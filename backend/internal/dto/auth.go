package dto

import (
	"time"

	"github.com/Nowap83/FrameRate/backend/internal/model"
)

// REQUESTS

type RegisterRequest struct {
	Username string `json:"username" binding:"required,username"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,strongpassword"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type UpdateProfileRequest struct {
	Username       *string       `json:"username,omitempty" binding:"omitempty,username"`
	Bio            *string       `json:"bio,omitempty" binding:"omitempty,max=500"`
	ProfilePicture *string       `json:"profile_picture_url,omitempty" binding:"omitempty,url"`
	GivenName      *string       `json:"given_name,omitempty" binding:"omitempty,max=100"`
	FamilyName     *string       `json:"family_name,omitempty" binding:"omitempty,max=100"`
	Location       *string       `json:"location,omitempty" binding:"omitempty,max=100"`
	Website        *string       `json:"website,omitempty" binding:"omitempty,max=255"`
	FavoriteFilms  []model.Movie `json:"favorite_films,omitempty"` // List of Movies
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,strongpassword"`
}

// RESPONSES

type UserResponse struct {
	ID             uint      `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	ProfilePicture *string   `json:"profile_picture_url,omitempty"`
	Bio            *string   `json:"bio,omitempty"`
	GivenName      *string   `json:"given_name,omitempty"`
	FamilyName     *string   `json:"family_name,omitempty"`
	Location       *string   `json:"location,omitempty"`
	Website        *string   `json:"website,omitempty"`
	IsVerified     bool      `json:"is_verified"`
	IsAdmin        bool      `json:"is_admin"`
	CreatedAt      time.Time `json:"created_at"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type PaginatedUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type VerifyEmailResponse struct {
	Token   string       `json:"token"`
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}

type UserStats struct {
	TotalFilms         int64          `json:"total_films"`
	MoviesThisYear     int64          `json:"movies_this_year"`
	Reviews            int64          `json:"reviews"`
	Following          int64          `json:"following"`
	Followers          int64          `json:"followers"`
	RatingDistribution map[string]int `json:"rating_distribution"`
}

type ProfileResponse struct {
	User           UserResponse  `json:"user"`
	Stats          *UserStats    `json:"stats,omitempty"`
	Favorites      []model.Movie `json:"favorites,omitempty"`
	RecentActivity []model.Movie `json:"recent_activity,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// CONVERTERS

func ToUserResponse(user *model.User) UserResponse {
	return UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePictureURL,
		Bio:            user.Bio,
		GivenName:      user.GivenName,
		FamilyName:     user.FamilyName,
		Location:       user.Location,
		Website:        user.Website,
		IsVerified:     user.IsVerified,
		IsAdmin:        user.IsAdmin,
		CreatedAt:      user.CreatedAt,
	}
}

// helpers de création de responses

func NewLoginResponse(token string, user *model.User) *LoginResponse {
	return &LoginResponse{
		Token: token,
		User:  ToUserResponse(user),
	}
}

func NewVerifyEmailResponse(token string, user *model.User, message string) *VerifyEmailResponse {
	return &VerifyEmailResponse{
		Token:   token,
		User:    ToUserResponse(user),
		Message: message,
	}
}
