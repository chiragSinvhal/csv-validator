package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"csv-validator/internal/models"
	"csv-validator/internal/utils"
	"csv-validator/pkg/logger"
)

// CSVService handles CSV file processing
type CSVService struct {
	fileService *FileService
	jobService  *JobService
}

// NewCSVService creates a new CSV service
func NewCSVService(fileService *FileService, jobService *JobService) *CSVService {
	return &CSVService{
		fileService: fileService,
		jobService:  jobService,
	}
}

// ProcessFile processes a CSV file asynchronously
func (cs *CSVService) ProcessFile(jobID string) {
	// Start processing in a goroutine
	go func() {
		if err := cs.processFileSync(jobID); err != nil {
			logger.Error(fmt.Sprintf("Failed to process file for job %s: %v", jobID, err))
			cs.jobService.UpdateJobError(jobID, err.Error())
		}
	}()
}

// processFileSync processes a CSV file synchronously
func (cs *CSVService) processFileSync(jobID string) error {
	// Update job status to processing
	if err := cs.jobService.UpdateJobStatus(jobID, models.JobStatusProcessing); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	// Get job details
	job, exists := cs.jobService.GetJob(jobID)
	if !exists {
		return ErrJobNotFound
	}

	// Open original file
	file, err := os.Open(job.OriginalFile)
	if err != nil {
		return fmt.Errorf("failed to open original file: %w", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("CSV file is empty")
	}

	// Process records
	processedRecords := cs.processRecords(records)

	// Create processed file path
	processedFileName := fmt.Sprintf("processed_%s", filepath.Base(job.OriginalFile))
	processedFilePath := filepath.Join(cs.fileService.GetDownloadDir(), processedFileName)

	// Write processed CSV
	if err := cs.writeProcessedCSV(processedFilePath, processedRecords); err != nil {
		return fmt.Errorf("failed to write processed CSV: %w", err)
	}

	// Update job with processed file path
	if err := cs.jobService.UpdateJobProcessedFile(jobID, processedFilePath); err != nil {
		return fmt.Errorf("failed to update job processed file: %w", err)
	}

	// Mark job as completed
	if err := cs.jobService.UpdateJobStatus(jobID, models.JobStatusCompleted); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	logger.Info(fmt.Sprintf("Successfully processed file for job %s", jobID))
	return nil
}

// processRecords adds email validation flag to CSV records
func (cs *CSVService) processRecords(records [][]string) [][]string {
	if len(records) == 0 {
		return records
	}

	processedRecords := make([][]string, len(records))

	// Process header row
	if len(records) > 0 {
		header := make([]string, len(records[0])+1)
		copy(header, records[0])
		header[len(header)-1] = "has_email"
		processedRecords[0] = header
	}

	// Process data rows
	for i := 1; i < len(records); i++ {
		record := records[i]

		// Skip empty rows
		if cs.isEmptyRow(record) {
			processedRecords[i] = record
			continue
		}

		// Create new record with email flag
		newRecord := make([]string, len(record)+1)
		copy(newRecord, record)

		// Check if any field contains a valid email
		hasEmail := false
		for _, field := range record {
			if utils.IsValidEmail(strings.TrimSpace(field)) {
				hasEmail = true
				break
			}
		}

		// Add email flag
		if hasEmail {
			newRecord[len(newRecord)-1] = "true"
		} else {
			newRecord[len(newRecord)-1] = "false"
		}

		processedRecords[i] = newRecord
	}

	return processedRecords
}

// isEmptyRow checks if a CSV row is empty or contains only whitespace
func (cs *CSVService) isEmptyRow(record []string) bool {
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}
	return true
}

// writeProcessedCSV writes processed records to a CSV file
func (cs *CSVService) writeProcessedCSV(filePath string, records [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create processed file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}
