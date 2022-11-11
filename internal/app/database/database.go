package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type PostgresProvider struct {
	Config PostgresConfig
	DB     *sql.DB
}

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
	SslMode  string
}

func (p *PostgresConfig) getConnectionString() string {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Name,
		p.SslMode,
	)
	return connectionString
}

func NewPostgresProvider(config PostgresConfig) (*PostgresProvider, error) {
	connectionString := config.getConnectionString()
	db, err := sql.Open("postgres", connectionString)
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
