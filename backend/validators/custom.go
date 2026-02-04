package validators

import (
    "regexp"
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
