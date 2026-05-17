package main

import (
	"net/http"
	"url-shortener/internal/config/db"
	"url-shortener/internal/config/srv"
	"url-shortener/internal/handler"
	"url-shortener/internal/logger"
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

	h := handler.NewHandler(serverConfig, service.UrlService{Storage: storage})

	router := chi.NewRouter()
	router.Get("/{id}", handler.LoggerHandler(h.UrlHandlerDecoder))
	router.Post("/", handler.LoggerHandler(h.UrlHandlerEncoder))
	router.Post("/api/shorten", handler.LoggerHandler(h.JsonUrlHandler))

	err := http.ListenAndServe(serverConfig.Addr, router)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
