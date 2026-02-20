package middleware

import (
	"net/http"

	"github.com/Nowap83/FrameRate/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// AdminRequired checks if the authenticated user is an admin
func AdminRequired(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		user, err := userRepo.GetByID(userID.(uint))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
