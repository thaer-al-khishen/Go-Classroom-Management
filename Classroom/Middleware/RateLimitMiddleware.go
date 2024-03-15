package Middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiters map[string]*limiterInfo
	mu       sync.Mutex
}

type limiterInfo struct {
	limiter *rate.Limiter
	timer   *time.Timer // Timer to track expiration
}

// Define the rate limit: 30 requests per minute (0.5 requests per second)
var (
	r = rate.Limit(0.5) // 30 requests per minute
	b = 5               // Burst size
)

var limiterInstance = newIPLimiter()

func newIPLimiter() *ipLimiter {
	return &ipLimiter{
		limiters: make(map[string]*limiterInfo),
	}
}

func (l *ipLimiter) getLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	if info, exists := l.limiters[ip]; exists {
		// Limiter exists, reset the deletion timer
		info.timer.Reset(10 * time.Minute)
		return info.limiter
	}

	// Create a new limiter and schedule its deletion
	lim := rate.NewLimiter(r, b)
	timer := time.AfterFunc(10*time.Minute, func() {
		l.removeLimiter(ip)
	})
	l.limiters[ip] = &limiterInfo{lim, timer}
	return lim
}

func (l *ipLimiter) removeLimiter(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.limiters, ip)
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if limiter := limiterInstance.getLimiter(ip); !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}
