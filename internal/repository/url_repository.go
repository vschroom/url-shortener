package repository

import (
	"fmt"
)

type ShortUrlDuplicateKeyError struct {
	ShortUrl string
	BaseUrl  string
}

func (err ShortUrlDuplicateKeyError) Error() string {
	return fmt.Sprintf("Short Url = %s already exists for base Url = %s", err.ShortUrl, err.BaseUrl)
}

func (e ShortUrlDuplicateKeyError) Is(err error) bool {
	_, is := err.(ShortUrlDuplicateKeyError)

	return is
}

type Storage struct {
	Store map[string]string
}

func (storage *Storage) StoreUrl(shortUrl string, baseUrl string) error {
	existBaseUrl := storage.Store[shortUrl]
	if existBaseUrl != "" && existBaseUrl != baseUrl {
		return ShortUrlDuplicateKeyError{
			ShortUrl: shortUrl,
			BaseUrl:  baseUrl,
		}
	}

	storage.Store[shortUrl] = baseUrl

	return nil
}

func (storage *Storage) GetUrl(shortUrl string) string {
	return storage.Store[shortUrl]
}
