package models

import "time"

type Track struct {
	UserID      uint       `gorm:"primaryKey"`
	MovieID     uint       `gorm:"primaryKey"`
	IsWatched   bool       `gorm:"default:false;index"`
	IsFavorite  bool       `gorm:"default:false;index"`
	IsWatchlist bool       `gorm:"default:false;index"`
	WatchedDate *time.Time `gorm:"index"`
	UpdatedAt   time.Time
	CreatedAt   time.Time

	User  User  `gorm:"foreignKey:UserID"`
	Movie Movie `gorm:"foreignKey:MovieID"`
}

type Rate struct {
	UserID    uint    `gorm:"primaryKey"`
	MovieID   uint    `gorm:"primaryKey"`
	Rating    float32 `gorm:"type:decimal(2,1);check:rating >= 0 AND rating <= 5"`
	CreatedAt time.Time
	UpdatedAt time.Time

	User  User  `gorm:"foreignKey:UserID"`
	Movie Movie `gorm:"foreignKey:MovieID"`
}
