package services

import (
	"errors"
	"time"

	"github.com/Nowap83/FrameRate/backend/dto"
	"github.com/Nowap83/FrameRate/backend/models"
	"github.com/Nowap83/FrameRate/backend/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db           *gorm.DB
	emailService *utils.EmailService
}

func NewAuthService(db *gorm.DB, emailService *utils.EmailService) *AuthService {
	return &AuthService{
		db:           db,
		emailService: emailService,
	}
}

//
// REGISTER
//

func (s *AuthService) Register(input dto.RegisterRequest) (*dto.RegisterResponse, error) {

	var existingUser models.User
	// verification si email ou username existe deja
	err := s.db.Where("email = ? OR username = ?", input.Email, input.Username).First(&existingUser).Error

	if err == nil {
		// user trouvé = déjà pris
		if existingUser.Email == input.Email {
			return nil, errors.New("email already exists")
		}
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
	user := models.User{
		Username:          input.Username,
		Email:             input.Email,
		PasswordHash:      hashedPassword,
		VerificationToken: &verificationToken,
		TokenExpiresAt:    &expiresAt,
		IsVerified:        false,
		IsAdmin:           false,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	// go routine envoi email
	go func() {
		if err := s.emailService.SendVerificationEmail(
			user.Email,
			user.Username,
			verificationToken,
		); err != nil {
			// TODO: Logger avec zap?
			println("Failed to send verification email:", err.Error())
		} else {
			println("Verification email sent to", user.Email)
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
	var user models.User
	result := s.db.Where("email = ? OR username = ?", input.Login, input.Login).First(&user)
	if result.Error != nil {
		return nil, errors.New("invalid credentials")
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

	return dto.NewLoginResponse(token, &user), nil
}

//
// VERIFY EMAIL
//

func (s *AuthService) VerifyEmail(token string) (*dto.VerifyEmailResponse, error) {
	var user models.User

	// find user
	if err := s.db.Where("verification_token = ? AND token_expires_at > ?", token, time.Now()).
		First(&user).Error; err != nil {
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
			&user,
			"Email already verified",
		), nil
	}

	// marque comme verif
	user.IsVerified = true
	user.VerificationToken = nil
	user.TokenExpiresAt = nil

	if err := s.db.Save(&user).Error; err != nil {
		return nil, errors.New("failed to verify email")
	}

	// gen jwt
	jwtToken, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return dto.NewVerifyEmailResponse(
		jwtToken,
		&user,
		"Email verified successfully! You are now logged in.",
	), nil
}

//
// GET USER BY ID
//

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

//
// UPDATE PROFILE
//

func (s *AuthService) UpdateProfile(userID uint, input dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	// recup user actuel
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// check nouveau username
	if input.Username != nil && *input.Username != user.Username {
		var existing models.User
		if err := s.db.Where("username = ? AND id != ?", *input.Username, userID).First(&existing).Error; err == nil {
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

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, errors.New("failed to update profile")
	}

	return dto.NewProfileResponse(&user), nil
}

//
// CHANGE PASSWORD
//

func (s *AuthService) ChangePassword(userID uint, input dto.ChangePasswordRequest) error {
	// recup user actuel
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// check password actuel
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// hash nouveau password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// update password
	if err := s.db.Model(&user).Update("password_hash", string(hashedPassword)).Error; err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

//
// DELETE ACCOUNT
//

func (s *AuthService) DeleteAccount(userID uint) error {
	if err := s.db.Delete(&models.User{}, userID).Error; err != nil {
		return errors.New("failed to delete account")
	}
	return nil
}
