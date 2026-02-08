package dto

type SearchMoviesRequest struct {
	Query    string `form:"q" binding:"required,min=2"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Year     int    `form:"year" binding:"omitempty,min=1900"`
	GenreID  int    `form:"genre_id" binding:"omitempty"`
	Language string `form:"language" binding:"omitempty,len=5"`
}

type SearchMoviesResponse struct {
	Results      []MovieSearchResult `json:"results"`
	Page         int                 `json:"page"`
	TotalPages   int                 `json:"total_pages"`
	TotalResults int                 `json:"total_results"`
}

type MovieSearchResult struct {
	TmdbID           int     `json:"tmdb_id"`
	Title            string  `json:"title"`
	Overview         string  `json:"overview"`
	PosterPath       *string `json:"poster_path"`
	BackdropPath     *string `json:"backdrop_path"`
	ReleaseDate      string  `json:"release_date"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	IsInDatabase     bool    `json:"is_in_database"`
	LocalRating      float32 `json:"local_rating"`
	LocalRatingCount int     `json:"local_rating_count"`
}
