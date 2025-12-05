package domain

import "context"

// JobRepository defines the interface for job persistence
type JobRepository interface {
	Create(ctx context.Context, job *ImportExportJob) error
	Update(ctx context.Context, job *ImportExportJob) error
	FindByID(ctx context.Context, id int64) (*ImportExportJob, error)
	FindByType(ctx context.Context, jobType JobType, status JobStatus, limit int) ([]*ImportExportJob, error)
	FindByStatus(ctx context.Context, status JobStatus, limit int) ([]*ImportExportJob, error)
	FindRecent(ctx context.Context, limit int) ([]*ImportExportJob, error)
	Delete(ctx context.Context, id int64) error
}
