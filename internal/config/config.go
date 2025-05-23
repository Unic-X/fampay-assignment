package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config holds all configuration for the application
type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	YouTubeAPIKeys []string
	SearchQuery    string
	ServerPort     string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Error loading .env file, using environment variables")
	}

	config := &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "youtube_videos"),
		YouTubeAPIKeys: strings.Split(getEnv("YOUTUBE_API_KEYS", ""), ","),
		SearchQuery:    getEnv("SEARCH_QUERY", "cricket"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
	}
	fmt.Println(config);
	return config, nil
}

// getEnv gets environment variable with fallback
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDatabaseURL returns database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
