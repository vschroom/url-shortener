package handler

import (
	"encoding/base64"
	"io"
	"net/http"
)

func UrlHandlerEncoder(rw http.ResponseWriter, r *http.Request) {
	if !(r.Method == http.MethodPost && r.Header.Get("Content-Type") == "text/plain") {
		rw.WriteHeader(http.StatusBadRequest)
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "error while reading body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		parsedUrl := string(body)
		encodeUrl := "http://" + r.Host + base64.URLEncoding.EncodeToString([]byte(parsedUrl))

		rw.WriteHeader(http.StatusCreated)
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(encodeUrl))
	}
}

func UrlHandlerDecoder(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusBadRequest)
	} else {
		encodeUrl := r.PathValue("id")
		originalUrl, err := base64.URLEncoding.DecodeString(encodeUrl)
		if err != nil {
			panic(err)
		}

		rw.Header().Set("Location", string(originalUrl))
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}
