package main

import (
	"net/http"
	"url-shortener/internal/handler"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()
	router.Get("/{id}", handler.UrlHandlerDecoder)
	router.Post("/", handler.UrlHandlerEncoder)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
