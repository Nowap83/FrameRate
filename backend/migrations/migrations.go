package migrations

import (
    "framerate/backend/config"
    "framerate/backend/models"
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
