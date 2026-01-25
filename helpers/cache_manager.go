package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheManager provides advanced cache management with monitoring and statistics
type CacheManager struct {
	cache  CacheInterface
	mu     sync.RWMutex
	stats  *CacheStats
	logger CacheLogger
}

// CacheStats tracks cache performance metrics
type CacheStats struct {
	Hits      int64
	Misses    int64
	Errors    int64
	Evictions int64
	LastReset time.Time
}

// CacheLogger interface for customizable logging
type CacheLogger interface {
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Info(msg string, args ...interface{})
}

// DefaultLogger is a simple logger implementation
type DefaultLogger struct{}

func (d *DefaultLogger) Debug(msg string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+msg+"\n", args...)
}

func (d *DefaultLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("[ERROR] "+msg+"\n", args...)
}

func (d *DefaultLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("[INFO] "+msg+"\n", args...)
}

// NewCacheManager creates a new cache manager instance
func NewCacheManager(cache CacheInterface) *CacheManager {
	return &CacheManager{
		cache:  cache,
		stats:  &CacheStats{LastReset: time.Now()},
		logger: &DefaultLogger{},
	}
}

// NewCacheManagerWithLogger creates a new cache manager with custom logger
func NewCacheManagerWithLogger(cache CacheInterface, logger CacheLogger) *CacheManager {
	return &CacheManager{
		cache:  cache,
		stats:  &CacheStats{LastReset: time.Now()},
		logger: logger,
	}
}

// GetOrSet retrieves value from cache, or sets it if not found (atomic operation)
func (cm *CacheManager) GetOrSet(
	ctx context.Context,
	key string,
	ttl time.Duration,
	getter func() (interface{}, error),
	dest interface{},
) error {
	// Try to get from cache first
	err := cm.cache.Get(ctx, key, dest)
	if err == nil {
		cm.mu.Lock()
		cm.stats.Hits++
		cm.mu.Unlock()
		cm.logger.Debug("Cache hit for key: %s", key)
		return nil
	}

	// Cache miss or error
	cm.mu.Lock()
	cm.stats.Misses++
	cm.mu.Unlock()

	// Execute getter function to fetch fresh data
	cm.logger.Debug("Cache miss for key: %s, fetching fresh data", key)
	value, err := getter()
	if err != nil {
		cm.mu.Lock()
		cm.stats.Errors++
		cm.mu.Unlock()
		cm.logger.Error("Error fetching value for key %s: %v", key, err)
		return err
	}

	// Set in cache
	if err := cm.cache.Set(ctx, key, value, ttl); err != nil {
		cm.mu.Lock()
		cm.stats.Errors++
		cm.mu.Unlock()
		cm.logger.Error("Error setting cache for key %s: %v", key, err)
		// Don't fail operation, just log error
	}

	// Marshal value to dest
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling value: %w", err)
	}

	return json.Unmarshal(data, dest)
}

// Invalidate removes keys from cache
func (cm *CacheManager) Invalidate(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	err := cm.cache.Delete(ctx, keys...)
	if err != nil && err != redis.Nil {
		cm.mu.Lock()
		cm.stats.Errors++
		cm.mu.Unlock()
		cm.logger.Error("Error invalidating cache keys: %v", err)
		return err
	}

	cm.logger.Info("Invalidated %d cache keys", len(keys))
	return nil
}

// InvalidatePattern invalidates all keys matching a pattern (for Redis)
func (cm *CacheManager) InvalidatePattern(ctx context.Context, pattern string) error {
	// This is a Redis-specific operation
	if redisClient, ok := cm.cache.(*RedisClient); ok {
		return redisClient.DeleteByPattern(ctx, pattern)
	}

	cm.logger.Debug("Pattern invalidation not supported for this cache backend")
	return nil
}

// GetStats returns current cache statistics
func (cm *CacheManager) GetStats() *CacheStats {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Create a copy to prevent external modification
	stats := *cm.stats
	return &stats
}

// ResetStats resets cache statistics
func (cm *CacheManager) ResetStats() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.stats = &CacheStats{LastReset: time.Now()}
	cm.logger.Info("Cache statistics reset")
}

// PrintStats logs current cache statistics
func (cm *CacheManager) PrintStats() {
	stats := cm.GetStats()
	total := stats.Hits + stats.Misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(stats.Hits) / float64(total) * 100
	}

	cm.logger.Info(
		"Cache Stats - Hits: %d, Misses: %d, Errors: %d, Evictions: %d, Hit Rate: %.2f%%, Since: %v",
		stats.Hits, stats.Misses, stats.Errors, stats.Evictions, hitRate, stats.LastReset,
	)
}

// Close closes the cache manager and underlying cache
func (cm *CacheManager) Close() error {
	return cm.cache.Close()
}
