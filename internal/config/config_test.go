package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clear any existing env vars
	os.Clearenv()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check defaults
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "./uploads", cfg.UploadDir)
	assert.Equal(t, "./downloads", cfg.DownloadDir)
	assert.Equal(t, int64(10*1024*1024), cfg.MaxFileSize)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "release", cfg.GinMode)
	assert.Equal(t, "*", cfg.AllowedOrigins)
}

func TestLoad_CustomValues(t *testing.T) {
	os.Clearenv()

	os.Setenv("PORT", "9000")
	os.Setenv("UPLOAD_DIR", "/tmp/uploads")
	os.Setenv("DOWNLOAD_DIR", "/tmp/downloads")
	os.Setenv("MAX_FILE_SIZE", "5242880") // 5MB
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("GIN_MODE", "debug")
	os.Setenv("ALLOWED_ORIGINS", "https://example.com")

	cfg, err := Load()
	assert.NoError(t, err)

	assert.Equal(t, "9000", cfg.Port)
	assert.Equal(t, "/tmp/uploads", cfg.UploadDir)
	assert.Equal(t, "/tmp/downloads", cfg.DownloadDir)
	assert.Equal(t, int64(5242880), cfg.MaxFileSize)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "debug", cfg.GinMode)
	assert.Equal(t, "https://example.com", cfg.AllowedOrigins)

	os.Clearenv()
}

func TestGetEnv(t *testing.T) {
	os.Clearenv()

	// Test fallback
	val := getEnv("NONEXISTENT", "fallback")
	assert.Equal(t, "fallback", val)

	// Test actual value
	os.Setenv("TEST_VAR", "actual")
	val = getEnv("TEST_VAR", "fallback")
	assert.Equal(t, "actual", val)

	os.Clearenv()
}

func TestGetEnvAsInt64(t *testing.T) {
	os.Clearenv()

	// Test fallback
	val := getEnvAsInt64("NONEXISTENT", 999)
	assert.Equal(t, int64(999), val)

	// Test valid int
	os.Setenv("TEST_INT", "12345")
	val = getEnvAsInt64("TEST_INT", 999)
	assert.Equal(t, int64(12345), val)

	// Test invalid int (should use fallback)
	os.Setenv("TEST_INT", "not-a-number")
	val = getEnvAsInt64("TEST_INT", 999)
	assert.Equal(t, int64(999), val)

	os.Clearenv()
}

func TestLoad_PartialEnvVars(t *testing.T) {
	os.Clearenv()

	// Set only some vars
	os.Setenv("PORT", "3000")
	os.Setenv("LOG_LEVEL", "warn")

	cfg, err := Load()
	assert.NoError(t, err)

	// Custom values
	assert.Equal(t, "3000", cfg.Port)
	assert.Equal(t, "warn", cfg.LogLevel)

	// Defaults for the rest
	assert.Equal(t, "./uploads", cfg.UploadDir)
	assert.Equal(t, "release", cfg.GinMode)

	os.Clearenv()
}
