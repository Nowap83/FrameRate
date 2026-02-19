package model

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	TmdbID              int            `gorm:"uniqueIndex;not null" json:"tmdb_id"`
	Title               string         `gorm:"type:varchar(255);not null" json:"title"`
	OriginalTitle       string         `gorm:"type:varchar(255)" json:"original_title"`
	ReleaseYear         int            `gorm:"index" json:"release_year"`
	DurationMinutes     int            `json:"duration_minutes"`
	Synopsis            string         `gorm:"type:text" json:"synopsis"`
	PosterURL           string         `gorm:"type:varchar(500)" json:"poster_url"`
	BackdropURL         string         `gorm:"type:varchar(500)" json:"backdrop_url"`
	TrailerURL          string         `gorm:"type:varchar(500)" json:"trailer_url"`
	Budget              int64          `json:"budget"`
	Revenue             int64          `json:"revenue"`
	ImdbRating          float32        `gorm:"type:decimal(3,1)" json:"imdb_rating"`
	MetacriticScore     int            `json:"metacritic_score"`
	RottenTomatoesScore int            `json:"rotten_tomatoes_score"`
	Language            string         `gorm:"type:varchar(10)" json:"language"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	// relations
	Genres    []Genre     `gorm:"many2many:movie_genres;"`
	Countries []Country   `gorm:"many2many:movie_countries;"`
	Cast      []MovieCast `gorm:"foreignKey:MovieID"` // actors
	Crew      []MovieCrew `gorm:"foreignKey:MovieID"` // directors / writers / producers
}
