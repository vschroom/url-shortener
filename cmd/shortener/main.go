package main

import (
	"net/http"
	"url-shortener/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.UrlHandlerEncoder)
	mux.HandleFunc("/{id}", handler.UrlHandlerDecoder)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
