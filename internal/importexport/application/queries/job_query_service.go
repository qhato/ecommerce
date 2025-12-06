package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/importexport/domain"
)

type JobQueryService struct {
	jobRepo domain.JobRepository
}

func NewJobQueryService(jobRepo domain.JobRepository) *JobQueryService {
	return &JobQueryService{jobRepo: jobRepo}
}

func (s *JobQueryService) GetJob(ctx context.Context, id int64) (*JobDTO, error) {
	job, err := s.jobRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find job: %w", err)
	}
	if job == nil {
		return nil, domain.ErrJobNotFound
	}

	return ToJobDTO(job), nil
}

func (s *JobQueryService) GetJobsByType(ctx context.Context, jobType string, status string, limit int) ([]*JobDTO, error) {
	jobs, err := s.jobRepo.FindByType(ctx, domain.JobType(jobType), domain.JobStatus(status), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find jobs: %w", err)
	}

	dtos := make([]*JobDTO, len(jobs))
	for i, job := range jobs {
		dtos[i] = ToJobDTO(job)
	}

	return dtos, nil
}

func (s *JobQueryService) GetJobsByStatus(ctx context.Context, status string, limit int) ([]*JobDTO, error) {
	jobs, err := s.jobRepo.FindByStatus(ctx, domain.JobStatus(status), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find jobs: %w", err)
	}

	dtos := make([]*JobDTO, len(jobs))
	for i, job := range jobs {
		dtos[i] = ToJobDTO(job)
	}

	return dtos, nil
}

func (s *JobQueryService) GetRecentJobs(ctx context.Context, limit int) ([]*JobDTO, error) {
	jobs, err := s.jobRepo.FindRecent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find recent jobs: %w", err)
	}

	dtos := make([]*JobDTO, len(jobs))
	for i, job := range jobs {
		dtos[i] = ToJobDTO(job)
	}

	return dtos, nil
}
