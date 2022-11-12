package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Needed by golang-migrate
	_ "github.com/lib/pq"                                // Needed by golang-migrate
)

type Database interface {
	GetUrl(string) (string, error)
	GetToken(string) (string, error)
	Save(string, string) error
}

type PostgresProvider struct {
	Config PostgresConfig
	DB     *sql.DB
}

type RedisProvider struct {
	Config RedisConfig
	DB     *redis.Client
}

type PostgresConfig struct {
	User           string
	Password       string
	Host           string
	Port           int
	Name           string
	SslMode        string
	MigrationsPath string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Timeout  int
}

func (p *PostgresConfig) getUrl() *url.URL {
	pgURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.User, p.Password),
		Path:   p.Host + ":" + strconv.Itoa(p.Port) + "/" + p.Name,
	}

	query := pgURL.Query()
	query.Add("sslmode", p.SslMode)
	pgURL.RawQuery = query.Encode()
	return pgURL
}

func NewPostgresProvider(config PostgresConfig) (*PostgresProvider, error) {
	pgURL := config.getUrl()
	db, err := sql.Open("postgres", pgURL.String())
	if err != nil {
		return nil, err
	}
	return &PostgresProvider{
		Config: config,
		DB:     db,
	}, nil
}

func (p *PostgresProvider) Close() {
	p.DB.Close()
}

func (p *PostgresProvider) GetUrl(token string) (string, error) {
	query := "SELECT url FROM linkshortener WHERE token=$1"
	var url string
	if err := p.DB.QueryRow(query, token).Scan(&url); err != nil {
		if err == sql.ErrNoRows {
			return "", ErrTokenNotFound
		}
		return "", err
	}
	return url, nil
}

func (p *PostgresProvider) GetToken(url string) (string, error) {
	query := "SELECT token FROM linkshortener WHERE url=$1"
	var token string
	if err := p.DB.QueryRow(query, url).Scan(&token); err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUrlNotFound
		}
		return "", err
	}
	return token, nil
}

func (p *PostgresProvider) Save(url, token string) error {
	insertStatement := "INSERT INTO linkshortener(url, token) VALUES ($1, $2)"
	_, err := p.DB.Exec(insertStatement, url, token)
	return err
}

func (p *PostgresProvider) Migrate() error {
	driver, err := postgres.WithInstance(p.DB, &postgres.Config{
		MultiStatementEnabled: true,
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		p.Config.MigrationsPath,
		"linkshortener", driver,
	)
	if err != nil {
		return err
	}

	// do the migration
	var migrationError error
	for migrationError == nil {
		migrationError = m.Steps(1)
	}

	log.Println("successfully applied the migrations")
	return nil
}

func NewRedisProvider(config RedisConfig) *RedisProvider {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	}
	db := redis.NewClient(opts)
	return &RedisProvider{
		Config: config,
		DB:     db,
	}
}

func (p *RedisProvider) GetUrl(token string) (string, error) {
	url, err := p.DB.Get(context.Background(), token).Result()
	if err == redis.Nil {
		return "", ErrUrlNotFound
	} else if err != nil {
		return "", err
	}
	return url, nil
}

func (p *RedisProvider) GetToken(url string) (string, error) {
	token, err := p.DB.Get(context.Background(), url).Result()
	if err == redis.Nil {
		return "", ErrTokenNotFound
	} else if err != nil {
		return "", err
	}
	return token, nil
}

func (p *RedisProvider) Save(url, token string) error {
	_, err := p.DB.Set(context.Background(), token, url, time.Duration(p.Config.Timeout)*time.Second).Result()
	if err != nil {
		return err
	}
	_, err = p.DB.Set(context.Background(), url, token, time.Duration(p.Config.Timeout)*time.Second).Result()
	if err != nil {
		return err
	}
	return nil
}
