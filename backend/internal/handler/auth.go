package handler

import (
	"net/http"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/service"
	internalValidator "github.com/Nowap83/FrameRate/backend/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

//
// REGISTER
//

func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": internalValidator.FormatValidationErrors(validationErr),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(input)
	if err != nil {
		switch err.Error() {
		case "email already exists":
			c.JSON(http.StatusConflict, gin.H{
				"errors": map[string]string{"email": "Email already exists"},
			})
		case "username already exists":
			c.JSON(http.StatusConflict, gin.H{
				"errors": map[string]string{"username": "Username already exists"},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Registration failed",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

//
// LOGIN
//

func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": internalValidator.FormatValidationErrors(validationErr),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON format",
		})
		return
	}
	response, err := h.authService.Login(input)
	if err != nil {
		switch err.Error() {
		case "invalid credentials":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email/username or password"})
		case "email not verified. please check your inbox":
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Please verify your email before logging in. Check your inbox.",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

//
// VERIFY EMAIL
//

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	response, err := h.authService.VerifyEmail(token)
	if err != nil {
		switch err.Error() {
		case "invalid or expired verification token":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification token"})
		case "failed to generate token":
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate login token"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email verification failed"})
		}
		return
	}

	c.JSON(http.StatusOK, response.Message)
}
