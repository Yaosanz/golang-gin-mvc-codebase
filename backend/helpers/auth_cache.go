package helpers

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserAuthData represents cached authentication data for fast login
type UserAuthData struct {
	UserID       string   `json:"user_id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	PasswordHash string   `json:"password_hash"`
	IsActive     bool     `json:"is_active"`
	Roles        []string `json:"roles"`
	Permissions  []string `json:"permissions"`
}

// NewUserAuthData creates authentication cache data
func NewUserAuthData(userID uuid.UUID, username, email, passwordHash string, isActive bool, roles, permissions []string) *UserAuthData {
	return &UserAuthData{
		UserID:       userID.String(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		IsActive:     isActive,
		Roles:        roles,
		Permissions:  permissions,
	}
}

// ToAuthData converts cache data for authentication
func (uad *UserAuthData) ToAuthData() (uuid.UUID, string, string, string, bool, []string, []string, error) {
	userID, err := uuid.Parse(uad.UserID)
	if err != nil {
		return uuid.Nil, "", "", "", false, nil, nil, err
	}

	return userID, uad.Username, uad.Email, uad.PasswordHash, uad.IsActive, uad.Roles, uad.Permissions, nil
}

// AuthCacheService provides optimized caching for authentication
type AuthCacheService struct {
	cache CacheInterface
}

// NewAuthCacheService creates a new auth cache service
func NewAuthCacheService(cache CacheInterface) *AuthCacheService {
	return &AuthCacheService{cache: cache}
}

// CacheUserAuth caches user authentication data
func (acs *AuthCacheService) CacheUserAuth(ctx context.Context, userID uuid.UUID, username, email, passwordHash string, isActive bool, roles, permissions []string) error {
	authData := NewUserAuthData(userID, username, email, passwordHash, isActive, roles, permissions)
	cacheKey := UserAuthCacheKey(username)
	return acs.cache.Set(ctx, cacheKey, authData, 10*time.Minute) // Shorter TTL for auth data
}

// GetUserAuth retrieves cached authentication data
func (acs *AuthCacheService) GetUserAuth(ctx context.Context, username string) (*UserAuthData, error) {
	cacheKey := UserAuthCacheKey(username)
	var authData UserAuthData
	err := acs.cache.Get(ctx, cacheKey, &authData)
	if err != nil {
		return nil, err
	}
	return &authData, nil
}

// InvalidateUserAuth removes user auth cache
func (acs *AuthCacheService) InvalidateUserAuth(ctx context.Context, username string) error {
	cacheKey := UserAuthCacheKey(username)
	return acs.cache.Delete(ctx, cacheKey)
}

// Cache keys for auth
func UserAuthCacheKey(username string) string {
	return NewCacheKeyBuilder(CachePrefixAuth).Build("user", username)
}
