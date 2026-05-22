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

	holder, errHolder := repository.NewUrlFileHolder(serverConfig.FileStoragePath, "test.json")
	if errHolder != nil {
		log.Fatal(errHolder)
	}
	h := handler.NewHandler(serverConfig, service.UrlService{
		Storage:       storage,
		UrlFileHolder: *holder,
	})

	defer holder.Close()

	router := chi.NewRouter()
	router.Get("/{id}", handler.LoggerHandler(handler.GzipHandler(h.UrlHandlerDecoder)))
	router.Post("/", handler.LoggerHandler(handler.GzipHandler(h.UrlHandlerEncoder)))
	router.Post("/api/shorten", handler.LoggerHandler(handler.GzipHandler(h.JsonUrlHandler)))

	err := http.ListenAndServe(serverConfig.Addr, router)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
