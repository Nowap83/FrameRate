package main

import (
	"log"
	"os"
	"time"

	"github.com/Nowap83/FrameRate/backend/config"
  "github.com/Nowap83/FrameRate/backend/migrations"
	"github.com/Nowap83/FrameRate/backend/routes"
	"github.com/Nowap83/FrameRate/backend/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/Nowap83/FrameRate/backend/validators"
	"github.com/gin-gonic/gin/binding"               
  "github.com/go-playground/validator/v10"
)

func main() {

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Error loading .env file")
	}

	config.ConnectDB()
    defer func() {
        sqlDB, err := config.DB.DB()
        if err != nil {
            log.Printf("Failed to get DB instance: %v", err)
            return
        }
        if err := sqlDB.Close(); err != nil {
            log.Printf("Error closing database: %v", err)
        } else {
            log.Println("Database connection closed gracefully")
        }
    }()

	
	migrations.AutoMigrateAll()
	
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
  		validators.RegisterCustomValidators(v)
  }

	emailService := utils.NewEmailService()
	
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	

  routes.SetupRoutes(r, config.DB, emailService)
	

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s\n", port)
  if err := r.Run(":" + port); err != nil {
  	log.Fatal("Failed to start server:", err)
	}
}
