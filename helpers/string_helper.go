package helpers

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// StringToPtr convert string to *string
//
//	Usage:
//	- use this helper if you need to convert a string to string pointer:
//	Params:
//	- s: string
//	Returns:
//	- *string
func StringToPtr(s string) *string {
	return &s
}

// PtrToString convert *string to string
//
//	Usage:
//	- use this helper if you need to convert string pointer to string:
//	Params:
//	- s: *string
//	Returns:
//	- string
func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// PtrToStringOrDefault convert *string to string
// if nil return default value ("")
//
//	Usage:
//	- use this helper if you need to convert string pointer to string, and return default value if nil:
//	Params:
//	- s: *string
//	Returns:
//	- string
func PtrToStringOrDefault(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

// IsPtrEqualsToString compares a string pointer (*string) with a regular string value without causing a panic.
//
//	Usage:
//	- use this helper If you need to compare a pointer string value (*string) with a regular string value (string):
//	Params:
//	- s: *string
//	- v: string
//	Returns:
//	- bool
func IsPtrEqualsToString(s *string, v string) bool {
	if s == nil {
		return false // if nil, cannot compare to any value
	}
	return *s == v
}

// GenerateSlug generate slug from string
//
// Usage:
// - use this helper to generate slug from string
// Params:
// - input: string
// Returns:
// - string
func GenerateSlug(input string) string {
	slug := strings.ToLower(input)

	// erase all non alphanumeric or space characters
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	slug = reg.ReplaceAllString(slug, "")

	// change space by hyphen
	slug = strings.ReplaceAll(slug, " ", "-")

	// remove hyphen from beginning and ending
	slug = strings.Trim(slug, "-")

	return slug
}

// ParamInt64 extracts a path parameter by key and converts it to int64.
// Returns (value, error).
func ParamInt64(ctx *gin.Context, key string) (int64, error) {
	param := ctx.Param(key)
	if param == "" {
		return 0, strconv.ErrSyntax
	}

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// HashPassword hashes the plaintext password and returns the hashed string.
// recommendedCost: bcrypt.DefaultCost (10) or 12 for more work factor on servers.
func HashPassword(password string, cost int) (string, error) {
	if password == "" {
		return "", errors.New("password is empty")
	}
	if cost == 0 {
		cost = bcrypt.DefaultCost // or 12
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPasswordHash CheckPassword verifies a plaintext password against a bcrypt hashed password.
// Returns nil if matched, or an error.
func CheckPasswordHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
