package routes

import (
    "github.com/Nowap83/FrameRate/backend/config"
    "github.com/Nowap83/FrameRate/backend/middleware"
		"net/http"

    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    // Health check (verif serveur)
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
        })
    })

    // Ping (verif co DB)
    r.GET("/ping", func(c *gin.Context) {
        sqlDB, err := config.DB.DB()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message":  "pong",
                "database": "error",
                "error":    err.Error(),
            })
            return
        }

        if err := sqlDB.Ping(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message":  "pong",
                "database": "disconnected",
                "error":    err.Error(),
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message":  "pong",
            "database": "connected",
        })
    })

    // Groupe API (ref swagger)
    api := r.Group("/api")
    {
        // Auth (publiques)
        auth := api.Group("/auth")
        {
            // ROUTE TEMPO DE TEST A SUPPR APRES AVOIR CODE REGISTER/LOGIN
            auth.GET("/generate-test-token", func(c *gin.Context) {
                token, err := middleware.GenerateToken(2) 
                if err != nil {
                    c.JSON(500, gin.H{"error": "Failed to generate token"})
                    return
                }
                c.JSON(200, gin.H{
                    "token": token,
                    "message": "test token !",
                    "user_id": 2,
                })
            })

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

