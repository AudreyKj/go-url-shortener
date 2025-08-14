package tests

import (
	"fmt"
	"testing"

	"go-url-shortner/utils"

	"github.com/stretchr/testify/assert"
)

func TestShortHash_ConsistentOutput(t *testing.T) {
	// Test that the same input always produces the same output
	input := "https://example.com"

	hash1 := utils.ShortHash(input)
	hash2 := utils.ShortHash(input)

	assert.Equal(t, hash1, hash2, "Same input should produce same hash")
	assert.NotEmpty(t, hash1, "Hash should not be empty")
}

func TestShortHash_DifferentInputs(t *testing.T) {
	// Test that different inputs produce different hashes
	input1 := "https://example.com"
	input2 := "https://example.org"
	input3 := "https://github.com"

	hash1 := utils.ShortHash(input1)
	hash2 := utils.ShortHash(input2)
	hash3 := utils.ShortHash(input3)

	assert.NotEqual(t, hash1, hash2, "Different inputs should produce different hashes")
	assert.NotEqual(t, hash1, hash3, "Different inputs should produce different hashes")
	assert.NotEqual(t, hash2, hash3, "Different inputs should produce different hashes")
}

func TestShortHash_EmptyInput(t *testing.T) {
	// Test that an empty string input produces a valid hash
	hash := utils.ShortHash("")

	assert.NotEmpty(t, hash, "Empty input should still produce a hash")
}

func TestShortHash_WhitespaceInput(t *testing.T) {
	// Test that whitespace-only inputs produce valid hashes
	hash1 := utils.ShortHash("   ")
	hash2 := utils.ShortHash("\t\n\r")

	assert.NotEmpty(t, hash1, "Whitespace input should produce a hash")
	assert.NotEmpty(t, hash2, "Whitespace input should produce a hash")
}

func TestShortHash_SpecialCharacters(t *testing.T) {
	// Test that inputs with special characters produce valid hashes
	inputs := []string{
		"https://example.com/path?param=value&other=123",
		"https://example.com/path#fragment",
		"https://user:pass@example.com:8080/path",
		"https://example.com/path with spaces",
		"https://example.com/path-with-underscores_and-dashes",
		"https://example.com/path/with/unicode/ñáéíóú",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			hash := utils.ShortHash(input)
			assert.NotEmpty(t, hash, "Special characters should not break hash generation")
		})
	}
}

func TestShortHash_Length(t *testing.T) {
	// Test that the hash length is within a reasonable range
	input := "https://example.com"
	hash := utils.ShortHash(input)

	// Hash should be reasonably short (typically 6-12 characters)
	assert.GreaterOrEqual(t, len(hash), 6, "Hash should be at least 6 characters")
	assert.LessOrEqual(t, len(hash), 12, "Hash should be at most 12 characters")
}

func TestShortHash_UnicodeInput(t *testing.T) {
	// Test that Unicode inputs produce valid hashes
	inputs := []string{
		"https://münchen.de",
		"https://café.com",
		"https://résumé.net",
		"https://example.com/路径/文件",
		"https://example.com/путь/файл",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			hash := utils.ShortHash(input)
			assert.NotEmpty(t, hash, "Unicode input should produce a hash")
		})
	}
}

func TestBase62Encode_Zero(t *testing.T) {
	// Test that encoding zero produces "0"
	result := utils.Base62Encode(0)
	assert.Equal(t, "0", result)
}

