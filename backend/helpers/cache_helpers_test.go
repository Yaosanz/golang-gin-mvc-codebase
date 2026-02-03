package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNoOpCache_Set(t *testing.T) {
	cache := NewNoOpCache()
	ctx := context.Background()

	err := cache.Set(ctx, "test-key", "test-value", 1*time.Minute)
	assert.NoError(t, err)
}

func TestNoOpCache_Get(t *testing.T) {
	cache := NewNoOpCache()
	ctx := context.Background()

	var result string
	err := cache.Get(ctx, "test-key", &result)
	assert.Error(t, err) // Always returns error (cache miss)
}

func TestNoOpCache_Delete(t *testing.T) {
	cache := NewNoOpCache()
	ctx := context.Background()

	err := cache.Delete(ctx, "test-key")
	assert.NoError(t, err)
}

func TestNoOpCache_Exists(t *testing.T) {
	cache := NewNoOpCache()
	ctx := context.Background()

	count, err := cache.Exists(ctx, "test-key")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count) // Always returns 0
}

func TestAuthCacheService_CacheUserAuth(t *testing.T) {
	cache := NewNoOpCache()
	authCache := NewAuthCacheService(cache)
	ctx := context.Background()

	userID := uuid.New()
	err := authCache.CacheUserAuth(
		ctx,
		userID,
		"testuser",
		"test@example.com",
		"hashedpassword",
		true,
		[]string{"user"},
		[]string{"user:read"},
	)

	assert.NoError(t, err)
}

func TestAuthorizationService_CheckRole(t *testing.T) {
	cache := NewNoOpCache()
	authzService := NewAuthorizationService(cache)
	ctx := context.Background()

	userID := uuid.New()

	// Cache roles first
	authzService.CacheUserRoles(ctx, userID, []string{"admin", "user"})

	// Check role - with NoOpCache it will fail
	hasRole, err := authzService.CheckRole(ctx, userID, "admin")
	assert.Error(t, err) // NoOpCache always fails
	assert.False(t, hasRole)
}

func TestAuthorizationService_CheckPermission(t *testing.T) {
	cache := NewNoOpCache()
	authzService := NewAuthorizationService(cache)
	ctx := context.Background()

	userID := uuid.New()

	// Cache permissions first
	authzService.CacheUserPermissions(ctx, userID, []string{"user:read", "user:write"})

	// Check permission - with NoOpCache it will fail
	hasPerm, err := authzService.CheckPermission(ctx, userID, "user:read")
	assert.Error(t, err) // NoOpCache always fails
	assert.False(t, hasPerm)
}

