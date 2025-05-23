package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Unic-X/fampay-assignment/internal/database"
)

type VideoHandler struct {
	db *database.DB
}

func NewVideoHandler(db *database.DB) *VideoHandler {
	return &VideoHandler{db: db}
}

// GetVideos handles GET /api/videos
func (h *VideoHandler) GetVideos(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	videos, err := h.db.GetVideos(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
		return
	}

	c.JSON(http.StatusOK, videos)
} 