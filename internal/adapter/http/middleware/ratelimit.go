
package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"mymodule/internal/config"
)

// IPRateLimiter holds the rate limiters for each IP address.
type IPRateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	config   config.RateLimitConfig
}

// NewIPRateLimiter creates a new IPRateLimiter.
func NewIPRateLimiter(cfg config.RateLimitConfig) *IPRateLimiter {
	return &IPRateLimiter{
		visitors: make(map[string]*rate.Limiter),
		config:   cfg,
	}
}

// getVisitor returns the rate limiter for the given IP address.
func (rl *IPRateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check after acquiring the lock.
		if limiter, exists = rl.visitors[ip]; !exists {
			limiter = rate.NewLimiter(rate.Limit(rl.config.RPS), rl.config.Burst)
			rl.visitors[ip] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

// RateLimitMiddleware applies rate limiting based on IP address.
func (rl *IPRateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.config.Enabled {
			c.Next()
			return
		}

		ip := c.ClientIP()
		limiter := rl.getVisitor(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
