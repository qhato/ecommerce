package domain

import "time"

// JobType represents the type of import/export job
type JobType string

const (
	JobTypeImport JobType = "IMPORT"
	JobTypeExport JobType = "EXPORT"
)

// JobStatus represents the status of a job
type JobStatus string

const (
	JobStatusPending    JobStatus = "PENDING"
	JobStatusProcessing JobStatus = "PROCESSING"
	JobStatusCompleted  JobStatus = "COMPLETED"
	JobStatusFailed     JobStatus = "FAILED"
	JobStatusCancelled  JobStatus = "CANCELLED"
)

// EntityType represents the entity being imported/exported
type EntityType string

const (
	EntityTypeProduct  EntityType = "PRODUCT"
	EntityTypeCategory EntityType = "CATEGORY"
	EntityTypeCustomer EntityType = "CUSTOMER"
	EntityTypeOrder    EntityType = "ORDER"
	EntityTypeContent  EntityType = "CONTENT"
)

// FileFormat represents the file format
type FileFormat string

const (
	FileFormatCSV  FileFormat = "CSV"
	FileFormatJSON FileFormat = "JSON"
	FileFormatXML  FileFormat = "XML"
)

// ImportExportJob represents an import or export job
type ImportExportJob struct {
	ID             int64
	Type           JobType
	EntityType     EntityType
	Format         FileFormat
	Status         JobStatus
	FilePath       string
	FileName       string
	TotalRecords   int
	ProcessedRecords int
	SuccessRecords int
	FailedRecords  int
	ErrorLog       string
	StartedAt      *time.Time
	CompletedAt    *time.Time
	CreatedBy      int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewImportJob creates a new import job
func NewImportJob(entityType EntityType, format FileFormat, filePath string, createdBy int64) (*ImportExportJob, error) {
	if filePath == "" {
		return nil, ErrFilePathRequired
	}

	now := time.Now()
	return &ImportExportJob{
		Type:       JobTypeImport,
		EntityType: entityType,
		Format:     format,
		Status:     JobStatusPending,
		FilePath:   filePath,
		CreatedBy:  createdBy,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// NewExportJob creates a new export job
func NewExportJob(entityType EntityType, format FileFormat, createdBy int64) (*ImportExportJob, error) {
	now := time.Now()
	return &ImportExportJob{
		Type:       JobTypeExport,
		EntityType: entityType,
		Format:     format,
		Status:     JobStatusPending,
		CreatedBy:  createdBy,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// Start starts the job
func (j *ImportExportJob) Start() {
	now := time.Now()
	j.Status = JobStatusProcessing
	j.StartedAt = &now
	j.UpdatedAt = now
}

// Complete completes the job
func (j *ImportExportJob) Complete() {
	now := time.Now()
	j.Status = JobStatusCompleted
	j.CompletedAt = &now
	j.UpdatedAt = now
}

// Fail marks the job as failed
func (j *ImportExportJob) Fail(errorLog string) {
	now := time.Now()
	j.Status = JobStatusFailed
	j.ErrorLog = errorLog
	j.CompletedAt = &now
	j.UpdatedAt = now
}

// UpdateProgress updates the job progress
func (j *ImportExportJob) UpdateProgress(processed, success, failed int) {
	j.ProcessedRecords = processed
	j.SuccessRecords = success
	j.FailedRecords = failed
	j.UpdatedAt = time.Now()
}

// GetProgress returns the progress percentage
func (j *ImportExportJob) GetProgress() int {
	if j.TotalRecords == 0 {
		return 0
	}
	return (j.ProcessedRecords * 100) / j.TotalRecords
}
