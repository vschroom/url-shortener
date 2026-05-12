package db

import "url-shortener/internal/repository"

func InitStore() repository.Storage {
	store := make(map[string]string)
	return repository.Storage{
		Store: store,
	}
}
