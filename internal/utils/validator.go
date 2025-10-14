package utils

import (
	"mime/multipart"
	"strings"
)

// ValidateCSVFile validates if the uploaded file is a valid CSV
func ValidateCSVFile(fileHeader *multipart.FileHeader) error {
	// Check file extension
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".csv") {
		return ErrInvalidFileExtension
	}

	// Check file size (basic check here, detailed check in service layer)
	if fileHeader.Size == 0 {
		return ErrEmptyFile
	}

	return nil
}

// SanitizeFilename removes potentially dangerous characters from filename
func SanitizeFilename(filename string) string {
	// Replace dangerous characters with underscores
	dangerous := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	sanitized := filename

	// First replace .. with underscores
	sanitized = strings.ReplaceAll(sanitized, "..", "_")

	for _, char := range dangerous {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}

	// Ensure filename is not empty after sanitization
	trimmed := strings.TrimSpace(sanitized)
	if trimmed == "" || trimmed == "_" || strings.Trim(trimmed, "_") == "" {
		sanitized = "file.csv"
	}

	return sanitized
}

// IsValidJobID checks if a job ID has the correct format (UUID)
func IsValidJobID(jobID string) bool {
	if len(jobID) != 36 {
		return false
	}

	// Basic UUID format check: 8-4-4-4-12 characters separated by hyphens
	parts := strings.Split(jobID, "-")
	if len(parts) != 5 {
		return false
	}

	expectedLengths := []int{8, 4, 4, 4, 12}
	for i, part := range parts {
		if len(part) != expectedLengths[i] {
			return false
		}
	}

	return true
}
