package tests

import (
	"context"
	"testing"

	"go-url-shortner/services"
)

func TestAISlugService_GenerateSlug(t *testing.T) {
	// Test with empty API key (should return nil service)
	service := services.NewAISlugService("")
	if service != nil {
		t.Error("Expected nil service when API key is empty")
	}

	// Test with valid API key
	service = services.NewAISlugService("test-key")
	if service == nil {
		t.Error("Expected service to be created with valid API key")
	}

	// Test slug generation
	ctx := context.Background()
	slug, err := service.GenerateSlug(ctx, "https://github.com")
	// We expect an error here since we don't have a real API key
	if err == nil {
		t.Error("Expected error when calling OpenAI API with invalid key")
	} else {
		t.Logf("Received expected error: %v", err)
	}

	// Validate slug 
	if slug != "" {
		t.Errorf("Expected empty slug due to invalid API key, got: %s", slug)
	}
}
