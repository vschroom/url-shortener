package repository

import (
	"errors"
)

var ErrShortUrlDuplicateKey = errors.New("short URL already exists")

type Storage struct {
	Store map[string]string
}

func (storage *Storage) StoreUrl(shortUrl string, baseUrl string) error {
	existBaseUrl := storage.Store[shortUrl]
	if existBaseUrl != "" && existBaseUrl != baseUrl {
		return ErrShortUrlDuplicateKey
	}

	storage.Store[shortUrl] = baseUrl

	return nil
}

func (storage *Storage) GetUrl(shortUrl string) string {
	return storage.Store[shortUrl]
}
