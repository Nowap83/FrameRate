package repository

import (
	"fmt"
	"testing"

	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
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

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	err := repo.Create(user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Errorf("expected user ID to be set")
	}

	// Verify it was actually saved
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 user in db, got %d", count)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	db.Create(user)

	foundUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if foundUser.Username != user.Username {
		t.Errorf("expected username %s, got %s", user.Username, foundUser.Username)
	}

	// Test non-existent ID
	_, err = repo.GetByID(999)
	if err == nil {
		t.Errorf("expected error for non-existent user")
	}
}

func TestUserRepository_GetByEmailOrUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	db.Create(user)

	tests := []struct {
		name    string
		login   string
		wantErr bool
	}{
		{"By Email", "test@example.com", false},
		{"By Username", "testuser", false},
		{"Non-existent", "doesnotexist@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, err := repo.GetByEmailOrUsername(tt.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByEmailOrUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && foundUser.ID != user.ID {
				t.Errorf("expected user ID %d, got %d", user.ID, foundUser.ID)
			}
		})
	}
}

func TestUserRepository_GetMethods(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	token := "token123"
	user := &model.User{
		Username:          "testuser",
		Email:             "test@example.com",
		VerificationToken: &token,
	}
	db.Create(user)

	// Test GetByEmail
	t.Run("GetByEmail", func(t *testing.T) {
		found, err := repo.GetByEmail("test@example.com")
		if err != nil || found.ID != user.ID {
			t.Errorf("GetByEmail failed")
		}
	})

	// Test GetByUsername
	t.Run("GetByUsername", func(t *testing.T) {
		found, err := repo.GetByUsername("testuser")
		if err != nil || found.ID != user.ID {
			t.Errorf("GetByUsername failed")
		}
	})

	// Test GetByVerificationToken
	t.Run("GetByVerificationToken", func(t *testing.T) {
		found, err := repo.GetByVerificationToken("token123")
		if err != nil || found.ID != user.ID {
			t.Errorf("GetByVerificationToken failed")
		}
	})
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{Username: "oldname", Email: "test@example.com"}
	db.Create(user)

	user.Username = "newname"
	err := repo.Update(user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	var updatedUser model.User
	db.First(&updatedUser, user.ID)
	if updatedUser.Username != "newname" {
		t.Errorf("expected updated username 'newname', got %s", updatedUser.Username)
	}
}

func TestUserRepository_UpdateFields(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{Username: "oldname", Email: "test@example.com", IsVerified: false}
	db.Create(user)

	updates := map[string]interface{}{
		"is_verified": true,
		"username":    "updatedname",
	}

	err := repo.UpdateFields(user.ID, updates)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	var updatedUser model.User
	db.First(&updatedUser, user.ID)
	if !updatedUser.IsVerified || updatedUser.Username != "updatedname" {
		t.Errorf("fields not properly updated")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{Username: "todelete", Email: "delete@example.com"}
	db.Create(user)

	err := repo.Delete(user.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	var count int64
	db.Model(&model.User{}).Where("id = ?", user.ID).Count(&count)
	if count != 0 {
		t.Errorf("expected user to be deleted (even soft-deleted from normal queries)")
	}
}

func TestUserRepository_GetAllUsers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	for i := 1; i <= 5; i++ {
		db.Create(&model.User{
			Username:     fmt.Sprintf("user%d", i),
			Email:        fmt.Sprintf("test%d@example.com", i),
			PasswordHash: "hash",
		})
	}

	users, total, err := repo.GetAllUsers(1, 3)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users per page, got %d", len(users))
	}

	users2, total2, err := repo.GetAllUsers(2, 3)
	if err != nil || total2 != 5 || len(users2) != 2 {
		t.Errorf("expected second page to have 2 users, got %d", len(users2))
	}
}
