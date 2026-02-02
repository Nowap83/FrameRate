package migrations

import (
    "github.com/Nowap83/FrameRate/backend/config"
    "github.com/Nowap83/FrameRate/backend/models"
    "log"
)

// migre auto tout les models
func AutoMigrateAll() {
    log.Println("Running database migrations...")

    err := config.DB.AutoMigrate(
        &models.User{},
    )

    if err != nil {
        log.Fatal("Migration failed:", err)
    }

    log.Println("Database migrated successfully")
}
