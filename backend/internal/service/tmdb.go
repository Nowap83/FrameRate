package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
)

type TMDBService struct {
	apiKey       string
	baseURL      string
	imageBaseURL string
	client       *http.Client
	cache        *CacheService
}

func NewTMDBService(cache *CacheService) *TMDBService {
	return &TMDBService{
		apiKey:       os.Getenv("TMDB_API_KEY"),
		baseURL:      os.Getenv("TMDB_BASE_URL"),
		imageBaseURL: os.Getenv("TMDB_IMAGE_BASE_URL"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: cache,
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

	// clé du cache
	cacheKey := fmt.Sprintf("tmdb:search:%s:%d:%s", query.Query, query.Page, query.Language)
	var cachedResult dto.TMDBSearchResponse

	if s.cache != nil {
		found, err := s.cache.Get(context.Background(), cacheKey, &cachedResult)
		if err == nil && found {
			utils.Log.Info(fmt.Sprintf("Cache hit for search: %s", query.Query))
			return &cachedResult, nil
		}
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

	// store dans cache (15min pour les recherches)
	if s.cache != nil {
		_ = s.cache.Set(context.Background(), cacheKey, result, 15*time.Minute)
	}

	return &result, nil
}

func (s *TMDBService) GetMovieDetails(tmdbID int, language string) (*dto.TMDBMovieDetails, error) {
	// valeurs par défaut
	if language == "" {
		language = "en-US"
	}

	// clé du cache
	cacheKey := fmt.Sprintf("tmdb:movie:%d:%s", tmdbID, language)
	var cachedResult dto.TMDBMovieDetails

	if s.cache != nil {
		found, err := s.cache.Get(context.Background(), cacheKey, &cachedResult)
		if err == nil && found {
			utils.Log.Info(fmt.Sprintf("Cache hit for aggregated movie details: %d", tmdbID))
			return &cachedResult, nil
		}
	}

	// build de l'url avec append_to_response
	url := fmt.Sprintf("%s/movie/%d?language=%s&append_to_response=credits", s.baseURL, tmdbID, language)

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

	// store dans cache (24h pour les détails)
	if s.cache != nil {
		_ = s.cache.Set(context.Background(), cacheKey, result, 24*time.Hour)
	}

	return &result, nil
}

func (s *TMDBService) GetMovieCredits(tmdbID int) (*dto.TMDBCredits, error) {
	// clé du cache
	cacheKey := fmt.Sprintf("tmdb:credits:%d", tmdbID)
	var cachedResult dto.TMDBCredits

	if s.cache != nil {
		found, err := s.cache.Get(context.Background(), cacheKey, &cachedResult)
		if err == nil && found {
			utils.Log.Info(fmt.Sprintf("Cache hit for movie credits: %d", tmdbID))
			return &cachedResult, nil
		}
	}

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

	// store dans cache (24h pour les crédits)
	if s.cache != nil {
		_ = s.cache.Set(context.Background(), cacheKey, result, 24*time.Hour)
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
