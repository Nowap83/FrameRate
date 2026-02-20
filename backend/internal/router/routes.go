package router

import (
	"github.com/Nowap83/FrameRate/backend/internal/handler"
	"github.com/Nowap83/FrameRate/backend/internal/middleware"
	"github.com/Nowap83/FrameRate/backend/internal/repository"
	"github.com/Nowap83/FrameRate/backend/internal/service"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"github.com/gin-gonic/gin"

	"net/http"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, rdb *redis.Client, emailService *utils.EmailService) {

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, emailService)

	authHandler := handler.NewAuthHandler(authService)

	cacheService := service.NewCacheService(rdb)
	tmdbService := service.NewTMDBService(cacheService)

	tmdbHandler := handler.NewTMDBHandler(tmdbService)

	movieRepo := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepo, tmdbService)
	movieHandler := handler.NewMovieHandler(movieService)

	userService := service.NewUserService(userRepo, movieRepo)
	userHandler := handler.NewUserHandler(userService)

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
			tmdb.GET("/popular", tmdbHandler.GetPopularMovies)
			tmdb.GET("/movie/:id", tmdbHandler.GetMovieDetails)
			tmdb.GET("/movie/:id/credits", tmdbHandler.GetMovieCredits)
			tmdb.GET("/movie/:id/videos", tmdbHandler.GetMovieVideos)
			tmdb.GET("/person/:id", tmdbHandler.GetPersonDetails)
			tmdb.GET("/person/:id/credits", tmdbHandler.GetPersonMovieCredits)
			tmdb.GET("/image", tmdbHandler.GetImageURL)
		}

		// Routes protégées
		protected := api.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// Users
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetProfile)
				users.PUT("/me", userHandler.UpdateProfile)
				users.POST("/me/avatar", userHandler.UploadAvatar)
				users.PUT("/me/password", userHandler.ChangePassword)
				users.DELETE("/me", userHandler.DeleteAccount)
				users.GET("/check-username", userHandler.CheckUsername)
			}

			// Movies (tracking, rating, review)
			movies := protected.Group("/movies")
			{
				movies.POST("/:tmdb_id/track", movieHandler.TrackMovie)
				movies.POST("/:tmdb_id/rate", movieHandler.RateMovie)
				movies.POST("/:tmdb_id/review", movieHandler.ReviewMovie)
			}
		}
	}
}
