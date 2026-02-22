package repository

import (
	"testing"
	"time"

	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupMovieTestDB(t *testing.T) *gorm.DB {
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

func TestMovieRepository_UpsertMovie(t *testing.T) {
	db := setupMovieTestDB(t)
	repo := NewMovieRepository(db)

	movie := &model.Movie{
		TmdbID:      123,
		Title:       "Test Movie",
		ReleaseYear: 2024,
	}

	err := repo.UpsertMovie(movie)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Upsert again with updated title
	movie.Title = "Updated Movie Title"
	err = repo.UpsertMovie(movie)
	if err != nil {
		t.Errorf("expected no error on second upsert, got %v", err)
	}

	var count int64
	db.Model(&model.Movie{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 movie in db, got %d", count)
	}

	savedMovie, _ := repo.GetMovieByTmdbID(123)
	if savedMovie.Title != "Updated Movie Title" {
		t.Errorf("expected title to be updated")
	}
}

func TestMovieRepository_GetMovieByTmdbID(t *testing.T) {
	db := setupMovieTestDB(t)
	repo := NewMovieRepository(db)

	movie := &model.Movie{TmdbID: 456, Title: "Another Movie"}
	repo.UpsertMovie(movie)

	found, err := repo.GetMovieByTmdbID(456)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if found.Title != "Another Movie" {
		t.Errorf("expected title 'Another Movie', got %s", found.Title)
	}

	_, err = repo.GetMovieByTmdbID(999)
	if err == nil {
		t.Errorf("expected error for non-existent movie")
	}
}

func TestMovieRepository_UpsertInteractions(t *testing.T) {
	db := setupMovieTestDB(t)
	repo := NewMovieRepository(db)

	user := &model.User{Username: "testuser", Email: "test@example.com"}
	db.Create(user)

	movie := &model.Movie{TmdbID: 101, Title: "Interaction Movie"}
	repo.UpsertMovie(movie)

	// Test Track
	trackDate := time.Now()
	track := &model.Track{
		UserID:      user.ID,
		MovieID:     movie.ID,
		IsWatched:   true,
		WatchedDate: &trackDate,
	}
	err := repo.UpsertTrack(track)
	if err != nil {
		t.Errorf("expected no error on UpsertTrack, got %v", err)
	}

	// Test Rate
	rate := &model.Rate{
		UserID:  user.ID,
		MovieID: movie.ID,
		Rating:  4.5,
	}
	err = repo.UpsertRate(rate)
	if err != nil {
		t.Errorf("expected no error on UpsertRate, got %v", err)
	}

	// Test Review
	review := &model.Review{
		UserID:  user.ID,
		MovieID: movie.ID,
		Content: "Great movie!",
	}
	err = repo.UpsertReview(review)
	if err != nil {
		t.Errorf("expected no error on UpsertReview, got %v", err)
	}

	// Verify UserInteraction
	tTrack, tRate, tReview, err := repo.GetUserInteraction(user.ID, movie.ID)
	if err != nil {
		t.Errorf("expected no error on GetUserInteraction, got %v", err)
	}
	if !tTrack.IsWatched || tRate.Rating != 4.5 || tReview.Content != "Great movie!" {
		t.Errorf("interactions not saved correctly")
	}
}

func TestMovieRepository_Counts(t *testing.T) {
	db := setupMovieTestDB(t)
	repo := NewMovieRepository(db)

	user := &model.User{Username: "userCounts", Email: "counts@example.com"}
	db.Create(user)

	movie1 := &model.Movie{TmdbID: 1, Title: "Movie 1"}
	movie2 := &model.Movie{TmdbID: 2, Title: "Movie 2"}
	repo.UpsertMovie(movie1)
	repo.UpsertMovie(movie2)

	// Watch two movies
	now := time.Now()
	repo.UpsertTrack(&model.Track{UserID: user.ID, MovieID: movie1.ID, IsWatched: true, WatchedDate: &now})
	repo.UpsertTrack(&model.Track{UserID: user.ID, MovieID: movie2.ID, IsWatched: true, WatchedDate: &now})
	repo.UpsertReview(&model.Review{UserID: user.ID, MovieID: movie1.ID, Content: "Review 1"})

	watched, _ := repo.CountWatched(user.ID)
	if watched != 2 {
		t.Errorf("expected 2 watched movies, got %d", watched)
	}

	watchedThisYear, _ := repo.CountWatchedThisYear(user.ID)
	if watchedThisYear != 2 {
		t.Errorf("expected 2 watched movies this year, got %d", watchedThisYear)
	}

	reviews, _ := repo.CountReviews(user.ID)
	if reviews != 1 {
		t.Errorf("expected 1 review, got %d", reviews)
	}
}

func TestMovieRepository_GetLists(t *testing.T) {
	db := setupMovieTestDB(t)
	repo := NewMovieRepository(db)

	user := &model.User{Username: "userLists", Email: "lists@example.com"}
	db.Create(user)

	m1 := model.Movie{TmdbID: 10, Title: "Fav 1"}
	m2 := model.Movie{TmdbID: 20, Title: "Fav 2"}

	err := repo.UpdateFavoriteFilms(user.ID, []model.Movie{m1, m2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	favs, err := repo.GetFavoriteMovies(user.ID, 5)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(favs) != 2 {
		t.Errorf("expected 2 favorite movies, got %d", len(favs))
	}

	repo.UpsertTrack(&model.Track{UserID: user.ID, MovieID: favs[0].ID, IsWatched: true, WatchedDate: &time.Time{}})
	recent, _ := repo.GetRecentWatched(user.ID, 5)
	if len(recent) != 1 {
		t.Errorf("expected 1 recent watched movie, got %d", len(recent))
	}
}

func TestMovieRepository_GetRatingDistribution(t *testing.T) {
	db := setupMovieTestDB(t)
	repo := NewMovieRepository(db)

	user := &model.User{Username: "userDist", Email: "dist@example.com"}
	db.Create(user)

	m1 := &model.Movie{TmdbID: 30, Title: "Rating 1"}
	m2 := &model.Movie{TmdbID: 40, Title: "Rating 2"}
	repo.UpsertMovie(m1)
	repo.UpsertMovie(m2)

	repo.UpsertRate(&model.Rate{UserID: user.ID, MovieID: m1.ID, Rating: 4.5})
	repo.UpsertRate(&model.Rate{UserID: user.ID, MovieID: m2.ID, Rating: 5.0})

	dist, err := repo.GetRatingDistribution(user.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if dist["4.5"] != 1 || dist["5.0"] != 1 {
		t.Errorf("rating distribution incorrect, got %v", dist)
	}
}
