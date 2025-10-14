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

func setupTestHandler(t *testing.T) (*Handler, string) {
	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "csv-validator-test")
	require.NoError(t, err)

	// Create test config
	cfg := &config.Config{
		Port:        "8080",
		UploadDir:   tempDir,
		MaxFileSize: 10 * 1024 * 1024, // 10MB
		LogLevel:    "info",
		GinMode:     "test",
	}

	// Initialize services
	fileService := services.NewFileService(cfg.UploadDir, cfg.UploadDir+"-downloads")
	jobService := services.NewJobService()
	csvService := services.NewCSVService(fileService, jobService)

	// Create handler
	handler := NewHandler(csvService, jobService, fileService, cfg)

	return handler, tempDir
}

func createTestCSVFile(t *testing.T, content string) *os.File {
	tmpFile, err := os.CreateTemp("", "test*.csv")
	require.NoError(t, err)

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)

	_, err = tmpFile.Seek(0, 0)
	require.NoError(t, err)

	return tmpFile
}

func createMultipartRequest(t *testing.T, filename string, content string) (*http.Request, string) {
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
	return req, writer.FormDataContentType()
}

func TestHandler_UploadFile_Success(t *testing.T) {
	handler, tempDir := setupTestHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	csvContent := "name,email,age\nChirag,Chirag@example.com,30\nYash,Yash@test.com,25"
	req, _ := createMultipartRequest(t, "test.csv", csvContent)

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

func TestHandler_UploadFile_NoFile(t *testing.T) {
	handler, tempDir := setupTestHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	req, err := http.NewRequest("POST", "/api/upload", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UploadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "No file provided")
}

func TestHandler_UploadFile_InvalidFileType(t *testing.T) {
	handler, tempDir := setupTestHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	txtContent := "This is a text file, not CSV"
	req, _ := createMultipartRequest(t, "test.txt", txtContent)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UploadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Only CSV files are allowed")
}

func TestHandler_UploadFile_FileTooLarge(t *testing.T) {
	handler, tempDir := setupTestHandler(t)
	defer os.RemoveAll(tempDir)

	// Set a very small max file size for testing
	handler.config.MaxFileSize = 10

	gin.SetMode(gin.TestMode)

	csvContent := "name,email,age\nChirag,Chirag@example.com,30\nYash,Yash@test.com,25"
	req, _ := createMultipartRequest(t, "test.csv", csvContent)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UploadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "exceeds maximum allowed size")
}

func TestHandler_DownloadFile_InvalidJobID(t *testing.T) {
	handler, tempDir := setupTestHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/download/invalid-id", nil)
	c.Params = []gin.Param{{Key: "id", Value: "invalid-id"}}

	handler.DownloadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Invalid job ID format")
}

func TestHandler_DownloadFile_JobNotFound(t *testing.T) {
	handler, tempDir := setupTestHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	validUUID := "a225eb00-0907-4273-92ca-5faadeefae5f"
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/download/"+validUUID, nil)
	c.Params = []gin.Param{{Key: "id", Value: validUUID}}

	handler.DownloadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Invalid job ID")
}

func TestHandler_DownloadFile_JobInProgress(t *testing.T) {
	handler, tempDir := setupTestHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	// Create a job in progress
	job := handler.jobService.CreateJob("test.csv")
	handler.jobService.UpdateJobStatus(job.ID, models.JobStatusProcessing)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/download/"+job.ID, nil)
	c.Params = []gin.Param{{Key: "id", Value: job.ID}}

	handler.DownloadFile(c)

	assert.Equal(t, http.StatusLocked, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Job is still in progress")
}

func TestHandler_DownloadFile_JobFailed(t *testing.T) {
	handler, tempDir := setupTestHandler(t)
	defer os.RemoveAll(tempDir)

	gin.SetMode(gin.TestMode)

	// Create a failed job
	job := handler.jobService.CreateJob("test.csv")
	handler.jobService.UpdateJobError(job.ID, "Processing failed")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/download/"+job.ID, nil)
	c.Params = []gin.Param{{Key: "id", Value: job.ID}}

	handler.DownloadFile(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Job processing failed")
}
