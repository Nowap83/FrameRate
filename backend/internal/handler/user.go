package handler

import (
	"errors"
	"net/http"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/service"
	"github.com/Nowap83/FrameRate/backend/internal/utils" // Import utils for logging
	internalValidator "github.com/Nowap83/FrameRate/backend/internal/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// returns the user's profile with stats
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	response, err := h.userService.GetProfile(userID.(uint))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// updates the user's profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
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

	utils.Log.Info("UpdateProfile Request",
		zap.Int("FavoriteFilmsCount", len(input.FavoriteFilms)),
	)
	if len(input.FavoriteFilms) > 0 {
		utils.Log.Info("First Favorite Film",
			zap.String("Title", input.FavoriteFilms[0].Title),
			zap.Int("TMDB_ID", input.FavoriteFilms[0].TmdbID),
		)
	}

	response, err := h.userService.UpdateProfile(userID.(uint), input)
	if err != nil {
		if errors.Is(err, service.ErrUsernameTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// updates the user's password
func (h *UserHandler) ChangePassword(c *gin.Context) {
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

	if err := h.userService.ChangePassword(userID.(uint), input); err != nil {
		if errors.Is(err, service.ErrPasswordIncorrect) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// deletes the user's account
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.userService.DeleteAccount(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}

// uploads the user's avatar
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// limit file size to 5MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5<<20)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// open the file to check its type
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}
	defer src.Close()

	// sniff the content type
	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	contentType := http.DetectContentType(buffer)

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed."})
		return
	}

	avatarURL, err := h.userService.UpdateAvatar(userID.(uint), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Avatar updated successfully",
		"url":     avatarURL,
	})
}

// checks if a username is available
func (h *UserHandler) CheckUsername(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	available, err := h.userService.CheckUsernameAvailability(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  username,
		"available": available,
	})
}
