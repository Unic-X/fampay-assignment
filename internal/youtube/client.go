package youtube

import (
	"context"
	"fmt"
	"strings"
	"time"

	_ "github.com/Unic-X/fampay-assignment/internal/database"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Video represents a YouTube video in our application
type Video struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	PublishedAt  time.Time `json:"published_at"`
	ThumbnailURL string    `json:"thumbnail_url"`
}

type Client struct {
	service     *youtube.Service
	apiKeys     []string
	currentKey  int
	searchQuery string
}

func New(apiKeys []string, searchQuery string) (*Client, error) {
	if len(apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided")
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKeys[0]))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube service: %v", err)
	}

	return &Client{
		service:     service,
		apiKeys:     apiKeys,
		currentKey:  0,
		searchQuery: searchQuery,
	}, nil
}

func (c *Client) FetchLatestVideos() ([]Video, error) {
	publishedAfter := time.Now().Add(-24 * time.Hour).Format(time.RFC3339) // 24 hours ago

	call := c.service.Search.List([]string{"snippet"}).
		Q(c.searchQuery).
		Type("video").
		Order("date").
		MaxResults(10).
		PublishedAfter(publishedAfter)

	response, err := call.Do()
	if err != nil {
		// Check if the error is due to quota exceeded
		if strings.Contains(err.Error(), "quotaExceeded") {
			logrus.Warnf("Quota exceeded for API key %d, trying next key", c.currentKey)
			return c.retryWithNextKey()
		}
		return nil, fmt.Errorf("error fetching videos: %v", err)
	}

	var videos []Video
	for _, item := range response.Items {
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			logrus.Warnf("Error parsing published date for video %s: %v", item.Id.VideoId, err)
			continue
		}

		video := Video{
			ID:           item.Id.VideoId,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishedAt:  publishedAt,
			ThumbnailURL: item.Snippet.Thumbnails.Default.Url,
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (c *Client) retryWithNextKey() ([]Video, error) {
	// Try the next API key
	c.currentKey = (c.currentKey + 1) % len(c.apiKeys)
	
	// If we've tried all keys, return an error
	if c.currentKey == 0 {
		return nil, fmt.Errorf("all API keys have exceeded their quota")
	}

	logrus.Infof("Switching to API key %d", c.currentKey)
	
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(c.apiKeys[c.currentKey]))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube service with new key: %v", err)
	}
	
	c.service = service
	return c.FetchLatestVideos()
}

