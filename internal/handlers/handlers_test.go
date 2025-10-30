package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

	csvData := "name,email\nJohn,Chirag@test.com"
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

func TestUploadInvalidFileType(t *testing.T) {
	handler, tempDir := setupHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	// Try uploading a .txt file
	req := createRequest(t, "test.txt", "not a csv")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UploadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid file type")
}

func TestUploadFileTooLarge(t *testing.T) {
	handler, tempDir := setupHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	// Create large content
	largeData := strings.Repeat("a", 2*1024*1024) // 2MB
	req := createRequest(t, "large.csv", largeData)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UploadFile(c)

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
}

func TestDownloadJobNotFound(t *testing.T) {
	handler, tempDir := setupHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "a225eb00-0907-4273-92ca-5faadeefae5f"}}

	handler.DownloadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Job not found")
}

func TestDownloadJobStillProcessing(t *testing.T) {
	handler, tempDir := setupHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	// Create a job manually
	job := handler.jobService.CreateJob("test.csv")
	handler.jobService.UpdateJobStatus(job.ID, models.JobStatusProcessing)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: job.ID}}

	handler.DownloadFile(c)

	assert.Equal(t, http.StatusLocked, w.Code)
	assert.Contains(t, w.Body.String(), "Still processing")
}

func TestDownloadJobFailed(t *testing.T) {
	handler, tempDir := setupHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	job := handler.jobService.CreateJob("test.csv")
	handler.jobService.UpdateJobError(job.ID, "something broke")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: job.ID}}

	handler.DownloadFile(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Processing failed")
}
