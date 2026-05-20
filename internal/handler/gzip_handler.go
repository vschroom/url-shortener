package handler

import (
	"net/http"
	"strings"
	"url-shortener/internal/compress"
)

func GzipHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		originalResp := resp

		acceptEncoding := req.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := compress.NewCompressWriter(resp)
			originalResp = cw
			defer cw.Close()
		}

		contentEncoding := req.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := compress.NewCompressReader(req.Body)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(originalResp, req)
	}
}
