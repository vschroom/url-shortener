package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"url-shortener/internal/config/cons"
	"url-shortener/internal/service"

	"github.com/go-chi/chi"

	"crypto/sha256"
	"net/url"
)

var store = make(map[int]string)

const shortUrlLength = 8

type Handler struct {
	SrvConsArg cons.ServerConsoleArg
	Storage    service.Storage
}

func (handler *Handler) UrlHandlerEncoder(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain" {
		rw.WriteHeader(http.StatusBadRequest)
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "error while reading body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		parsedUrl := string(body)
		urlHash := sha256.Sum256([]byte(parsedUrl))
		shortUrl := fmt.Sprintf("%x", urlHash)[:shortUrlLength]

		handler.Storage.StoreUrl(shortUrl, parsedUrl)

		resultUrl, err := url.JoinPath(handler.SrvConsArg.BaseShortAddr, shortUrl)
		if err != nil {
			log.Fatal(err)
		}

		rw.WriteHeader(http.StatusCreated)
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(resultUrl))
	}
}

func (handler *Handler) UrlHandlerDecoder(rw http.ResponseWriter, r *http.Request) {
	encodeUrl := chi.URLParam(r, "id")
	baseUrl := handler.Storage.GetUrl(encodeUrl)

	rw.Header().Set("Location", baseUrl)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
