package database

import (
	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// migre auto tout les models
func AutoMigrateAll(db *gorm.DB) {
	utils.Log.Info("Running database migrations...")

	err := db.AutoMigrate(
		&model.User{},

		// Movie models
		&model.Movie{},
		&model.Genre{},
		&model.Country{},
		&model.Person{},

		// Junction tables
		&model.MovieCast{},
		&model.MovieCrew{},
		&model.Track{},
		&model.Rate{},
	)

	if err != nil {
		utils.Log.Fatal("Migration failed", zap.Error(err))
	}

	utils.Log.Info("Database migrated successfully")
}
