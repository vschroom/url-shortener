package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

// для сжатия данных, реализует http.ResponseWriter
type compressWriter struct {
	writer     http.ResponseWriter
	gzipWriter *gzip.Writer
}

func NewCompressWriter(writer http.ResponseWriter) *compressWriter {
	return &compressWriter{
		writer:     writer,
		gzipWriter: gzip.NewWriter(writer),
	}
}

func (compressWriter *compressWriter) Header() http.Header {
	return compressWriter.writer.Header()
}

func (compressWriter *compressWriter) Write(p []byte) (int, error) {
	return compressWriter.gzipWriter.Write(p)
}

func (compressWriter *compressWriter) WriteHeader(statusCode int) {
	compressWriter.writer.Header().Set("Content-Encoding", "gzip")
	compressWriter.writer.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (compressWriter *compressWriter) Close() error {
	return compressWriter.gzipWriter.Close()
}

// распаковывает сжатые данные, реализует io.ReadCloser
type compressReader struct {
	readCloser io.ReadCloser
	gzipReader *gzip.Reader
}

func NewCompressReader(readCloser io.ReadCloser) (*compressReader, error) {
	gzipReader, err := gzip.NewReader(readCloser)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		readCloser: readCloser,
		gzipReader: gzipReader,
	}, nil
}

func (compressReader compressReader) Read(data []byte) (n int, err error) {
	return compressReader.gzipReader.Read(data)
}

func (compressReader *compressReader) Close() error {
	if err := compressReader.readCloser.Close(); err != nil {
		return err
	}
	return compressReader.gzipReader.Close()
}
