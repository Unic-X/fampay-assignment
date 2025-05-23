package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

// New creates a new database connection
func New(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &DB{db}, nil
}

// SaveVideo saves a video to the database
func (db *DB) SaveVideo(video *Video) error {
	query := `
		INSERT INTO videos (id, title, description, published_at, thumbnail_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE
		SET title = $2, description = $3, published_at = $4, thumbnail_url = $5, updated_at = $7
	`

	now := time.Now()
	_, err := db.Exec(query,
		video.ID,
		video.Title,
		video.Description,
		video.PublishedAt,
		video.ThumbnailURL,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("error saving video: %v", err)
	}

	return nil
}

// GetVideos returns paginated videos sorted by the specified field and order
func (db *DB) GetVideos(page, limit int, sortBy, sortOrder string) (*PaginatedResponse, error) {
	offset := (page - 1) * limit

	// Get total count
	var totalCount int
	err := db.QueryRow("SELECT COUNT(*) FROM videos").Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("error getting total count: %v", err)
	}

	// Build the query with dynamic sorting
	query := fmt.Sprintf(`
		SELECT id, title, description, published_at, thumbnail_url, created_at, updated_at
		FROM videos
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sortBy, sortOrder)

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying videos: %v", err)
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var v Video
		err := rows.Scan(
			&v.ID,
			&v.Title,
			&v.Description,
			&v.PublishedAt,
			&v.ThumbnailURL,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning video: %v", err)
		}
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating videos: %v", err)
	}

	totalPages := (totalCount + limit - 1) / limit

	return &PaginatedResponse{
		Videos:     videos,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

