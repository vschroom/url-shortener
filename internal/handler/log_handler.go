package handler

import (
	"net/http"
	"url-shortener/internal/logger"

	"time"

	"go.uber.org/zap"
)

func LoggerHandler(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h(&lw, r)

		duration := time.Since(start)

		logger.Log.Info("Log request and response",
			zap.String("Method", r.Method),
			zap.String("URI", r.RequestURI),
			zap.Duration("Duration", duration),
			zap.Int("Status", responseData.status),
			zap.Int("Response size", responseData.size),
		)
	})
}

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
