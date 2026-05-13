package main

import (
	"net/http"
	"url-shortener/internal/config/cons"
	"url-shortener/internal/config/db"
	"url-shortener/internal/handler"
	"url-shortener/internal/service"

	"log"

	"errors"

	"github.com/go-chi/chi"
)

func main() {
	serverArgs := cons.ParseServerFlags()
	storage := db.InitStore()

	h := handler.NewHandler(serverArgs, service.UrlService{Storage: storage})

	router := chi.NewRouter()
	router.Get("/{id}", h.UrlHandlerDecoder)
	router.Post("/", h.UrlHandlerEncoder)

	err := http.ListenAndServe(serverArgs.Addr, router)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
