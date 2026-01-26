package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"go-starter-app/app/services"
	"go-starter-app/config"
	"go-starter-app/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SecureAuthMiddleware provides secure JWT authentication with server-side authorization
type SecureAuthMiddleware struct {
	authService *services.SecureAuthService
	jwtConfig   config.JwtConfig
}

// NewSecureAuthMiddleware creates secure auth middleware
func NewSecureAuthMiddleware(authService *services.SecureAuthService, jwtConfig config.JwtConfig) *SecureAuthMiddleware {
	return &SecureAuthMiddleware{
		authService: authService,
		jwtConfig:   jwtConfig,
	}
}

// AuthRequired validates JWT token and provides user context
func (m *SecureAuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			return
		}

		// Extract Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			return
		}

		tokenString := tokenParts[1]

		// Validate secure JWT token
		claims, err := helpers.ValidateSecureJWT(tokenString, m.jwtConfig.Secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		// Check if session is valid (not blacklisted)
		isValid, err := m.authService.IsSessionValid(c.Request.Context(), claims.SessionID)
		if err != nil || !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Session expired",
			})
			return
		}

		// Get user ID from token
		userID, err := claims.GetUserID()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID in token",
			})
			return
		}

		// Check if user has admin role from JWT - admin bypasses active check
		fmt.Printf("User roles: %v\n", claims.Roles)
		if claims.HasRole("admin") {
			// Admin is always considered active
			fmt.Printf("Admin user detected, bypassing active check\n")
		} else {
			// Check if user is active (server-side verification) for non-admin users
			isActive, err := m.authService.IsUserActive(c.Request.Context(), userID)
			if err != nil || !isActive {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Account is inactive",
				})
				return
			}
		}

		// Store user context in Gin context for use in handlers
		c.Set("user_id", userID)
		c.Set("user_type", claims.TokenType)
		c.Set("session_id", claims.SessionID)
		c.Set("claims", claims) // For permission middleware

		c.Next()
	}
}

// PermissionRequired creates middleware that checks for specific permission
func (m *SecureAuthMiddleware) PermissionRequired(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated first
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			return
		}

		userID, ok := userIDVal.(uuid.UUID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user context",
			})
			return
		}

		// Check if user has admin role - admin bypasses all permission checks
		hasAdminRole, err := m.authService.CheckRole(c.Request.Context(), userID, "admin")
		if err == nil && hasAdminRole {
			// Admin has all permissions, proceed
			c.Next()
			return
		}

		// Check permission server-side for non-admin users
		hasPermission, err := m.authService.CheckPermission(c.Request.Context(), userID, permission)
		if err != nil {
			// Permission check error usually means cache expired - require re-login
			log.Printf("[PERMISSION ERROR] Failed to check permission '%s' for user %s: %v", permission, userID, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Session expired, please login again",
			})
			return
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			return
		}

		c.Next()
	}
}

// RoleRequired creates middleware that checks for specific role
func (m *SecureAuthMiddleware) RoleRequired(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated first
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			return
		}

		userID, ok := userIDVal.(uuid.UUID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user context",
			})
			return
		}

		// Check role server-side
		hasRole, err := m.authService.CheckRole(c.Request.Context(), userID, role)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Role check failed",
			})
			return
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Insufficient role",
			})
			return
		}

		c.Next()
	}
}

// AdminRequired is a convenience middleware for admin role
func (m *SecureAuthMiddleware) AdminRequired() gin.HandlerFunc {
	return m.RoleRequired("admin")
}

// CMSRequired is a convenience middleware for CMS role
func (m *SecureAuthMiddleware) CMSRequired() gin.HandlerFunc {
	return m.RoleRequired("cms")
}

// UserRequired is a convenience middleware for user role
func (m *SecureAuthMiddleware) UserRequired() gin.HandlerFunc {
	return m.RoleRequired("user")
}

// Helper functions for handlers

// GetUserID extracts user ID from Gin context
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user not authenticated")
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid user ID in context")
	}

	return userID, nil
}

// GetSessionID extracts session ID from Gin context
func GetSessionID(c *gin.Context) (string, error) {
	sessionIDVal, exists := c.Get("session_id")
	if !exists {
		return "", errors.New("session not found")
	}

	sessionID, ok := sessionIDVal.(string)
	if !ok {
		return "", errors.New("invalid session ID in context")
	}

	return sessionID, nil
}

// GetUserType extracts user type from Gin context
func GetUserType(c *gin.Context) (string, error) {
	userTypeVal, exists := c.Get("user_type")
	if !exists {
		return "", errors.New("user type not found")
	}

	userType, ok := userTypeVal.(string)
	if !ok {
		return "", errors.New("invalid user type in context")
	}

	return userType, nil
}

// MustGetUserID panics if user ID cannot be retrieved (use with caution)
func MustGetUserID(c *gin.Context) uuid.UUID {
	userID, err := GetUserID(c)
	if err != nil {
		panic(err)
	}
	return userID
}
