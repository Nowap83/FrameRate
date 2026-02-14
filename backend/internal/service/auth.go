package service

import (
	"errors"
	"time"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/Nowap83/FrameRate/backend/internal/repository"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo     repository.UserRepository
	emailService *utils.EmailService
}

func NewAuthService(userRepo repository.UserRepository, emailService *utils.EmailService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

//
// REGISTER
//

func (s *AuthService) Register(input dto.RegisterRequest) (*dto.RegisterResponse, error) {

	// verification si email ou username existe deja
	_, err := s.userRepo.GetByEmail(input.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("database error")
	}

	_, err = s.userRepo.GetByUsername(input.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("database error")
	}

	// hashage password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// generation du token mail
	verificationToken, err := utils.GenerateVerificationToken()
	if err != nil {
		return nil, errors.New("failed to generate verification token")
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	// creation user
	user := model.User{
		Username:          input.Username,
		Email:             input.Email,
		PasswordHash:      hashedPassword,
		VerificationToken: &verificationToken,
		TokenExpiresAt:    &expiresAt,
		IsVerified:        false,
		IsAdmin:           false,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// go routine envoi email
	go func() {
		if err := s.emailService.SendVerificationEmail(
			user.Email,
			user.Username,
			verificationToken,
		); err != nil {
			utils.Log.Error("Failed to send verification email",
				zap.Uint("user_id", user.ID),
				zap.String("email", user.Email),
				zap.Error(err),
			)
		} else {
			utils.Log.Info("Verification email sent",
				zap.Uint("user_id", user.ID),
				zap.String("email", user.Email),
			)
		}
	}()

	return &dto.RegisterResponse{
		Message: "Registration successful! Please check your email to verify your account.",
	}, nil
}

//
// LOGIN
//

func (s *AuthService) Login(input dto.LoginRequest) (*dto.LoginResponse, error) {
	// cherche user par mail ou username
	user, err := s.userRepo.GetByEmailOrUsername(input.Login)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("database error")
	}

	// verif du mdp
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return nil, errors.New("email not verified. please check your inbox")
	}

	// genere jwt
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return dto.NewLoginResponse(token, user), nil
}

//
// VERIFY EMAIL
//

func (s *AuthService) VerifyEmail(token string) (*dto.VerifyEmailResponse, error) {
	// find user
	user, err := s.userRepo.GetByVerificationToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired verification token")
	}

	if user.TokenExpiresAt != nil && user.TokenExpiresAt.Before(time.Now()) {
		return nil, errors.New("invalid or expired verification token")
	}

	// deja verif ?
	if user.IsVerified {
		// genere jwt quand même
		jwtToken, err := utils.GenerateToken(user.ID)
		if err != nil {
			return nil, errors.New("failed to generate token")
		}

		return dto.NewVerifyEmailResponse(
			jwtToken,
			user,
			"Email already verified",
		), nil
	}

	// marque comme verif
	user.IsVerified = true
	user.VerificationToken = nil
	user.TokenExpiresAt = nil

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to verify email")
	}

	// gen jwt
	jwtToken, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return dto.NewVerifyEmailResponse(
		jwtToken,
		user,
		"Email verified successfully! You are now logged in.",
	), nil
}

//
// GET USER BY ID
//

func (s *AuthService) GetUserByID(userID uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

//
// UPDATE PROFILE
//

func (s *AuthService) UpdateProfile(userID uint, input dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	// recup user actuel
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// check nouveau username
	if input.Username != nil && *input.Username != user.Username {
		existing, err := s.userRepo.GetByUsername(*input.Username)
		if err == nil && existing.ID != userID {
			return nil, errors.New("username already taken")
		}
	}

	// update les champs
	updates := make(map[string]interface{})
	if input.Username != nil {
		updates["username"] = *input.Username
	}
	if input.Bio != nil {
		updates["bio"] = *input.Bio
	}
	if input.ProfilePictureURL != nil {
		updates["profile_picture_url"] = *input.ProfilePictureURL
	}

	if err := s.userRepo.UpdateFields(userID, updates); err != nil {
		return nil, errors.New("failed to update profile")
	}

	// refresh user data for response
	updatedUser, _ := s.userRepo.GetByID(userID)

	return dto.NewProfileResponse(updatedUser), nil
}

//
// CHANGE PASSWORD
//

func (s *AuthService) ChangePassword(userID uint, input dto.ChangePasswordRequest) error {
	// recup user actuel
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// check password actuel
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// hash nouveau password
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// update password
	if err := s.userRepo.UpdateFields(userID, map[string]interface{}{"password_hash": string(hashedPassword)}); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

//
// DELETE ACCOUNT
//

func (s *AuthService) DeleteAccount(userID uint) error {
	if err := s.userRepo.Delete(userID); err != nil {
		return errors.New("failed to delete account")
	}
	return nil
}
