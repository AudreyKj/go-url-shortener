package tests

import (
	"context"

	"go-url-shortner/services"

	"github.com/stretchr/testify/mock"
)

// MockURLService is a mock implementation that can be used where services.URLServiceInterface is expected
type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) CreateShortURL(ctx context.Context, req services.URLRequest) (*services.URLResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.URLResponse), args.Error(1)
}

func (m *MockURLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	args := m.Called(ctx, shortCode)
	return args.String(0), args.Error(1)
}

// MockAISlugService is a mock implementation that can be used where services.AISlugServiceInterface is expected
type MockAISlugService struct {
	mock.Mock
}

func (m *MockAISlugService) GenerateSlug(ctx context.Context, originalURL string) (string, error) {
	args := m.Called(ctx, originalURL)
	return args.String(0), args.Error(1)
}
