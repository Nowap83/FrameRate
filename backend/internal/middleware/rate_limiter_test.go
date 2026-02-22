package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.InitLogger() // required for logging inside RateLimiter

	// Create a tracker for testing that allows 2 requests per second with burst 2.
	testTracker := &IPTracker{
		rate:  rate.Limit(2),
		burst: 2,
	}

	r := gin.New()
	r.Use(RateLimiter(testTracker))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("Allow Requests Within Burst", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:1234"
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("Block Excess Requests", func(t *testing.T) {
		// Next request immediately should be blocked (burst was 2 and we already did 2)
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Contains(t, w.Body.String(), "Too many requests. Please try again later.")
	})

	t.Run("Allow Different IP", func(t *testing.T) {
		// A different IP should be allowed
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:5678" // New IP
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Wait and allow again", func(t *testing.T) {
		// Sleep enough time (more than 1/rate second) to replenish tokens
		// Since rate=2/sec, sleeping ~600ms gives us >1 token.
		time.Sleep(600 * time.Millisecond)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthRateLimiterConfig(t *testing.T) {
	// Let's just ensure the constructors return valid middleware functions
	middleware := AuthRateLimiter()
	assert.NotNil(t, middleware)
}

func TestAPIRateLimiterConfig(t *testing.T) {
	middleware := APIRateLimiter()
	assert.NotNil(t, middleware)
}
