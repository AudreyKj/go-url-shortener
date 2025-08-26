package tests

import (
	"context"
	"testing"

	"go-url-shortner/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestURLService_CreateShortURL_WithAI(t *testing.T) {
	// Setup
	mockStorage := NewMockRedisStorage()
	mockAI := new(MockAISlugService)

	service := services.NewURLService(mockStorage, mockAI, "localhost", "8080")

	req := services.URLRequest{URL: "https://github.com"}
	aiSlug := "ghub"

	// Mock AI service returning a slug
	mockAI.On("GenerateSlug", mock.Anything, req.URL).Return(aiSlug, nil)

	// Ensure AI slug is not empty
	assert.NotEmpty(t, aiSlug)

	// Mock storage check - slug is available
	mockStorage.On("GetURL", mock.Anything, aiSlug).Return("", assert.AnError)

	// Mock storage store
	mockStorage.On("StoreURL", mock.Anything, aiSlug, req.URL).Return(nil)

	// Execute
	response, err := service.CreateShortURL(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.URL, response.OriginalURL)
	assert.Equal(t, aiSlug, response.ShortCode)
	assert.Equal(t, "ai_generated", response.SlugType)
	assert.Equal(t, "http://localhost:8080/"+aiSlug, response.ShortURL)

	mockStorage.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

func TestURLService_CreateShortURL_WithoutAI(t *testing.T) {
	// Setup
	mockStorage := NewMockRedisStorage()

	service := services.NewURLService(mockStorage, nil, "localhost", "8080")

	req := services.URLRequest{URL: "https://example.com"}

	// Mock storage store
	mockStorage.On("StoreURL", mock.Anything, mock.AnythingOfType("string"), req.URL).Return(nil)

	// Execute
	response, err := service.CreateShortURL(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.URL, response.OriginalURL)
	assert.NotEmpty(t, response.ShortCode)
	assert.Equal(t, "hash_based", response.SlugType)

	mockStorage.AssertExpectations(t)
}

func TestURLService_CreateShortURL_EmptyURL(t *testing.T) {
	// Setup
	mockStorage := NewMockRedisStorage()
	mockAI := new(MockAISlugService)

	service := services.NewURLService(mockStorage, mockAI, "localhost", "8080")

	req := services.URLRequest{URL: ""}

	// Execute
	response, err := service.CreateShortURL(context.Background(), req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "URL is required")
}

func TestURLService_CreateShortURL_StorageError(t *testing.T) {
	// Setup
	mockStorage := NewMockRedisStorage()
	mockAI := new(MockAISlugService)

	service := services.NewURLService(mockStorage, mockAI, "localhost", "8080")

	req := services.URLRequest{URL: "https://example.com"}

	// Mock AI service to fail (so it falls back to hash-based slug)
	mockAI.On("GenerateSlug", mock.Anything, req.URL).Return("", assert.AnError)

	// Mock storage store error
	mockStorage.On("StoreURL", mock.Anything, mock.AnythingOfType("string"), req.URL).Return(assert.AnError)

	// Execute
	response, err := service.CreateShortURL(context.Background(), req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to store URL")

	mockStorage.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

func TestURLService_GetOriginalURL(t *testing.T) {
	// Setup
	mockStorage := NewMockRedisStorage()
	mockAI := new(MockAISlugService)

	service := services.NewURLService(mockStorage, mockAI, "localhost", "8080")

	shortCode := "abc123"
	originalURL := "https://example.com"

	mockStorage.On("GetURL", mock.Anything, shortCode).Return(originalURL, nil)

	// Execute
	result, err := service.GetOriginalURL(context.Background(), shortCode)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, originalURL, result)

	mockStorage.AssertExpectations(t)
}

func TestURLService_GetOriginalURL_NotFound(t *testing.T) {
	// Setup
	mockStorage := NewMockRedisStorage()
	mockAI := new(MockAISlugService)

	service := services.NewURLService(mockStorage, mockAI, "localhost", "8080")

	shortCode := "nonexistent"

	mockStorage.On("GetURL", mock.Anything, shortCode).Return("", assert.AnError)

	// Execute
	result, err := service.GetOriginalURL(context.Background(), shortCode)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, result)

	mockStorage.AssertExpectations(t)
}
