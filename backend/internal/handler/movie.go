package handler

import (
	"net/http"
	"strconv"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	movieService *service.MovieService
}

func NewMovieHandler(movieService *service.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

func (h *MovieHandler) TrackMovie(c *gin.Context) {
	tmdbID, err := strconv.Atoi(c.Param("tmdb_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	userID, _ := c.Get("userID")

	var req dto.TrackMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.movieService.TrackMovie(userID.(uint), tmdbID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track movie", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie tracking updated successfully"})
}

func (h *MovieHandler) RateMovie(c *gin.Context) {
	tmdbID, err := strconv.Atoi(c.Param("tmdb_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	userID, _ := c.Get("userID")

	var req dto.RateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.movieService.RateMovie(userID.(uint), tmdbID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rate movie", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie rating updated successfully"})
}

func (h *MovieHandler) ReviewMovie(c *gin.Context) {
	tmdbID, err := strconv.Atoi(c.Param("tmdb_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	userID, _ := c.Get("userID")

	var req dto.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.movieService.ReviewMovie(userID.(uint), tmdbID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to review movie", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie review updated successfully"})
}
