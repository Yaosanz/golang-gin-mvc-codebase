package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-starter-app/helpers"
	"go-starter-app/interfaces"
)

func JWTAuthMiddleware(app interfaces.IAppDependencies) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header missing",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization format",
			})
			return
		}

		tokenString := parts[1]

		// ✅ PAKAI JWT YANG SUDAH ADA DI PROJECT
		secret := app.GetConfig().Jwt().Secret
		claims, err := helpers.ValidateJWT(tokenString, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			return
		}

		// ✅ SET CONTEXT
		c.Set("user_id", claims.UserID)
		c.Set("role_id", claims.RoleID)
		c.Set("app", app)

		c.Next()
	}
}
