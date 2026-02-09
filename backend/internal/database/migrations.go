package migrations

import (
	"log"

	"github.com/Nowap83/FrameRate/backend/config"
	"github.com/Nowap83/FrameRate/backend/models"
)

// migre auto tout les models
func AutoMigrateAll() {
	log.Println("Running database migrations...")

	err := config.DB.AutoMigrate(
		&models.User{},

		// Movie models
		&models.Movie{},
		&models.Genre{},
		&models.Country{},
		&models.Person{},

		// Junction tables
		&models.MovieCast{},
		&models.MovieCrew{},
		&models.Track{},
		&models.Rate{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migrated successfully")
}
