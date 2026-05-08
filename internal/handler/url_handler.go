package handler

import (
	"io"
	"net/http"

	"url-shortener/internal/config/cons"

	"github.com/go-chi/chi"
	"github.com/speps/go-hashids/v2"
)

var store = make(map[int]string)

const CustomSalt string = "randomsalt"
const ShortUrlLength = 8

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
		count := len(store) + 1
		store[count] = parsedUrl

		hd := hashids.NewData()
		hd.Salt = CustomSalt
		hd.MinLength = ShortUrlLength
		h, _ := hashids.NewWithData(hd)
		encodeUrl, _ := h.Encode([]int{count})

		resultUrl := cons.ServerConsoleArg.BaseShortAddr + "/" + encodeUrl

		rw.WriteHeader(http.StatusCreated)
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(resultUrl))
	}
}

func UrlHandlerDecoder(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusBadRequest)
	} else {
		hd := hashids.NewData()
		hd.Salt = CustomSalt
		hd.MinLength = ShortUrlLength
		h, _ := hashids.NewWithData(hd)

		encodeUrl := chi.URLParam(r, "id")
		d, _ := h.DecodeWithError(encodeUrl)
		key := d[0]

		rw.Header().Set("Location", store[key])
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}
