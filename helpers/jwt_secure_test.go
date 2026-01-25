package helpers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSecureUserToken(t *testing.T) {
	userID := uuid.New()
	sessionID := GenerateSessionID()
	roles := []string{"user", "cms"}

	token, err := GenerateSecureUserToken(
		userID,
		"user",
		sessionID,
		roles,
		"test-secret",
		24*time.Hour,
		"test-issuer",
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token
	claims, err := ValidateSecureJWT(token, "test-secret")
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), claims.UserID)
	assert.Equal(t, "user", claims.TokenType)
	assert.Equal(t, sessionID, claims.SessionID)
	assert.ElementsMatch(t, roles, claims.Roles)
}

func TestValidateSecureJWT(t *testing.T) {
	userID := uuid.New()
	sessionID := GenerateSessionID()
	roles := []string{"admin"}

	token, _ := GenerateSecureUserToken(
		userID,
		"cms",
		sessionID,
		roles,
		"test-secret",
		24*time.Hour,
		"test-issuer",
	)

	claims, err := ValidateSecureJWT(token, "test-secret")
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), claims.UserID)
	assert.True(t, claims.IsCMSToken())
	assert.False(t, claims.IsUserToken())
}

func TestSecureJwtClaims_HasRole(t *testing.T) {
	claims := &SecureJwtClaims{
		Roles: []string{"admin", "user"},
	}

	assert.True(t, claims.HasRole("admin"))
	assert.True(t, claims.HasRole("user"))
	assert.False(t, claims.HasRole("cms"))
}

func TestGenerateSessionID(t *testing.T) {
	sessionID1 := GenerateSessionID()
	sessionID2 := GenerateSessionID()

	assert.NotEmpty(t, sessionID1)
	assert.NotEmpty(t, sessionID2)
	assert.NotEqual(t, sessionID1, sessionID2)
}

func TestInvalidTokenSecret(t *testing.T) {
	userID := uuid.New()
	sessionID := GenerateSessionID()

	token, _ := GenerateSecureUserToken(
		userID,
		"user",
		sessionID,
		[]string{"user"},
		"correct-secret",
		24*time.Hour,
		"test-issuer",
	)

	// Try to validate with wrong secret
	_, err := ValidateSecureJWT(token, "wrong-secret")
	assert.Error(t, err)
}

func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	sessionID := GenerateSessionID()

	// Create token with negative duration (already expired)
	token, _ := GenerateSecureUserToken(
		userID,
		"user",
		sessionID,
		[]string{"user"},
		"test-secret",
		-1*time.Hour,
		"test-issuer",
	)

	// Try to validate expired token
	_, err := ValidateSecureJWT(token, "test-secret")
	assert.Error(t, err)
}
