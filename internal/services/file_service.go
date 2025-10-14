package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"csv-validator/internal/models"
)

// FileService handles file operations
type FileService struct {
	uploadDir   string
	downloadDir string
}

// NewFileService creates a new file service
func NewFileService(uploadDir, downloadDir string) *FileService {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	// Create download directory if it doesn't exist
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create download directory: %v", err))
	}

	return &FileService{
		uploadDir:   uploadDir,
		downloadDir: downloadDir,
	}
}

// SaveFile saves an uploaded file to the upload directory
func (fs *FileService) SaveFile(file *multipart.FileHeader, jobID string) (string, error) {
	// Create a unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d_%s", jobID, timestamp, file.Filename)
	filepath := filepath.Join(fs.uploadDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return filepath, nil
}

// GetFile returns the file path if it exists
func (fs *FileService) GetFile(filename string) (string, error) {
	filepath := filepath.Join(fs.uploadDir, filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", filename)
	}

	return filepath, nil
}

// DeleteFile removes a file from the filesystem
func (fs *FileService) DeleteFile(filename string) error {
	filepath := filepath.Join(fs.uploadDir, filename)
	return os.Remove(filepath)
}

// ValidateFile validates the uploaded file
func (fs *FileService) ValidateFile(file *multipart.FileHeader, maxSize int64) (*models.FileInfo, error) {
	// Check file size
	if file.Size > maxSize {
		return nil, fmt.Errorf("file size (%d bytes) exceeds maximum allowed size (%d bytes)", file.Size, maxSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".csv" {
		return nil, fmt.Errorf("invalid file format. Only CSV files are allowed")
	}

	// Get MIME type by reading file header
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file for validation: %w", err)
	}
	defer src.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read file header: %w", err)
	}

	// Reset file pointer
	src.Close()

	// Additional validation: check if it's a text file
	mimeType := "text/csv"
	if !isTextFile(buffer[:n]) {
		return nil, fmt.Errorf("file does not appear to be a valid text/CSV file")
	}

	return &models.FileInfo{
		Filename: file.Filename,
		Size:     file.Size,
		MimeType: mimeType,
	}, nil
}

// isTextFile checks if the file content appears to be text
func isTextFile(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// Check for null bytes (common in binary files)
	for _, b := range data {
		if b == 0 {
			return false
		}
	}

	// Check if it contains mostly printable ASCII characters
	printableCount := 0
	for _, b := range data {
		if (b >= 32 && b <= 126) || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}

	// At least 95% should be printable
	ratio := float64(printableCount) / float64(len(data))
	return ratio >= 0.95
}

// GetDownloadDir returns the download directory path
func (fs *FileService) GetDownloadDir() string {
	return fs.downloadDir
}
