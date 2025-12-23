package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// OidcClaims represents the claims extracted from a JWT token
type OidcClaims struct {
	Exp            int      `json:"exp"`
	Iat            int      `json:"iat"`
	AuthTime       int      `json:"auth_time"`
	Jti            string   `json:"jti"`
	Iss            string   `json:"iss"`
	Sub            string   `json:"sub"`
	Typ            string   `json:"typ"`
	Azp            string   `json:"azp"`
	Nonce          string   `json:"nonce"`
	SessionState   string   `json:"session_state"`
	Acr            string   `json:"acr"`
	AllowedOrigins []string `json:"allowed-origins"`
	ResourceAccess struct {
		Access struct {
			Roles []string `json:"roles"`
		} `json:"access"`
	} `json:"resource_access"`
	Scope             string `json:"scope"`
	Sid               string `json:"sid"`
	EmailVerified     bool   `json:"email_verified"`
	ShadowExpire      string `json:"shadowExpire"`
	Name              string `json:"name"`
	KodeIdentitas     string `json:"kodeIdentitas"`
	PreferredUsername string `json:"preferred_username"`
	Civitas           string `json:"civitas"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Email             string `json:"email"`
}

// JWTDecode decodes a JWT token's payload and maps it to the given struct
func JWTDecode(tokenString string) (*OidcClaims, error) {
	// Split the token into its parts (header, payload, signature)
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT token format")
	}

	// Decode the payload (second part of the JWT)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("failed to decode JWT payload")
	}

	// Unmarshal the payload into the struct
	var claims OidcClaims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return nil, errors.New("failed to map JWT payload to struct")
	}
	return &claims, nil
}

// JWTMapClaims map jwt claim
func JWTMapClaims(ctx *gin.Context) (map[string]interface{}, error) {
	claims, exists := ctx.Get("claims")
	if !exists {
		return nil, errors.New("claims not found in context")
	}

	// Assert the type
	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid claims format")
	}

	return claimsMap, nil
}
