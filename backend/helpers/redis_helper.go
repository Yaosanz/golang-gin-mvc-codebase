package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps Redis operations with context
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client wrapper with optimized settings
func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		// Optimized connection pooling
		PoolSize:     10,  // Number of connections
		MinIdleConns: 5,   // Minimum idle connections
		MaxIdleConns: 10,  // Maximum idle connections
		// Connection timeouts
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		// Retry logic
		MaxRetries: 3,
		// Connection validation
		PoolTimeout: 4 * time.Second,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		// Log error but don't fail - allow fallback to no-op cache
		fmt.Printf("Redis connection failed: %v\n", err)
		return &RedisClient{client: nil}
	}

	return &RedisClient{
		client: rdb,
	}
}

// CacheInterface defines caching operations
type CacheInterface interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
	FlushDB(ctx context.Context) error
	Close() error
}

// NoOpCache is a cache implementation that does nothing (used when Redis is unavailable)
type NoOpCache struct{}

// NewNoOpCache creates a new no-op cache
func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

// Set does nothing and returns nil
func (n *NoOpCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

// Get always returns an error (cache miss)
func (n *NoOpCache) Get(ctx context.Context, key string, dest interface{}) error {
	return redis.Nil
}

// Delete does nothing and returns nil
func (n *NoOpCache) Delete(ctx context.Context, keys ...string) error {
	return nil
}

// Exists always returns 0 (no keys exist)
func (n *NoOpCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return 0, nil
}

// Expire does nothing and returns nil
func (n *NoOpCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return nil
}

// FlushDB does nothing and returns nil
func (n *NoOpCache) FlushDB(ctx context.Context) error {
	return nil
}

// Close does nothing and returns nil
func (n *NoOpCache) Close() error {
	return nil
}

// Set stores a value in cache with TTL
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value from cache
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete removes keys from cache
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists checks if keys exist in cache
func (r *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Expire sets expiration on a key
func (r *RedisClient) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

// FlushDB clears the entire database
func (r *RedisClient) FlushDB(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// DeleteByPattern deletes keys matching a pattern (Redis-specific)
func (r *RedisClient) DeleteByPattern(ctx context.Context, pattern string) error {
	if r.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	// Use SCAN to find keys matching pattern (non-blocking alternative to KEYS)
	var cursor uint64
	var keys []string

	for {
		scanResult, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("scan error: %w", err)
		}

		keys = append(keys, scanResult...)
		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}

// Incr increments a counter
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// IncrBy increments by specific value
func (r *RedisClient) IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	return r.client.IncrBy(ctx, key, increment).Result()
}

// GetInt retrieves and returns as int64
func (r *RedisClient) GetInt(ctx context.Context, key string) (int64, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	var value int64
	if err := json.Unmarshal([]byte(result), &value); err != nil {
		// Try simple string conversion
		_, err := r.client.Get(ctx, key).Int64()
		return 0, err
	}

	return value, nil
}

// Ping checks if Redis is healthy
func (r *RedisClient) Ping(ctx context.Context) error {
	if r.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return r.client.Ping(ctx).Err()
}

// CacheKeyBuilder helps build cache keys
type CacheKeyBuilder struct {
	prefix string
}

// NewCacheKeyBuilder creates a new cache key builder
func NewCacheKeyBuilder(prefix string) *CacheKeyBuilder {
	return &CacheKeyBuilder{prefix: prefix}
}

// Build creates a cache key with the configured prefix
func (b *CacheKeyBuilder) Build(parts ...string) string {
	key := b.prefix
	for _, part := range parts {
		key += ":" + part
	}
	return key
}

// Common cache key patterns
const (
	CachePrefixUser     = "user"
	CachePrefixShorten  = "shorten"
	CachePrefixAuth     = "auth"
)

// User cache keys
func UserCacheKey(userID string) string {
	return NewCacheKeyBuilder(CachePrefixUser).Build(userID)
}

func UserByUsernameCacheKey(username string) string {
	return NewCacheKeyBuilder(CachePrefixUser).Build("username", username)
}

func UserByEmailCacheKey(email string) string {
	return NewCacheKeyBuilder(CachePrefixUser).Build("email", email)
}

// Shorten link cache keys
func ShortenLinkCacheKey(id string) string {
	return NewCacheKeyBuilder(CachePrefixShorten).Build("link", id)
}

func ShortenLinkByCodeCacheKey(code string) string {
	return NewCacheKeyBuilder(CachePrefixShorten).Build("code", code)
}

func UserShortenLinksCacheKey(userID string) string {
	return NewCacheKeyBuilder(CachePrefixShorten).Build("user", userID, "links")
}

// Auth cache keys
func UserPermissionsCacheKey(userID string) string {
	return NewCacheKeyBuilder(CachePrefixAuth).Build("permissions", userID)
}

func UserRolesCacheKey(userID string) string {
	return NewCacheKeyBuilder(CachePrefixAuth).Build("roles", userID)
}
