package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// ThrottleMiddleware simple rate limiting middleware
func ThrottleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Since(start)
		log.Printf("[MIDDLEWARE] Request took %v", elapsed)
	}
}
