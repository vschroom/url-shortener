package main

import (
	"net/http"
	"url-shortener/internal/config/db"
	"url-shortener/internal/config/srv"
	"url-shortener/internal/handler"
	"url-shortener/internal/service"

	"log"

	"errors"

	"github.com/go-chi/chi"
)

func main() {
	serverConfig := srv.InitServerConfig()
	storage := db.InitStore()

	h := handler.NewHandler(serverConfig, service.UrlService{Storage: storage})

	router := chi.NewRouter()
	router.Get("/{id}", h.UrlHandlerDecoder)
	router.Post("/", h.UrlHandlerEncoder)

	err := http.ListenAndServe(serverConfig.Addr, router)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
