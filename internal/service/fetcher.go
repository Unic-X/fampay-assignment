package service

import (
	"time"

	"github.com/Unic-X/fampay-assignment/internal/database"
	"github.com/Unic-X/fampay-assignment/internal/youtube"
	"github.com/sirupsen/logrus"
)

type Fetcher struct {
	youtubeClient *youtube.Client
	db            *database.DB
	interval      time.Duration
	stopChan      chan struct{}
}

func NewFetcher(youtubeClient *youtube.Client, db *database.DB, interval time.Duration) *Fetcher {
	return &Fetcher{
		youtubeClient: youtubeClient,
		db:            db,
		interval:      interval,
		stopChan:      make(chan struct{}),
	}
}

func (f *Fetcher) Start() {
	ticker := time.NewTicker(f.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := f.fetchAndStore(); err != nil {
				logrus.Errorf("Error fetching videos: %v", err)
			}
		case <-f.stopChan:
			return
		}
	}
}

func (f *Fetcher) Stop() {
	close(f.stopChan)
}

func (f *Fetcher) fetchAndStore() error {
	videos, err := f.youtubeClient.FetchLatestVideos()
	if err != nil {
		return err
	}

	for _, video := range videos {
		dbVideo := &database.Video{
			ID:           video.ID,
			Title:        video.Title,
			Description:  video.Description,
			PublishedAt:  video.PublishedAt,
			ThumbnailURL: video.ThumbnailURL,
		}
		if err := f.db.SaveVideo(dbVideo); err != nil {
			logrus.Errorf("Error saving video %s: %v", video.ID, err)
			continue
		}
	}

	return nil
}
