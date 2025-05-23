package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Unic-X/fampay-assignment/internal/database"
	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	db *database.DB
}

func NewVideoHandler(db *database.DB) *VideoHandler {
	return &VideoHandler{db: db}
}

// GetVideos godoc
// @Summary Get paginated list of videos
// @Description Get a paginated list of videos sorted by published date in descending order
// @Tags videos
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 50)"
// @Param sort_by query string false "Field to sort by (published_at, title, created_at) (default: published_at)"
// @Param sort_order query string false "Sort order (asc, desc) (default: desc)"
// @Success 200 {object} database.PaginatedResponse
// @Failure 500 {object} map[string]string
// @Router /api/videos [get]
func (h *VideoHandler) GetVideos(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	// Get sorting parameters
	sortBy := strings.ToLower(c.DefaultQuery("sort_by", "published_at"))
	sortOrder := strings.ToLower(c.DefaultQuery("sort_order", "desc"))

	// Validate sort_by parameter
	validSortFields := map[string]string{
		"published_at": "published_at",
		"title":        "title",
		"created_at":   "created_at",
	}

	if _, valid := validSortFields[sortBy]; !valid {
		sortBy = "published_at"
	}

	// Validate sort_order parameter
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	videos, err := h.db.GetVideos(page, limit, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
		return
	}

	c.JSON(http.StatusOK, videos)
}

