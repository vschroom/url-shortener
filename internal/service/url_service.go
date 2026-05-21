package service

import (
	"errors"
	"math/rand/v2"
	"url-shortener/internal/logger"
	"url-shortener/internal/model"
	"url-shortener/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UrlService struct {
	Storage       repository.Storage
	UrlFileReader repository.UrlFileReader
	UrlFileWriter repository.UrlFileWriter
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

		urlInfo := &model.UrlFileEntity{
			Id:          uuid.New(),
			ShortUrl:    randShortUrl,
			OriginalUrl: baseUrl,
		}

		err := urlService.UrlFileWriter.StoreUrlInfo(urlInfo)
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
	urlInfo, err := urlService.UrlFileReader.GetUrlInfo()
	if err != nil {
		panic(err)
	}

	for _, info := range *urlInfo {
		if target := info.ShortUrl; target == shortUrl {
			return info.OriginalUrl
		}
	}
	return ""
}

func randomString(length int) string {
	rowBytes := make([]byte, length)

	for i := range rowBytes {
		randByteIdx := rand.IntN(len(charset))
		rowBytes[i] = charset[randByteIdx]
	}

	return string(rowBytes)
}
