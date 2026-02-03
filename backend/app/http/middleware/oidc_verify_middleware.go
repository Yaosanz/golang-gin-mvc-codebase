package middleware

import (
	"fmt"
	"go-starter-app/app/http/utils"
	"go-starter-app/pkg/provider"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Key for context
const (
	ClaimsKey = "claims"
	RolesKey  = "roles"
	UserKey   = "user"
)

// CustomClaim represents the parsed claims from OIDC token
type CustomClaim struct {
	ResourceAccess map[string]interface{} `json:"resource_access"`
	Roles          []string               `json:"roles"`
}

func OidcVerifyMiddleware(oidcVerifier *provider.OidcProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("[middleware:oidc] error: %s", "Authorization header missing")
			utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Printf("[middleware:oidc] error: %s", "Bearer token missing")
			utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
			c.Abort()
			return
		}

		// verifikasi token string
		idToken, err := oidcVerifier.VerifyToken(c, tokenString)
		if err != nil {
			log.Printf("[middleware:oidc] error: %s", "Invalid token")
			utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
			c.Abort()
			return
		}

		// parse token claim
		var claims map[string]interface{}
		err = oidcVerifier.ParseClaims(idToken, &claims)
		if err != nil {
			utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
			return
		}

		claimMapper, err := mapClaims(claims)
		if err != nil {
			utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
			return
		}

		// pass the claims in case of another user
		c.Set(ClaimsKey, claims)
		c.Set(RolesKey, claimMapper.Roles)
		c.Next()
	}
}

func mapClaims(claims map[string]interface{}) (*CustomClaim, error) {
	resourceAccess, ok := claims["resource_access"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid resource_access format")
	}

	access, ok := resourceAccess["access"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid resource access format")
	}

	rawRoles, ok := access["roles"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid roles format")
	}

	roles := make([]string, len(rawRoles))
	for i, role := range rawRoles {
		if roleStr, ok := role.(string); ok {
			roles[i] = roleStr
		}
	}

	return &CustomClaim{
		ResourceAccess: resourceAccess,
		Roles:          roles,
	}, nil
}
