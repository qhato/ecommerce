package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/workflow/domain"
)

// WorkflowDTO represents a workflow definition for API responses
type WorkflowDTO struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description"`
	Type            string                   `json:"type"`
	Version         string                   `json:"version"`
	IsActive        bool                     `json:"is_active"`
	Activities      []ActivityDTO            `json:"activities"`
	Transitions     []TransitionDTO          `json:"transitions"`
	StartActivityID string                   `json:"start_activity_id"`
	EndActivityIDs  []string                 `json:"end_activity_ids"`
	Metadata        map[string]interface{}   `json:"metadata,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

// WorkflowExecutionDTO represents a workflow execution for API responses
type WorkflowExecutionDTO struct {
	ID                string                   `json:"id"`
	WorkflowID        string                   `json:"workflow_id"`
	WorkflowVersion   string                   `json:"workflow_version"`
	Status            string                   `json:"status"`
	Context           map[string]interface{}   `json:"context,omitempty"`
	InputData         map[string]interface{}   `json:"input_data,omitempty"`
	OutputData        map[string]interface{}   `json:"output_data,omitempty"`
	CurrentActivityID string                   `json:"current_activity_id,omitempty"`
	ActivityHistory   []ActivityExecutionDTO   `json:"activity_history,omitempty"`
	ErrorMessage      *string                  `json:"error_message,omitempty"`
	RetryCount        int                      `json:"retry_count"`
	StartedBy         string                   `json:"started_by"`
	StartedAt         time.Time                `json:"started_at"`
	CompletedAt       *time.Time               `json:"completed_at,omitempty"`
	LastHeartbeat     time.Time                `json:"last_heartbeat"`
	EntityType        string                   `json:"entity_type,omitempty"`
	EntityID          string                   `json:"entity_id,omitempty"`
	Metadata          map[string]interface{}   `json:"metadata,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

// ActivityDTO represents an activity for API responses
type ActivityDTO struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type"`
	Config      ActivityConfigDTO      `json:"config"`
	Rollback    *RollbackConfigDTO     `json:"rollback,omitempty"`
	Timeout     *int64                 `json:"timeout,omitempty"` // milliseconds
	RetryPolicy *RetryPolicyDTO        `json:"retry_policy,omitempty"`
	IsAsync     bool                   `json:"is_async"`
	Order       int                    `json:"order"`
}

// ActivityConfigDTO represents activity configuration
type ActivityConfigDTO struct {
	Handler       string                 `json:"handler"`
	InputMapping  map[string]string      `json:"input_mapping,omitempty"`
	OutputMapping map[string]string      `json:"output_mapping,omitempty"`
	Parameters    map[string]interface{} `json:"parameters,omitempty"`
	Conditions    []ConditionDTO         `json:"conditions,omitempty"`
}

// ConditionDTO represents a condition
type ConditionDTO struct {
	Expression     string `json:"expression"`
	NextActivityID string `json:"next_activity_id"`
}

// RollbackConfigDTO represents rollback configuration
type RollbackConfigDTO struct {
	Handler    string                 `json:"handler"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// RetryPolicyDTO represents retry policy
type RetryPolicyDTO struct {
	MaxAttempts     int     `json:"max_attempts"`
	InitialInterval int64   `json:"initial_interval"` // milliseconds
	MaxInterval     int64   `json:"max_interval"`     // milliseconds
	Multiplier      float64 `json:"multiplier"`
}

// TransitionDTO represents a transition
type TransitionDTO struct {
	FromActivityID string `json:"from_activity_id"`
	ToActivityID   string `json:"to_activity_id"`
	Condition      string `json:"condition,omitempty"`
	Priority       int    `json:"priority"`
}

// ActivityExecutionDTO represents an activity execution
type ActivityExecutionDTO struct {
	ID           string                 `json:"id"`
	ActivityID   string                 `json:"activity_id"`
	ActivityName string                 `json:"activity_name"`
	Status       string                 `json:"status"`
	InputData    map[string]interface{} `json:"input_data,omitempty"`
	OutputData   map[string]interface{} `json:"output_data,omitempty"`
	ErrorMessage *string                `json:"error_message,omitempty"`
	RetryCount   int                    `json:"retry_count"`
	StartedAt    time.Time              `json:"started_at"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty"`
	Duration     *int64                 `json:"duration,omitempty"` // milliseconds
}

