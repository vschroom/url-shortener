package handler

import (
	"io"
	"log"
	"net/http"

	"url-shortener/internal/config/cons"
	"url-shortener/internal/service"

	"github.com/go-chi/chi"

	"net/url"
)

const shortUrlLength = 8

type Handler struct {
	SrvConsArg cons.ServerConsoleArg
	UrlService service.UrlService
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

		shortUrl, shortUrlErr := handler.UrlService.StoreUrl(parsedUrl)
		if shortUrlErr != nil {
			log.Default().Println(shortUrlErr.Error())

			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			resultUrl, err := url.JoinPath(handler.SrvConsArg.BaseShortAddr, shortUrl)
			if err != nil {
				log.Default().Println(err.Error())
				rw.WriteHeader(http.StatusInternalServerError)
			} else {
				rw.WriteHeader(http.StatusCreated)
				rw.Header().Set("Content-Type", "text/plain")
				rw.Write([]byte(resultUrl))
			}
		}
	}
}

func (handler *Handler) UrlHandlerDecoder(rw http.ResponseWriter, r *http.Request) {
	encodeUrl := chi.URLParam(r, "id")
	baseUrl := handler.UrlService.GetUrl(encodeUrl)

	rw.Header().Set("Location", baseUrl)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
