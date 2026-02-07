package middleware

import (
	"net/http"
	"strings"

	"github.com/Nowap83/FrameRate/backend/utils"
	"github.com/gin-gonic/gin"
)

// vérifie le JWT
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// recup le header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// format => "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// parse et valide le token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		// garde user id dans le contexte
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
