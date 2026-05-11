package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"url-shortener/internal/config/cons"
	"url-shortener/internal/config/db"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUrlHandlerEncoder(t *testing.T) {
	// init
	serverConsoleArgs := cons.ParseServerFlags()
	storage := db.InitStore()
	h := &Handler{
		SrvConsArg: serverConsoleArgs,
		Storage:    storage,
	}

	type want struct {
		code               int
		request            string
		requestMethod      string
		requestContentType string
		response           string
		contentType        string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "#1 should return created for https://yandex.ru/",
			want: want{
				code:               201,
				request:            "https://yandex.ru/",
				requestMethod:      http.MethodPost,
				requestContentType: "text/plain",
				response:           "http://localhost:8080/77fca595",
				contentType:        "text/plain",
			},
		},
		{
			name: "#2 should return created for https://ya.ru/",
			want: want{
				code:               201,
				request:            "https://ya.ru/",
				requestMethod:      http.MethodPost,
				requestContentType: "text/plain",
				response:           "http://localhost:8080/e12f5f6c",
				contentType:        "text/plain",
			},
		},
		{
			name: "#3 should return bad request with wrong http method",
			want: want{
				code:               400,
				request:            "https://ya.ru/",
				requestMethod:      http.MethodGet,
				requestContentType: "text/plain",
				response:           "",
				contentType:        "",
			},
		},
		{
			name: "#4 should return bad request with wrong content type header",
			want: want{
				code:               400,
				request:            "https://ya.ru/",
				requestMethod:      http.MethodPost,
				requestContentType: "application/json",
				response:           "",
				contentType:        "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.want.requestMethod, "http://localhost:8080/", strings.NewReader(test.want.request))
			request.Header.Set("Content-Type", test.want.requestContentType)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.UrlHandlerEncoder(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, w.Header().Get("Content-Type"))
		})
	}
}

func TestUrlHandlerDecoder(t *testing.T) {
	serverConsoleArgs := cons.ParseServerFlags()
	storage := db.InitStore()
	h := &Handler{
		SrvConsArg: serverConsoleArgs,
		Storage:    storage,
	}

	type want struct {
		code               int
		request            string
		requestMethod      string
		requestContentType string
		headerLocation     string
		contentType        string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "#1 should return 307 for 77fca595",
			want: want{
				code:               307,
				request:            "77fca595",
				requestMethod:      http.MethodGet,
				requestContentType: "text/plain",
				headerLocation:     "https://yandex.ru/",
				contentType:        "text/plain",
			},
		},
		{
			name: "#2 should return 307 for e12f5f6c",
			want: want{
				code:               307,
				request:            "e12f5f6c",
				requestMethod:      http.MethodGet,
				requestContentType: "text/plain",
				headerLocation:     "https://ya.ru/",
				contentType:        "text/plain",
			},
		},
		{
			name: "#3 should return bad request with wrong http method",
			want: want{
				code:               405,
				request:            "lGBPMKLe",
				requestMethod:      http.MethodPost,
				requestContentType: "text/plain",
				headerLocation:     "",
				contentType:        "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			postReq := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader(test.want.headerLocation))
			postReq.Header.Set("Content-Type", "text/plain")
			h.UrlHandlerEncoder(httptest.NewRecorder(), postReq)

			request := httptest.NewRequest(test.want.requestMethod, "/"+test.want.request, nil)
			// request.SetPathValue("id", test.want.request)
			request.Header.Set("Content-Type", test.want.requestContentType)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// UrlHandlerDecoder(w, request)
			r := chi.NewRouter()
			r.Get("/{id}", h.UrlHandlerDecoder)
			r.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.headerLocation, w.Header().Get("Location"))
		})
	}
}
