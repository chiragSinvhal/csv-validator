package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"csv-validator/internal/config"
	"csv-validator/internal/models"
	"csv-validator/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupHandler(t *testing.T) (*Handler, string) {
	tempDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)

	cfg := &config.Config{
		Port:        "8080",
		UploadDir:   tempDir,
		MaxFileSize: 1024 * 1024, // 1MB
		LogLevel:    "info",
		GinMode:     "test",
	}

	fileService := services.NewFileService(cfg.UploadDir, cfg.UploadDir+"-downloads")
	jobService := services.NewJobService()
	csvService := services.NewCSVService(fileService, jobService)

	handler := NewHandler(csvService, jobService, fileService, cfg)
	return handler, tempDir
}

func createRequest(t *testing.T, filename, content string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	require.NoError(t, err)

	_, err = part.Write([]byte(content))
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/upload", body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestUpload(t *testing.T) {
	handler, tempDir := setupHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	csvData := "name,email\nJohn,john@test.com"
	req := createRequest(t, "test.csv", csvData)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UploadFile(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.UploadResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
}

func TestUploadNoFile(t *testing.T) {
	handler, tempDir := setupHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	req, _ := http.NewRequest("POST", "/api/upload", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UploadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDownloadInvalidID(t *testing.T) {
	handler, tempDir := setupHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	handler.DownloadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
