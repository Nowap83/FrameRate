package repository

import (
	"github.com/Nowap83/FrameRate/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) GetMovieByTmdbID(tmdbID int) (*model.Movie, error) {
	var movie model.Movie
	err := r.db.Where("tmdb_id = ?", tmdbID).First(&movie).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) UpsertMovie(movie *model.Movie) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tmdb_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "original_title", "release_year", "duration_minutes", "synopsis", "poster_url", "backdrop_url", "language", "updated_at"}),
	}).Create(movie).Error
}

func (r *MovieRepository) UpsertTrack(track *model.Track) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "movie_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_watched", "is_favorite", "is_watchlist", "watched_date", "updated_at"}),
	}).Create(track).Error
}

func (r *MovieRepository) UpsertRate(rate *model.Rate) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "movie_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"rating", "updated_at"}),
	}).Create(rate).Error
}

func (r *MovieRepository) UpsertReview(review *model.Review) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "movie_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"content", "is_spoiler", "updated_at"}),
	}).Create(review).Error
}

func (r *MovieRepository) GetUserInteraction(userID uint, movieID uint) (*model.Track, *model.Rate, *model.Review, error) {
	var track model.Track
	var rate model.Rate
	var review model.Review

	_ = r.db.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&track).Error
	_ = r.db.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&rate).Error
	_ = r.db.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&review).Error

	return &track, &rate, &review, nil
}
