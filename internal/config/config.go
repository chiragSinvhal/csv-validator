package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Port           string
	UploadDir      string
	DownloadDir    string
	MaxFileSize    int64
	LogLevel       string
	GinMode        string
	AllowedOrigins string
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		UploadDir:      getEnv("UPLOAD_DIR", "./uploads"),
		DownloadDir:    getEnv("DOWNLOAD_DIR", "./downloads"),
		MaxFileSize:    getEnvAsInt64("MAX_FILE_SIZE", 10*1024*1024), // 10MB default
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		GinMode:        getEnv("GIN_MODE", "release"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),
	}

	return cfg, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt64 gets an environment variable as int64 with a fallback value
func getEnvAsInt64(key string, fallback int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return fallback
}
