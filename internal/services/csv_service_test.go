package services

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"time"

	"csv-validator/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSVService_ProcessRecords(t *testing.T) {
	fileService := NewFileService("./test-uploads")
	jobService := NewJobService()
	csvService := NewCSVService(fileService, jobService)

	tests := []struct {
		name     string
		input    [][]string
		expected [][]string
	}{
		{
			name: "CSV with valid emails",
			input: [][]string{
				{"name", "email", "age"},
				{"Chirag", "Chirag@example.com", "30"},
				{"Yash", "Yash@test.com", "25"},
				{"Rohan", "not-an-email", "35"},
			},
			expected: [][]string{
				{"name", "email", "age", "has_email"},
				{"Chirag", "Chirag@example.com", "30", "true"},
				{"Yash", "Yash@test.com", "25", "true"},
				{"Rohan", "not-an-email", "35", "false"},
			},
		},
		{
			name: "CSV with no emails",
			input: [][]string{
				{"name", "phone", "age"},
				{"Chirag", "123-456-7890", "30"},
				{"Yash", "987-654-3210", "25"},
			},
			expected: [][]string{
				{"name", "phone", "age", "has_email"},
				{"Chirag", "123-456-7890", "30", "false"},
				{"Yash", "987-654-3210", "25", "false"},
			},
		},
		{
			name: "CSV with empty rows",
			input: [][]string{
				{"name", "email"},
				{"Chirag", "Chirag@example.com"},
				{"", ""},
				{"Yash", "Yash@test.com"},
			},
			expected: [][]string{
				{"name", "email", "has_email"},
				{"Chirag", "Chirag@example.com", "true"},
				{"", ""},
				{"Yash", "Yash@test.com", "true"},
			},
		},
		{
			name:     "Empty CSV",
			input:    [][]string{},
			expected: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := csvService.processRecords(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCSVService_IsEmptyRow(t *testing.T) {
	fileService := NewFileService("./test-uploads")
	jobService := NewJobService()
	csvService := NewCSVService(fileService, jobService)

	tests := []struct {
		name     string
		input    []string
		expected bool
	}{
		{
			name:     "Non-empty row",
			input:    []string{"Chirag", "Chirag@example.com", "30"},
			expected: false,
		},
		{
			name:     "Empty row",
			input:    []string{"", "", ""},
			expected: true,
		},
		{
			name:     "Row with whitespace",
			input:    []string{"  ", "  ", "  "},
			expected: true,
		},
		{
			name:     "Mixed empty row",
			input:    []string{"", "  ", ""},
			expected: true,
		},
		{
			name:     "Partially empty row",
			input:    []string{"Chirag", "", ""},
			expected: false,
		},
		{
			name:     "Empty slice",
			input:    []string{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := csvService.isEmptyRow(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCSVService_WriteProcessedCSV(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "csv-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fileService := NewFileService(tempDir)
	jobService := NewJobService()
	csvService := NewCSVService(fileService, jobService)

	records := [][]string{
		{"name", "email", "has_email"},
		{"Chirag", "Chirag@example.com", "true"},
		{"Yash", "Yash@test.com", "true"},
		{"Rohan", "not-an-email", "false"},
	}

	filePath := filepath.Join(tempDir, "test_output.csv")
	err = csvService.writeProcessedCSV(filePath, records)
	require.NoError(t, err)

	// Read back the file and verify content
	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	readRecords, err := reader.ReadAll()
	require.NoError(t, err)

	assert.Equal(t, records, readRecords)
}

func TestCSVService_ProcessFileSync_Success(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "csv-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test CSV file
	csvContent := "name,email,age\nChirag,Chirag@example.com,30\nYash,invalid-email,25\nRohan,Rohan@test.org,35"
	testFile := filepath.Join(tempDir, "test.csv")
	err = os.WriteFile(testFile, []byte(csvContent), 0644)
	require.NoError(t, err)

	fileService := NewFileService(tempDir)
	jobService := NewJobService()
	csvService := NewCSVService(fileService, jobService)

	// Create job
	job := jobService.CreateJob(testFile)

	// Process file
	err = csvService.processFileSync(job.ID)
	require.NoError(t, err)

	// Verify job status
	updatedJob, exists := jobService.GetJob(job.ID)
	require.True(t, exists)
	assert.Equal(t, models.JobStatusCompleted, updatedJob.Status)
	assert.NotEmpty(t, updatedJob.ProcessedFile)

	// Verify processed file exists and has correct content
	processedFile, err := os.Open(updatedJob.ProcessedFile)
	require.NoError(t, err)
	defer processedFile.Close()

	reader := csv.NewReader(processedFile)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Check header
	expectedHeader := []string{"name", "email", "age", "has_email"}
	assert.Equal(t, expectedHeader, records[0])

	// Check data rows
	assert.Equal(t, "true", records[1][3])  // Chirag has valid email
	assert.Equal(t, "false", records[2][3]) // Yash has invalid email
	assert.Equal(t, "true", records[3][3])  // Rohan has valid email
}

func TestCSVService_ProcessFileSync_FileNotFound(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "csv-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fileService := NewFileService(tempDir)
	jobService := NewJobService()
	csvService := NewCSVService(fileService, jobService)

	// Create job with non-existent file
	job := jobService.CreateJob("non-existent.csv")

	// Process file should fail
	err = csvService.processFileSync(job.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open original file")
}

func TestCSVService_ProcessFileSync_EmptyFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "csv-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create empty CSV file
	testFile := filepath.Join(tempDir, "empty.csv")
	err = os.WriteFile(testFile, []byte(""), 0644)
	require.NoError(t, err)

	fileService := NewFileService(tempDir)
	jobService := NewJobService()
	csvService := NewCSVService(fileService, jobService)

	// Create job
	job := jobService.CreateJob(testFile)

	// Process file should fail
	err = csvService.processFileSync(job.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CSV file is empty")
}

func TestCSVService_ProcessFile_Async(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "csv-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test CSV file
	csvContent := "name,email\nChirag,Chirag@example.com\nYash,Yash@test.com"
	testFile := filepath.Join(tempDir, "test.csv")
	err = os.WriteFile(testFile, []byte(csvContent), 0644)
	require.NoError(t, err)

	fileService := NewFileService(tempDir)
	jobService := NewJobService()
	csvService := NewCSVService(fileService, jobService)

	// Create job
	job := jobService.CreateJob(testFile)

	// Start async processing
	csvService.ProcessFile(job.ID)

	// Wait for processing to complete (with timeout)
	timeout := time.After(5 * time.Second)
	tick := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-timeout:
			t.Fatal("Processing timed out")
		case <-tick:
			updatedJob, exists := jobService.GetJob(job.ID)
			require.True(t, exists)

			if updatedJob.Status == models.JobStatusCompleted {
				assert.NotEmpty(t, updatedJob.ProcessedFile)
				return
			} else if updatedJob.Status == models.JobStatusFailed {
				t.Fatalf("Job failed: %s", updatedJob.ErrorMessage)
			}
		}
	}
}
