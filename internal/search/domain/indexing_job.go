package domain

import "time"

// IndexingJobStatus represents the status of an indexing job
type IndexingJobStatus string

const (
	IndexingJobStatusPending    IndexingJobStatus = "PENDING"
	IndexingJobStatusRunning    IndexingJobStatus = "RUNNING"
	IndexingJobStatusCompleted  IndexingJobStatus = "COMPLETED"
	IndexingJobStatusFailed     IndexingJobStatus = "FAILED"
	IndexingJobStatusCancelled  IndexingJobStatus = "CANCELLED"
)

// IndexingJobType represents the type of indexing job
type IndexingJobType string

const (
	IndexingJobTypeFull        IndexingJobType = "FULL"        // Reindexación completa
	IndexingJobTypeIncremental IndexingJobType = "INCREMENTAL" // Indexación incremental
	IndexingJobTypeSingle      IndexingJobType = "SINGLE"      // Indexación de un solo item
)

// IndexingJob represents an indexing job
// Business Logic: Gestiona trabajos de indexación (full reindex, incremental)
type IndexingJob struct {
	ID               int64
	Type             IndexingJobType
	Status           IndexingJobStatus
	EntityType       string // "product", "category", etc.
	TotalItems       int
	ProcessedItems   int
	FailedItems      int
	ErrorMessage     string
	StartedAt        *time.Time
	CompletedAt      *time.Time
	CreatedBy        int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// IndexingJobRepository defines the interface for indexing job persistence
type IndexingJobRepository interface {
	Create(job *IndexingJob) error
	Update(job *IndexingJob) error
	FindByID(id int64) (*IndexingJob, error)
	FindRecent(limit int) ([]*IndexingJob, error)
	FindRunning() ([]*IndexingJob, error)
}

// Start marks the job as running
func (j *IndexingJob) Start() {
	j.Status = IndexingJobStatusRunning
	now := time.Now()
	j.StartedAt = &now
	j.UpdatedAt = now
}

// Complete marks the job as completed
func (j *IndexingJob) Complete() {
	j.Status = IndexingJobStatusCompleted
	now := time.Now()
	j.CompletedAt = &now
	j.UpdatedAt = now
}

// Fail marks the job as failed
func (j *IndexingJob) Fail(errorMessage string) {
	j.Status = IndexingJobStatusFailed
	j.ErrorMessage = errorMessage
	now := time.Now()
	j.CompletedAt = &now
	j.UpdatedAt = now
}

// Cancel marks the job as cancelled
func (j *IndexingJob) Cancel() {
	j.Status = IndexingJobStatusCancelled
	now := time.Now()
	j.CompletedAt = &now
	j.UpdatedAt = now
}

// IncrementProcessed increments the processed items counter
func (j *IndexingJob) IncrementProcessed() {
	j.ProcessedItems++
	j.UpdatedAt = time.Now()
}

// IncrementFailed increments the failed items counter
func (j *IndexingJob) IncrementFailed() {
	j.FailedItems++
	j.ProcessedItems++
	j.UpdatedAt = time.Now()
}

// GetProgress returns the progress percentage
func (j *IndexingJob) GetProgress() float64 {
	if j.TotalItems == 0 {
		return 0
	}
	return float64(j.ProcessedItems) / float64(j.TotalItems) * 100
}

// IsCompleted checks if the job is completed
func (j *IndexingJob) IsCompleted() bool {
	return j.Status == IndexingJobStatusCompleted ||
		j.Status == IndexingJobStatusFailed ||
		j.Status == IndexingJobStatusCancelled
}

// GetDuration returns the duration of the job
func (j *IndexingJob) GetDuration() time.Duration {
	if j.StartedAt == nil {
		return 0
	}
	if j.CompletedAt != nil {
		return j.CompletedAt.Sub(*j.StartedAt)
	}
	return time.Since(*j.StartedAt)
}
