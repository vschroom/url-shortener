package main

import (
	"net/http"
	"url-shortener/internal/config/db"
	"url-shortener/internal/config/srv"
	"url-shortener/internal/handler"
	"url-shortener/internal/logger"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	"log"

	"errors"

	"github.com/go-chi/chi"
)

func main() {
	serverConfig := srv.InitServerConfig()
	storage := db.InitStore()

	logErr := logger.Initialize(serverConfig.LoggerLevel)
	if logErr != nil {
		log.Fatal(logErr)
	}

	reader, errReader := repository.NewUrlFileReader(serverConfig.FileStoragePath + "test.json")
	if errReader != nil {
		log.Fatal(errReader)
	}
	writer, errWriter := repository.NewUrlFileWriter(serverConfig.FileStoragePath + "test.json")
	if errWriter != nil {
		log.Fatal(errWriter)
	}
	h := handler.NewHandler(serverConfig, service.UrlService{
		Storage:       storage,
		UrlFileReader: *reader,
		UrlFileWriter: *writer,
	})

	router := chi.NewRouter()
	router.Get("/{id}", handler.LoggerHandler(handler.GzipHandler(h.UrlHandlerDecoder)))
	router.Post("/", handler.LoggerHandler(handler.GzipHandler(h.UrlHandlerEncoder)))
	router.Post("/api/shorten", handler.LoggerHandler(handler.GzipHandler(h.JsonUrlHandler)))

	err := http.ListenAndServe(serverConfig.Addr, router)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
