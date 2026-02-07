package models

import (
	"time"

	"gorm.io/gorm"
)

type Gender int

const (
	GenderNotSet Gender = iota
	GenderFemale
	GenderMale
	GenderNonBinary
)

func (g Gender) String() string {
	switch g {
	case GenderFemale:
		return "Female"
	case GenderMale:
		return "Male"
	case GenderNonBinary:
		return "Non-binary"
	default:
		return "Not set"
	}
}

func (g Gender) IsValid() bool {
	return g >= GenderNotSet && g <= GenderNonBinary
}

// Person : Unifie Actor et Director
type Person struct {
	ID                uint   `gorm:"primaryKey"`
	TmdbID            int    `gorm:"uniqueIndex"`
	Name              string `gorm:"type:varchar(255);not null;index"`
	ProfilePictureURL string `gorm:"type:varchar(500)"`
	Biography         string `gorm:"type:text"`
	BirthDate         *time.Time
	BirthPlace        string `gorm:"type:varchar(255)"`
	DeathDate         *time.Time
	Gender            Gender `gorm:"type:smallint; default:0"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

// MovieCast : Table de jointure pour les acteurs (avec rôle)
type MovieCast struct {
	MovieID       uint   `gorm:"primaryKey"`
	PersonID      uint   `gorm:"primaryKey"`
	CharacterName string `gorm:"type:varchar(255)"`
	CastOrder     int    `gorm:"index"` // Pour trier par importance
	CreatedAt     time.Time

	Movie  Movie  `gorm:"foreignKey:MovieID"`
	Person Person `gorm:"foreignKey:PersonID"`
}

// MovieCrew : Table de jointure pour l'équipe technique (directors, writers, etc.)
type MovieCrew struct {
	MovieID    uint   `gorm:"primaryKey"`
	PersonID   uint   `gorm:"primaryKey"`
	Job        string `gorm:"primaryKey;type:varchar(100)"` // "Director", "Writer", "Producer"...
	Department string `gorm:"type:varchar(100);index"`      // "Directing", "Writing", "Production"...
	CreatedAt  time.Time

	Movie  Movie  `gorm:"foreignKey:MovieID"`
	Person Person `gorm:"foreignKey:PersonID"`
}
