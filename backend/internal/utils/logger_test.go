package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	t.Run("Development Mode", func(t *testing.T) {
		os.Setenv("ENV", "development")
		InitLogger()
		assert.NotNil(t, Log)
		Log.Info("Dev logger test")
	})

	t.Run("Production Mode", func(t *testing.T) {
		os.Setenv("ENV", "production")
		InitLogger()
		assert.NotNil(t, Log)
		Log.Info("Prod logger test")
	})
}
