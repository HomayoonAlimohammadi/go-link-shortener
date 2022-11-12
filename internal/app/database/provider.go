package database

import (
	"database/sql"
	"log"
)

type LinkShortenerProvider struct {
	Providers []Database
}

func (p *LinkShortenerProvider) GetUrl(url string) (string, error) {
	var err error
	var result string
	for _, provider := range p.Providers {
		result, err = provider.GetUrl(url)
		if result != "" {
			return result, nil
		}
	}
	if err == sql.ErrNoRows {
		return "", ErrUrlNotFound
	}
	return "", err
}

func (p *LinkShortenerProvider) GetToken(token string) (string, error) {
	var err error
	var result string
	for _, provider := range p.Providers {
		result, err = provider.GetToken(token)
		if result != "" {
			return result, nil
		}
	}
	if err == sql.ErrNoRows {
		return "", ErrTokenNotFound
	}
	return "", err
}

func (p *LinkShortenerProvider) Save(url, token string) error {
	var err error
	for _, provider := range p.Providers {
		err = provider.Save(url, token)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (p *LinkShortenerProvider) AddProvider(provider Database) {
	p.Providers = append(p.Providers, provider)
}
