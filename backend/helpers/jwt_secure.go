package helpers

import (
	"context"
	"errors"
	"go-starter-app/app/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// SecureJwtClaims - Minimal, secure JWT claims following best practices
type SecureJwtClaims struct {
	UserID    string   `json:"user_id"`    // Only essential identifier
	TokenType string   `json:"token_type"` // "user" or "cms"
	SessionID string   `json:"sid"`        // Session identifier for revocation
	Roles     []string `json:"roles"`      // User roles for permission checking

	jwt.RegisteredClaims
}

// JwtClaims alias for backward compatibility
type JwtClaims = SecureJwtClaims

// GenerateSecureUserToken creates a minimal, secure JWT token
// Only includes essential data - all authorization done server-side
func GenerateSecureUserToken(
	userID uuid.UUID,
	tokenType string,
	sessionID string,
	roles []string,
	secret string,
	expired time.Duration,
	issuer string,
) (string, error) {

	claims := SecureJwtClaims{
		UserID:    userID.String(),
		TokenType: tokenType,
		SessionID: sessionID,
		Roles:     roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expired)),
			Subject:   userID.String(),
			ID:        sessionID, // Use session ID as JWT ID for revocation
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateSecureJWT validates and parses secure JWT token
func ValidateSecureJWT(tokenString string, secret string) (*SecureJwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SecureJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SecureJwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GetUserID returns the user ID as UUID
func (c *SecureJwtClaims) GetUserID() (uuid.UUID, error) {
	return uuid.Parse(c.UserID)
}

// IsUserToken checks if this is a user token
func (c *SecureJwtClaims) IsUserToken() bool {
	return c.TokenType == "user"
}

// IsCMSToken checks if this is a CMS token
func (c *SecureJwtClaims) IsCMSToken() bool {
	return c.TokenType == "cms"
}

// HasPermission - For backward compatibility, always returns false (server-side only)
func (c *SecureJwtClaims) HasPermission(permission string) bool {
	return false // Permissions checked server-side only
}

// HasRole checks if user has the specified role
func (c *SecureJwtClaims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// GenerateSessionID creates a unique session identifier
func GenerateSessionID() string {
	return uuid.New().String()
}

// AuthorizationService - Server-side authorization (recommended approach)
type AuthorizationService struct {
	cache CacheInterface
}

// NewAuthorizationService creates authorization service
func NewAuthorizationService(cache CacheInterface) *AuthorizationService {
	return &AuthorizationService{cache: cache}
}

// CacheUserPermissions caches user permissions for authorization
func (as *AuthorizationService) CacheUserPermissions(ctx context.Context, userID uuid.UUID, permissions []string) error {
	cacheKey := UserPermissionsCacheKey(userID.String())
	return as.cache.Set(ctx, cacheKey, permissions, 10*time.Minute)
}

// CacheUserRoles caches user roles for authorization
func (as *AuthorizationService) CacheUserRoles(ctx context.Context, userID uuid.UUID, roles []string) error {
	cacheKey := UserRolesCacheKey(userID.String())
	return as.cache.Set(ctx, cacheKey, roles, 10*time.Minute)
}

// CheckPermission verifies if user has permission (server-side only)
func (as *AuthorizationService) CheckPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	cacheKey := UserPermissionsCacheKey(userID.String())

	var permissions []string
	err := as.cache.Get(ctx, cacheKey, &permissions)
	if err != nil {
		// Cache miss - this should not happen in normal operation
		// Permissions should be cached during login
		return false, errors.New("permissions not cached - please login again")
	}

	for _, p := range permissions {
		if p == permission {
			return true, nil
		}
	}
	return false, nil
}

// CheckRole verifies if user has role (server-side only)
func (as *AuthorizationService) CheckRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error) {
	cacheKey := UserRolesCacheKey(userID.String())

	var roles []string
	err := as.cache.Get(ctx, cacheKey, &roles)
	if err != nil {
		return false, errors.New("unable to verify roles")
	}

	for _, r := range roles {
		if r == roleName {
			return true, nil
		}
	}
	return false, nil
}

// GetUserPermissions returns all user permissions (server-side only)
func (as *AuthorizationService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	cacheKey := UserPermissionsCacheKey(userID.String())

	var permissions []string
	err := as.cache.Get(ctx, cacheKey, &permissions)
	if err != nil {
		return nil, errors.New("unable to retrieve permissions")
	}

	return permissions, nil
}

// GetUserRoles returns all user roles (server-side only)
func (as *AuthorizationService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	cacheKey := UserRolesCacheKey(userID.String())

	var roles []string
	err := as.cache.Get(ctx, cacheKey, &roles)
	if err != nil {
		return nil, errors.New("unable to retrieve roles")
	}

	return roles, nil
}

// IsUserActive checks if user account is active (server-side only)
// Returns cache miss error that should be handled by caller with DB fallback
func (as *AuthorizationService) IsUserActive(ctx context.Context, userID uuid.UUID) (bool, error) {
	cacheKey := UserCacheKey(userID.String())

	var user models.UserCacheData
	err := as.cache.Get(ctx, cacheKey, &user)
	if err != nil {
		// Return error so caller can fallback to database
		return false, err
	}

	return user.IsActive, nil
}

// Middleware helper functions for secure authorization

// RequirePermission middleware helper
func RequirePermission(authSvc *AuthorizationService, permission string) func(ctx context.Context, userID uuid.UUID) error {
	return func(ctx context.Context, userID uuid.UUID) error {
		hasPermission, err := authSvc.CheckPermission(ctx, userID, permission)
		if err != nil {
			return err
		}
		if !hasPermission {
			return errors.New("insufficient permissions")
		}
		return nil
	}
}

// RequireRole middleware helper
func RequireRole(authSvc *AuthorizationService, role string) func(ctx context.Context, userID uuid.UUID) error {
	return func(ctx context.Context, userID uuid.UUID) error {
		hasRole, err := authSvc.CheckRole(ctx, userID, role)
		if err != nil {
			return err
		}
		if !hasRole {
			return errors.New("insufficient role")
		}
		return nil
	}
}

// RequireActiveUser middleware helper
func RequireActiveUser(authSvc *AuthorizationService) func(ctx context.Context, userID uuid.UUID) error {
	return func(ctx context.Context, userID uuid.UUID) error {
		isActive, err := authSvc.IsUserActive(ctx, userID)
		if err != nil {
			return err
		}
		if !isActive {
			return errors.New("account is inactive")
		}
		return nil
	}
}
