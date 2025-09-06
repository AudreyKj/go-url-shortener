package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-url-shortner/handlers"
	"go-url-shortner/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestCreateShortURL_Success(t *testing.T) {
	// Test the successful creation of a short URL.
	// Setup
	router := setupTestRouter()
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)

	requestBody := services.URLRequest{URL: "https://example.com"}
	expectedResponse := &services.URLResponse{
		OriginalURL: "https://example.com",
		ShortCode:   "abc123",
		ShortURL:    "http://localhost:8080/abc123",
		SlugType:    "hash_based",
	}

	mockService.On("CreateShortURL", mock.Anything, requestBody).Return(expectedResponse, nil)

	router.POST("/api/urls", handler.CreateShortURL)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 for successful URL creation")

	var response services.URLResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.OriginalURL, response.OriginalURL)
	assert.Equal(t, expectedResponse.ShortCode, response.ShortCode, "ShortCode mismatch")

	mockService.AssertExpectations(t)
}

func TestCreateShortURL_InvalidJSON(t *testing.T) {
	// Test handling of invalid JSON input.
	// Setup
	router := setupTestRouter()
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)

	router.POST("/api/urls", handler.CreateShortURL)

	// Create request with invalid JSON
	req := httptest.NewRequest("POST", "/api/urls", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP status 400 for invalid JSON input")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid JSON", response["error"], "Error message mismatch for invalid JSON")
}

func TestCreateShortURL_ServiceError(t *testing.T) {
	// Test handling of service errors during URL creation.
	// Setup
	router := setupTestRouter()
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)

	requestBody := services.URLRequest{URL: "https://example.com"}
	mockService.On("CreateShortURL", mock.Anything, requestBody).Return(nil, assert.AnError)

	router.POST("/api/urls", handler.CreateShortURL)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected HTTP status 500 for service error")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, assert.AnError.Error(), response["error"], "Error message mismatch for service error")

	mockService.AssertExpectations(t)
}

func TestRedirectToURL_Success(t *testing.T) {
	// Test successful redirection to the original URL.
	// Setup
	router := setupTestRouter()
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)

	shortCode := "abc123"
	originalURL := "https://example.com"

	mockService.On("GetOriginalURL", mock.Anything, shortCode).Return(originalURL, nil)

	router.GET("/:shortCode", handler.RedirectToURL)

	// Create request
	req := httptest.NewRequest("GET", "/"+shortCode, nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusMovedPermanently, w.Code, "Expected HTTP status 301 for successful redirection")
	assert.Equal(t, originalURL, w.Header().Get("Location"), "Location header mismatch")

	mockService.AssertExpectations(t)
}

func TestRedirectToURL_NotFound(t *testing.T) {
	// Test handling of non-existent short codes.
	// Setup
	router := setupTestRouter()
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)

	shortCode := "nonexistent"
	mockService.On("GetOriginalURL", mock.Anything, shortCode).Return("", assert.AnError)

	router.GET("/:shortCode", handler.RedirectToURL)

	// Create request
	req := httptest.NewRequest("GET", "/"+shortCode, nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP status 404 for non-existent short code")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Short URL not found", response["error"], "Error message mismatch for non-existent short code")

	mockService.AssertExpectations(t)
}

func TestHealthCheck(t *testing.T) {
	// Test the health check endpoint.
	// Setup
	router := setupTestRouter()
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)

	router.GET("/health", handler.HealthCheck)

	// Create request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 for health check")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"], "Health status mismatch")
	assert.Equal(t, "redis", response["storage"], "Storage type mismatch")
}

func TestCreateShortURL_InvalidURL(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)

	router.POST("/api/urls", handler.CreateShortURL)

	// Invalid URL (no domain)
	requestBody := services.URLRequest{URL: "not_a_url"}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP status 400 for invalid URL")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "valid domain", "Error message should mention domain validation")
}

func TestCreateShortURL_NormalizeURL_Success(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)

	// Input without scheme
	requestBody := services.URLRequest{URL: "google.com"}
	// The normalized URL should be https://google.com
	expectedResponse := &services.URLResponse{
		OriginalURL: "https://google.com",
		ShortCode:   "abc123",
		ShortURL:    "http://localhost:8080/abc123",
		SlugType:    "hash_based",
	}

	mockService.On("CreateShortURL", mock.Anything, mock.MatchedBy(func(req services.URLRequest) bool {
		return req.URL == "https://google.com"
	})).Return(expectedResponse, nil)

	router.POST("/api/urls", handler.CreateShortURL)

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 for normalized URL input")

	var response services.URLResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.OriginalURL, response.OriginalURL)
	assert.Equal(t, expectedResponse.ShortCode, response.ShortCode)
	mockService.AssertExpectations(t)
}
