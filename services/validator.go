package services

import (
	"fmt"
	"net/url"
	"strings"
)

type URLValidator struct{}

func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

func (v *URLValidator) ValidateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL is required")
	}

	// Check for invalid patterns
	if strings.HasSuffix(urlStr, ".") || strings.HasPrefix(urlStr, ".") {
		return fmt.Errorf("URL cannot start or end with a dot")
	}

	// Add scheme if missing
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check if host is present
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a valid host")
	}

	// Check if host has at least one dot (basic domain validation)
	if !strings.Contains(parsedURL.Host, ".") {
		return fmt.Errorf("URL must have a valid domain")
	}

	// Additional validation for edge cases
	if strings.HasSuffix(parsedURL.Host, ".") || strings.HasPrefix(parsedURL.Host, ".") {
		return fmt.Errorf("URL host cannot start or end with a dot")
	}

	// Check for @ character in invalid position (not in userinfo format)
	if strings.Contains(urlStr, "@") {
		// Check if @ is in the correct userinfo format (username:password@host)
		schemeEnd := strings.Index(urlStr, "://")
		if schemeEnd == -1 {
			return fmt.Errorf("URL contains @ character in invalid position")
		}

		atIndex := strings.Index(urlStr, "@")
		if atIndex < schemeEnd+3 {
			return fmt.Errorf("URL contains @ character in invalid position")
		}

		// Check if there's a colon before @ (username:password@host)
		userInfo := urlStr[schemeEnd+3 : atIndex]
		if !strings.Contains(userInfo, ":") {
			return fmt.Errorf("URL contains @ character in invalid position")
		}
	}

	return nil
}

func (v *URLValidator) NormalizeURL(urlStr string) (string, error) {
	if err := v.ValidateURL(urlStr); err != nil {
		return "", err
	}

	// Add scheme if missing
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	return urlStr, nil
}
