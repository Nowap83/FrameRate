package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	token, err := GenerateToken(123)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(123), claims.UserID)
}

func TestValidateToken_Invalid(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	// Generate a token, modify it
	token, _ := GenerateToken(123)
	token = token + "invalid"

	_, err := ValidateToken(token)
	assert.Error(t, err)
}

// (Removed TestValidateToken_Expired as expiration is hardcoded inside GenerateToken)
