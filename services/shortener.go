package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go-url-shortner/utils"
)

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	SlugType    string `json:"slug_type"`
}

type URLService struct {
	storage    StorageInterface
	aiService  AISlugServiceInterface
	serverHost string
	serverPort string
}

func NewURLService(storage StorageInterface, aiService AISlugServiceInterface, serverHost, serverPort string) *URLService {
	return &URLService{
		storage:    storage,
		aiService:  aiService,
		serverHost: serverHost,
		serverPort: serverPort,
	}
}

func (s *URLService) CreateShortURL(ctx context.Context, req URLRequest) (*URLResponse, error) {
	if req.URL == "" {
		return nil, fmt.Errorf("URL is required")
	}

	var shortCode string
	var err error

	if s.aiService != nil {
		aiSlug, aiErr := s.aiService.GenerateSlug(ctx, req.URL)

		if aiErr == nil && aiSlug != "" {
			if s.isSlugAvailable(ctx, aiSlug) {
				shortCode = aiSlug
				log.Printf("Using AI-generated slug: %s", shortCode)
			} else {
				log.Printf("AI-generated slug '%s' already exists, falling back to hash", aiSlug)
			}
		} else {
			log.Printf("AI slug generation failed: %v, falling back to hash", aiErr)
		}
	}

	// Fallback to hash-based slug if AI failed or slug is unavailable
	if shortCode == "" {
		shortCode = utils.ShortHash(req.URL)
		log.Printf("Using hash-based slug: %s", shortCode)
	}

	err = s.storage.StoreURL(ctx, shortCode, req.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to store URL: %w", err)
	}

	log.Printf("Stored URL mapping - Short: %s, Original: %s", shortCode, req.URL)

	return &URLResponse{
		OriginalURL: req.URL,
		ShortCode:   shortCode,
		ShortURL:    fmt.Sprintf("http://%s:%s/%s", s.serverHost, s.serverPort, shortCode),
		SlugType:    s.getSlugType(shortCode),
	}, nil
}

func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	return s.storage.GetURL(ctx, shortCode)
}

func (s *URLService) isSlugAvailable(ctx context.Context, slug string) bool {
	_, err := s.storage.GetURL(ctx, slug)
	return err != nil 
}


func (s *URLService) getSlugType(slug string) string {
	if s.aiService != nil {
		// Simple heuristic: AI slugs are typically shorter and more readable
		if len(slug) <= 8 && s.isReadableSlug(slug) {
			return "ai_generated"
		}
	}
	return "hash_based"
}

func (s *URLService) isReadableSlug(slug string) bool {
	consonantCount := 0
	vowelCount := 0

	for _, char := range slug {
		if strings.ContainsRune("bcdfghjklmnpqrstvwxz", char) {
			consonantCount++
		} else if strings.ContainsRune("aeiou", char) {
			vowelCount++
		}
	}

	// AI slugs typically have a good balance of consonants and vowels
	return consonantCount > 0 && vowelCount > 0 && consonantCount+vowelCount >= 3
}
