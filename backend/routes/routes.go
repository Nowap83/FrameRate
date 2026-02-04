package routes

import (
    "github.com/Nowap83/FrameRate/backend/middleware"
		"github.com/Nowap83/FrameRate/backend/handlers"
		"net/http"
		"gorm.io/gorm"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

		authHandler := handlers.NewAuthHandler(db)	

    // Health check (verif serveur)
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
        })
    })

    // Groupe API (ref swagger)
    api := r.Group("/api")
    {
        // Auth (publiques)
        auth := api.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)

        }

        // Routes protégées 
        protected := api.Group("")
        protected.Use(middleware.AuthRequired())
        {
            // Users
            users := protected.Group("/users")
            {            
          			users.GET("/test", func(c *gin.Context) {
                    userID := c.GetUint("userID")
                    c.JSON(200, gin.H{
                        "message": "protected route ok",
                        "user_id": userID,
                    })
                })
            }
        }
    }
}

