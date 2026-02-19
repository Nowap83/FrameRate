package repository

import (
	"fmt"
	"time"

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

func (r *MovieRepository) CountWatched(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Track{}).Where("user_id = ? AND is_watched = ?", userID, true).Count(&count).Error
	return count, err
}

func (r *MovieRepository) CountWatchedThisYear(userID uint) (int64, error) {
	var count int64
	currentYear := time.Now().Year()
	startOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.Model(&model.Track{}).
		Where("user_id = ? AND is_watched = ? AND watched_date >= ?", userID, true, startOfYear).
		Count(&count).Error
	return count, err
}

func (r *MovieRepository) CountReviews(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Review{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *MovieRepository) GetFavoriteMovies(userID uint, limit int) ([]model.Movie, error) {
	var movies []model.Movie
	err := r.db.Joins("JOIN tracks ON tracks.movie_id = movies.id").
		Where("tracks.user_id = ? AND tracks.is_favorite = ?", userID, true).
		Limit(limit).
		Find(&movies).Error
	return movies, err
}

func (r *MovieRepository) GetRecentWatched(userID uint, limit int) ([]model.Movie, error) {
	var movies []model.Movie
	err := r.db.Joins("JOIN tracks ON tracks.movie_id = movies.id").
		Where("tracks.user_id = ? AND tracks.is_watched = ?", userID, true).
		Order("tracks.watched_date DESC").
		Limit(limit).
		Find(&movies).Error
	return movies, err
}

func (r *MovieRepository) GetRatingDistribution(userID uint) (map[string]int, error) {
	rows, err := r.db.Model(&model.Rate{}).
		Select("rating, count(*) as count").
		Where("user_id = ?", userID).
		Group("rating").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	distribution := make(map[string]int)
	for rows.Next() {
		var rating float64
		var count int
		if err := rows.Scan(&rating, &count); err != nil {
			continue
		}
		key := fmt.Sprintf("%.1f", rating)
		distribution[key] = count
	}
	return distribution, nil
}

func (r *MovieRepository) UpdateFavoriteFilms(userID uint, movies []model.Movie) error {
	var savedMovies []model.Movie

	for _, movie := range movies {
		if err := r.UpsertMovie(&movie); err != nil {
			return err
		}
		var saved model.Movie
		if err := r.db.Where("tmdb_id = ?", movie.TmdbID).First(&saved).Error; err != nil {
			return err
		}
		savedMovies = append(savedMovies, saved)
	}

	var user model.User
	user.ID = userID

	return r.db.Model(&user).Association("FavoriteFilms").Replace(&savedMovies)
}
