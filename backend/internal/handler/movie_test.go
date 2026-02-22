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
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func setupMovieHandlerTest() (*gin.Engine, *gorm.DB) {
	utils.Log = zap.NewNop()
	gin.SetMode(gin.TestMode)

	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{}, &model.Movie{}, &model.Track{}, &model.Rate{}, &model.Review{})

	movieRepo := repository.NewMovieRepository(db)
	tmdbService := service.NewTMDBService(nil) // It's fine if we pre-populate movies

	movieService := service.NewMovieService(movieRepo, tmdbService)
	movieHandler := NewMovieHandler(movieService)

	r := gin.New()

	// Mock Auth Middleware
	mockAuth := func(c *gin.Context) {
		c.Set("userID", uint(1)) // User ID = 1
		c.Next()
	}

	api := r.Group("/movie")
	api.Use(mockAuth)
	{
		api.POST("/:tmdb_id/track", movieHandler.TrackMovie)
		api.POST("/:tmdb_id/rate", movieHandler.RateMovie)
		api.POST("/:tmdb_id/review", movieHandler.ReviewMovie)
	}

	return r, db
}

func TestMovieHandler_TrackMovie(t *testing.T) {
	r, db := setupMovieHandlerTest()

	// Pre-populate data
	user := &model.User{ID: 1, Username: "trackuser", Email: "track@example.com"}
	db.Create(user)
	db.Create(&model.Movie{TmdbID: 100, Title: "Test Movie 100"})

	tWatched := true
	reqBody := dto.TrackMovieRequest{IsWatched: &tWatched}
	body, _ := json.Marshal(reqBody)

	// Valid request
	req, _ := http.NewRequest("POST", "/movie/100/track", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var track model.Track
	db.First(&track, "user_id = ? AND movie_id = ?", 1, 1)
	if !track.IsWatched {
		t.Errorf("expected movie to be tracked as watched")
	}

	// Invalid tmdb_id
	req2, _ := http.NewRequest("POST", "/movie/abc/track", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for invalid ID, got %d", w2.Code)
	}
}

func TestMovieHandler_RateMovie(t *testing.T) {
	r, db := setupMovieHandlerTest()

	user := &model.User{ID: 1, Username: "rateuser", Email: "rate@example.com"}
	db.Create(user)
	db.Create(&model.Movie{TmdbID: 200, Title: "Test Movie 200"})

	reqBody := dto.RateMovieRequest{Rating: 4.5}
	body, _ := json.Marshal(reqBody)

	// Valid Request
	req, _ := http.NewRequest("POST", "/movie/200/rate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var rate model.Rate
	db.First(&rate, "user_id = ? AND movie_id = ?", 1, 1)
	if rate.Rating != 4.5 {
		t.Errorf("expected rating 4.5, got %f", rate.Rating)
	}
}

func TestMovieHandler_ReviewMovie(t *testing.T) {
	r, db := setupMovieHandlerTest()

	user := &model.User{ID: 1, Username: "reviewuser", Email: "review@example.com"}
	db.Create(user)
	db.Create(&model.Movie{TmdbID: 300, Title: "Test Movie 300"})

	reqBody := dto.ReviewRequest{Content: "Great!", IsSpoiler: false}
	body, _ := json.Marshal(reqBody)

	// Valid Request
	req, _ := http.NewRequest("POST", "/movie/300/review", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var review model.Review
	db.First(&review, "user_id = ? AND movie_id = ?", 1, 1)
	if review.Content != "Great!" {
		t.Errorf("expected review 'Great!', got %s", review.Content)
	}
}
