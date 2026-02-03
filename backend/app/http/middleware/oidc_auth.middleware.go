package middleware

import (
	"context"
	"fmt"
	"go-starter-app/app/http/utils"
	"go-starter-app/config"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
)

type ResourceAccess struct {
	Roles []string `json:"roles"`
}

type CustomClaims struct {
	Exp               int64                     `json:"exp"`
	Iat               int64                     `json:"iat"`
	Jti               string                    `json:"jti"`
	Iss               string                    `json:"iss"`
	Sub               string                    `json:"sub"`
	Typ               string                    `json:"typ"`
	Azp               string                    `json:"azp"`
	Name              string                    `json:"name"`
	Email             string                    `json:"email"`
	PreferredUsername string                    `json:"preferred_username"`
	Civitas           string                    `json:"civitas"`
	ResourceAccess    map[string]ResourceAccess `json:"resource_access"`
}

// OIDCAuthMiddleware verify access_token received from identity provider (e.x KeyCloak)
func OIDCAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	provider, err := oidc.NewProvider(context.Background(), cfg.Oidc().Issuer)
	if err != nil {
		panic(fmt.Sprintf("Failed to create OIDC provider instance: %v", err))
	}

	// verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	verifier := provider.Verifier(&oidc.Config{
		ClientID:          cfg.Oidc().ClientID,
		SkipClientIDCheck: cfg.Oidc().SkipClientIDCheck, // Skipping default client ID check
	})

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

		// Verify the token
		token, err := verifier.Verify(c.Request.Context(), tokenString)
		if err != nil {
			log.Printf("[middleware:oidc] error: %s", "Invalid token")
			utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
			c.Abort()
			return
		}

		// Extract claims if needed
		var claims map[string]interface{}
		if err := token.Claims(&claims); err != nil {
			log.Printf("[middleware:oidc] error: %s", "Unable to parse claims")
			utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
			c.Abort()
			return
		}
		// Navigate to `roles`
		resourceAccess, ok := claims["resource_access"].(map[string]interface{})
		if !ok {
			fmt.Println("Invalid resource_access format")
		}

		access, ok := resourceAccess["access"].(map[string]interface{})
		if !ok {
			fmt.Println("Invalid resource access format")
		}

		roles, ok := access["roles"].([]interface{})
		if !ok {
			fmt.Println("Invalid roles format")
		}

		// Set claims or token info in the context if necessary
		c.Set("claims", claims)
		c.Set("roles", roles)
		c.Next()
	}
}
