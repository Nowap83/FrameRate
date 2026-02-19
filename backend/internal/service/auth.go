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
