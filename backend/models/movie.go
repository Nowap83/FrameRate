package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID                  uint   `gorm:"primaryKey"`
	TmdbID              int    `gorm:"uniqueIndex;not null"`
	Title               string `gorm:"type:varchar(255);not null"`
	OriginalTitle       string `gorm:"type:varchar(255)"`
	ReleaseYear         int    `gorm:"index"`
	DurationMinutes     int
	Synopsis            string `gorm:"type:text"`
	PosterURL           string `gorm:"type:varchar(500)"`
	BackdropURL         string `gorm:"type:varchar(500)"`
	TrailerURL          string `gorm:"type:varchar(500)"`
	Budget              int64
	Revenue             int64
	ImdbRating          float32 `gorm:"type:decimal(3,1)"`
	MetacriticScore     int
	RottenTomatoesScore int
	Language            string `gorm:"type:varchar(10)"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`

	// relations
	Genres    []Genre     `gorm:"many2many:movie_genres;"`
	Countries []Country   `gorm:"many2many:movie_countries;"`
	Cast      []MovieCast `gorm:"foreignKey:MovieID"` // actors
	Crew      []MovieCrew `gorm:"foreignKey:MovieID"` // directors / writers / producers
}
