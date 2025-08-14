package storage

import (
	"context"
)

type URLStorage interface {
	StoreURL(ctx context.Context, shortCode, originalURL string) error
	GetURL(ctx context.Context, shortCode string) (string, error)
	DeleteURL(ctx context.Context, shortCode string) error
	Close() error
}

type URLMapping struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	CreatedAt   int64  `json:"created_at"`
	ExpiresAt   int64  `json:"expires_at"`
}

type StorageStats struct {
	TotalURLs   int64 `json:"total_urls"`
	ActiveURLs  int64 `json:"active_urls"`
	ExpiredURLs int64 `json:"expired_urls"`
	StorageSize int64 `json:"storage_size"`
}
