package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/Nowap83/FrameRate/backend/internal/repository"
	"gorm.io/gorm"
)

type MovieService struct {
	movieRepo   *repository.MovieRepository
	tmdbService *TMDBService
}

func NewMovieService(movieRepo *repository.MovieRepository, tmdbService *TMDBService) *MovieService {
	return &MovieService{
		movieRepo:   movieRepo,
		tmdbService: tmdbService,
	}
}

func (s *MovieService) ensureMovieExists(tmdbID int) (*model.Movie, error) {
	movie, err := s.movieRepo.GetMovieByTmdbID(tmdbID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Fetch from TMDB
			tmdbMovie, err := s.tmdbService.GetMovieDetails(tmdbID, "fr-FR")
			if err != nil {
				return nil, err
			}

			movie = &model.Movie{
				TmdbID:          tmdbMovie.ID,
				Title:           tmdbMovie.Title,
				OriginalTitle:   tmdbMovie.OriginalTitle,
				DurationMinutes: tmdbMovie.Runtime,
				Synopsis:        tmdbMovie.Overview,
				Language:        tmdbMovie.OriginalLanguage,
			}

			if tmdbMovie.PosterPath != nil {
				movie.PosterURL = *tmdbMovie.PosterPath
			}
			if tmdbMovie.BackdropPath != nil {
				movie.BackdropURL = *tmdbMovie.BackdropPath
			}

			if tmdbMovie.ReleaseDate != "" {
				// Simple year parsing
				if len(tmdbMovie.ReleaseDate) >= 4 {
					var year int
					_, _ = fmt.Sscanf(tmdbMovie.ReleaseDate[:4], "%d", &year)
					movie.ReleaseYear = year
				}
			}

			err = s.movieRepo.UpsertMovie(movie)
			if err != nil {
				return nil, err
			}
			return movie, nil
		}
		return nil, err
	}
	return movie, nil
}

func (s *MovieService) TrackMovie(userID uint, tmdbID int, req dto.TrackMovieRequest) error {
	movie, err := s.ensureMovieExists(tmdbID)
	if err != nil {
		return err
	}

	track := &model.Track{
		UserID:  userID,
		MovieID: movie.ID,
	}

	if req.IsWatched != nil {
		track.IsWatched = *req.IsWatched
	}
	if req.IsFavorite != nil {
		track.IsFavorite = *req.IsFavorite
	}
	if req.IsWatchlist != nil {
		track.IsWatchlist = *req.IsWatchlist
	}
	if req.WatchedDate != nil {
		track.WatchedDate = req.WatchedDate
	}

	return s.movieRepo.UpsertTrack(track)
}

func (s *MovieService) RateMovie(userID uint, tmdbID int, req dto.RateMovieRequest) error {
	movie, err := s.ensureMovieExists(tmdbID)
	if err != nil {
		return err
	}

	// validation de la note => palier de 0.5
	if float64(req.Rating*2) != float64(int(req.Rating*2)) {
		return errors.New("Rating must be in increments of 0.5")
	}

	rate := &model.Rate{
		UserID:  userID,
		MovieID: movie.ID,
		Rating:  req.Rating,
	}

	err = s.movieRepo.UpsertRate(rate)
	if err != nil {
		return err
	}

	// auto-mark as watched when rated
	now := time.Now()
	track := &model.Track{
		UserID:      userID,
		MovieID:     movie.ID,
		IsWatched:   true,
		WatchedDate: &now,
	}
	return s.movieRepo.UpsertTrack(track)
}

func (s *MovieService) ReviewMovie(userID uint, tmdbID int, req dto.ReviewRequest) error {
	movie, err := s.ensureMovieExists(tmdbID)
	if err != nil {
		return err
	}

	review := &model.Review{
		UserID:    userID,
		MovieID:   movie.ID,
		Content:   req.Content,
		IsSpoiler: req.IsSpoiler,
	}

	return s.movieRepo.UpsertReview(review)
}

func (s *MovieService) GetMovieInteraction(userID uint, tmdbID int) (*dto.UserInteractionResponse, error) {
	movie, err := s.movieRepo.GetMovieByTmdbID(tmdbID)
	if err != nil {
		// no interaction if movie doesn't exist
		return &dto.UserInteractionResponse{
			IsWatched:   false,
			IsFavorite:  false,
			IsWatchlist: false,
		}, nil
	}

	track, rate, review, _ := s.movieRepo.GetUserInteraction(userID, movie.ID)

	response := &dto.UserInteractionResponse{
		IsWatched:   false,
		IsFavorite:  false,
		IsWatchlist: false,
	}

	if track != nil {
		response.IsWatched = track.IsWatched
		response.IsFavorite = track.IsFavorite
		response.IsWatchlist = track.IsWatchlist
		response.WatchedDate = track.WatchedDate
	}

	if rate != nil && rate.Rating > 0 {
		r := rate.Rating
		response.UserRating = &r
	}

	if review != nil && review.Content != "" {
		response.UserReview = &dto.ReviewResponse{
			Content:   review.Content,
			IsSpoiler: review.IsSpoiler,
			CreatedAt: review.CreatedAt,
		}
	}

	return response, nil
}
