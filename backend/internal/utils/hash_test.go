package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "my_secret_password"
	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	isValid := CheckPassword(password, hash)
	assert.True(t, isValid)
}

func TestCheckPassword_Invalid(t *testing.T) {
	password := "my_secret_password"
	hash, _ := HashPassword(password)

	isValid := CheckPassword("wrong_password", hash)
	assert.False(t, isValid)
}
