package validator

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// enregistre tous les validateurs perso
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("strongpassword", validateStrongPassword)
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	match, _ := regexp.MatchString(`^[a-zA-Z0-9_-]{3,50}$`, username)
	return match
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// formate les erreurs de validation en map lisible
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// erreur générique
		errors["error"] = "invalid request data"
		return errors
	}

	// erreur de validation,
	// ? check plus tard si min/max pertinent
	for _, e := range validationErrors {
		field := strings.ToLower(e.Field())

		switch e.Tag() {
		case "required":
			errors[field] = field + " is required"
		case "email":
			errors[field] = "invalid email format"
		case "min":
			errors[field] = field + " must be at least " + e.Param() + " characters"
		case "max":
			errors[field] = field + " must be at most " + e.Param() + " characters"
		case "username":
			errors[field] = "username must be 3-50 characters (letters, numbers, _ and - only)"
		case "strongpassword":
			errors[field] = "password must contain at least 8 characters, 1 uppercase, 1 lowercase, 1 number and 1 special character"
		default:
			errors[field] = "invalid " + field
		}
	}

	return errors
}
