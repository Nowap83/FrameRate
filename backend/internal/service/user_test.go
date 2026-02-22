package service

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/Nowap83/FrameRate/backend/internal/repository"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupUserServiceTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Movie{}, &model.Track{}, &model.Rate{}, &model.Review{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	return db
}

func TestUserService_GetUserByID(t *testing.T) {
	utils.Log = zap.NewNop()
	db := setupUserServiceTestDB(t)
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	userService := NewUserService(userRepo, movieRepo)

	user := &model.User{Username: "test1", Email: "test1@example.com"}
	db.Create(user)

	found, err := userService.GetUserByID(user.ID)
	if err != nil || found.ID != user.ID {
		t.Fatalf("expected to find user, got error %v", err)
	}

	_, err = userService.GetUserByID(999)
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	utils.Log = zap.NewNop()
	db := setupUserServiceTestDB(t)
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	userService := NewUserService(userRepo, movieRepo)

	for i := 1; i <= 5; i++ {
		db.Create(&model.User{Username: "user" + string(rune(i)), Email: "test" + string(rune(i)) + "@test.com"})
	}

	resp, err := userService.GetAllUsers(1, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Users) != 3 || resp.Total != 5 || resp.TotalPages != 2 {
		t.Errorf("pagination incorrect: got %v", resp)
	}
}

func TestUserService_GetProfile(t *testing.T) {
	utils.Log = zap.NewNop()
	db := setupUserServiceTestDB(t)
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	userService := NewUserService(userRepo, movieRepo)

	user := &model.User{Username: "profileuser", Email: "profile@example.com"}
	db.Create(user)

	movie := &model.Movie{TmdbID: 10, Title: "Stat Movie"}
	movieRepo.UpsertMovie(movie)

	now := time.Now()
	movieRepo.UpsertTrack(&model.Track{UserID: user.ID, MovieID: movie.ID, IsWatched: true, WatchedDate: &now})

	profile, err := userService.GetProfile(user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if profile.Stats.TotalFilms != 1 || profile.User.Username != "profileuser" {
		t.Errorf("profile stats incorrect")
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	utils.Log = zap.NewNop()
	db := setupUserServiceTestDB(t)
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	userService := NewUserService(userRepo, movieRepo)

	user := &model.User{Username: "olduser", Email: "update@example.com"}
	db.Create(user)

	newUsername := "newuser"
	req := dto.UpdateProfileRequest{Username: &newUsername}

	profile, err := userService.UpdateProfile(user.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if profile.User.Username != "newuser" {
		t.Errorf("expected updated username 'newuser', got %s", profile.User.Username)
	}
}

func TestUserService_ChangePassword(t *testing.T) {
	utils.Log = zap.NewNop()
	db := setupUserServiceTestDB(t)
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	userService := NewUserService(userRepo, movieRepo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.DefaultCost)
	user := &model.User{Username: "passuser", Email: "pass@example.com", PasswordHash: string(hash)}
	db.Create(user)

	req := dto.ChangePasswordRequest{CurrentPassword: "oldpass", NewPassword: "newpass"}
	err := userService.ChangePassword(user.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var updated model.User
	db.First(&updated, user.ID)
	if err := bcrypt.CompareHashAndPassword([]byte(updated.PasswordHash), []byte("newpass")); err != nil {
		t.Errorf("expected new password to be valid")
	}
}

func TestUserService_CheckUsernameAvailability(t *testing.T) {
	utils.Log = zap.NewNop()
	db := setupUserServiceTestDB(t)
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	userService := NewUserService(userRepo, movieRepo)

	user := &model.User{Username: "taken", Email: "taken@example.com"}
	db.Create(user)

	avail1, _ := userService.CheckUsernameAvailability("taken")
	if avail1 {
		t.Errorf("expected 'taken' to be false")
	}

	avail2, _ := userService.CheckUsernameAvailability("free")
	if !avail2 {
		t.Errorf("expected 'free' to be true")
	}
}

func TestUserService_DeleteAccount(t *testing.T) {
	utils.Log = zap.NewNop()
	db := setupUserServiceTestDB(t)
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	userService := NewUserService(userRepo, movieRepo)

	user := &model.User{Username: "todelete", Email: "delete@example.com"}
	db.Create(user)

	err := userService.DeleteAccount(user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = userService.GetUserByID(user.ID)
	if err != ErrUserNotFound {
		t.Errorf("expected user to be deleted/not found")
	}
}

func TestUserService_UpdateAvatar(t *testing.T) {
	utils.Log = zap.NewNop()
	db := setupUserServiceTestDB(t)
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	userService := NewUserService(userRepo, movieRepo)

	user := &model.User{Username: "avataruser", Email: "avatar@example.com"}
	db.Create(user)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
	part.Write([]byte("fake image data"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/dummy", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_ = req.ParseMultipartForm(10 << 20)
	_, fileHeader, _ := req.FormFile("avatar")

	fileURL, err := userService.UpdateAvatar(user.ID, fileHeader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := userService.GetUserByID(user.ID)
	if updated.ProfilePictureURL == nil || *updated.ProfilePictureURL != fileURL {
		t.Errorf("expected avatar to be updated")
	}

	os.RemoveAll("./uploads/avatars")
}
