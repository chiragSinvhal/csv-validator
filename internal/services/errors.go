package services

import "errors"

// Common service errors
var (
	ErrJobNotFound     = errors.New("job not found")
	ErrJobInProgress   = errors.New("job is still in progress")
	ErrJobFailed       = errors.New("job processing failed")
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file size exceeds limit")
)
