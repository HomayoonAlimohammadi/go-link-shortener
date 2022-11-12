package database

import (
	"database/sql"
	"errors"
	"log"
	"net/url"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Needed by golang-migrate
	_ "github.com/lib/pq"                                // Needed by golang-migrate
)

type PostgresProvider struct {
	Config PostgresConfig
	DB     *sql.DB
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

func (p *PostgresConfig) getURL() *url.URL {
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
	pgURL := config.getURL()
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
			return "", errors.New("no matching rows for the given token")
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
			return "", errors.New("no matching rows for the given url")
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
