package database

import "time"

// Video represents a YouTube video in our database
type Video struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	ThumbnailURL string   `json:"thumbnail_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PaginatedResponse represents a paginated response of videos
type PaginatedResponse struct {
	Videos     []Video `json:"videos"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
} 