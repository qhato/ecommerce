package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/importexport/domain"
)

type PostgresJobRepository struct {
	db *sql.DB
}

func NewPostgresJobRepository(db *sql.DB) *PostgresJobRepository {
	return &PostgresJobRepository{db: db}
}

func (r *PostgresJobRepository) Create(ctx context.Context, job *domain.ImportExportJob) error {
	query := `INSERT INTO blc_import_export_job (
		type, entity_type, format, status, file_path, file_name,
		total_records, processed_records, success_records, failed_records,
		error_log, started_at, completed_at, created_by, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		job.Type, job.EntityType, job.Format, job.Status, job.FilePath, job.FileName,
		job.TotalRecords, job.ProcessedRecords, job.SuccessRecords, job.FailedRecords,
		job.ErrorLog, job.StartedAt, job.CompletedAt, job.CreatedBy, job.CreatedAt, job.UpdatedAt,
	).Scan(&job.ID)
}

func (r *PostgresJobRepository) Update(ctx context.Context, job *domain.ImportExportJob) error {
	query := `UPDATE blc_import_export_job SET
		status = $1, file_path = $2, file_name = $3,
		total_records = $4, processed_records = $5, success_records = $6, failed_records = $7,
		error_log = $8, started_at = $9, completed_at = $10, updated_at = $11
	WHERE id = $12`

	_, err := r.db.ExecContext(ctx, query,
		job.Status, job.FilePath, job.FileName,
		job.TotalRecords, job.ProcessedRecords, job.SuccessRecords, job.FailedRecords,
		job.ErrorLog, job.StartedAt, job.CompletedAt, job.UpdatedAt, job.ID,
	)
	return err
}

func (r *PostgresJobRepository) FindByID(ctx context.Context, id int64) (*domain.ImportExportJob, error) {
	query := `SELECT id, type, entity_type, format, status, file_path, file_name,
		total_records, processed_records, success_records, failed_records,
		error_log, started_at, completed_at, created_by, created_at, updated_at
	FROM blc_import_export_job WHERE id = $1`

	return r.scanJob(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresJobRepository) FindByType(ctx context.Context, jobType domain.JobType, status domain.JobStatus, limit int) ([]*domain.ImportExportJob, error) {
	query := `SELECT id, type, entity_type, format, status, file_path, file_name,
		total_records, processed_records, success_records, failed_records,
		error_log, started_at, completed_at, created_by, created_at, updated_at
	FROM blc_import_export_job WHERE type = $1 AND status = $2 ORDER BY created_at DESC LIMIT $3`

	return r.queryJobs(ctx, query, jobType, status, limit)
}

func (r *PostgresJobRepository) FindByStatus(ctx context.Context, status domain.JobStatus, limit int) ([]*domain.ImportExportJob, error) {
	query := `SELECT id, type, entity_type, format, status, file_path, file_name,
		total_records, processed_records, success_records, failed_records,
		error_log, started_at, completed_at, created_by, created_at, updated_at
	FROM blc_import_export_job WHERE status = $1 ORDER BY created_at DESC LIMIT $2`

	return r.queryJobs(ctx, query, status, limit)
}

func (r *PostgresJobRepository) FindRecent(ctx context.Context, limit int) ([]*domain.ImportExportJob, error) {
	query := `SELECT id, type, entity_type, format, status, file_path, file_name,
		total_records, processed_records, success_records, failed_records,
		error_log, started_at, completed_at, created_by, created_at, updated_at
	FROM blc_import_export_job ORDER BY created_at DESC LIMIT $1`

	return r.queryJobs(ctx, query, limit)
}

func (r *PostgresJobRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_import_export_job WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresJobRepository) scanJob(row interface {
	Scan(dest ...interface{}) error
}) (*domain.ImportExportJob, error) {
	job := &domain.ImportExportJob{}

	err := row.Scan(
		&job.ID, &job.Type, &job.EntityType, &job.Format, &job.Status,
		&job.FilePath, &job.FileName, &job.TotalRecords, &job.ProcessedRecords,
		&job.SuccessRecords, &job.FailedRecords, &job.ErrorLog,
		&job.StartedAt, &job.CompletedAt, &job.CreatedBy, &job.CreatedAt, &job.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (r *PostgresJobRepository) queryJobs(ctx context.Context, query string, args ...interface{}) ([]*domain.ImportExportJob, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]*domain.ImportExportJob, 0)
	for rows.Next() {
		job := &domain.ImportExportJob{}

		if err := rows.Scan(
			&job.ID, &job.Type, &job.EntityType, &job.Format, &job.Status,
			&job.FilePath, &job.FileName, &job.TotalRecords, &job.ProcessedRecords,
			&job.SuccessRecords, &job.FailedRecords, &job.ErrorLog,
			&job.StartedAt, &job.CompletedAt, &job.CreatedBy, &job.CreatedAt, &job.UpdatedAt,
		); err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