func TestBase62Encode_SmallNumbers(t *testing.T) {
	// Test encoding of small numbers into base62
	tests := []struct {
		input    uint64
		expected string
	}{
		{1, "1"},
		{10, "A"},  // 10 in base62 is 'A' (10th character in base62Chars)
		{35, "Z"},  // 35 in base62 is 'Z' (35th character in base62Chars)
		{36, "a"},  // 36 in base62 is 'a' (36th character in base62Chars)
		{61, "z"},  // 61 in base62 is 'z' (61st character in base62Chars)
		{62, "10"}, // 62 in base62 is "10" (62 = 1*62^1 + 0*62^0)
		{63, "11"}, // 63 in base62 is "11" (63 = 1*62^1 + 1*62^0)
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.input), func(t *testing.T) {
			result := utils.Base62Encode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBase62Encode_LargeNumbers(t *testing.T) {
	// Test encoding of large numbers into base62
	tests := []struct {
		input    uint64
		expected string
	}{
		{62, "10"},       // 62 = 1*62^1 + 0*62^0
		{63, "11"},       // 63 = 1*62^1 + 1*62^0
		{124, "20"},      // 124 = 2*62^1 + 0*62^0
		{3844, "100"},    // 3844 = 1*62^2 + 0*62^1 + 0*62^0
		{238328, "1000"}, // 238328 = 1*62^3 + 0*62^2 + 0*62^1 + 0*62^0
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.input), func(t *testing.T) {
			result := utils.Base62Encode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBase62Encode_EdgeCases(t *testing.T) {
	// Test edge cases for base62 encoding
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{35, "Z"},     // 35 = 'Z' (35th character)
		{36, "a"},     // 36 = 'a' (36th character)
		{61, "z"},     // 61 = 'z' (61st character)
		{62, "10"},    // 62 = 1*62^1 + 0*62^0
		{63, "11"},    // 63 = 1*62^1 + 1*62^0
		{100, "1c"},   // 100 = 1*62^1 + 38*62^0, 38 = 'c' (lowercase)
		{3844, "100"}, // 3844 = 1*62^2 + 0*62^1 + 0*62^0
		{1000, "G8"},  // 1000 = 16*62^1 + 8*62^0, 16 = 'G' (uppercase), 8 = '8'
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.input), func(t *testing.T) {
			result := utils.Base62Encode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBase62Encode_Consistency(t *testing.T) {
	// Test that encoding the same number multiple times produces consistent results
	input := uint64(12345)

	for i := 0; i < 100; i++ {
		result := utils.Base62Encode(input)
		assert.Equal(t, "3D7", result)
	}
}

func TestBase62Encode_CharacterSet(t *testing.T) {
	// Test that all characters in the base62 character set are used
	usedChars := make(map[byte]bool)

	// Test a range of numbers to see what characters are used
	for i := uint64(0); i < 10000; i++ {
		result := utils.Base62Encode(i)
		for j := 0; j < len(result); j++ {
			usedChars[result[j]] = true
		}
	}

	// Should have used most of the base62 character set
	assert.Greater(t, len(usedChars), 50, "Should use most of the base62 character set")
}

func TestBase62Encode_Reversibility(t *testing.T) {
	// Test that base62 encoding produces valid reversible strings
	// Note: This is not a full base62 decode test, just validation that the encoding produces valid base62 strings

	testCases := []uint64{0, 1, 10, 35, 36, 61, 62, 63, 100, 1000, 10000}

	for _, input := range testCases {
		t.Run(string(rune(input)), func(t *testing.T) {
			encoded := utils.Base62Encode(input)

			// Check that the encoded string only contains valid base62 characters
			for _, char := range encoded {
				assert.True(t,
					(char >= '0' && char <= '9') ||
						(char >= 'A' && char <= 'Z') ||
						(char >= 'a' && char <= 'z'),
					"Encoded string should only contain base62 characters: %c", char)
			}
		})
	}
}

func TestShortHash_Base62Compliance(t *testing.T) {
	// Test that ShortHash produces base62-compliant strings
	input := "https://example.com"
	hash := utils.ShortHash(input)

	// Check that the hash only contains valid base62 characters
	for _, char := range hash {
		assert.True(t,
			(char >= '0' && char <= '9') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= 'a' && char <= 'z'),
			"Hash should only contain base62 characters: %c", char)
	}
}
