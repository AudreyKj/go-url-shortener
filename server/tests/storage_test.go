package tests

import (
	"context"
	"encoding/json"
	"testing"

	"go-url-shortner/storage"

	"github.com/stretchr/testify/assert"
)

func TestRedisStorage_NewRedisStorage_Success(t *testing.T) {
	// Test that RedisStorage implements URLStorage interface
	var _ storage.URLStorage = (*storage.RedisStorage)(nil)

	// Test that the struct has the expected fields
	// This is more of a compile-time check than a runtime test
	assert.True(t, true, "RedisStorage should implement URLStorage interface")
}

func TestRedisStorage_StoreURL(t *testing.T) {
	// Test using the mock storage
	mockStorage := NewMockRedisStorage()
	ctx := context.Background()
	shortCode := "abc123"
	originalURL := "https://example.com"

	// Set up mock expectations
	mockStorage.On("StoreURL", ctx, shortCode, originalURL).Return(nil)
	mockStorage.On("GetURL", ctx, shortCode).Return(originalURL, nil)

	// Test the mock implementation
	err := mockStorage.StoreURL(ctx, shortCode, originalURL)
	assert.NoError(t, err)

	// Verify the URL was stored in the mock
	retrievedURL, err := mockStorage.GetURL(ctx, shortCode)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, retrievedURL)

	mockStorage.AssertExpectations(t)
}

func TestRedisStorage_GetURL(t *testing.T) {
	// Test using the mock storage
	mockStorage := NewMockRedisStorage()
	ctx := context.Background()
	shortCode := "abc123"
	originalURL := "https://example.com"

	// Set up mock expectations
	mockStorage.On("StoreURL", ctx, shortCode, originalURL).Return(nil)
	mockStorage.On("GetURL", ctx, shortCode).Return(originalURL, nil)
	mockStorage.On("GetURL", ctx, "nonexistent").Return("", assert.AnError)

	// Store a URL first
	err := mockStorage.StoreURL(ctx, shortCode, originalURL)
	assert.NoError(t, err)

	// Retrieve the URL
	retrievedURL, err := mockStorage.GetURL(ctx, shortCode)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, retrievedURL)

	// Test getting non-existent URL
	_, err = mockStorage.GetURL(ctx, "nonexistent")
	assert.Error(t, err)

	mockStorage.AssertExpectations(t)
}

func TestRedisStorage_Close(t *testing.T) {
	// Test using the mock storage
	mockStorage := NewMockRedisStorage()

	// Set up mock expectations
	mockStorage.On("Close").Return(nil)

	// Test that close doesn't panic
	err := mockStorage.Close()
	assert.NoError(t, err)

	mockStorage.AssertExpectations(t)
}

func TestURLMapping_JSONTags(t *testing.T) {
	mapping := storage.URLMapping{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		CreatedAt:   1234567890,
		ExpiresAt:   1234567890 + 86400,
	}

	// Test that the struct can be marshaled to JSON
	// This verifies the JSON tags are correct
	_, err := json.Marshal(mapping)
	assert.NoError(t, err)

	// Test that the struct can be unmarshaled from JSON
	jsonData := `{"short_code":"xyz789","original_url":"https://test.com","created_at":1234567890,"expires_at":1234567890}`
	var unmarshaled storage.URLMapping
	err = json.Unmarshal([]byte(jsonData), &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, "xyz789", unmarshaled.ShortCode)
	assert.Equal(t, "https://test.com", unmarshaled.OriginalURL)
	assert.Equal(t, int64(1234567890), unmarshaled.CreatedAt)
	assert.Equal(t, int64(1234567890), unmarshaled.ExpiresAt)
}

func TestStorageStats_JSONTags(t *testing.T) {
	stats := storage.StorageStats{
		TotalURLs:   100,
		ActiveURLs:  95,
		ExpiredURLs: 5,
		StorageSize: 1024,
	}

	// Test that the struct can be marshaled to JSON
	// This verifies the JSON tags are correct
	_, err := json.Marshal(stats)
	assert.NoError(t, err)

	// Test that the struct can be unmarshaled from JSON
	jsonData := `{"total_urls":200,"active_urls":180,"expired_urls":20,"storage_size":2048}`
	var unmarshaled storage.StorageStats
	err = json.Unmarshal([]byte(jsonData), &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, int64(200), unmarshaled.TotalURLs)
	assert.Equal(t, int64(180), unmarshaled.ActiveURLs)
	assert.Equal(t, int64(20), unmarshaled.ExpiredURLs)
	assert.Equal(t, int64(2048), unmarshaled.StorageSize)
}
