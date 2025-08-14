package tests

import (
	"context"
	"sync"

	"github.com/stretchr/testify/mock"
)

type MockRedisStorage struct {
	mock.Mock
	urls map[string]string
	sync.RWMutex
}

func NewMockRedisStorage() *MockRedisStorage {
	return &MockRedisStorage{
		urls: make(map[string]string),
	}
}

// StoreURL mocks storing a shortCode → originalURL mapping
func (m *MockRedisStorage) StoreURL(ctx context.Context, shortCode, originalURL string) error {
	m.Lock()
	defer m.Unlock()

	args := m.Called(ctx, shortCode, originalURL)

	if args.Error(0) == nil {
		m.urls[shortCode] = originalURL
	}

	return args.Error(0)
}

// GetURL mocks retrieving the originalURL for a given shortCode
func (m *MockRedisStorage) GetURL(ctx context.Context, shortCode string) (string, error) {
	m.RLock()
	defer m.RUnlock()

	args := m.Called(ctx, shortCode)

	// If the in-memory map has it, return it (still honoring the mocked error)
	if url, exists := m.urls[shortCode]; exists {
		return url, args.Error(1)
	}

	// Otherwise return the mocked return value
	if args.Get(0) != nil {
		return args.String(0), args.Error(1)
	}

	return "", args.Error(1)
}

// DeleteURL mocks removing a shortCode mapping
func (m *MockRedisStorage) DeleteURL(ctx context.Context, shortCode string) error {
	m.Lock()
	defer m.Unlock()

	args := m.Called(ctx, shortCode)

	if args.Error(0) == nil {
		delete(m.urls, shortCode)
	}

	return args.Error(0)
}

// Close mocks closing the Redis connection
func (m *MockRedisStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}
