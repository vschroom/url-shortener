package repository

import (
	"encoding/json"
	"io"
	"os"
	"url-shortener/internal/logger"
	"url-shortener/internal/model"

	"path/filepath"

	"errors"

	"go.uber.org/zap"
)

type FileHolder struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewUrlFileHolder(filePath string, filename string) (*FileHolder, error) {
	path := filepath.Join("./", filePath)
	_, err := os.Stat(path)

	filename = filepath.Join(path, "/", filename)
	var file *os.File
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		f, err := os.OpenFile(filename, os.O_RDWR, 0777)
		if err != nil {
			return nil, err
		}
		file = f
	} else {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			return nil, errDir
		}
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return nil, err
		}
		file = f
	}

	jsonEncoder := json.NewEncoder(file)
	jsonDecoder := json.NewDecoder(file)

	return &FileHolder{
		file:    file,
		encoder: jsonEncoder,
		decoder: jsonDecoder,
	}, nil
}

func (fileHolder *FileHolder) StoreUrlInfo(urlInfo *model.UrlFileEntity) error {
	if info, err := fileHolder.file.Stat(); err == nil {
		if info.Size() == 0 {
			fileHolder.file.Write([]byte("[\r\n"))
			data, jsonErr := json.Marshal(&urlInfo)
			fileHolder.file.Write(data)
			fileHolder.file.Write([]byte("\r\n]"))
			return jsonErr
		}
	}

	data, jsonErr := json.Marshal(&urlInfo)
	if jsonErr != nil {
		return jsonErr
	}
	_, seekErr := fileHolder.file.Seek(-3, io.SeekEnd)
	if seekErr != nil {
		logger.Log.Error("Seek err", zap.Error(seekErr))
	}
	fileHolder.file.Write([]byte(",\r\n"))
	fileHolder.file.Write(data)
	fileHolder.file.Write([]byte("\r\n]"))

	return jsonErr
}

func (fileHolder *FileHolder) GetUrlInfo() (*[]model.UrlFileEntity, error) {
	fileHolder.file.Seek(0, io.SeekStart)
	var urlInfo []model.UrlFileEntity
	if err := fileHolder.decoder.Decode(&urlInfo); err != nil {
		return nil, err
	}

	return &urlInfo, nil
}

func (fileHolder *FileHolder) Close() error {
	return fileHolder.file.Close()
}

func (fileHolder *FileHolder) Remove() error {
	return os.Remove(fileHolder.file.Name())
}
