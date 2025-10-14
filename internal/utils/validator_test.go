package utils

import (
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCSVFile(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		size        int64
		expectError bool
	}{
		{
			name:        "Valid CSV file",
			filename:    "test.csv",
			size:        1024,
			expectError: false,
		},
		{
			name:        "Valid CSV file uppercase",
			filename:    "test.CSV",
			size:        1024,
			expectError: false,
		},
		{
			name:        "Invalid file extension - txt",
			filename:    "test.txt",
			size:        1024,
			expectError: true,
		},
		{
			name:        "Invalid file extension - xlsx",
			filename:    "test.xlsx",
			size:        1024,
			expectError: true,
		},
		{
			name:        "Empty file",
			filename:    "test.csv",
			size:        0,
			expectError: true,
		},
		{
			name:        "No extension",
			filename:    "test",
			size:        1024,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileHeader := &multipart.FileHeader{
				Filename: tt.filename,
				Size:     tt.size,
			}

			err := ValidateCSVFile(fileHeader)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal filename",
			input:    "test.csv",
			expected: "test.csv",
		},
		{
			name:     "Filename with dangerous characters",
			input:    "test/file\\name:*.csv",
			expected: "test_file_name__.csv",
		},
		{
			name:     "Filename with path traversal",
			input:    "../../../etc/passwd",
			expected: "______etc_passwd",
		},
		{
			name:     "Empty filename",
			input:    "",
			expected: "file.csv",
		},
		{
			name:     "Filename with only dangerous characters",
			input:    "/\\:*?\"<>|",
			expected: "file.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFilename(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidJobID(t *testing.T) {
	tests := []struct {
		name     string
		jobID    string
		expected bool
	}{
		{
			name:     "Valid UUID",
			jobID:    "a225eb00-0907-4273-92ca-5faadeefae5f",
			expected: true,
		},
		{
			name:     "Valid UUID - different format",
			jobID:    "123e4567-e89b-12d3-a456-426614174000",
			expected: true,
		},
		{
			name:     "Invalid UUID - too short",
			jobID:    "a225eb00-0907-4273-92ca-5faadeefa",
			expected: false,
		},
		{
			name:     "Invalid UUID - too long",
			jobID:    "a225eb00-0907-4273-92ca-5faadeefae5ff",
			expected: false,
		},
		{
			name:     "Invalid UUID - wrong format",
			jobID:    "a225eb000907427392ca5faadeefae5f",
			expected: false,
		},
		{
			name:     "Invalid UUID - missing hyphens",
			jobID:    "a225eb00090742739ca5faadeefae5f",
			expected: false,
		},
		{
			name:     "Empty string",
			jobID:    "",
			expected: false,
		},
		{
			name:     "Invalid UUID - wrong segment lengths",
			jobID:    "a225eb0-0907-4273-92ca-5faadeefae5f",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidJobID(tt.jobID)
			assert.Equal(t, tt.expected, result)
		})
	}
}
