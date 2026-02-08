package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Nowap83/FrameRate/backend/dto"
)

type TMDBService struct {
	apiKey       string
	baseURL      string
	imageBaseURL string
	client       *http.Client
}

func NewTMDBService() *TMDBService {
	return &TMDBService{
		apiKey:       os.Getenv("TMDB_API_KEY"),
		baseURL:      os.Getenv("TMDB_BASE_URL"),
		imageBaseURL: os.Getenv("TMDB_IMAGE_BASE_URL"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *TMDBService) SearchMovies(query dto.SearchMoviesRequest) (*dto.TMDBSearchResponse, error) {
	// valeurs par défaut
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Language == "" {
		query.Language = "fr-FR"
	}

	// build de l'url
	url := fmt.Sprintf("%s/search/movie?query=%s&page=%d&language=%s",
		s.baseURL,
		url.QueryEscape(query.Query), // caracteres speciaux
		query.Page,
		query.Language)

	// requete http
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	// execution de la requete
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("TMDB API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TMDB API returned status %d: %s",
			resp.StatusCode, string(bodyBytes))
	}

	// parser response json
	var result dto.TMDBSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode TMDB response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetMovieDetails(tmdbID int, language string) (*dto.TMDBMovieDetails, error) {
	// valeurs par défaut
	if language == "" {
		language = "en-US"
	}

	// build de l'url
	url := fmt.Sprintf("%s/movie/%d?language=%s", s.baseURL, tmdbID, language)

	// requete http
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	// execution de la requete
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("TMDB API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TMDB API returned status %d: %s",
			resp.StatusCode, string(bodyBytes))
	}

	// parser response json
	var result dto.TMDBMovieDetails
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode TMDB response: %w", err)
	}

	return &result, nil
}

func (s *TMDBService) GetMovieCredits(tmdbID int) (*dto.TMDBCredits, error) {
	// build de l'url
	url := fmt.Sprintf("%s/movie/%d/credits", s.baseURL, tmdbID)

	// requete http
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	// execution de la requete
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("TMDB API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TMDB API returned status %d: %s",
			resp.StatusCode, string(bodyBytes))
	}

	// parser response json
	var result dto.TMDBCredits
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode TMDB response: %w", err)
	}

	return &result, nil
}

// helper pour construire les URLs d'images
func (s *TMDBService) GetImageURL(path string, size string) string {
	if path == "" {
		return ""
	}
	// Sizes disponibles: w92, w154, w185, w342, w500, w780, original
	return fmt.Sprintf("%s/%s%s", s.imageBaseURL, size, path)
}

// helper pour valider la taille d'image
func (s *TMDBService) ValidateImageSize(size string) bool {
	validSizes := map[string]bool{
		"w92":      true,
		"w154":     true,
		"w185":     true,
		"w342":     true,
		"w500":     true,
		"w780":     true,
		"original": true,
	}
	return validSizes[size]
}
