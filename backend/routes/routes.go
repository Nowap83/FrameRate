package routes

import (
	"github.com/Nowap83/FrameRate/backend/handlers"
	"github.com/Nowap83/FrameRate/backend/middleware"
	"github.com/Nowap83/FrameRate/backend/services"
	"github.com/Nowap83/FrameRate/backend/utils"
	"github.com/gin-gonic/gin"

	"net/http"

	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, emailService *utils.EmailService) {

	authService := services.NewAuthService(db, emailService)

	authHandler := handlers.NewAuthHandler(authService)

	tmdbService := services.NewTMDBService()

	tmdbHandler := handlers.NewTMDBHandler(tmdbService)

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
			auth.GET("/verify-email", authHandler.VerifyEmail)
		}

		// TMDB
		tmdb := api.Group("/tmdb")
		{
			tmdb.GET("/search", tmdbHandler.SearchMovies)
			tmdb.GET("/movie/:id", tmdbHandler.GetMovieDetails)
			tmdb.GET("/movie/:id/credits", tmdbHandler.GetMovieCredits)
			tmdb.GET("/image", tmdbHandler.GetImageURL)
		}

		// Routes protégées
		protected := api.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// Users
			users := protected.Group("/users")
			{
				users.GET("/me", authHandler.GetProfile)
				users.PUT("/me", authHandler.UpdateProfile)
				users.PUT("/me/password", authHandler.ChangePassword)
				users.DELETE("/me", authHandler.DeleteAccount)
			}
		}
	}
}
