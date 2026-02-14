package handler

import (
	"net/http"
	"strconv"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type TMDBHandler struct {
	tmdbService *service.TMDBService
}

func NewTMDBHandler(tmdbService *service.TMDBService) *TMDBHandler {
	return &TMDBHandler{
		tmdbService: tmdbService,
	}
}

// * @param : ?query=inception&page=1&language=en-US
func (h *TMDBHandler) SearchMovies(c *gin.Context) {
	var req dto.SearchMoviesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid search parameters",
			"details": err.Error(),
		})
		return
	}

	results, err := h.tmdbService.SearchMovies(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search movies",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
	})
}

// * @param: ?language=fr-FR (optionnel)
func (h *TMDBHandler) GetMovieDetails(c *gin.Context) {

	tmdbID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid movie ID",
		})
		return
	}

	language := c.DefaultQuery("language", "en-US")

	details, err := h.tmdbService.GetMovieDetails(tmdbID, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch movie details",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    details,
	})
}

func (h *TMDBHandler) GetMovieCredits(c *gin.Context) {

	tmdbID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid movie ID",
		})
		return
	}

	credits, err := h.tmdbService.GetMovieCredits(tmdbID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch movie credits",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    credits,
	})
}

// * @param: ?path=/abc.jpg&size=w500
func (h *TMDBHandler) GetImageURL(c *gin.Context) {
	path := c.Query("path")
	size := c.DefaultQuery("size", "w500")

	if !h.tmdbService.ValidateImageSize(size) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Invalid image size",
			"valid_sizes": []string{"w92", "w154", "w185", "w342", "w500", "w780", "original"},
		})
		return
	}

	imageURL := h.tmdbService.GetImageURL(path, size)

	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing image path",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"url":     imageURL,
	})
}
