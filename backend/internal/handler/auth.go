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
				"errors": map[string]string{"email": "email already exists"},
			})
		case "username already exists":
			c.JSON(http.StatusConflict, gin.H{
				"errors": map[string]string{"username": "username already exists"},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "registration failed",
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

//
// GET PROFILE (
//

func (h *AuthHandler) GetProfile(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

//
// UPDATE PROFILE
//

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.UpdateProfile(userID.(uint), input)
	if err != nil {
		switch err.Error() {
		case "username already taken":
			c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

//
// CHANGE PASSWORD
//

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		if validationErrors := internalValidator.FormatValidationErrors(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ChangePassword(userID.(uint), input); err != nil {
		switch err.Error() {
		case "current password is incorrect":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

//
// DELETE ACCOUNT
//

func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.authService.DeleteAccount(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}
