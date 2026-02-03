package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRefreshToken creates a secure opaque refresh token using 256 bits of entropy
// Returned token is URL-safe and suitable for storage in HttpOnly cookies or headers
func GenerateRefreshToken() (string, error) {
    buf := make([]byte, 32) // 256-bit
    if _, err := rand.Read(buf); err != nil {
        return "", fmt.Errorf("failed to generate refresh token: %w", err)
    }
    return base64.RawURLEncoding.EncodeToString(buf), nil
}
