package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-starter-app/helpers"
	"go-starter-app/interfaces"
)

// RequirePermission checks if the user has the required permission
// Uses cached permissions from JWT claims for better performance
func RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsAny, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token claims not found",
			})
			return
		}

		claims, ok := claimsAny.(*helpers.JwtClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid token claims",
			})
			return
		}

		// Check permission from JWT claims first (faster)
		if claims.HasPermission(code) {
			c.Next()
			return
		}

		// Fallback to database check if not in claims
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

		// Check permissions for all user roles
		for _, roleName := range claims.Roles {
			hasPermission, err := app.GetService().GetPermissionService().
				HasPermission(c.Request.Context(), roleName, code)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error checking permissions",
				})
				return
			}

			if hasPermission {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden: insufficient permissions",
		})
	}
}

// RequireRole checks if the user has the required role
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsAny, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token claims not found",
			})
			return
		}

		claims, ok := claimsAny.(*helpers.JwtClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid token claims",
			})
			return
		}

		if !claims.HasRole(roleName) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden: role '" + roleName + "' required",
			})
			return
		}

		c.Next()
	}
}

// RequireAnyRole checks if the user has any of the required roles
func RequireAnyRole(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsAny, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token claims not found",
			})
			return
		}

		claims, ok := claimsAny.(*helpers.JwtClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid token claims",
			})
			return
		}

		for _, roleName := range roleNames {
			if claims.HasRole(roleName) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden: one of the required roles needed",
		})
	}
}
