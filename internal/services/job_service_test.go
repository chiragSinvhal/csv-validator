package services

import (
	"testing"
	"time"

	"csv-validator/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestJobService_CreateJob(t *testing.T) {
	js := NewJobService()

	originalFile := "test.csv"
	job := js.CreateJob(originalFile)

	assert.NotEmpty(t, job.ID)
	assert.Equal(t, originalFile, job.OriginalFile)
	assert.Equal(t, models.JobStatusPending, job.Status)
	assert.False(t, job.CreatedAt.IsZero())
	assert.Nil(t, job.CompletedAt)
}

func TestJobService_GetJob(t *testing.T) {
	js := NewJobService()

	// Test getting non-existent job
	job, exists := js.GetJob("non-existent")
	assert.False(t, exists)
	assert.Nil(t, job)

	// Create a job
	originalFile := "test.csv"
	createdJob := js.CreateJob(originalFile)

	// Test getting existing job
	job, exists = js.GetJob(createdJob.ID)
	assert.True(t, exists)
	assert.NotNil(t, job)
	assert.Equal(t, createdJob.ID, job.ID)
	assert.Equal(t, originalFile, job.OriginalFile)
}

func TestJobService_UpdateJobStatus(t *testing.T) {
	js := NewJobService()

	// Create a job
	job := js.CreateJob("test.csv")

	// Update to processing
	err := js.UpdateJobStatus(job.ID, models.JobStatusProcessing)
	assert.NoError(t, err)

	// Verify update
	updatedJob, exists := js.GetJob(job.ID)
	assert.True(t, exists)
	assert.Equal(t, models.JobStatusProcessing, updatedJob.Status)
	assert.Nil(t, updatedJob.CompletedAt)

	// Update to completed
	err = js.UpdateJobStatus(job.ID, models.JobStatusCompleted)
	assert.NoError(t, err)

	// Verify update with completion time
	updatedJob, exists = js.GetJob(job.ID)
	assert.True(t, exists)
	assert.Equal(t, models.JobStatusCompleted, updatedJob.Status)
	assert.NotNil(t, updatedJob.CompletedAt)

	// Test updating non-existent job
	err = js.UpdateJobStatus("non-existent", models.JobStatusCompleted)
	assert.Error(t, err)
	assert.Equal(t, ErrJobNotFound, err)
}

func TestJobService_UpdateJobProcessedFile(t *testing.T) {
	js := NewJobService()

	// Create a job
	job := js.CreateJob("test.csv")

	processedFile := "/path/to/processed.csv"
	err := js.UpdateJobProcessedFile(job.ID, processedFile)
	assert.NoError(t, err)

	// Verify update
	updatedJob, exists := js.GetJob(job.ID)
	assert.True(t, exists)
	assert.Equal(t, processedFile, updatedJob.ProcessedFile)

	// Test updating non-existent job
	err = js.UpdateJobProcessedFile("non-existent", processedFile)
	assert.Error(t, err)
	assert.Equal(t, ErrJobNotFound, err)
}

func TestJobService_UpdateJobError(t *testing.T) {
	js := NewJobService()

	// Create a job
	job := js.CreateJob("test.csv")

	errorMessage := "Processing failed"
	err := js.UpdateJobError(job.ID, errorMessage)
	assert.NoError(t, err)

	// Verify update
	updatedJob, exists := js.GetJob(job.ID)
	assert.True(t, exists)
	assert.Equal(t, errorMessage, updatedJob.ErrorMessage)
	assert.Equal(t, models.JobStatusFailed, updatedJob.Status)
	assert.NotNil(t, updatedJob.CompletedAt)

	// Test updating non-existent job
	err = js.UpdateJobError("non-existent", errorMessage)
	assert.Error(t, err)
	assert.Equal(t, ErrJobNotFound, err)
}

func TestJobService_ListJobs(t *testing.T) {
	js := NewJobService()

	// Initially empty
	jobs := js.ListJobs()
	assert.Empty(t, jobs)

	// Create some jobs
	job1 := js.CreateJob("test1.csv")
	job2 := js.CreateJob("test2.csv")

	// List jobs
	jobs = js.ListJobs()
	assert.Len(t, jobs, 2)

	// Find our jobs
	found1, found2 := false, false
	for _, job := range jobs {
		if job.ID == job1.ID {
			found1 = true
			assert.Equal(t, "test1.csv", job.OriginalFile)
		}
		if job.ID == job2.ID {
			found2 = true
			assert.Equal(t, "test2.csv", job.OriginalFile)
		}
	}
	assert.True(t, found1)
	assert.True(t, found2)
}

func TestJobService_CleanupOldJobs(t *testing.T) {
	js := NewJobService()

	// Create jobs with different timestamps
	job1 := js.CreateJob("old.csv")
	job2 := js.CreateJob("new.csv")

	// Manually set old timestamp
	js.mu.Lock()
	js.jobs[job1.ID].CreatedAt = time.Now().Add(-2 * time.Hour)
	js.mu.Unlock()

	// Cleanup jobs older than 1 hour
	removed := js.CleanupOldJobs(time.Hour)
	assert.Equal(t, 1, removed)

	// Verify only new job remains
	jobs := js.ListJobs()
	assert.Len(t, jobs, 1)
	assert.Equal(t, job2.ID, jobs[0].ID)
}
