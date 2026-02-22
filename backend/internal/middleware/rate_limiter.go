package middleware

import (
	"net/http"
	"sync"

	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// IPTracker maps IPs to rate limiters
type IPTracker struct {
	limiters sync.Map
	rate     rate.Limit // Requests per second
	burst    int        // Max burst size
}

// Global instance for auth endpoints (1 request per second, burst 5)
// This is strict to prevent brute force
var authLimiter = &IPTracker{
	rate:  rate.Limit(1),
	burst: 5,
}

// Global instance for general API endpoints (10 requests per second, burst 30)
var apiLimiter = &IPTracker{
	rate:  rate.Limit(10),
	burst: 30,
}

func (i *IPTracker) getLimiter(ip string) *rate.Limiter {
	limiter, exists := i.limiters.Load(ip)
	if !exists {
		newLimiter := rate.NewLimiter(i.rate, i.burst)
		i.limiters.Store(ip, newLimiter)
		return newLimiter
	}
	return limiter.(*rate.Limiter)
}

// RateLimiter creates a middleware with a specific rate tracking instance
func RateLimiter(tracker *IPTracker) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := tracker.getLimiter(ip)

		if !limiter.Allow() {
			utils.Log.Warn("Rate limit exceeded",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRateLimiter is specifically configured for login/register endpoints
func AuthRateLimiter() gin.HandlerFunc {
	return RateLimiter(authLimiter)
}

// APIRateLimiter is the general limiter for other API routes
func APIRateLimiter() gin.HandlerFunc {
	return RateLimiter(apiLimiter)
}
