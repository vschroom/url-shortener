package repository

import (
	"encoding/json"
	"io"
	"os"
	"url-shortener/internal/logger"
	"url-shortener/internal/model"

	"go.uber.org/zap"
)

type UrlFileWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func NewUrlFileWriter(filename string) (*UrlFileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	jsonEncoder := json.NewEncoder(file)

	return &UrlFileWriter{
		file:    file,
		encoder: jsonEncoder,
	}, nil
}

func (urlFileWriter *UrlFileWriter) StoreUrlInfo(urlInfo *model.UrlFileEntity) error {
	if info, err := urlFileWriter.file.Stat(); err == nil {
		if info.Size() == 0 {
			urlFileWriter.file.Write([]byte("[\r\n"))
			data, jsonErr := json.Marshal(&urlInfo)
			urlFileWriter.file.Write(data)
			urlFileWriter.file.Write([]byte("\r\n]"))
			return jsonErr
		}
	}

	data, jsonErr := json.Marshal(&urlInfo)
	if jsonErr != nil {
		return jsonErr
	}
	_, seekErr := urlFileWriter.file.Seek(-3, io.SeekEnd)
	if seekErr != nil {
		logger.Log.Error("Seek err", zap.Error(seekErr))
	}
	urlFileWriter.file.Write([]byte(",\r\n"))
	urlFileWriter.file.Write(data)
	urlFileWriter.file.Write([]byte("\r\n]"))

	return jsonErr
}

func (urlFileWriter *UrlFileWriter) Close() error {
	return urlFileWriter.file.Close()
}

func (urlFileWriter *UrlFileWriter) Remove() error {
	return os.Remove(urlFileWriter.file.Name())
}

type UrlFileReader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewUrlFileReader(filename string) (*UrlFileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	jsonDecoder := json.NewDecoder(file)

	return &UrlFileReader{
		file:    file,
		decoder: jsonDecoder,
	}, nil
}

func (c *UrlFileReader) GetUrlInfo() (*[]model.UrlFileEntity, error) {
	var urlInfo []model.UrlFileEntity
	if err := c.decoder.Decode(&urlInfo); err != nil {
		return nil, err
	}

	return &urlInfo, nil
}

func (c *UrlFileReader) Close() error {
	return c.file.Close()
}

func (c *UrlFileReader) Remove() error {
	return os.Remove(c.file.Name())
}
