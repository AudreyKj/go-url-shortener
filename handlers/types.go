package handlers

import (
	"go-url-shortner/services"
)

type URLHandler struct {
	urlService services.URLServiceInterface
}

func NewURLHandler(urlService services.URLServiceInterface) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}
