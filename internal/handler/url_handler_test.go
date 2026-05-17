package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"url-shortener/internal/config/db"
	"url-shortener/internal/config/srv"
	"url-shortener/internal/service"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"regexp"

	"github.com/go-resty/resty/v2"
)

func TestUrlHandlerEncoder(t *testing.T) {
	serverConfig := srv.InitServerConfig()
	storage := db.InitStore()
	h := NewHandler(serverConfig, service.UrlService{Storage: storage})

	urlHandlerEncoder := http.HandlerFunc(h.UrlHandlerEncoder)
	server := httptest.NewServer(urlHandlerEncoder)

	defer server.Close()

	type want struct {
		code               int
		request            string
		requestMethod      string
		requestContentType string
		responseRegexp     string
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
				responseRegexp:     "^" + serverConfig.BaseShortAddr + "[a-zA-Z0-9]{8}$",
				contentType:        "text/plain; charset=utf-8",
			},
		},
		{
			name: "#2 should return created for https://ya.ru/",
			want: want{
				code:               201,
				request:            "https://ya.ru/",
				requestMethod:      http.MethodPost,
				requestContentType: "text/plain",
				responseRegexp:     "^" + serverConfig.BaseShortAddr + "[a-zA-Z0-9]{8}$",
				contentType:        "text/plain; charset=utf-8",
			},
		},
		{
			name: "#3 should return bad request with wrong http method",
			want: want{
				code:               400,
				request:            "https://ya.ru/",
				requestMethod:      http.MethodGet,
				requestContentType: "text/plain; charset=utf-8",
				responseRegexp:     "",
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
				responseRegexp:     "",
				contentType:        "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := resty.New().R().
				SetHeader("Content-Type", test.want.requestContentType).
				SetBody(strings.NewReader(test.want.request)).
				Post(server.URL)

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, test.want.code, resp.StatusCode())

			rowBody := resp.Body()
			assert.Regexp(t, regexp.MustCompile(test.want.responseRegexp), string(rowBody))
			assert.Equal(t, test.want.contentType, resp.Header().Get("Content-Type"))
		})
	}
}

func TestUrlHandlerDecoder(t *testing.T) {
	storage := db.InitStore()
	h := &Handler{
		urlService: service.UrlService{Storage: storage},
	}

	urlHandlerDecoder := http.HandlerFunc(h.UrlHandlerDecoder)
	server := httptest.NewServer(urlHandlerDecoder)

	defer server.Close()

	type want struct {
		code               int
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
			name: "#1 should return 307 for https://yandex.ru/",
			want: want{
				code:               307,
				requestMethod:      http.MethodGet,
				requestContentType: "text/plain",
				headerLocation:     "https://yandex.ru/",
				contentType:        "text/plain",
			},
		},
		{
			name: "#2 should return 307 for https://ya.ru/",
			want: want{
				code:               307,
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
			rdr := httptest.NewRecorder()
			h.UrlHandlerEncoder(rdr, postReq)
			url, err := io.ReadAll(rdr.Body)
			assert.NoError(t, err)

			request := httptest.NewRequest(test.want.requestMethod, "http://localhost:8080/"+string(url), nil)
			request.Header.Set("Content-Type", test.want.requestContentType)

			// создаём новый Recorder
			w := httptest.NewRecorder()
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

func TestJsonUrlHandler(t *testing.T) {
	storage := db.InitStore()
	h := &Handler{
		urlService: service.UrlService{Storage: storage},
	}

	r := chi.NewRouter()
	r.Post("/api/shorten", h.JsonUrlHandler)
	server := httptest.NewServer(r)

	defer server.Close()

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "#1 method not allowed for GET",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "#2 method not allowed for PUT",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "#3 method not allowed for DELETE",
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "#4 method POST without body",
			method:       http.MethodPost,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
		{
			name:         "#5 method_post_unsupported_type",
			method:       http.MethodPost,
			body:         `[{"url":""}]`,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
		{
			name:         "method_post_success",
			method:       http.MethodPost,
			body:         `{"url":"https://ya.ru/"}`,
			expectedCode: http.StatusCreated,
			expectedBody: `{"result":""}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, server.URL+"/api/shorten", strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			resp := w.Result()

			assert.Equal(t, tc.expectedCode, resp.StatusCode, "Response code didn't match expected")
			if tc.expectedBody != "" {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err, "error read json body")
				assert.NotNil(t, tc.expectedBody, string(body))
			}
		})
	}
}
