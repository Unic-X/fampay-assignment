package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Unic-X/fampay-assignment/internal/config"
	"github.com/Unic-X/fampay-assignment/internal/database"
	"github.com/Unic-X/fampay-assignment/internal/handler"
	"github.com/Unic-X/fampay-assignment/internal/service"
	"github.com/Unic-X/fampay-assignment/internal/youtube"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database
	db, err := database.New(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Initialize YouTube client
	youtubeClient, err := youtube.New(cfg.YouTubeAPIKeys, cfg.SearchQuery)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	// Initialize video fetcher service
	fetcher := service.NewFetcher(youtubeClient, db, 10*time.Second)
	go fetcher.Start()
	defer fetcher.Stop()

	// Initialize router
	router := gin.Default()
	videoHandler := handler.NewVideoHandler(db)

	// Register routes
	router.GET("/api/videos", videoHandler.GetVideos)

	// Start server
	go func() {
		if err := router.Run(":" + cfg.ServerPort); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
} 