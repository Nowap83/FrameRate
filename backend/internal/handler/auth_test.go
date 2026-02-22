package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/Nowap83/FrameRate/backend/internal/repository"
	"github.com/Nowap83/FrameRate/backend/internal/service"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	internalValidator "github.com/Nowap83/FrameRate/backend/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/glebarez/sqlite"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// MockEmailSender for tests
type MockEmailSender struct{}

func (m *MockEmailSender) SendVerificationEmail(to, username, token string) error {
	return nil
}

func setupAuthHandlerTest() (*gin.Engine, *gorm.DB) {
	utils.Log = zap.NewNop()
	gin.SetMode(gin.TestMode)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		internalValidator.RegisterCustomValidators(v)
	}

	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{}, &model.Movie{}, &model.Track{}, &model.Rate{}, &model.Review{})

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, &MockEmailSender{})
	authHandler := NewAuthHandler(authService)

	r := gin.New()
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.GET("/verify-email", authHandler.VerifyEmail)

	return r, db
}

func TestAuthHandler_Register(t *testing.T) {
	r, db := setupAuthHandlerTest()

	// 1. Success
	reqBody := dto.RegisterRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "Password1!",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", w.Code)
	}
	var count int64
	db.Model(&model.User{}).Where("email = ?", "new@example.com").Count(&count)
	if count != 1 {
		t.Errorf("expected user to be in DB")
	}

	// 2. Conflict (Email Exists)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict, got %d", w2.Code)
	}

	// 3. Bad Request (Validation)
	badReq := map[string]string{"username": ""}
	badBody, _ := json.Marshal(badReq)
	req3, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(badBody))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for validation error, got %d", w3.Code)
	}
}

func TestAuthHandler_Login(t *testing.T) {
	r, db := setupAuthHandlerTest()
	os.Setenv("JWT_SECRET", "testsecret")

	// Create user
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user := &model.User{
		Username:     "loginuser",
		Email:        "login@example.com",
		IsVerified:   false,
		PasswordHash: string(hash),
	}
	db.Create(user)

	// 1. Unverified Login
	reqBody := dto.LoginRequest{Login: "loginuser", Password: "password"}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for unverified, got %d", w.Code)
	}

	// 2. Success Login
	db.Model(&user).Update("is_verified", true)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w2.Code)
	}

	// 3. Invalid credentials
	badReq := dto.LoginRequest{Login: "loginuser", Password: "wrongpassword"}
	badBody, _ := json.Marshal(badReq)
	req3, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(badBody))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", w3.Code)
	}
}

func TestAuthHandler_VerifyEmail(t *testing.T) {
	r, db := setupAuthHandlerTest()
	os.Setenv("JWT_SECRET", "testsecret")

	token := "verify123"
	user := &model.User{
		Username:          "verifyuser",
		Email:             "verify@example.com",
		VerificationToken: &token,
		IsVerified:        false,
	}
	db.Create(user)

	// 1. Success
	req, _ := http.NewRequest("GET", "/verify-email?token=verify123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var updatedUser model.User
	db.First(&updatedUser, user.ID)
	if !updatedUser.IsVerified {
		t.Errorf("expected user to be verified")
	}

	// 2. Missing token
	req2, _ := http.NewRequest("GET", "/verify-email", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for missing token, got %d", w2.Code)
	}

	// 3. Invalid token
	req3, _ := http.NewRequest("GET", "/verify-email?token=invalid", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for invalid token, got %d", w3.Code)
	}
}
