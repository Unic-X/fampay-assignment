package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	YouTube  YouTubeConfig
	Server   ServerConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// YouTubeConfig holds YouTube API configuration
type YouTubeConfig struct {
	APIKeys       []string
	SearchQuery   string
	FetchInterval time.Duration
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatal("Invalid DB_PORT value")
	}

	fetchInterval, err := time.ParseDuration(getEnv("FETCH_INTERVAL", "10s"))
	if err != nil {
		log.Fatal("Invalid FETCH_INTERVAL value")
	}

	apiKeysStr := getEnv("YOUTUBE_API_KEYS", "")
	if apiKeysStr == "" {
		log.Fatal("YOUTUBE_API_KEYS environment variable is required")
	}

	apiKeys := strings.Split(apiKeysStr, ",")
	for i, key := range apiKeys {
		apiKeys[i] = strings.TrimSpace(key)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "youtube_api"),
		},
		YouTube: YouTubeConfig{
			APIKeys:       apiKeys,
			SearchQuery:   getEnv("SEARCH_QUERY", "football"),
			FetchInterval: fetchInterval,
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// GetDatabaseURL returns database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.DBName)
}