// ToWorkflowDTO converts domain Workflow to WorkflowDTO
func ToWorkflowDTO(workflow *domain.Workflow) *WorkflowDTO {
	activities := make([]ActivityDTO, len(workflow.Activities))
	for i, a := range workflow.Activities {
		activities[i] = toActivityDTO(a)
	}

	transitions := make([]TransitionDTO, len(workflow.Transitions))
	for i, t := range workflow.Transitions {
		transitions[i] = toTransitionDTO(t)
	}

	return &WorkflowDTO{
		ID:              workflow.ID,
		Name:            workflow.Name,
		Description:     workflow.Description,
		Type:            string(workflow.Type),
		Version:         workflow.Version,
		IsActive:        workflow.IsActive,
		Activities:      activities,
		Transitions:     transitions,
		StartActivityID: workflow.StartActivityID,
		EndActivityIDs:  workflow.EndActivityIDs,
		Metadata:        workflow.Metadata,
		CreatedAt:       workflow.CreatedAt,
		UpdatedAt:       workflow.UpdatedAt,
	}
}

// ToWorkflowExecutionDTO converts domain WorkflowExecution to WorkflowExecutionDTO
func ToWorkflowExecutionDTO(execution *domain.WorkflowExecution) *WorkflowExecutionDTO {
	activityHistory := make([]ActivityExecutionDTO, len(execution.ActivityHistory))
	for i, ae := range execution.ActivityHistory {
		activityHistory[i] = toActivityExecutionDTO(ae)
	}

	return &WorkflowExecutionDTO{
		ID:                execution.ID,
		WorkflowID:        execution.WorkflowID,
		WorkflowVersion:   execution.WorkflowVersion,
		Status:            string(execution.Status),
		Context:           execution.Context,
		InputData:         execution.InputData,
		OutputData:        execution.OutputData,
		CurrentActivityID: execution.CurrentActivityID,
		ActivityHistory:   activityHistory,
		ErrorMessage:      execution.ErrorMessage,
		RetryCount:        execution.RetryCount,
		StartedBy:         execution.StartedBy,
		StartedAt:         execution.StartedAt,
		CompletedAt:       execution.CompletedAt,
		LastHeartbeat:     execution.LastHeartbeat,
		EntityType:        execution.EntityType,
		EntityID:          execution.EntityID,
		Metadata:          execution.Metadata,
		CreatedAt:         execution.CreatedAt,
		UpdatedAt:         execution.UpdatedAt,
	}
}

func toActivityDTO(activity domain.Activity) ActivityDTO {
	dto := ActivityDTO{
		ID:          activity.ID,
		Name:        activity.Name,
		Description: activity.Description,
		Type:        string(activity.Type),
		Config:      toActivityConfigDTO(activity.Config),
		IsAsync:     activity.IsAsync,
		Order:       activity.Order,
	}

	if activity.Rollback != nil {
		dto.Rollback = &RollbackConfigDTO{
			Handler:    activity.Rollback.Handler,
			Parameters: activity.Rollback.Parameters,
		}
	}

	if activity.Timeout != nil {
		timeout := activity.Timeout.Milliseconds()
		dto.Timeout = &timeout
	}

	if activity.RetryPolicy != nil {
		dto.RetryPolicy = &RetryPolicyDTO{
			MaxAttempts:     activity.RetryPolicy.MaxAttempts,
			InitialInterval: activity.RetryPolicy.InitialInterval.Milliseconds(),
			MaxInterval:     activity.RetryPolicy.MaxInterval.Milliseconds(),
			Multiplier:      activity.RetryPolicy.Multiplier,
		}
	}

	return dto
}

func toActivityConfigDTO(config domain.ActivityConfig) ActivityConfigDTO {
	conditions := make([]ConditionDTO, len(config.Conditions))
	for i, c := range config.Conditions {
		conditions[i] = ConditionDTO{
			Expression:     c.Expression,
			NextActivityID: c.NextActivityID,
		}
	}

	return ActivityConfigDTO{
		Handler:       config.Handler,
		InputMapping:  config.InputMapping,
		OutputMapping: config.OutputMapping,
		Parameters:    config.Parameters,
		Conditions:    conditions,
	}
}

func toTransitionDTO(transition domain.Transition) TransitionDTO {
	return TransitionDTO{
		FromActivityID: transition.FromActivityID,
		ToActivityID:   transition.ToActivityID,
		Condition:      transition.Condition,
		Priority:       transition.Priority,
	}
}

func toActivityExecutionDTO(execution domain.ActivityExecution) ActivityExecutionDTO {
	return ActivityExecutionDTO{
		ID:           execution.ID,
		ActivityID:   execution.ActivityID,
		ActivityName: execution.ActivityName,
		Status:       string(execution.Status),
		InputData:    execution.InputData,
		OutputData:   execution.OutputData,
		ErrorMessage: execution.ErrorMessage,
		RetryCount:   execution.RetryCount,
		StartedAt:    execution.StartedAt,
		CompletedAt:  execution.CompletedAt,
		Duration:     execution.Duration,
	}
}
