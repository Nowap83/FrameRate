package dto

import (
	"time"

	"github.com/Nowap83/FrameRate/backend/models"
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
	Username          *string `json:"username,omitempty" binding:"omitempty,username"`
	Bio               *string `json:"bio,omitempty" binding:"omitempty,max=500"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty" binding:"omitempty,url"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,strongpassword"`
}

// RESPONSES

type UserResponse struct {
	ID                uint      `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	ProfilePictureURL *string   `json:"profile_picture_url,omitempty"`
	Bio               *string   `json:"bio,omitempty"`
	IsVerified        bool      `json:"is_verified"`
	IsAdmin           bool      `json:"is_admin"`
	CreatedAt         time.Time `json:"created_at"`
}

type RegisterResponse struct {
	Message string `json:"message"`
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

type ProfileResponse struct {
	User UserResponse `json:"user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// CONVERTERS

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		ProfilePictureURL: user.ProfilePictureURL,
		Bio:               user.Bio,
		IsVerified:        user.IsVerified,
		IsAdmin:           user.IsAdmin,
		CreatedAt:         user.CreatedAt,
	}
}

// helpers de création de responses

func NewLoginResponse(token string, user *models.User) *LoginResponse {
	return &LoginResponse{
		Token: token,
		User:  ToUserResponse(user),
	}
}

func NewVerifyEmailResponse(token string, user *models.User, message string) *VerifyEmailResponse {
	return &VerifyEmailResponse{
		Token:   token,
		User:    ToUserResponse(user),
		Message: message,
	}
}

func NewProfileResponse(user *models.User) *ProfileResponse {
	return &ProfileResponse{
		User: ToUserResponse(user),
	}
}
