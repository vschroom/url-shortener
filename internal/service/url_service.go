package service

import (
	"errors"
	"math/rand/v2"
	"url-shortener/internal/logger"
	"url-shortener/internal/repository"

	"go.uber.org/zap"
)

type UrlService struct {
	Storage repository.Storage
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortUrlLength = 8
const maxShortUrlRetryCount = 5

func (urlService *UrlService) StoreUrl(baseUrl string) (string, error) {
	randShortUrl := randomString(shortUrlLength)
	for n := range maxShortUrlRetryCount {
		logger.Log.Info(
			"Try to store short Url",
			zap.Int("Try #", n+1),
		)
		err := urlService.Storage.StoreUrl(randShortUrl, baseUrl)
		if err != nil && errors.Is(err, repository.ErrShortUrlDuplicateKey) {
			logger.Log.Info("Short Url duplicate for different base urls. Try to generate another one")
			randShortUrl = randomString(shortUrlLength)
		} else if err != nil {
			return "", err
		} else {
			logger.Log.Info("Short Url successfully generated and store")
			return randShortUrl, nil
		}
	}

	return "", errors.New("Short Url generation failed")
}

func (urlService *UrlService) GetUrl(shortUrl string) string {
	return urlService.Storage.GetUrl(shortUrl)
}

func randomString(length int) string {
	rowBytes := make([]byte, length)

	for i := range rowBytes {
		randByteIdx := rand.IntN(len(charset))
		rowBytes[i] = charset[randByteIdx]
	}

	return string(rowBytes)
}
