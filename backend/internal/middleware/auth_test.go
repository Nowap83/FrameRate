package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthRequired_NoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a dummy request with no Authorization header
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	// Execute middleware
	AuthRequired()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
	assert.True(t, c.IsAborted())
}

func TestAuthRequired_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "InvalidFormatToken")

	AuthRequired()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid authorization format")
	assert.True(t, c.IsAborted())
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer faketoken123")

	AuthRequired()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestAuthRequired_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	token, _ := utils.GenerateToken(123)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	AuthRequired()(c)

	assert.False(t, c.IsAborted())

	// Check if userID was set in context
	userID, exists := c.Get("userID")
	assert.True(t, exists)
	assert.Equal(t, uint(123), userID)
}
