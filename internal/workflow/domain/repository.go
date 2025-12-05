package domain

import "context"

// WorkflowRepository defines the interface for workflow definition persistence
type WorkflowRepository interface {
	// Create creates a new workflow definition
	Create(ctx context.Context, workflow *Workflow) error

	// Update updates an existing workflow definition
	Update(ctx context.Context, workflow *Workflow) error

	// FindByID finds a workflow by ID
	FindByID(ctx context.Context, id string) (*Workflow, error)

	// FindByName finds workflows by name
	FindByName(ctx context.Context, name string) ([]*Workflow, error)

	// FindByType finds workflows by type
	FindByType(ctx context.Context, workflowType WorkflowType, activeOnly bool) ([]*Workflow, error)

	// FindAll finds all workflow definitions
	FindAll(ctx context.Context, activeOnly bool) ([]*Workflow, error)

	// Delete deletes a workflow definition
	Delete(ctx context.Context, id string) error

	// ExistsByName checks if a workflow with the given name exists
	ExistsByName(ctx context.Context, name string) (bool, error)
}

// WorkflowExecutionRepository defines the interface for workflow execution persistence
type WorkflowExecutionRepository interface {
	// Create creates a new workflow execution
	Create(ctx context.Context, execution *WorkflowExecution) error

	// Update updates an existing workflow execution
	Update(ctx context.Context, execution *WorkflowExecution) error

	// FindByID finds a workflow execution by ID
	FindByID(ctx context.Context, id string) (*WorkflowExecution, error)

	// FindByWorkflowID finds all executions of a specific workflow
	FindByWorkflowID(ctx context.Context, workflowID string, limit int) ([]*WorkflowExecution, error)

	// FindByStatus finds workflow executions by status
	FindByStatus(ctx context.Context, status WorkflowStatus, limit int) ([]*WorkflowExecution, error)

	// FindByEntityReference finds workflow executions by entity reference
	FindByEntityReference(ctx context.Context, entityType, entityID string) ([]*WorkflowExecution, error)

	// FindActiveExecutions finds all active (running or suspended) workflow executions
	FindActiveExecutions(ctx context.Context, limit int) ([]*WorkflowExecution, error)

	// FindStaleExecutions finds workflow executions that haven't had a heartbeat within the specified duration
	FindStaleExecutions(ctx context.Context, staleAfterMinutes int, limit int) ([]*WorkflowExecution, error)

	// Delete deletes a workflow execution
	Delete(ctx context.Context, id string) error

	// CountByStatus counts workflow executions by status
	CountByStatus(ctx context.Context, status WorkflowStatus) (int64, error)
}
