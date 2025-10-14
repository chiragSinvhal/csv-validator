package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"csv-validator/internal/config"
	"csv-validator/internal/models"
	"csv-validator/internal/services"
	"csv-validator/internal/utils"
	"csv-validator/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Handler contains all HTTP handlers
type Handler struct {
	csvService  *services.CSVService
	jobService  *services.JobService
	fileService *services.FileService
	config      *config.Config
}

// NewHandler creates a new handler instance
func NewHandler(csvService *services.CSVService, jobService *services.JobService, fileService *services.FileService, config *config.Config) *Handler {
	return &Handler{
		csvService:  csvService,
		jobService:  jobService,
		fileService: fileService,
		config:      config,
	}
}

// UploadFile handles CSV file uploads
func (h *Handler) UploadFile(c *gin.Context) {
	logger.Info("Received file upload request")

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get file from form: %v", err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "No file provided or invalid form data",
		})
		return
	}

	logger.Info(fmt.Sprintf("Processing file: %s (size: %d bytes)", file.Filename, file.Size))

	// Validate file
	_, err = h.fileService.ValidateFile(file, h.config.MaxFileSize)
	if err != nil {
		logger.Error(fmt.Sprintf("File validation failed: %v", err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Additional validation using utils
	if err := utils.ValidateCSVFile(file); err != nil {
		logger.Error(fmt.Sprintf("CSV validation failed: %v", err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Check file size against configured limit
	if file.Size > h.config.MaxFileSize {
		logger.Error(fmt.Sprintf("File size %d exceeds limit %d", file.Size, h.config.MaxFileSize))
		c.JSON(http.StatusRequestEntityTooLarge, models.ErrorResponse{
			Error: fmt.Sprintf("File size (%d bytes) exceeds maximum allowed size (%d bytes)", file.Size, h.config.MaxFileSize),
		})
		return
	}

	// Create job
	job := h.jobService.CreateJob(file.Filename)
	logger.Info(fmt.Sprintf("Created job %s for file %s", job.ID, file.Filename))

	// Save file
	savedPath, err := h.fileService.SaveFile(file, job.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to save file: %v", err))
		h.jobService.UpdateJobError(job.ID, fmt.Sprintf("Failed to save file: %v", err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to save uploaded file",
		})
		return
	}

	// Update job with saved file path
	if err := h.jobService.UpdateJobProcessedFile(job.ID, savedPath); err != nil {
		logger.Error(fmt.Sprintf("Failed to update job with file path: %v", err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update job information",
		})
		return
	}

	// Update job original file path
	job.OriginalFile = savedPath

	// Start async processing
	h.csvService.ProcessFile(job.ID)

	logger.Info(fmt.Sprintf("Successfully initiated processing for job %s", job.ID))

	// Return job ID
	c.JSON(http.StatusOK, models.UploadResponse{
		ID: job.ID,
	})
}

// DownloadFile handles file download requests
func (h *Handler) DownloadFile(c *gin.Context) {
	jobID := c.Param("id")
	logger.Info(fmt.Sprintf("Received download request for job %s", jobID))

	// Validate job ID format
	if !utils.IsValidJobID(jobID) {
		logger.Error(fmt.Sprintf("Invalid job ID format: %s", jobID))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid job ID format",
		})
		return
	}

	// Get job
	job, exists := h.jobService.GetJob(jobID)
	if !exists {
		logger.Error(fmt.Sprintf("Job not found: %s", jobID))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid job ID",
		})
		return
	}

	logger.Info(fmt.Sprintf("Job %s status: %s", jobID, job.Status))

	// Check job status
	switch job.Status {
	case models.JobStatusPending, models.JobStatusProcessing:
		// Job is still in progress - return 423 (Locked)
		logger.Info(fmt.Sprintf("Job %s is still in progress", jobID))
		c.JSON(http.StatusLocked, models.ErrorResponse{
			Error: "Job is still in progress",
		})
		return

	case models.JobStatusFailed:
		// Job failed
		logger.Error(fmt.Sprintf("Job %s failed: %s", jobID, job.ErrorMessage))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: fmt.Sprintf("Job processing failed: %s", job.ErrorMessage),
		})
		return

	case models.JobStatusCompleted:
		// Job completed successfully
		if job.ProcessedFile == "" {
			logger.Error(fmt.Sprintf("Job %s completed but no processed file found", jobID))
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Processed file not found",
			})
			return
		}

		// Check if processed file exists
		_, err := h.fileService.GetFile(filepath.Base(job.ProcessedFile))
		if err != nil {
			logger.Error(fmt.Sprintf("Processed file not found for job %s: %v", jobID, err))
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Processed file not found",
			})
			return
		}

		// Determine filename for download
		originalFilename := filepath.Base(job.OriginalFile)
		downloadFilename := fmt.Sprintf("processed_%s", originalFilename)

		// Clean filename to remove any job ID prefix
		if strings.Contains(downloadFilename, "_") {
			parts := strings.Split(downloadFilename, "_")
			if len(parts) >= 3 {
				// Remove job ID and timestamp from filename
				downloadFilename = strings.Join(parts[2:], "_")
				downloadFilename = "processed_" + downloadFilename
			}
		}

		logger.Info(fmt.Sprintf("Serving file %s for job %s", job.ProcessedFile, jobID))

		// Set headers for file download
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", downloadFilename))
		c.Header("Content-Type", "text/csv")

		// Serve the file
		c.File(job.ProcessedFile)
		return

	default:
		// Unknown status
		logger.Error(fmt.Sprintf("Job %s has unknown status: %s", jobID, job.Status))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Unknown job status",
		})
		return
	}
}
