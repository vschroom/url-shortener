package main

import (
	"net/http"
	"url-shortener/internal/config/cons"
	"url-shortener/internal/handler"

	"github.com/go-chi/chi"
)

func main() {
	cons.ParseServerFlags()
	serverArgs := cons.ServerConsoleArg

	router := chi.NewRouter()
	router.Get("/{id}", handler.UrlHandlerDecoder)
	router.Post("/", handler.UrlHandlerEncoder)

	err := http.ListenAndServe(serverArgs.Addr, router)
	if err != nil {
		panic(err)
	}
}
