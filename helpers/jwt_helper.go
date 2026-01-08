package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JwtClaims adalah custom claims untuk JWT
type JwtClaims struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	RoleID   uuid.UUID `json:"role_id"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT token dengan UUID user
func GenerateToken(
	userID uuid.UUID,
	username string,
	role string,
	roleID uuid.UUID,
	secret string,
	expired time.Duration,
	issuer string,
) (string, error) {

	claims := JwtClaims{
		UserID:   userID.String(),
		Username: username,
		Role:     role,
		RoleID:   roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expired)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT memvalidasi dan parse JWT token
func ValidateJWT(tokenString string, secret string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
