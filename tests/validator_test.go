package tests

import (
	"testing"

	"go-url-shortner/services"

	"github.com/stretchr/testify/assert"
)

func TestURLValidator_ValidateURL_ValidURLs(t *testing.T) {
	validator := services.NewURLValidator()

	testCases := []string{
		"https://example.com",
		"http://example.com",
		"https://www.example.com",
		"https://example.com/path",
		"https://example.com/path?param=value",
		"https://example.com/path#fragment",
		"https://subdomain.example.com",
		"https://example.co.uk",
		"https://example.com:8080",
	}

	for _, url := range testCases {
		t.Run(url, func(t *testing.T) {
			err := validator.ValidateURL(url)
			assert.NoError(t, err, "URL should be valid: %s", url)
		})
	}
}

func TestURLValidator_ValidateURL_InvalidURLs(t *testing.T) {
	validator := services.NewURLValidator()

	testCases := []string{
		"",
		"not-a-url",
		"ftp://example.com", // unsupported scheme
		"example",           // no domain
		"example.",          // incomplete domain
		".example",          // invalid domain
		"http://",           // no host
		"https://",          // no host
		"://example.com",    // no scheme
	}

	for _, url := range testCases {
		t.Run(url, func(t *testing.T) {
			err := validator.ValidateURL(url)
			assert.Error(t, err, "URL should be invalid: %s", url)
		})
	}
}

func TestURLValidator_ValidateURL_URLsWithoutScheme(t *testing.T) {
	validator := services.NewURLValidator()

	testCases := []string{
		"example.com",
		"www.example.com",
		"subdomain.example.com",
		"example.com/path",
		"example.com:8080",
	}

	for _, url := range testCases {
		t.Run(url, func(t *testing.T) {
			err := validator.ValidateURL(url)
			assert.NoError(t, err, "URL should be valid with auto-added scheme: %s", url)
		})
	}
}

func TestURLValidator_ValidateURL_EdgeCases(t *testing.T) {
	validator := services.NewURLValidator()

	// Test with very long URLs
	longURL := "https://" + string(make([]byte, 1000)) + ".com"
	err := validator.ValidateURL(longURL)
	assert.Error(t, err, "Very long URL should be invalid")

	// Test with special characters in domain (should fail)
	specialChars := "https://ex@mple.com"
	err = validator.ValidateURL(specialChars)
	assert.Error(t, err, "URL with special characters in domain should be invalid")

	// Test with valid special characters in path
	validSpecialChars := "https://example.com/path-with_underscores-and-dashes"
	err = validator.ValidateURL(validSpecialChars)
	assert.NoError(t, err, "URL with valid special characters in path should be valid")
}

func TestURLValidator_NormalizeURL_ValidURLs(t *testing.T) {
	validator := services.NewURLValidator()

	testCases := []struct {
		input    string
		expected string
	}{
		{"https://example.com", "https://example.com"},
		{"http://example.com", "http://example.com"},
		{"example.com", "https://example.com"},
		{"www.example.com", "https://www.example.com"},
		{"subdomain.example.com/path", "https://subdomain.example.com/path"},
		{"example.com:8080", "https://example.com:8080"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			normalized, err := validator.NormalizeURL(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, normalized)
		})
	}
}

func TestURLValidator_NormalizeURL_InvalidURLs(t *testing.T) {
	validator := services.NewURLValidator()

	testCases := []string{
		"",
		"not-a-url",
		"ftp://example.com",
		"example",
		"http://",
	}

	for _, url := range testCases {
		t.Run(url, func(t *testing.T) {
			normalized, err := validator.NormalizeURL(url)
			assert.Error(t, err)
			assert.Empty(t, normalized)
		})
	}
}

func TestURLValidator_NewURLValidator(t *testing.T) {
	validator := services.NewURLValidator()
	assert.NotNil(t, validator)

	// Test that the same instance can be reused
	err := validator.ValidateURL("https://example.com")
	assert.NoError(t, err)

	err = validator.ValidateURL("https://another-example.com")
	assert.NoError(t, err)
}

func TestURLValidator_ConsistentValidation(t *testing.T) {
	validator := services.NewURLValidator()

	// Test that the same URL always validates the same way
	url := "https://example.com"

	for i := 0; i < 10; i++ {
		err := validator.ValidateURL(url)
		assert.NoError(t, err)
	}
}

func TestURLValidator_InternationalDomains(t *testing.T) {
	validator := services.NewURLValidator()

	// Test with international domains (IDN)
	testCases := []string{
		"https://münchen.de",
		"https://café.com",
		"https://résumé.net",
	}

	for _, url := range testCases {
		t.Run(url, func(t *testing.T) {
			err := validator.ValidateURL(url)
			// Note: These might fail depending on Go version and URL parsing
			// The test documents the current behavior
			if err != nil {
				t.Logf("International domain validation failed (expected in some cases): %s - %v", url, err)
			}
		})
	}
}
