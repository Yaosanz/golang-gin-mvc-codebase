package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter implements sliding window rate limiting using Redis
type RateLimiter struct {
	redis     *redis.Client
	limit     int           // Max requests allowed
	window    time.Duration // Time window
	KeyFunc   func(*gin.Context) string // Custom key function (default: IP-based)
}

// NewRateLimiter creates a new rate limiter
// limit: maximum number of requests allowed
// window: time window duration (e.g., 1*time.Minute)
func NewRateLimiter(redisClient *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redis:  redisClient,
		limit:  limit,
		window: window,
		KeyFunc: func(c *gin.Context) string {
			// Default: rate limit by IP + path
			return fmt.Sprintf("rate_limit:ip:%s:%s", c.ClientIP(), c.Request.URL.Path)
		},
	}
}

// Limit returns the rate limiting middleware
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Graceful degradation: if Redis is unavailable, allow the request
		if rl.redis == nil {
		log.Println("[INFO] Rate limiter: Redis disabled - all requests allowed")
		}

		ctx := context.Background()
		key := rl.KeyFunc(c)
		now := time.Now().Unix()
		windowStart := now - int64(rl.window.Seconds())

		// Use Redis sorted set for sliding window counter
		pipe := rl.redis.Pipeline()

		// Remove old entries outside the window
		pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

		// Add current request with timestamp as score
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: fmt.Sprintf("%d", now),
		})

		// Count requests in current window
		pipe.ZCard(ctx, key)

		// Set expiration to window duration + buffer
		pipe.Expire(ctx, key, rl.window+time.Minute)

		// Execute pipeline
		cmds, err := pipe.Exec(ctx)
		if err != nil {
			log.Printf("[ERROR] Rate limiter: Redis error: %v, allowing request", err)
			c.Next()
			return
		}

		// Get request count from ZCard result
		count := cmds[2].(*redis.IntCmd).Val()
		remaining := rl.limit - int(count)
		if remaining < 0 {
			remaining = 0
		}

		// Calculate reset time
		resetTime := now + int64(rl.window.Seconds())

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		// Check if limit exceeded
		if count > int64(rl.limit) {
			retryAfter := int64(rl.window.Seconds()) - (now - windowStart)
			if retryAfter < 0 {
				retryAfter = 1
			}

			c.Header("Retry-After", strconv.FormatInt(retryAfter, 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": fmt.Sprintf("Rate limit exceeded. Try again in %d seconds", retryAfter),
				"error": gin.H{
					"code":         "RATE_LIMIT_EXCEEDED",
					"limit":        rl.limit,
					"window":       rl.window.String(),
					"retry_after":  retryAfter,
					"reset_at":     time.Unix(resetTime, 0).Format(time.RFC3339),
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ByUserID creates a rate limiter that limits by authenticated user ID
func (rl *RateLimiter) ByUserID() gin.HandlerFunc {
	rl.KeyFunc = func(c *gin.Context) string {
		userID, exists := c.Get("user_id")
		if !exists {
			// Fallback to IP if user not authenticated
			return fmt.Sprintf("rate_limit:ip:%s:%s", c.ClientIP(), c.Request.URL.Path)
		}
		return fmt.Sprintf("rate_limit:user:%v:%s", userID, c.Request.URL.Path)
	}
	return rl.Limit()
}

// Global creates a global rate limiter (all endpoints combined)
func (rl *RateLimiter) Global() gin.HandlerFunc {
	rl.KeyFunc = func(c *gin.Context) string {
		return fmt.Sprintf("rate_limit:global:ip:%s", c.ClientIP())
	}
	return rl.Limit()
}
