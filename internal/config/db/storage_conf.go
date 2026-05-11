package db

import "url-shortener/internal/service"

func InitStore() service.Storage {
	store := make(map[string]string)
	return service.Storage{
		Store: store,
	}
}
