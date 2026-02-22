package utils

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateVerificationToken(t *testing.T) {
	token, err := GenerateVerificationToken()

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Len(t, token, 64) // 32 bytes hex encoded = 64 chars

	// Check valid hex
	_, err = hex.DecodeString(token)
	assert.NoError(t, err)

	// Check uniqueness
	token2, err := GenerateVerificationToken()
	assert.NoError(t, err)
	assert.NotEqual(t, token, token2)
}
