package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockUserRepository is a manual stub for repository.UserRepository
type mockUserRepo struct {
	user *model.User
	err  error
}

func (m *mockUserRepo) Create(user *model.User) error                              { return nil }
func (m *mockUserRepo) GetByID(id uint) (*model.User, error)                       { return m.user, m.err }
func (m *mockUserRepo) GetAllUsers() ([]*model.User, error)                        { return nil, nil }
func (m *mockUserRepo) GetByEmailOrUsername(login string) (*model.User, error)     { return nil, nil }
func (m *mockUserRepo) GetByEmail(email string) (*model.User, error)               { return nil, nil }
func (m *mockUserRepo) GetByUsername(username string) (*model.User, error)         { return nil, nil }
func (m *mockUserRepo) GetByVerificationToken(token string) (*model.User, error)   { return nil, nil }
func (m *mockUserRepo) Update(user *model.User) error                              { return nil }
func (m *mockUserRepo) UpdateFields(id uint, updates map[string]interface{}) error { return nil }
func (m *mockUserRepo) Delete(id uint) error                                       { return nil }

func TestAdminRequired_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	repo := &mockUserRepo{}
	AdminRequired(repo)(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
	assert.True(t, c.IsAborted())
}

func TestAdminRequired_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Set("userID", uint(1))

	repo := &mockUserRepo{err: errors.New("not found")}
	AdminRequired(repo)(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
	assert.True(t, c.IsAborted())
}

func TestAdminRequired_NotAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Set("userID", uint(1))

	repo := &mockUserRepo{user: &model.User{IsAdmin: false}}
	AdminRequired(repo)(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Admin privileges required")
	assert.True(t, c.IsAborted())
}

func TestAdminRequired_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Set("userID", uint(1))

	repo := &mockUserRepo{user: &model.User{IsAdmin: true}}
	AdminRequired(repo)(c)

	assert.False(t, c.IsAborted())
}
