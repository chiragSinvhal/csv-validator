package services

import (
	"sync"
	"time"

	"csv-validator/internal/models"

	"github.com/google/uuid"
)

// JobService manages file processing jobs
type JobService struct {
	jobs map[string]*models.Job
	mu   sync.RWMutex
}

// NewJobService creates a new job service
func NewJobService() *JobService {
	return &JobService{
		jobs: make(map[string]*models.Job),
	}
}

// CreateJob creates a new job with pending status
func (js *JobService) CreateJob(originalFile string) *models.Job {
	js.mu.Lock()
	defer js.mu.Unlock()

	job := &models.Job{
		ID:           uuid.New().String(),
		Status:       models.JobStatusPending,
		OriginalFile: originalFile,
		CreatedAt:    time.Now(),
	}

	js.jobs[job.ID] = job
	return job
}

// GetJob retrieves a job by ID
func (js *JobService) GetJob(id string) (*models.Job, bool) {
	js.mu.RLock()
	defer js.mu.RUnlock()

	job, exists := js.jobs[id]
	if !exists {
		return nil, false
	}

	// Return a copy to avoid race conditions
	jobCopy := *job
	return &jobCopy, true
}

// UpdateJobStatus updates the status of a job
func (js *JobService) UpdateJobStatus(id string, status models.JobStatus) error {
	js.mu.Lock()
	defer js.mu.Unlock()

	job, exists := js.jobs[id]
	if !exists {
		return ErrJobNotFound
	}

	job.Status = status
	if status == models.JobStatusCompleted || status == models.JobStatusFailed {
		now := time.Now()
		job.CompletedAt = &now
	}

	return nil
}

// UpdateJobProcessedFile updates the processed file path for a job
func (js *JobService) UpdateJobProcessedFile(id string, processedFile string) error {
	js.mu.Lock()
	defer js.mu.Unlock()

	job, exists := js.jobs[id]
	if !exists {
		return ErrJobNotFound
	}

	job.ProcessedFile = processedFile
	return nil
}

// UpdateJobError updates the error message for a job
func (js *JobService) UpdateJobError(id string, errorMessage string) error {
	js.mu.Lock()
	defer js.mu.Unlock()

	job, exists := js.jobs[id]
	if !exists {
		return ErrJobNotFound
	}

	job.ErrorMessage = errorMessage
	job.Status = models.JobStatusFailed
	now := time.Now()
	job.CompletedAt = &now

	return nil
}

// ListJobs returns all jobs (for debugging/monitoring)
func (js *JobService) ListJobs() []*models.Job {
	js.mu.RLock()
	defer js.mu.RUnlock()

	jobs := make([]*models.Job, 0, len(js.jobs))
	for _, job := range js.jobs {
		jobCopy := *job
		jobs = append(jobs, &jobCopy)
	}

	return jobs
}

// CleanupOldJobs removes jobs older than the specified duration
func (js *JobService) CleanupOldJobs(maxAge time.Duration) int {
	js.mu.Lock()
	defer js.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	removed := 0

	for id, job := range js.jobs {
		if job.CreatedAt.Before(cutoff) {
			delete(js.jobs, id)
			removed++
		}
	}

	return removed
}
