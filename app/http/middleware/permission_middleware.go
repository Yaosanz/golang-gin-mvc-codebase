package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-starter-app/interfaces"
)

func RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ❌ JANGAN MustGet
		appAny, exists := c.Get("app")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Application context not found",
			})
			return
		}

		app, ok := appAny.(interfaces.IAppDependencies)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid application context",
			})
			return
		}

		roleIDAny, exists := c.Get("role_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Role not found in token",
			})
			return
		}

		roleID, ok := roleIDAny.(uuid.UUID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid role ID type",
			})
			return
		}

		hasPermission, err := app.GetService().PermissionService.
			HasPermission(c.Request.Context(), roleID, code)

		if err != nil || !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden",
			})
			return
		}

		c.Next()
	}
}
