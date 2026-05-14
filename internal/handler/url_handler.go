package handler

import (
	"io"
	"log"
	"net/http"

	"url-shortener/internal/config/srv"
	"url-shortener/internal/service"

	"github.com/go-chi/chi"

	"net/url"
)

type Handler struct {
	serverConfig srv.ServerConfig
	urlService   service.UrlService
}

func NewHandler(srvConsArg srv.ServerConfig, urlService service.UrlService) *Handler {
	return &Handler{
		serverConfig: srvConsArg,
		urlService:   urlService,
	}
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

		shortUrl, shortUrlErr := handler.urlService.StoreUrl(parsedUrl)
		if shortUrlErr != nil {
			log.Default().Println(shortUrlErr.Error())

			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			resultUrl, err := url.JoinPath(handler.serverConfig.BaseShortAddr, shortUrl)
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
	baseUrl := handler.urlService.GetUrl(encodeUrl)

	rw.Header().Set("Location", baseUrl)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
