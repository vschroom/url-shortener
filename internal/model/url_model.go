package model

import (
	"github.com/google/uuid"
)

type Request struct {
	Url string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type UrlFileEntity struct {
	Id          uuid.UUID `json:"id"`
	ShortUrl    string    `json:"short_url"`
	OriginalUrl string    `json:"original_url"`
}
