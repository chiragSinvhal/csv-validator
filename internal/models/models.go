package models

import (
	"time"
)

// JobStatus represents the status of a file processing job
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

// Job represents a file processing job
type Job struct {
	ID            string     `json:"id"`
	Status        JobStatus  `json:"status"`
	OriginalFile  string     `json:"original_file"`
	ProcessedFile string     `json:"processed_file,omitempty"`
	ErrorMessage  string     `json:"error_message,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
}

// UploadResponse represents the response for file upload
type UploadResponse struct {
	ID string `json:"id"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// FileInfo contains information about uploaded file
type FileInfo struct {
	Filename string
	Size     int64
	MimeType string
}
