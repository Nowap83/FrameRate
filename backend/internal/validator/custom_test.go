package validator

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestCustomValidators(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	type UserInput struct {
		Username string `validate:"username"`
		Password string `validate:"strongpassword"`
	}

	t.Run("Valid Input", func(t *testing.T) {
		input := UserInput{
			Username: "valid_name",
			Password: "Valid1Password!",
		}
		err := v.Struct(input)
		assert.NoError(t, err)
	})

	t.Run("Invalid Username", func(t *testing.T) {
		inputs := []string{
			"ab",            // too short
			"invalid name!", // invalid chars
		}

		for _, name := range inputs {
			err := v.Struct(UserInput{Username: name, Password: "Valid1Password!"})
			assert.Error(t, err)
			var valErr validator.ValidationErrors
			errors.As(err, &valErr)
			assert.Equal(t, "Username", valErr[0].Field())
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		inputs := []string{
			"short",            // length < 8
			"nouppercase1!",    // no upper
			"NOLOWERCASE1!",    // no lower
			"NoNumberPass!",    // no number
			"NoSpecialChar123", // no special
		}

		for _, pass := range inputs {
			err := v.Struct(UserInput{Username: "valid_name", Password: pass})
			assert.Error(t, err)
			var valErr validator.ValidationErrors
			errors.As(err, &valErr)
			assert.Equal(t, "Password", valErr[0].Field())
		}
	})
}

func TestFormatValidationErrors(t *testing.T) {
	v := validator.New()
	RegisterCustomValidators(v)

	type TestStruct struct {
		ReqField  string `validate:"required"`
		Email     string `validate:"email"`
		MinField  string `validate:"min=5"`
		MaxField  string `validate:"max=10"`
		Username  string `validate:"username"`
		Password  string `validate:"strongpassword"`
		SomeField string `validate:"url"` // to hit default case
	}

	t.Run("Test Formatting", func(t *testing.T) {
		input := TestStruct{
			ReqField:  "",
			Email:     "invalidemail",
			MinField:  "123",
			MaxField:  "thisiswaytoolong",
			Username:  "a", // invalid username
			Password:  "weak",
			SomeField: "not-a-url",
		}

		err := v.Struct(input)
		assert.Error(t, err)

		formatted := FormatValidationErrors(err)

		assert.Equal(t, "reqfield is required", formatted["reqfield"])
		assert.Equal(t, "invalid email format", formatted["email"])
		assert.Equal(t, "minfield must be at least 5 characters", formatted["minfield"])
		assert.Equal(t, "maxfield must be at most 10 characters", formatted["maxfield"])
		assert.Equal(t, "username must be 3-50 characters (letters, numbers, _ and - only)", formatted["username"])
		assert.Equal(t, "password must contain at least 8 characters, 1 uppercase, 1 lowercase, 1 number and 1 special character", formatted["password"])
		assert.Equal(t, "invalid somefield", formatted["somefield"])
	})

	t.Run("Generic Error Formatting", func(t *testing.T) {
		genericErr := errors.New("just a random error")
		formatted := FormatValidationErrors(genericErr)
		assert.Equal(t, "invalid request data", formatted["error"])
	})
}
