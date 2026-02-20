package dto

type TMDBMovie struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Overview     string  `json:"overview"`
	ReleaseDate  string  `json:"release_date"`
	PosterPath   *string `json:"poster_path"`
	BackdropPath *string `json:"backdrop_path"`
	VoteAverage  float64 `json:"vote_average"`
	VoteCount    int     `json:"vote_count"`
}

type TMDBSearchResponse struct {
	Page         int         `json:"page"`
	Results      []TMDBMovie `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}

type TMDBMovieDetails struct {
	ID               int          `json:"id"`
	Title            string       `json:"title"`
	OriginalTitle    string       `json:"original_title"`
	Overview         string       `json:"overview"`
	ReleaseDate      string       `json:"release_date"`
	Runtime          int          `json:"runtime"`
	Budget           int64        `json:"budget"`
	Revenue          int64        `json:"revenue"`
	PosterPath       *string      `json:"poster_path"`
	BackdropPath     *string      `json:"backdrop_path"`
	VoteAverage      float64      `json:"vote_average"`
	VoteCount        int          `json:"vote_count"`
	ImdbID           string       `json:"imdb_id"`
	OriginalLanguage string       `json:"original_language"`
	Genres           []TMDBGenre  `json:"genres"`
	Credits          *TMDBCredits `json:"credits,omitempty"`
}

type TMDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TMDBCredits struct {
	ID   int              `json:"id"`
	Cast []TMDBCastMember `json:"cast"`
	Crew []TMDBCrewMember `json:"crew"`
}

type TMDBCastMember struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Character   string  `json:"character"`
	ProfilePath *string `json:"profile_path"`
	Order       int     `json:"order"`
	Gender      int     `json:"gender"`
}

type TMDBCrewMember struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Job         string  `json:"job"`
	Department  string  `json:"department"`
	ProfilePath *string `json:"profile_path"`
	Gender      int     `json:"gender"`
}

type TMDBPersonDetails struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Biography    string  `json:"biography"`
	Birthday     *string `json:"birthday"`
	Deathday     *string `json:"deathday"`
	PlaceOfBirth *string `json:"place_of_birth"`
	ProfilePath  *string `json:"profile_path"`
	Gender       int     `json:"gender"`
}

type TMDBPersonCredits struct {
	ID   int                   `json:"id"`
	Cast []TMDBPersonCastMovie `json:"cast"`
	Crew []TMDBPersonCrewMovie `json:"crew"`
}

type TMDBPersonCastMovie struct {
	TMDBMovie
	Character string `json:"character"`
}

type TMDBPersonCrewMovie struct {
	TMDBMovie
	Job        string `json:"job"`
	Department string `json:"department"`
}

type TMDBVideo struct {
	ID          string `json:"id"`
	ISO639_1    string `json:"iso_639_1"`
	ISO3166_1   string `json:"iso_3166_1"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Site        string `json:"site"`
	Size        int    `json:"size"`
	Type        string `json:"type"`
	Official    bool   `json:"official"`
	PublishedAt string `json:"published_at"`
}

type TMDBVideoResponse struct {
	ID      int         `json:"id"`
	Results []TMDBVideo `json:"results"`
}
