package dto

import "time"

type MovieListResponse struct {
	ID                uint    `json:"id"`
	TmdbID            int     `json:"tmdb_id"`
	Title             string  `json:"title"`
	ReleaseYear       int     `json:"release_year"`
	PosterURL         string  `json:"poster_url"`
	AverageUserRating float32 `json:"average_user_rating"`
	TotalRatings      int     `json:"total_ratings"`
}

type MovieResponse struct {
	ID                uint             `json:"id"`
	TmdbID            int              `json:"tmdb_id"`
	Title             string           `json:"title"`
	OriginalTitle     string           `json:"original_title"`
	ReleaseYear       int              `json:"release_year"`
	DurationMinutes   int              `json:"duration_minutes"`
	Synopsis          string           `json:"synopsis"`
	PosterURL         string           `json:"poster_url"`
	BackdropURL       string           `json:"backdrop_url"`
	TrailerURL        string           `json:"trailer_url"`
	ImdbRating        float32          `json:"imdb_rating"`
	AverageUserRating float32          `json:"average_user_rating"`
	TotalRatings      int              `json:"total_ratings"`
	Genres            []GenreResponse  `json:"genres"`
	Directors         []PersonResponse `json:"directors"`
	TopCast           []CastResponse   `json:"top_cast"` // 5 premiers
	CreatedAt         time.Time        `json:"created_at"`
}

type MovieDetailResponse struct {
	MovieResponse
	Budget              int64                    `json:"budget"`
	Revenue             int64                    `json:"revenue"`
	MetacriticScore     int                      `json:"metacritic_score"`
	RottenTomatoesScore int                      `json:"rotten_tomatoes_score"`
	Language            string                   `json:"language"`
	Countries           []CountryResponse        `json:"countries"`
	FullCast            []CastResponse           `json:"full_cast"`
	Crew                []CrewResponse           `json:"crew"`
	UserInteraction     *UserInteractionResponse `json:"user_interaction,omitempty"` // Si connecté
}

type GenreResponse struct {
	ID     uint   `json:"id"`
	TmdbID int    `json:"tmdb_id"`
	Name   string `json:"name"`
}

type PersonResponse struct {
	ID                uint   `json:"id"`
	TmdbID            int    `json:"tmdb_id"`
	Name              string `json:"name"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

type CastResponse struct {
	PersonResponse
	CharacterName string `json:"character_name"`
	CastOrder     int    `json:"cast_order"`
}

type CrewResponse struct {
	PersonResponse
	Job        string `json:"job"`
	Department string `json:"department"`
}

type CountryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type UserInteractionResponse struct {
	IsWatched   bool            `json:"is_watched"`
	IsFavorite  bool            `json:"is_favorite"`
	IsWatchlist bool            `json:"is_watchlist"`
	UserRating  *float32        `json:"user_rating,omitempty"`
	UserReview  *ReviewResponse `json:"user_review,omitempty"`
	WatchedDate *time.Time      `json:"watched_date,omitempty"`
}

type MovieReviewResponse struct {
	Content   string    `json:"content"`
	IsSpoiler bool      `json:"is_spoiler"`
	CreatedAt time.Time `json:"created_at"`
}

type TrackMovieRequest struct {
	IsWatched   *bool      `json:"is_watched"`
	IsFavorite  *bool      `json:"is_favorite"`
	IsWatchlist *bool      `json:"is_watchlist"`
	WatchedDate *time.Time `json:"watched_date"`
}

type RateMovieRequest struct {
	Rating float32 `json:"rating" binding:"required,min=0,max=5"`
}

type ReviewRequest struct {
	Content   string `json:"content" binding:"required"`
	IsSpoiler bool   `json:"is_spoiler"`
}

type ReviewResponse struct {
	Content   string    `json:"content"`
	IsSpoiler bool      `json:"is_spoiler"`
	CreatedAt time.Time `json:"created_at"`
}
