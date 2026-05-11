package service

type Storage struct {
	Store map[string]string
}

func (storage *Storage) StoreUrl(shortUrl string, baseUrl string) {
	storage.Store[shortUrl] = baseUrl
}

func (storage *Storage) GetUrl(shortUrl string) string {
	return storage.Store[shortUrl]
}
