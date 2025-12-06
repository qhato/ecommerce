package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/importexport/domain"
)

type JobCommandHandler struct {
	jobRepo domain.JobRepository
}

func NewJobCommandHandler(jobRepo domain.JobRepository) *JobCommandHandler {
	return &JobCommandHandler{jobRepo: jobRepo}
}

func (h *JobCommandHandler) HandleCreateImportJob(ctx context.Context, cmd CreateImportJobCommand) (*domain.ImportExportJob, error) {
	job, err := domain.NewImportJob(
		domain.EntityType(cmd.EntityType),
		domain.FileFormat(cmd.Format),
		cmd.FilePath,
		cmd.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	job.FileName = cmd.FileName

	if err := h.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create import job: %w", err)
	}

	return job, nil
}

func (h *JobCommandHandler) HandleCreateExportJob(ctx context.Context, cmd CreateExportJobCommand) (*domain.ImportExportJob, error) {
	job, err := domain.NewExportJob(
		domain.EntityType(cmd.EntityType),
		domain.FileFormat(cmd.Format),
		cmd.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	if err := h.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create export job: %w", err)
	}

	return job, nil
}

func (h *JobCommandHandler) HandleStartJob(ctx context.Context, cmd StartJobCommand) (*domain.ImportExportJob, error) {
	job, err := h.jobRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find job: %w", err)
	}
	if job == nil {
		return nil, domain.ErrJobNotFound
	}

	if job.Status != domain.JobStatusPending {
		return nil, domain.ErrJobAlreadyRunning
	}

	job.Start()
	if err := h.jobRepo.Update(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to start job: %w", err)
	}

	return job, nil
}

func (h *JobCommandHandler) HandleCompleteJob(ctx context.Context, cmd CompleteJobCommand) (*domain.ImportExportJob, error) {
	job, err := h.jobRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find job: %w", err)
	}
	if job == nil {
		return nil, domain.ErrJobNotFound
	}

	if job.Status != domain.JobStatusProcessing {
		return nil, domain.ErrJobNotProcessing
	}

	job.Complete()
	if err := h.jobRepo.Update(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to complete job: %w", err)
	}

	return job, nil
}

func (h *JobCommandHandler) HandleFailJob(ctx context.Context, cmd FailJobCommand) (*domain.ImportExportJob, error) {
	job, err := h.jobRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find job: %w", err)
	}
	if job == nil {
		return nil, domain.ErrJobNotFound
	}

	job.Fail(cmd.ErrorLog)
	if err := h.jobRepo.Update(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to mark job as failed: %w", err)
	}

	return job, nil
}

func (h *JobCommandHandler) HandleUpdateProgress(ctx context.Context, cmd UpdateProgressCommand) (*domain.ImportExportJob, error) {
	job, err := h.jobRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find job: %w", err)
	}
	if job == nil {
		return nil, domain.ErrJobNotFound
	}

	job.TotalRecords = cmd.TotalRecords
	job.UpdateProgress(cmd.ProcessedRecords, cmd.SuccessRecords, cmd.FailedRecords)

	if err := h.jobRepo.Update(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to update job progress: %w", err)
	}

	return job, nil
}

func (h *JobCommandHandler) HandleCancelJob(ctx context.Context, cmd CancelJobCommand) (*domain.ImportExportJob, error) {
	job, err := h.jobRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find job: %w", err)
	}
	if job == nil {
		return nil, domain.ErrJobNotFound
	}

	job.Status = domain.JobStatusCancelled
	job.UpdatedAt = job.UpdatedAt // Update timestamp

	if err := h.jobRepo.Update(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to cancel job: %w", err)
	}

	return job, nil
}

func (h *JobCommandHandler) HandleDeleteJob(ctx context.Context, cmd DeleteJobCommand) error {
	return h.jobRepo.Delete(ctx, cmd.ID)
}
