package model

import "time"

type Genre struct {
	ID        int    `gorm:"primaryKey;autoIncrement:false" json:"id"` // TMDB ID
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt time.Time
}
