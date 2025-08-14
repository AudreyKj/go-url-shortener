package tests

import (
	"os"
	"testing"

	"go-url-shortner/utils"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load_DefaultValues(t *testing.T) {
	// Clear environment variables to test defaults
	os.Unsetenv("REDIS_ADDR")
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("REDIS_DB")
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("OPENAI_API_KEY")

	// Load configuration
	cfg := utils.Load()

	// Assert default values
	assert.Equal(t, "localhost:6379", cfg.RedisAddr)
	assert.Equal(t, "", cfg.RedisPassword)
	assert.Equal(t, 0, cfg.RedisDB)
	assert.Equal(t, "localhost", cfg.ServerHost)
	assert.Equal(t, "8080", cfg.ServerPort)
	assert.Equal(t, "", cfg.OpenAIAPIKey)
}

func TestConfig_Load_EnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("REDIS_ADDR", "redis.example.com:6380")
	os.Setenv("REDIS_PASSWORD", "secret123")
	os.Setenv("REDIS_DB", "5")
	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("OPENAI_API_KEY", "sk-test123")

	// Load configuration
	cfg := utils.Load()

	// Assert environment variable values
	assert.Equal(t, "redis.example.com:6380", cfg.RedisAddr)
	assert.Equal(t, "secret123", cfg.RedisPassword)
	assert.Equal(t, 5, cfg.RedisDB)
	assert.Equal(t, "0.0.0.0", cfg.ServerHost)
	assert.Equal(t, "9090", cfg.ServerPort)
	assert.Equal(t, "sk-test123", cfg.OpenAIAPIKey)

	// Clean up
	os.Unsetenv("REDIS_ADDR")
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("REDIS_DB")
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("OPENAI_API_KEY")
}

func TestConfig_Load_InvalidRedisDB(t *testing.T) {
	// Set invalid Redis DB value
	os.Setenv("REDIS_DB", "invalid")

	// Load configuration
	cfg := utils.Load()

	// Should fall back to default value
	assert.Equal(t, 0, cfg.RedisDB)

	// Clean up
	os.Unsetenv("REDIS_DB")
}

func TestConfig_Load_EmptyRedisDB(t *testing.T) {
	// Set empty Redis DB value
	os.Setenv("REDIS_DB", "")

	// Load configuration
	cfg := utils.Load()

	// Should fall back to default value
	assert.Equal(t, 0, cfg.RedisDB)

	// Clean up
	os.Unsetenv("REDIS_DB")
}

func TestConfig_Load_MixedValues(t *testing.T) {
	// Set some environment variables, leave others unset
	os.Setenv("REDIS_ADDR", "custom-redis:6379")
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("OPENAI_API_KEY", "sk-custom123")

	// Load configuration
	cfg := utils.Load()

	// Assert mixed values
	assert.Equal(t, "custom-redis:6379", cfg.RedisAddr)
	assert.Equal(t, "", cfg.RedisPassword)       // default
	assert.Equal(t, 0, cfg.RedisDB)              // default
	assert.Equal(t, "localhost", cfg.ServerHost) // default
	assert.Equal(t, "3000", cfg.ServerPort)
	assert.Equal(t, "sk-custom123", cfg.OpenAIAPIKey)

	// Clean up
	os.Unsetenv("REDIS_ADDR")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("OPENAI_API_KEY")
}

func TestConfig_Load_RedisDBEdgeCases(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"0", 0},
		{"1", 1},
		{"15", 15},
		{"-1", 0},    // invalid, should default to 0
		{"abc", 0},   // invalid, should default to 0
		{"", 0},      // empty, should default to 0
		{"999", 999}, // large number
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			os.Setenv("REDIS_DB", tc.input)
			cfg := utils.Load()
			assert.Equal(t, tc.expected, cfg.RedisDB)
			os.Unsetenv("REDIS_DB")
		})
	}
}

func TestConfig_Load_ConsistentBehavior(t *testing.T) {
	// Test that loading config multiple times gives consistent results
	os.Setenv("REDIS_ADDR", "test-redis:6379")
	os.Setenv("SERVER_PORT", "1234")

	cfg1 := utils.Load()
	cfg2 := utils.Load()

	// Should be equal
	assert.Equal(t, cfg1.RedisAddr, cfg2.RedisAddr)
	assert.Equal(t, cfg1.ServerPort, cfg2.ServerPort)

	// Clean up
	os.Unsetenv("REDIS_ADDR")
	os.Unsetenv("SERVER_PORT")
}

func TestConfig_Load_EnvironmentVariablePriority(t *testing.T) {
	// Test that environment variables take priority over defaults
	os.Setenv("REDIS_ADDR", "env-redis:6379")

	cfg := utils.Load()
	assert.Equal(t, "env-redis:6379", cfg.RedisAddr)

	// Clean up
	os.Unsetenv("REDIS_ADDR")
}

func TestConfig_Load_EmptyStringEnvironmentVariables(t *testing.T) {
	// Test that empty string environment variables are handled correctly
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("OPENAI_API_KEY", "")

	cfg := utils.Load()
	assert.Equal(t, "", cfg.RedisPassword)
	assert.Equal(t, "", cfg.OpenAIAPIKey)

	// Clean up
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("OPENAI_API_KEY")
}

func TestConfig_Load_WhitespaceEnvironmentVariables(t *testing.T) {
	// Test that whitespace in environment variables is preserved
	os.Setenv("REDIS_PASSWORD", "  password with spaces  ")
	os.Setenv("SERVER_HOST", "  host with spaces  ")

	cfg := utils.Load()
	assert.Equal(t, "  password with spaces  ", cfg.RedisPassword)
	assert.Equal(t, "  host with spaces  ", cfg.ServerHost)

	// Clean up
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("SERVER_HOST")
}

func TestConfig_Load_UnicodeEnvironmentVariables(t *testing.T) {
	// Test that Unicode characters in environment variables are handled correctly
	os.Setenv("REDIS_PASSWORD", "p@ssw0rd_ñáéíóú")
	os.Setenv("SERVER_HOST", "host-ñáéíóú.com")

	cfg := utils.Load()
	assert.Equal(t, "p@ssw0rd_ñáéíóú", cfg.RedisPassword)
	assert.Equal(t, "host-ñáéíóú.com", cfg.ServerHost)

	// Clean up
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("SERVER_HOST")
}
