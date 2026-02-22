package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func setupUserHandlerTest() (*gin.Engine, *gorm.DB) {
	utils.Log = zap.NewNop()
	gin.SetMode(gin.TestMode)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		internalValidator.RegisterCustomValidators(v)
	}

	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{}, &model.Movie{}, &model.Track{}, &model.Rate{}, &model.Review{})

	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	userService := service.NewUserService(userRepo, movieRepo)
	userHandler := NewUserHandler(userService)

	r := gin.New()

	// Middleware mockup for authenticated routes
	mockAuth := func(c *gin.Context) {
		// Mock userID = 1
		c.Set("userID", uint(1))
		c.Next()
	}

	api := r.Group("/user")
	api.Use(mockAuth)
	{
		api.GET("/profile", userHandler.GetProfile)
		api.PUT("/profile", userHandler.UpdateProfile)
		api.PUT("/password", userHandler.ChangePassword)
		api.DELETE("/account", userHandler.DeleteAccount)
	}

	r.GET("/check-username", userHandler.CheckUsername)

	admin := r.Group("/admin")
	{
		admin.GET("/users", userHandler.GetAllUsers)
		admin.DELETE("/users/:id", userHandler.DeleteUserAdmin)
	}

	return r, db
}

func TestUserHandler_GetProfile(t *testing.T) {
	r, db := setupUserHandlerTest()

	user := &model.User{ID: 1, Username: "testuser", Email: "test@example.com"}
	db.Create(user)

	req, _ := http.NewRequest("GET", "/user/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var resp dto.ProfileResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.User.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %s", resp.User.Username)
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	r, db := setupUserHandlerTest()

	user := &model.User{ID: 1, Username: "olduser", Email: "test@example.com"}
	db.Create(user)

	newUsername := "newuser"
	reqBody := dto.UpdateProfileRequest{Username: &newUsername}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/user/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var updatedUser model.User
	db.First(&updatedUser, 1)
	if updatedUser.Username != "newuser" {
		t.Errorf("expected updated username 'newuser', got %s", updatedUser.Username)
	}
}

func TestUserHandler_ChangePassword(t *testing.T) {
	r, db := setupUserHandlerTest()

	hash, _ := bcrypt.GenerateFromPassword([]byte("oldpassword!"), bcrypt.DefaultCost)
	user := &model.User{ID: 1, Username: "passuser", Email: "pass@example.com", PasswordHash: string(hash)}
	db.Create(user)

	reqBody := dto.ChangePasswordRequest{
		CurrentPassword: "oldpassword!",
		NewPassword:     "NewPassword1!", // Must match strong password validator
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/user/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestUserHandler_DeleteAccount(t *testing.T) {
	r, db := setupUserHandlerTest()

	user := &model.User{ID: 1, Username: "deleteuser", Email: "del@example.com"}
	db.Create(user)

	req, _ := http.NewRequest("DELETE", "/user/account", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var count int64
	db.Model(&model.User{}).Where("id = 1").Count(&count)
	if count != 0 {
		t.Errorf("expected user to be deleted")
	}
}

func TestUserHandler_CheckUsername(t *testing.T) {
	r, db := setupUserHandlerTest()

	user := &model.User{Username: "taken", Email: "taken@example.com"}
	db.Create(user)

	// Available
	req, _ := http.NewRequest("GET", "/check-username?username=free", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte(`"available":true`)) {
		t.Errorf("expected available: true")
	}

	// Taken
	req2, _ := http.NewRequest("GET", "/check-username?username=taken", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w2.Code)
	}
	if !bytes.Contains(w2.Body.Bytes(), []byte(`"available":false`)) {
		t.Errorf("expected available: false")
	}
}

func TestUserHandler_AdminFunctions(t *testing.T) {
	r, db := setupUserHandlerTest()

	user1 := &model.User{ID: 10, Username: "u1", Email: "1@t.com"}
	user2 := &model.User{ID: 11, Username: "u2", Email: "2@t.com"}
	db.Create(user1)
	db.Create(user2)

	// Get all
	req, _ := http.NewRequest("GET", "/admin/users?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	// Delete
	req2, _ := http.NewRequest("DELETE", "/admin/users/10", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for admin delete, got %d", w2.Code)
	}
}
