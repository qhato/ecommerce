package commands

import (
	"time"

	"github.com/qhato/ecommerce/internal/workflow/domain"
)

// CreateWorkflowCommand creates a new workflow definition
type CreateWorkflowCommand struct {
	Name        string
	Description string
	Type        string
	Activities  []ActivityDTO
	Transitions []TransitionDTO
	StartActivityID string
	EndActivityIDs []string
	Metadata    map[string]interface{}
	CreatedBy   string
}

// UpdateWorkflowCommand updates an existing workflow definition
type UpdateWorkflowCommand struct {
	ID          string
	Name        string
	Description string
	Activities  []ActivityDTO
	Transitions []TransitionDTO
	StartActivityID string
	EndActivityIDs []string
	Metadata    map[string]interface{}
	UpdatedBy   string
}

// ActivateWorkflowCommand activates a workflow
type ActivateWorkflowCommand struct {
	ID string
}

// DeactivateWorkflowCommand deactivates a workflow
type DeactivateWorkflowCommand struct {
	ID string
}

// DeleteWorkflowCommand deletes a workflow
type DeleteWorkflowCommand struct {
	ID string
}

// StartWorkflowExecutionCommand starts a new workflow execution
type StartWorkflowExecutionCommand struct {
	WorkflowID  string
	InputData   map[string]interface{}
	StartedBy   string
	EntityType  string
	EntityID    string
}

// CompleteActivityCommand completes an activity in a workflow execution
type CompleteActivityCommand struct {
	ExecutionID string
	ActivityID  string
	OutputData  map[string]interface{}
}

// FailActivityCommand marks an activity as failed
type FailActivityCommand struct {
	ExecutionID  string
	ActivityID   string
	ErrorMessage string
}

// SuspendWorkflowCommand suspends a workflow execution
type SuspendWorkflowCommand struct {
	ExecutionID   string
	SuspendReason string
}

// ResumeWorkflowCommand resumes a suspended workflow execution
type ResumeWorkflowCommand struct {
	ExecutionID string
	ResumedBy   string
}

// CancelWorkflowCommand cancels a workflow execution
type CancelWorkflowCommand struct {
	ExecutionID  string
	CancelReason string
	CancelledBy  string
}

// SetWorkflowContextCommand sets a value in workflow context
type SetWorkflowContextCommand struct {
	ExecutionID string
	Key         string
	Value       interface{}
}

// ActivityDTO represents an activity in commands
type ActivityDTO struct {
	ID          string
	Name        string
	Description string
	Type        string
	Config      ActivityConfigDTO
	Rollback    *RollbackConfigDTO
	Timeout     *int64 // milliseconds
	RetryPolicy *RetryPolicyDTO
	IsAsync     bool
	Order       int
}

// ActivityConfigDTO represents activity configuration
type ActivityConfigDTO struct {
	Handler       string
	InputMapping  map[string]string
	OutputMapping map[string]string
	Parameters    map[string]interface{}
	Conditions    []ConditionDTO
}

// ConditionDTO represents a condition
type ConditionDTO struct {
	Expression     string
	NextActivityID string
}

// RollbackConfigDTO represents rollback configuration
type RollbackConfigDTO struct {
	Handler    string
	Parameters map[string]interface{}
}

// RetryPolicyDTO represents retry policy
type RetryPolicyDTO struct {
	MaxAttempts     int
	InitialInterval int64 // milliseconds
	MaxInterval     int64 // milliseconds
	Multiplier      float64
}

// TransitionDTO represents a transition
type TransitionDTO struct {
	FromActivityID string
	ToActivityID   string
	Condition      string
	Priority       int
}

// ToActivity converts ActivityDTO to domain Activity
func (dto ActivityDTO) ToActivity() domain.Activity {
	activity := domain.Activity{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Type:        domain.ActivityType(dto.Type),
		Config:      dto.Config.ToActivityConfig(),
		IsAsync:     dto.IsAsync,
		Order:       dto.Order,
	}

	if dto.Rollback != nil {
		rollback := domain.RollbackConfig{
			Handler:    dto.Rollback.Handler,
			Parameters: dto.Rollback.Parameters,
		}
		activity.Rollback = &rollback
	}

	if dto.Timeout != nil {
		timeout := time.Duration(*dto.Timeout) * time.Millisecond
		activity.Timeout = &timeout
	}

	if dto.RetryPolicy != nil {
		retryPolicy := domain.RetryPolicy{
			MaxAttempts:     dto.RetryPolicy.MaxAttempts,
			InitialInterval: time.Duration(dto.RetryPolicy.InitialInterval) * time.Millisecond,
			MaxInterval:     time.Duration(dto.RetryPolicy.MaxInterval) * time.Millisecond,
			Multiplier:      dto.RetryPolicy.Multiplier,
		}
		activity.RetryPolicy = &retryPolicy
	}

	return activity
}

// ToActivityConfig converts ActivityConfigDTO to domain ActivityConfig
func (dto ActivityConfigDTO) ToActivityConfig() domain.ActivityConfig {
	conditions := make([]domain.Condition, len(dto.Conditions))
	for i, c := range dto.Conditions {
		conditions[i] = domain.Condition{
			Expression:     c.Expression,
			NextActivityID: c.NextActivityID,
		}
	}

	return domain.ActivityConfig{
		Handler:       dto.Handler,
		InputMapping:  dto.InputMapping,
		OutputMapping: dto.OutputMapping,
		Parameters:    dto.Parameters,
		Conditions:    conditions,
	}
}

// ToTransition converts TransitionDTO to domain Transition
func (dto TransitionDTO) ToTransition() domain.Transition {
	return domain.Transition{
		FromActivityID: dto.FromActivityID,
		ToActivityID:   dto.ToActivityID,
		Condition:      dto.Condition,
		Priority:       dto.Priority,
	}
}
