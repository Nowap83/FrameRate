package models

import "time"

type Genre struct {
	ID        uint   `gorm:"primaryKey"`
	TmdbID    int    `gorm:"uniqueIndex"`
	Name      string `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedAt time.Time
}
