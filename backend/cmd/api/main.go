package main

import (
	"os"
	"strings"
	"time"

	"github.com/Nowap83/FrameRate/backend/internal/config"
	"github.com/Nowap83/FrameRate/backend/internal/database"
	"github.com/Nowap83/FrameRate/backend/internal/router"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	internalValidator "github.com/Nowap83/FrameRate/backend/internal/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	utils.InitLogger()
	defer utils.Log.Sync()

	if err := godotenv.Load("../.env"); err != nil {
		utils.Log.Warn("Error loading .env file")
	}

	if err := config.ValidateEnvironment(); err != nil {
		utils.Log.Fatal("Configuration error", zap.Error(err))
	}

	// Ensure uploads directory exists
	if err := os.MkdirAll("./uploads/avatars", 0755); err != nil {
		utils.Log.Fatal("Failed to create uploads directory", zap.Error(err))
	}

	db, err := database.ConnectDB()
	if err != nil {
		utils.Log.Fatal("Database initialization failed", zap.Error(err))
	}

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			utils.Log.Error("Failed to get DB instance", zap.Error(err))
			return
		}
		if err := sqlDB.Close(); err != nil {
			utils.Log.Error("Error closing database", zap.Error(err))
		} else {
			utils.Log.Info("Database connection closed gracefully")
		}
	}()

	database.AutoMigrateAll(db)

	rdb, err := database.ConnectRedis()
	if err != nil {
		utils.Log.Warn("Redis connection failed, continuing without cache", zap.Error(err))
	} else {
		defer rdb.Close()
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		internalValidator.RegisterCustomValidators(v)
	}

	emailService := utils.NewEmailService()

	r := gin.Default()
	allowedOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
	} else {
		allowedOrigins = []string{"http://localhost:5173"}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve static files
	r.Static("/uploads", "./uploads")

	router.SetupRoutes(r, db, rdb, emailService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.Log.Info("Server starting",
		zap.String("port", port),
		zap.String("env", os.Getenv("ENV")),
	)

	if err := r.Run(":" + port); err != nil {
		utils.Log.Fatal("Failed to start server", zap.Error(err))
	}
}
