package utils

import (
	"crypto/sha1"
	"encoding/binary"
	"strings"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// ShortHash generates a hash-based short code from input string
func ShortHash(input string) string {
	hash := sha1.Sum([]byte(input))
	num := binary.BigEndian.Uint64(hash[:8])
	return Base62Encode(num)
}

// Base62Encode encodes a number to base62
func Base62Encode(num uint64) string {
	if num == 0 {
		return "0"
	}

	var result strings.Builder
	base := uint64(len(base62Chars))

	for num > 0 {
		result.WriteByte(base62Chars[num%base])
		num /= base
	}

	// Reverse the string
	str := result.String()
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
