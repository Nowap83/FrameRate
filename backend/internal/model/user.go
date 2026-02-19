package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Username          string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email             string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash      string         `gorm:"not null" json:"-"`
	ProfilePictureURL *string        `gorm:"size:500" json:"profile_picture,omitempty"`
	Bio               *string        `gorm:"size:500" json:"bio,omitempty"`
	GivenName         *string        `gorm:"size:100" json:"given_name,omitempty"`
	FamilyName        *string        `gorm:"size:100" json:"family_name,omitempty"`
	Location          *string        `gorm:"size:100" json:"location,omitempty"`
	Website           *string        `gorm:"size:255" json:"website,omitempty"`
	IsVerified        bool           `gorm:"default:false" json:"is_verified"`
	VerificationToken *string        `json:"-"`
	TokenExpiresAt    *time.Time     `json:"-"`
	IsAdmin           bool           `gorm:"default:false" json:"is_admin"`
	FavoriteFilms     []Movie        `gorm:"many2many:user_favorite_films;" json:"favorite_films,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// hook GORM juste avant insert
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Validation supplémentaire si besoin
	if u.Username == "" || u.Email == "" {
		return gorm.ErrInvalidValue
	}
	return nil
}
