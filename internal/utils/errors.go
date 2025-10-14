package utils

import "errors"

// Validation errors
var (
	ErrInvalidFileExtension = errors.New("invalid file extension, only .csv files are allowed")
	ErrEmptyFile            = errors.New("file is empty")
	ErrInvalidJobID         = errors.New("invalid job ID format")
)
