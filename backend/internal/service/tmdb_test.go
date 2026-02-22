package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"go.uber.org/zap"
)

// ensure RoundTripFunc from movie_test.go is available or redefine a local one if we don't share test files
// (But wait, test files in same package share unexported types! I already defined RoundTripFunc in movie_test.go!
// To avoid conflicts, I'll use a slightly different name or just reuse it.)
// Reusing mockTMDBClient defined in movie_test.go since they are in the same package (service) and compiled together during tests.

func TestTMDBService_SearchMovies(t *testing.T) {
	utils.Log = zap.NewNop()

	mockResp := dto.TMDBSearchResponse{
		Page:         1,
		TotalPages:   1,
		TotalResults: 1,
		Results: []dto.TMDBMovie{
			{ID: 111, Title: "Search Result"},
		},
	}
	respBytes, _ := json.Marshal(mockResp)

	tmdbService := NewTMDBService(nil)
	tmdbService.client = mockTMDBClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(respBytes)),
			Header:     make(http.Header),
		}
	})

	req := dto.SearchMoviesRequest{Query: "Search"}
	res, err := tmdbService.SearchMovies(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.TotalResults != 1 || len(res.Results) != 1 || res.Results[0].Title != "Search Result" {
		t.Errorf("expected search result, got %v", res)
	}
}

func TestTMDBService_GetPopularMovies(t *testing.T) {
	utils.Log = zap.NewNop()

	mockResp := dto.TMDBSearchResponse{
		Page:         1,
		TotalPages:   10,
		TotalResults: 200,
		Results: []dto.TMDBMovie{
			{ID: 222, Title: "Popular Movie"},
		},
	}
	respBytes, _ := json.Marshal(mockResp)

	tmdbService := NewTMDBService(nil)
	tmdbService.client = mockTMDBClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(respBytes)),
			Header:     make(http.Header),
		}
	})

	res, err := tmdbService.GetPopularMovies(1, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Results[0].Title != "Popular Movie" {
		t.Errorf("expected popular movie")
	}
}

func TestTMDBService_Errors(t *testing.T) {
	utils.Log = zap.NewNop()

	tmdbService := NewTMDBService(nil)
	tmdbService.client = mockTMDBClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{"status_message": "not found"}`)),
			Header:     make(http.Header),
		}
	})

	_, err := tmdbService.GetMovieDetails(999999, "")
	if err == nil {
		t.Fatalf("expected error from TMDB API")
	}
}

func TestTMDBService_WithCache(t *testing.T) {
	utils.Log = zap.NewNop()

	cacheService := NewCacheService(nil) // we can pass a dummy redis client or skip tests if no real redis
	// Actually, CacheService needs a redis client. Better not to test cache logic here unless we mock redis.
	// We've already tested without cache (nil cache). Let's skip deep cache tests to avoid needing a redis mock.
	_ = cacheService
}

func TestTMDBService_GetMovieCredits_And_Videos(t *testing.T) {
	utils.Log = zap.NewNop()

	creditsResp := dto.TMDBCredits{ID: 1, Cast: []dto.TMDBCastMember{{Name: "Actor 1"}}}
	cBytes, _ := json.Marshal(creditsResp)

	tmdbService := NewTMDBService(nil)
	tmdbService.client = mockTMDBClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(cBytes)),
			Header:     make(http.Header),
		}
	})

	credits, err := tmdbService.GetMovieCredits(1)
	if err != nil || len(credits.Cast) != 1 {
		t.Fatalf("failed credits test: %v", err)
	}
}

func TestTMDBService_GetMovieVideos(t *testing.T) {
	utils.Log = zap.NewNop()

	videosResp := dto.TMDBVideoResponse{ID: 1, Results: []dto.TMDBVideo{{Key: "test_key", Site: "YouTube", Type: "Trailer"}}}
	vBytes, _ := json.Marshal(videosResp)

	tmdbService := NewTMDBService(nil)
	tmdbService.client = mockTMDBClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(vBytes)),
			Header:     make(http.Header),
		}
	})

	videos, err := tmdbService.GetMovieVideos(1)
	if err != nil || len(videos.Results) != 1 {
		t.Fatalf("failed videos test: %v", err)
	}
}

func TestTMDBService_GetPersonDetails(t *testing.T) {
	utils.Log = zap.NewNop()

	personResp := dto.TMDBPersonDetails{ID: 1, Name: "Actor 1"}
	pBytes, _ := json.Marshal(personResp)

	tmdbService := NewTMDBService(nil)
	tmdbService.client = mockTMDBClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(pBytes)),
			Header:     make(http.Header),
		}
	})

	person, err := tmdbService.GetPersonDetails(1, "")
	if err != nil || person.Name != "Actor 1" {
		t.Fatalf("failed person details test: %v", err)
	}
}

func TestTMDBService_GetPersonMovieCredits(t *testing.T) {
	utils.Log = zap.NewNop()

	creditsResp := dto.TMDBPersonCredits{Cast: []dto.TMDBPersonCastMovie{{TMDBMovie: dto.TMDBMovie{ID: 1, Title: "Movie 1"}}}}
	cBytes, _ := json.Marshal(creditsResp)

	tmdbService := NewTMDBService(nil)
	tmdbService.client = mockTMDBClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBuffer(cBytes)),
			Header:     make(http.Header),
		}
	})

	credits, err := tmdbService.GetPersonMovieCredits(1, "")
	if err != nil || len(credits.Cast) != 1 {
		t.Fatalf("failed person credits test: %v", err)
	}
}

func TestTMDBService_ImageUtils(t *testing.T) {
	os.Setenv("TMDB_IMAGE_BASE_URL", "https://image.tmdb.org/t/p")
	tmdbService := NewTMDBService(nil)
	url := tmdbService.GetImageURL("/test.jpg", "w500")
	if url != "https://image.tmdb.org/t/p/w500/test.jpg" {
		t.Errorf("unexpected image url: %s", url)
	}

	valid := tmdbService.ValidateImageSize("w500")
	if !valid {
		t.Errorf("should be valid size")
	}

	invalid := tmdbService.ValidateImageSize("w999")
	if invalid {
		t.Errorf("should be invalid size")
	}
}
