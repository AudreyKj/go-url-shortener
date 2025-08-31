package services

import (
	"context"
)

type URLServiceInterface interface {
	CreateShortURL(ctx context.Context, req URLRequest) (*URLResponse, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
}

type AISlugServiceInterface interface {
	GenerateSlug(ctx context.Context, originalURL string) (string, error)
}

type StorageInterface interface {
	StoreURL(ctx context.Context, shortCode, originalURL string) error
	GetURL(ctx context.Context, shortCode string) (string, error)
	Close() error
}
