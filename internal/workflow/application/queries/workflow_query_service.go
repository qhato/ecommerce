package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/workflow/domain"
)

// WorkflowQueryService handles workflow-related queries
type WorkflowQueryService struct {
	workflowRepo  domain.WorkflowRepository
	executionRepo domain.WorkflowExecutionRepository
}

// NewWorkflowQueryService creates a new workflow query service
func NewWorkflowQueryService(
	workflowRepo domain.WorkflowRepository,
	executionRepo domain.WorkflowExecutionRepository,
) *WorkflowQueryService {
	return &WorkflowQueryService{
		workflowRepo:  workflowRepo,
		executionRepo: executionRepo,
	}
}

// GetWorkflow retrieves a workflow definition by ID
func (s *WorkflowQueryService) GetWorkflow(ctx context.Context, id string) (*WorkflowDTO, error) {
	workflow, err := s.workflowRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find workflow: %w", err)
	}
	if workflow == nil {
		return nil, domain.ErrWorkflowNotFound
	}

	return ToWorkflowDTO(workflow), nil
}

// GetWorkflowsByName retrieves workflows by name
func (s *WorkflowQueryService) GetWorkflowsByName(ctx context.Context, name string) ([]*WorkflowDTO, error) {
	workflows, err := s.workflowRepo.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to find workflows: %w", err)
	}

	dtos := make([]*WorkflowDTO, len(workflows))
	for i, w := range workflows {
		dtos[i] = ToWorkflowDTO(w)
	}

	return dtos, nil
}

// GetWorkflowsByType retrieves workflows by type
func (s *WorkflowQueryService) GetWorkflowsByType(ctx context.Context, workflowType string, activeOnly bool) ([]*WorkflowDTO, error) {
	workflows, err := s.workflowRepo.FindByType(ctx, domain.WorkflowType(workflowType), activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find workflows: %w", err)
	}

	dtos := make([]*WorkflowDTO, len(workflows))
	for i, w := range workflows {
		dtos[i] = ToWorkflowDTO(w)
	}

	return dtos, nil
}

// GetAllWorkflows retrieves all workflow definitions
func (s *WorkflowQueryService) GetAllWorkflows(ctx context.Context, activeOnly bool) ([]*WorkflowDTO, error) {
	workflows, err := s.workflowRepo.FindAll(ctx, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find workflows: %w", err)
	}

	dtos := make([]*WorkflowDTO, len(workflows))
	for i, w := range workflows {
		dtos[i] = ToWorkflowDTO(w)
	}

	return dtos, nil
}

// GetWorkflowExecution retrieves a workflow execution by ID
func (s *WorkflowQueryService) GetWorkflowExecution(ctx context.Context, id string) (*WorkflowExecutionDTO, error) {
	execution, err := s.executionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find execution: %w", err)
	}
	if execution == nil {
		return nil, domain.ErrWorkflowExecutionNotFound
	}

	return ToWorkflowExecutionDTO(execution), nil
}

// GetWorkflowExecutionsByWorkflowID retrieves all executions of a workflow
func (s *WorkflowQueryService) GetWorkflowExecutionsByWorkflowID(ctx context.Context, workflowID string, limit int) ([]*WorkflowExecutionDTO, error) {
	executions, err := s.executionRepo.FindByWorkflowID(ctx, workflowID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find executions: %w", err)
	}

	dtos := make([]*WorkflowExecutionDTO, len(executions))
	for i, e := range executions {
		dtos[i] = ToWorkflowExecutionDTO(e)
	}

	return dtos, nil
}

// GetWorkflowExecutionsByStatus retrieves workflow executions by status
func (s *WorkflowQueryService) GetWorkflowExecutionsByStatus(ctx context.Context, status string, limit int) ([]*WorkflowExecutionDTO, error) {
	executions, err := s.executionRepo.FindByStatus(ctx, domain.WorkflowStatus(status), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find executions: %w", err)
	}

	dtos := make([]*WorkflowExecutionDTO, len(executions))
	for i, e := range executions {
		dtos[i] = ToWorkflowExecutionDTO(e)
	}

	return dtos, nil
}

// GetWorkflowExecutionsByEntityReference retrieves workflow executions by entity reference
func (s *WorkflowQueryService) GetWorkflowExecutionsByEntityReference(ctx context.Context, entityType, entityID string) ([]*WorkflowExecutionDTO, error) {
	executions, err := s.executionRepo.FindByEntityReference(ctx, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to find executions: %w", err)
	}

	dtos := make([]*WorkflowExecutionDTO, len(executions))
	for i, e := range executions {
		dtos[i] = ToWorkflowExecutionDTO(e)
	}

	return dtos, nil
}

// GetActiveWorkflowExecutions retrieves all active workflow executions
func (s *WorkflowQueryService) GetActiveWorkflowExecutions(ctx context.Context, limit int) ([]*WorkflowExecutionDTO, error) {
	executions, err := s.executionRepo.FindActiveExecutions(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find active executions: %w", err)
	}

	dtos := make([]*WorkflowExecutionDTO, len(executions))
	for i, e := range executions {
		dtos[i] = ToWorkflowExecutionDTO(e)
	}

	return dtos, nil
}

// GetStaleWorkflowExecutions retrieves stale workflow executions
func (s *WorkflowQueryService) GetStaleWorkflowExecutions(ctx context.Context, staleAfterMinutes, limit int) ([]*WorkflowExecutionDTO, error) {
	executions, err := s.executionRepo.FindStaleExecutions(ctx, staleAfterMinutes, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find stale executions: %w", err)
	}

	dtos := make([]*WorkflowExecutionDTO, len(executions))
	for i, e := range executions {
		dtos[i] = ToWorkflowExecutionDTO(e)
	}

	return dtos, nil
}

// CountWorkflowExecutionsByStatus counts workflow executions by status
func (s *WorkflowQueryService) CountWorkflowExecutionsByStatus(ctx context.Context, status string) (int64, error) {
	count, err := s.executionRepo.CountByStatus(ctx, domain.WorkflowStatus(status))
	if err != nil {
		return 0, fmt.Errorf("failed to count executions: %w", err)
	}

	return count, nil
}
