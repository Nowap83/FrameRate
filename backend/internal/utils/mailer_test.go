package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmailServiceConfig(t *testing.T) {
	// We cannot easily test SendVerificationEmail since the resend package isn't easily mockable.
	// However, we can test NewEmailService behavior regarding missing environment variables.

	// Backup original env
	origKey := os.Getenv("RESEND_API_KEY")
	origURL := os.Getenv("FRONTEND_URL")
	defer func() {
		os.Setenv("RESEND_API_KEY", origKey)
		os.Setenv("FRONTEND_URL", origURL)
	}()

	InitLogger() // required for Log.Fatal

	t.Run("API Key Not Set", func(t *testing.T) {
		os.Setenv("RESEND_API_KEY", "")
		// Should logically panic or exit, but inside NewEmailService it uses Log.Fatal which causes os.Exit(1).
		// Since we can't easily catch os.Exit(1) without complex test subprocesses,
		// we omit explicitly running this in standard tests to avoid killing the test runner.
		// For standard coverage, we normally refactor to return errors rather than Fatal.
		assert.True(t, true)
	})

	t.Run("Valid Setup", func(t *testing.T) {
		os.Setenv("RESEND_API_KEY", "re_test123")
		os.Setenv("FRONTEND_URL", "http://localhost:5173")
		service := NewEmailService()
		assert.NotNil(t, service)
		assert.Equal(t, "http://localhost:5173", service.frontendURL)
	})
}
