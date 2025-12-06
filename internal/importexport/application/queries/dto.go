package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/importexport/domain"
)

type JobDTO struct {
	ID               int64      `json:"id"`
	Type             string     `json:"type"`
	EntityType       string     `json:"entity_type"`
	Format           string     `json:"format"`
	Status           string     `json:"status"`
	FilePath         string     `json:"file_path,omitempty"`
	FileName         string     `json:"file_name,omitempty"`
	TotalRecords     int        `json:"total_records"`
	ProcessedRecords int        `json:"processed_records"`
	SuccessRecords   int        `json:"success_records"`
	FailedRecords    int        `json:"failed_records"`
	Progress         int        `json:"progress_percentage"`
	ErrorLog         string     `json:"error_log,omitempty"`
	StartedAt        *time.Time `json:"started_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	CreatedBy        int64      `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func ToJobDTO(job *domain.ImportExportJob) *JobDTO {
	return &JobDTO{
		ID:               job.ID,
		Type:             string(job.Type),
		EntityType:       string(job.EntityType),
		Format:           string(job.Format),
		Status:           string(job.Status),
		FilePath:         job.FilePath,
		FileName:         job.FileName,
		TotalRecords:     job.TotalRecords,
		ProcessedRecords: job.ProcessedRecords,
		SuccessRecords:   job.SuccessRecords,
		FailedRecords:    job.FailedRecords,
		Progress:         job.GetProgress(),
		ErrorLog:         job.ErrorLog,
		StartedAt:        job.StartedAt,
		CompletedAt:      job.CompletedAt,
		CreatedBy:        job.CreatedBy,
		CreatedAt:        job.CreatedAt,
		UpdatedAt:        job.UpdatedAt,
	}
}
