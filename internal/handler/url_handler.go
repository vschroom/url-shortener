package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"url-shortener/internal/logger"
	"url-shortener/internal/model"

	"url-shortener/internal/config/srv"
	"url-shortener/internal/service"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

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
			logger.Log.Info(shortUrlErr.Error())

			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			resultUrl, err := url.JoinPath(handler.serverConfig.BaseShortAddr, shortUrl)
			if err != nil {
				logger.Log.Info(err.Error())
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
	baseUrl, err := handler.urlService.GetUrl(encodeUrl)
	if err != nil {
		http.Error(rw, "error while getting result url", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Location", baseUrl)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func (handler *Handler) JsonUrlHandler(rw http.ResponseWriter, r *http.Request) {
	var req model.Request
	jsonDec := json.NewDecoder(r.Body)
	if err := jsonDec.Decode(&req); err != nil {
		logger.Log.Error("Cannot parse request json body", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortUrl, shortUrlErr := handler.urlService.StoreUrl(req.Url)
	if shortUrlErr != nil {
		logger.Log.Error(shortUrlErr.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultUrl, joinPathErr := url.JoinPath(handler.serverConfig.BaseShortAddr, shortUrl)
	if joinPathErr != nil {
		logger.Log.Error(joinPathErr.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	resp := model.Response{Result: resultUrl}
	jsonEnc := json.NewEncoder(rw)
	jsonEnc.Encode(resp)

	logger.Log.Info("Successfully processed")
}
