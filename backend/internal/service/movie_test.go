package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/Nowap83/FrameRate/backend/internal/repository"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func setupMovieServiceTestDB(t *testing.T) *gorm.DB {
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

func mockTMDBClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestMovieService_TrackMovie(t *testing.T) {
	utils.Log = zap.NewNop()
	db := setupMovieServiceTestDB(t)
	repo := repository.NewMovieRepository(db)

	mockResponse := dto.TMDBMovieDetails{
		ID:            123,
		Title:         "Inception",
		OriginalTitle: "Inception",
		Runtime:       148,
		Overview:      "A dream within a dream.",
		ReleaseDate:   "2010-07-16",
	}
	respBytes, _ := json.Marshal(mockResponse)

	tmdbService := NewTMDBService(nil)
	tmdbService.client = mockTMDBClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(respBytes)),
			Header:     make(http.Header),
		}
	})

	movieService := NewMovieService(repo, tmdbService)

	user := &model.User{Username: "test", Email: "test@example.com"}
	db.Create(user)

	tWatched := true
	req := dto.TrackMovieRequest{IsWatched: &tWatched}

	err := movieService.TrackMovie(user.ID, 123, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify track is inserted
	var track model.Track
	db.First(&track, "user_id = ? AND movie_id = ?", user.ID, 1) // First movie gets ID=1
	if !track.IsWatched {
		t.Errorf("expected movie to be tracked as watched")
	}

	// Verify movie is inserted
	var movie model.Movie
	db.First(&movie, 1)
	if movie.Title != "Inception" || movie.ReleaseYear != 2010 {
		t.Errorf("expected movie 'Inception' (2010), got %v", movie)
	}
}

func TestMovieService_RateMovie(t *testing.T) {
	utils.Log = zap.NewNop()
	db := setupMovieServiceTestDB(t)
	repo := repository.NewMovieRepository(db)

	tmdbService := NewTMDBService(nil)
	// Mock returning existing movie, so it won't hit TMDB.
	// But ensureMovieExists will still hit TMDB if not found, so we seed the DB.
	repo.UpsertMovie(&model.Movie{TmdbID: 456, Title: "Interstellar"})

	movieService := NewMovieService(repo, tmdbService)

	user := &model.User{Username: "test2", Email: "test2@example.com"}
	db.Create(user)

	// Valid rating
	req := dto.RateMovieRequest{Rating: 4.5}
	err := movieService.RateMovie(user.ID, 456, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Invalid rating increments
	req2 := dto.RateMovieRequest{Rating: 4.7}
	err = movieService.RateMovie(user.ID, 456, req2)
	if err == nil || err.Error() != "Rating must be in increments of 0.5" {
		t.Fatalf("expected error for invalid rating increment, got %v", err)
	}
}

func TestMovieService_ReviewMovie(t *testing.T) {
	utils.Log = zap.NewNop()
	db := setupMovieServiceTestDB(t)
	repo := repository.NewMovieRepository(db)

	tmdbService := NewTMDBService(nil)
	repo.UpsertMovie(&model.Movie{TmdbID: 789, Title: "The Dark Knight"})

	movieService := NewMovieService(repo, tmdbService)

	user := &model.User{Username: "test3", Email: "test3@example.com"}
	db.Create(user)

	req := dto.ReviewRequest{Content: "Masterpiece.", IsSpoiler: false}
	err := movieService.ReviewMovie(user.ID, 789, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var review model.Review
	db.First(&review)
	if review.Content != "Masterpiece." || review.IsSpoiler {
		t.Errorf("expected review correctly inserted")
	}
}
