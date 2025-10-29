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

type Handler struct {
	csvService  *services.CSVService
	jobService  *services.JobService
	fileService *services.FileService
	config      *config.Config
}

func NewHandler(csvService *services.CSVService, jobService *services.JobService, fileService *services.FileService, config *config.Config) *Handler {
	return &Handler{
		csvService:  csvService,
		jobService:  jobService,
		fileService: fileService,
		config:      config,
	}
}

// UploadFile handles file uploads
func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "No file provided",
		})
		return
	}

	// Basic checks
	if file.Size > h.config.MaxFileSize {
		c.JSON(http.StatusRequestEntityTooLarge, models.ErrorResponse{
			Error: "File too big",
		})
		return
	}

	if err := utils.ValidateCSVFile(file); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid file type",
		})
		return
	}

	job := h.jobService.CreateJob(file.Filename)
	filePath, err := h.fileService.SaveFile(file, job.ID)
	if err != nil {
		h.jobService.UpdateJobError(job.ID, "Save failed")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Could not save file",
		})
		return
	}

	job.OriginalFile = filePath
	h.csvService.ProcessFile(job.ID)

	c.JSON(http.StatusOK, models.UploadResponse{ID: job.ID})
}

// DownloadFile handles file downloads
func (h *Handler) DownloadFile(c *gin.Context) {
	jobID := c.Param("id")
	logger.Info(fmt.Sprintf("Download request for %s", jobID))

	if !utils.IsValidJobID(jobID) {
		logger.Error(fmt.Sprintf("Bad job ID: %s", jobID))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid job ID",
		})
		return
	}

	job, exists := h.jobService.GetJob(jobID)
	if !exists {
		logger.Error(fmt.Sprintf("Job not found: %s", jobID))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Job not found",
		})
		return
	}

	switch job.Status {
	case models.JobStatusPending, models.JobStatusProcessing:
		logger.Info(fmt.Sprintf("Job %s still processing", jobID))
		c.JSON(http.StatusLocked, models.ErrorResponse{
			Error: "Still processing",
		})
		return

	case models.JobStatusFailed:
		logger.Error(fmt.Sprintf("Job %s failed: %s", jobID, job.ErrorMessage))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Processing failed",
		})
		return

	case models.JobStatusCompleted:
		if job.ProcessedFile == "" {
			logger.Error(fmt.Sprintf("No processed file for job %s", jobID))
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "No processed file",
			})
			return
		}

		// Clean up filename for download
		originalName := filepath.Base(job.OriginalFile)
		downloadName := fmt.Sprintf("processed_%s", originalName)

		// Remove job ID prefix if present
		if strings.Contains(downloadName, "_") {
			parts := strings.Split(downloadName, "_")
			if len(parts) >= 3 {
				downloadName = "processed_" + strings.Join(parts[2:], "_")
			}
		}

		logger.Info(fmt.Sprintf("Serving %s for job %s", job.ProcessedFile, jobID))

		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", downloadName))
		c.File(job.ProcessedFile)
		return

	default:
		logger.Error(fmt.Sprintf("Unknown status %s for job %s", job.Status, jobID))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Unknown status",
		})
		return
	}
}
