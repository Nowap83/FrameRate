package routes

import (
    "github.com/Nowap83/FrameRate/backend/config"
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
        api.GET("/test", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "message": "API is working!",
            })
        })
    }
}
