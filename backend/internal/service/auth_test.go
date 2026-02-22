package service

import (
	"os"
	"testing"
	"time"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// MockUserRepository
type MockUserRepository struct {
	User *model.User
	Err  error

	CreateFn                 func(user *model.User) error
	GetByIDFn                func(id uint) (*model.User, error)
	GetAllUsersFn            func(page, limit int) ([]*model.User, int64, error)
	GetByEmailOrUsernameFn   func(login string) (*model.User, error)
	GetByEmailFn             func(email string) (*model.User, error)
	GetByUsernameFn          func(username string) (*model.User, error)
	GetByVerificationTokenFn func(token string) (*model.User, error)
	UpdateFn                 func(user *model.User) error
	UpdateFieldsFn           func(id uint, updates map[string]interface{}) error
	DeleteFn                 func(id uint) error
}

func (m *MockUserRepository) Create(user *model.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(user)
	}
	return m.Err
}
func (m *MockUserRepository) GetByID(id uint) (*model.User, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return m.User, m.Err
}
func (m *MockUserRepository) GetAllUsers(page, limit int) ([]*model.User, int64, error) {
	if m.GetAllUsersFn != nil {
		return m.GetAllUsersFn(page, limit)
	}
	return nil, 0, m.Err
}
func (m *MockUserRepository) GetByEmailOrUsername(login string) (*model.User, error) {
	if m.GetByEmailOrUsernameFn != nil {
		return m.GetByEmailOrUsernameFn(login)
	}
	return m.User, m.Err
}
func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(email)
	}
	return m.User, m.Err
}
func (m *MockUserRepository) GetByUsername(username string) (*model.User, error) {
	if m.GetByUsernameFn != nil {
		return m.GetByUsernameFn(username)
	}
	return m.User, m.Err
}
func (m *MockUserRepository) GetByVerificationToken(token string) (*model.User, error) {
	if m.GetByVerificationTokenFn != nil {
		return m.GetByVerificationTokenFn(token)
	}
	return m.User, m.Err
}
func (m *MockUserRepository) Update(user *model.User) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(user)
	}
	return m.Err
}
func (m *MockUserRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	if m.UpdateFieldsFn != nil {
		return m.UpdateFieldsFn(id, updates)
	}
	return m.Err
}
func (m *MockUserRepository) Delete(id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return m.Err
}

// MockEmailSender
type MockEmailSender struct {
	SendVerificationEmailFn func(to, username, token string) error
	Sent                    bool
}

func (m *MockEmailSender) SendVerificationEmail(to, username, token string) error {
	m.Sent = true
	if m.SendVerificationEmailFn != nil {
		return m.SendVerificationEmailFn(to, username, token)
	}
	return nil
}

func TestAuthService_Register_Success(t *testing.T) {
	utils.Log = zap.NewNop()
	userRepo := &MockUserRepository{
		GetByEmailFn: func(email string) (*model.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		GetByUsernameFn: func(username string) (*model.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		CreateFn: func(user *model.User) error {
			user.ID = 1
			return nil
		},
	}
	emailSender := &MockEmailSender{}
	authService := NewAuthService(userRepo, emailSender)

	req := dto.RegisterRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	resp, err := authService.Register(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.Message == "" {
		t.Fatalf("expected response with message")
	}

	// Wait briefly for goroutine to finish (hacky but sufficient for manual mock without channels)
	time.Sleep(10 * time.Millisecond)
	if !emailSender.Sent {
		t.Errorf("expected verification email to be sent")
	}
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	utils.Log = zap.NewNop()
	userRepo := &MockUserRepository{
		GetByEmailFn: func(email string) (*model.User, error) {
			return &model.User{}, nil
		},
	}
	emailSender := &MockEmailSender{}
	authService := NewAuthService(userRepo, emailSender)

	req := dto.RegisterRequest{
		Username: "newuser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	_, err := authService.Register(req)
	if err == nil || err.Error() != "email already exists" {
		t.Fatalf("expected error 'email already exists', got %v", err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	utils.Log = zap.NewNop()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	token := "tok"
	expires := time.Now().Add(1 * time.Hour)
	userRepo := &MockUserRepository{
		GetByEmailOrUsernameFn: func(login string) (*model.User, error) {
			return &model.User{
				ID:                1,
				Username:          "testuser",
				Email:             "testuser@example.com",
				PasswordHash:      string(hashedPassword),
				IsVerified:        true,
				VerificationToken: &token,
				TokenExpiresAt:    &expires,
			}, nil
		},
	}
	authService := NewAuthService(userRepo, &MockEmailSender{})

	os.Setenv("JWT_SECRET", "testsecret")

	req := dto.LoginRequest{
		Login:    "testuser",
		Password: "password123",
	}

	resp, err := authService.Login(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Errorf("expected token in response")
	}
}

func TestAuthService_Login_NotVerified(t *testing.T) {
	utils.Log = zap.NewNop()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	userRepo := &MockUserRepository{
		GetByEmailOrUsernameFn: func(login string) (*model.User, error) {
			return &model.User{
				ID:           1,
				Username:     "testuser",
				Email:        "testuser@example.com",
				PasswordHash: string(hashedPassword),
				IsVerified:   false,
			}, nil
		},
	}
	authService := NewAuthService(userRepo, &MockEmailSender{})

	req := dto.LoginRequest{
		Login:    "testuser",
		Password: "password123",
	}

	_, err := authService.Login(req)
	if err == nil || err.Error() != "email not verified. please check your inbox" {
		t.Fatalf("expected error 'email not verified...', got %v", err)
	}
}

func TestAuthService_VerifyEmail_Success(t *testing.T) {
	utils.Log = zap.NewNop()
	expires := time.Now().Add(1 * time.Hour)
	userRepo := &MockUserRepository{
		GetByVerificationTokenFn: func(token string) (*model.User, error) {
			return &model.User{
				ID:             1,
				IsVerified:     false,
				TokenExpiresAt: &expires,
			}, nil
		},
		UpdateFn: func(user *model.User) error {
			return nil
		},
	}
	authService := NewAuthService(userRepo, &MockEmailSender{})

	os.Setenv("JWT_SECRET", "testsecret")

	resp, err := authService.VerifyEmail("validtoken")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Message != "Email verified successfully! You are now logged in." {
		t.Errorf("expected success message, got %s", resp.Message)
	}
}

func TestAuthService_VerifyEmail_Expired(t *testing.T) {
	utils.Log = zap.NewNop()
	expires := time.Now().Add(-1 * time.Hour)
	userRepo := &MockUserRepository{
		GetByVerificationTokenFn: func(token string) (*model.User, error) {
			return &model.User{
				ID:             1,
				IsVerified:     false,
				TokenExpiresAt: &expires,
			}, nil
		},
	}
	authService := NewAuthService(userRepo, &MockEmailSender{})

	_, err := authService.VerifyEmail("expiredtoken")
	if err == nil || err.Error() != "invalid or expired verification token" {
		t.Fatalf("expected error for expired token, got %v", err)
	}
}
