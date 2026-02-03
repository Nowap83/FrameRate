package models

import (
    "time"

    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type User struct {
    ID             uint      `gorm:"primaryKey" json:"id"`
    Username       string    `gorm:"uniqueIndex;not null;size:50" json:"username" binding:"required,min=3,max=50"`
    Email          string    `gorm:"uniqueIndex;not null;size:255" json:"email" binding:"required,email"`
    PasswordHash   string    `gorm:"not null" json:"-"` 
    ProfilePicture string    `gorm:"size:500" json:"profile_picture"`
    Bio            string    `gorm:"size:500" json:"bio"`
    IsVerified     bool      `gorm:"default:false" json:"is_verified"`
    IsAdmin        bool      `gorm:"default:false" json:"is_admin"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

type UserResponse struct {
    ID             uint      `json:"id"`
    Username       string    `json:"username"`
    Email          string    `json:"email"`
    ProfilePicture string    `json:"profile_picture,omitempty"`
    Bio            string    `json:"bio,omitempty"`
    IsVerified     bool      `json:"is_verified"`
    IsAdmin        bool      `json:"is_admin"`
    CreatedAt      time.Time `json:"created_at"`
}

// convertisseur
func (u *User) ToResponse() UserResponse {
    return UserResponse{
        ID:             u.ID,
        Username:       u.Username,
        Email:          u.Email,
        ProfilePicture: u.ProfilePicture,
        Bio:            u.Bio,
        IsVerified:     u.IsVerified,
        IsAdmin:        u.IsAdmin,
        CreatedAt:      u.CreatedAt,
    }
}

func (u *User) HashPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.PasswordHash = string(hashedPassword)
    return nil
}

func (u *User) ValidatePassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
    return err == nil
}

// hook GORM juste avant insert
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Validation supplémentaire si besoin
    if u.Username == "" || u.Email == "" {
        return gorm.ErrInvalidValue
    }
    return nil
}
