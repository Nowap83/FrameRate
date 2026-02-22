package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/service"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupTMDBHandlerTest(mockResp interface{}, statusCode int) (*gin.Engine, *httptest.Server) {
	utils.Log = zap.NewNop()
	gin.SetMode(gin.TestMode)

	respBytes, _ := json.Marshal(mockResp)

	// Create a mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(respBytes)
	}))

	// Override the TMDB_BASE_URL
	os.Setenv("TMDB_BASE_URL", ts.URL)
	os.Setenv("TMDB_API_KEY", "testkey")

	tmdbService := service.NewTMDBService(nil)
	tmdbHandler := NewTMDBHandler(tmdbService)

	r := gin.New()
	api := r.Group("/tmdb")
	{
		api.GET("/search", tmdbHandler.SearchMovies)
		api.GET("/popular", tmdbHandler.GetPopularMovies)
		api.GET("/movie/:id", tmdbHandler.GetMovieDetails)
		api.GET("/movie/:id/credits", tmdbHandler.GetMovieCredits)
		api.GET("/movie/:id/videos", tmdbHandler.GetMovieVideos)
		api.GET("/person/:id", tmdbHandler.GetPersonDetails)
		api.GET("/person/:id/movie_credits", tmdbHandler.GetPersonMovieCredits)
		api.GET("/image", tmdbHandler.GetImageURL)
	}

	return r, ts
}

func TestTMDBHandler_SearchMovies(t *testing.T) {
	mockResp := dto.TMDBSearchResponse{
		Results: []dto.TMDBMovie{{ID: 1, Title: "Test Search"}},
	}
	r, ts := setupTMDBHandlerTest(mockResp, 200)
	defer ts.Close()

	req, _ := http.NewRequest("GET", "/tmdb/search?q=test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Test Search")) {
		t.Errorf("expected to find 'Test Search' in response")
	}
}

func TestTMDBHandler_GetPopularMovies(t *testing.T) {
	mockResp := dto.TMDBSearchResponse{
		Results: []dto.TMDBMovie{{ID: 2, Title: "Popular"}},
	}
	r, ts := setupTMDBHandlerTest(mockResp, 200)
	defer ts.Close()

	req, _ := http.NewRequest("GET", "/tmdb/popular?page=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestTMDBHandler_GetMovieDetails_And_Others(t *testing.T) {
	mockResp := dto.TMDBMovieDetails{ID: 3, Title: "Details"}
	r, ts := setupTMDBHandlerTest(mockResp, 200)
	defer ts.Close()

	endpoints := []string{
		"/tmdb/movie/3",
		"/tmdb/movie/3/credits",
		"/tmdb/movie/3/videos",
		"/tmdb/person/4",
		"/tmdb/person/4/movie_credits",
	}

	for _, endpoint := range endpoints {
		req, _ := http.NewRequest("GET", endpoint, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200 OK for %s, got %d", endpoint, w.Code)
		}
	}
}

func TestTMDBHandler_GetImageURL(t *testing.T) {
	r, ts := setupTMDBHandlerTest(nil, 200)
	defer ts.Close()

	req, _ := http.NewRequest("GET", "/tmdb/image?path=/test.jpg&size=w500", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	// Invalid size
	req2, _ := http.NewRequest("GET", "/tmdb/image?path=/test.jpg&size=invalid", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for invalid size, got %d", w2.Code)
	}
}
