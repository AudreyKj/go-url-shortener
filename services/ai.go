package services

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type AISlugService struct {
	client *openai.Client
}

func NewAISlugService(apiKey string) *AISlugService {
	if apiKey == "" {
		return nil
	}

	client := openai.NewClient(apiKey)
	return &AISlugService{
		client: client,
	}
}

// Generates an AI-powered slug for a given URL
func (s *AISlugService) GenerateSlug(ctx context.Context, originalURL string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}

	// Extract domain and path from URL for better context
	domain := extractDomain(originalURL)

	prompt := fmt.Sprintf(`Generate exactly one short, catchy, and memorable URL slug for this website:
	URL: %s
	Domain: %s

	Requirements:
	- 3 to 8 characters
	- Memorable and relevant to the website's name or purpose
	- Use only lowercase letters, numbers, and hyphens
	- No spaces, underscores, or special characters
	- Avoid generic or overused slugs

	Examples:
	- For "github.com" -> "ghub" or "git"
	- For "stackoverflow.com" -> "stack" or "so"
	- For "reddit.com" -> "reddit" or "rdt"

	Output:
	Only return the slug itself with no explanation or formatting.`, originalURL, domain)

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   20,
			Temperature: 0.7,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate AI slug: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	slug := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Clean and validate the slug
	cleanSlug := cleanSlug(slug)

	if cleanSlug == "" {
		return "", fmt.Errorf("generated slug is empty after cleaning")
	}

	log.Printf("AI generated slug '%s' for URL: %s", cleanSlug, originalURL)
	return cleanSlug, nil
}

// Extracts the domain from a URL
func extractDomain(url string) string {
	if strings.HasPrefix(url, "http://") {
		url = url[7:]
	} else if strings.HasPrefix(url, "https://") {
		url = url[8:]
	}

	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}

	return url
}

func cleanSlug(slug string) string {
	// Remove any non-alphanumeric characters except hyphens
	re := regexp.MustCompile(`[^a-z0-9-]`)
	clean := re.ReplaceAllString(strings.ToLower(slug), "")

	// Remove leading/trailing hyphens
	clean = strings.Trim(clean, "-")

	// Ensure it's not too long
	if len(clean) > 8 {
		clean = clean[:8]
	}

	// Ensure it's not too short
	if len(clean) < 3 {
		return ""
	}

	return clean
}
