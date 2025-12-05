package domain

import (
	"encoding/json"
	"time"
)

// WorkflowExecution represents a runtime instance of a workflow
// Business Logic: Track execution state and progress of workflow instances
type WorkflowExecution struct {
	ID              string
	WorkflowID      string
	WorkflowVersion string
	Status          WorkflowStatus

	// Context and data
	Context         map[string]interface{} // Workflow-level variables and data
	InputData       map[string]interface{} // Initial input data
	OutputData      map[string]interface{} // Final output data

	// Execution tracking
	CurrentActivityID string
	ActivityHistory   []ActivityExecution
	ErrorMessage      *string
	RetryCount        int

	// Metadata
	StartedBy       string
	StartedAt       time.Time
	CompletedAt     *time.Time
	LastHeartbeat   time.Time

	// Associated entity (e.g., checkout session, order)
	EntityType      string
	EntityID        string

	// Metadata
	Metadata        map[string]interface{}
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ActivityExecution represents the execution of a single activity
type ActivityExecution struct {
	ID              string
	ActivityID      string
	ActivityName    string
	Status          ActivityStatus

	// Execution details
	InputData       map[string]interface{}
	OutputData      map[string]interface{}
	ErrorMessage    *string
	RetryCount      int

	// Timing
	StartedAt       time.Time
	CompletedAt     *time.Time
	Duration        *int64 // Duration in milliseconds
}

// ActivityStatus represents the status of an activity execution
type ActivityStatus string

const (
	ActivityStatusPending    ActivityStatus = "PENDING"
	ActivityStatusRunning    ActivityStatus = "RUNNING"
	ActivityStatusCompleted  ActivityStatus = "COMPLETED"
	ActivityStatusFailed     ActivityStatus = "FAILED"
	ActivityStatusSkipped    ActivityStatus = "SKIPPED"
	ActivityStatusWaiting    ActivityStatus = "WAITING" // Waiting for external event
)

// NewWorkflowExecution creates a new workflow execution instance
func NewWorkflowExecution(workflowID, workflowVersion, startedBy string, inputData map[string]interface{}) *WorkflowExecution {
	now := time.Now()
	return &WorkflowExecution{
		ID:              generateExecutionID(),
		WorkflowID:      workflowID,
		WorkflowVersion: workflowVersion,
		Status:          WorkflowStatusPending,
		Context:         make(map[string]interface{}),
		InputData:       inputData,
		OutputData:      make(map[string]interface{}),
		ActivityHistory: make([]ActivityExecution, 0),
		StartedBy:       startedBy,
		StartedAt:       now,
		LastHeartbeat:   now,
		Metadata:        make(map[string]interface{}),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Start starts the workflow execution
func (we *WorkflowExecution) Start(startActivityID string) error {
	if we.Status != WorkflowStatusPending {
		return ErrWorkflowNotInPendingState
	}

	we.Status = WorkflowStatusRunning
	we.CurrentActivityID = startActivityID
	we.UpdatedAt = time.Now()
	we.LastHeartbeat = time.Now()

	return nil
}

// StartActivity starts an activity execution
func (we *WorkflowExecution) StartActivity(activityID, activityName string, inputData map[string]interface{}) error {
	if we.Status != WorkflowStatusRunning {
		return ErrWorkflowNotRunning
	}

	activityExec := ActivityExecution{
		ID:           generateActivityExecutionID(),
		ActivityID:   activityID,
		ActivityName: activityName,
		Status:       ActivityStatusRunning,
		InputData:    inputData,
		OutputData:   make(map[string]interface{}),
		StartedAt:    time.Now(),
	}

	we.ActivityHistory = append(we.ActivityHistory, activityExec)
	we.CurrentActivityID = activityID
	we.UpdatedAt = time.Now()
	we.LastHeartbeat = time.Now()

	return nil
}

// CompleteActivity completes an activity execution
func (we *WorkflowExecution) CompleteActivity(activityID string, outputData map[string]interface{}) error {
	if we.Status != WorkflowStatusRunning {
		return ErrWorkflowNotRunning
	}

	// Find the activity execution
	activityExec, err := we.findLastActivityExecution(activityID)
	if err != nil {
		return err
	}

	now := time.Now()
	duration := now.Sub(activityExec.StartedAt).Milliseconds()

	activityExec.Status = ActivityStatusCompleted
	activityExec.OutputData = outputData
	activityExec.CompletedAt = &now
	activityExec.Duration = &duration

	// Merge output data into workflow context
	we.mergeOutputToContext(outputData)

	we.UpdatedAt = now
	we.LastHeartbeat = now

	return nil
}

// FailActivity marks an activity as failed
func (we *WorkflowExecution) FailActivity(activityID string, errorMsg string) error {
	if we.Status != WorkflowStatusRunning {
		return ErrWorkflowNotRunning
	}

	// Find the activity execution
	activityExec, err := we.findLastActivityExecution(activityID)
	if err != nil {
		return err
	}

	now := time.Now()
	duration := now.Sub(activityExec.StartedAt).Milliseconds()

	activityExec.Status = ActivityStatusFailed
	activityExec.ErrorMessage = &errorMsg
	activityExec.CompletedAt = &now
	activityExec.Duration = &duration
	activityExec.RetryCount++

	we.UpdatedAt = now
	we.LastHeartbeat = now

	return nil
}

// MoveToNextActivity moves execution to the next activity
func (we *WorkflowExecution) MoveToNextActivity(nextActivityID string) error {
	if we.Status != WorkflowStatusRunning {
		return ErrWorkflowNotRunning
	}

	we.CurrentActivityID = nextActivityID
	we.UpdatedAt = time.Now()
	we.LastHeartbeat = time.Now()

	return nil
}

// Complete marks the workflow execution as completed
func (we *WorkflowExecution) Complete(outputData map[string]interface{}) error {
	if we.Status != WorkflowStatusRunning {
		return ErrWorkflowNotRunning
	}

	now := time.Now()
	we.Status = WorkflowStatusCompleted
	we.OutputData = outputData
	we.CompletedAt = &now
	we.UpdatedAt = now

	return nil
}

// Fail marks the workflow execution as failed
func (we *WorkflowExecution) Fail(errorMsg string) error {
	if we.Status == WorkflowStatusCompleted || we.Status == WorkflowStatusCancelled {
		return ErrWorkflowAlreadyFinished
	}

	now := time.Now()
	we.Status = WorkflowStatusFailed
	we.ErrorMessage = &errorMsg
	we.CompletedAt = &now
	we.UpdatedAt = now

	return nil
}

// Suspend suspends the workflow execution (e.g., waiting for external event)
func (we *WorkflowExecution) Suspend() error {
	if we.Status != WorkflowStatusRunning {
		return ErrWorkflowNotRunning
	}

	we.Status = WorkflowStatusSuspended
	we.UpdatedAt = time.Now()

	return nil
}

// Resume resumes a suspended workflow execution
func (we *WorkflowExecution) Resume() error {
	if we.Status != WorkflowStatusSuspended {
		return ErrWorkflowNotSuspended
	}

	we.Status = WorkflowStatusRunning
	we.UpdatedAt = time.Now()
	we.LastHeartbeat = time.Now()

	return nil
}

// Cancel cancels the workflow execution
func (we *WorkflowExecution) Cancel() error {
	if we.Status == WorkflowStatusCompleted {
		return ErrWorkflowAlreadyCompleted
	}
	if we.Status == WorkflowStatusCancelled {
		return ErrWorkflowAlreadyCancelled
	}

	now := time.Now()
	we.Status = WorkflowStatusCancelled
	we.CompletedAt = &now
	we.UpdatedAt = now

	return nil
}

// Heartbeat updates the last heartbeat timestamp
func (we *WorkflowExecution) Heartbeat() {
	we.LastHeartbeat = time.Now()
	we.UpdatedAt = time.Now()
}

// SetContext sets a value in the workflow context
func (we *WorkflowExecution) SetContext(key string, value interface{}) {
	we.Context[key] = value
	we.UpdatedAt = time.Now()
}

// GetContext gets a value from the workflow context
func (we *WorkflowExecution) GetContext(key string) (interface{}, bool) {
	value, exists := we.Context[key]
	return value, exists
}

// SetEntityReference sets the entity reference for this workflow
func (we *WorkflowExecution) SetEntityReference(entityType, entityID string) {
	we.EntityType = entityType
	we.EntityID = entityID
	we.UpdatedAt = time.Now()
}

// GetLastActivityExecution gets the last execution of an activity
func (we *WorkflowExecution) GetLastActivityExecution(activityID string) (*ActivityExecution, error) {
	return we.findLastActivityExecution(activityID)
}

// GetActivityExecutions gets all executions of a specific activity
func (we *WorkflowExecution) GetActivityExecutions(activityID string) []ActivityExecution {
	executions := make([]ActivityExecution, 0)
	for _, exec := range we.ActivityHistory {
		if exec.ActivityID == activityID {
			executions = append(executions, exec)
		}
	}
	return executions
}

// IsCompleted checks if the workflow is completed
func (we *WorkflowExecution) IsCompleted() bool {
	return we.Status == WorkflowStatusCompleted
}

// IsFailed checks if the workflow has failed
func (we *WorkflowExecution) IsFailed() bool {
	return we.Status == WorkflowStatusFailed
}

// IsRunning checks if the workflow is running
func (we *WorkflowExecution) IsRunning() bool {
	return we.Status == WorkflowStatusRunning
}

// IsSuspended checks if the workflow is suspended
func (we *WorkflowExecution) IsSuspended() bool {
	return we.Status == WorkflowStatusSuspended
}

// GetDuration gets the total duration of the workflow execution
func (we *WorkflowExecution) GetDuration() *time.Duration {
	if we.CompletedAt == nil {
		return nil
	}
	duration := we.CompletedAt.Sub(we.StartedAt)
	return &duration
}

// MarshalContextToJSON marshals the context to JSON
func (we *WorkflowExecution) MarshalContextToJSON() ([]byte, error) {
	return json.Marshal(we.Context)
}

// UnmarshalContextFromJSON unmarshals the context from JSON
func (we *WorkflowExecution) UnmarshalContextFromJSON(data []byte) error {
	return json.Unmarshal(data, &we.Context)
}

// Private helper methods

func (we *WorkflowExecution) findLastActivityExecution(activityID string) (*ActivityExecution, error) {
	// Search from the end to find the most recent execution
	for i := len(we.ActivityHistory) - 1; i >= 0; i-- {
		if we.ActivityHistory[i].ActivityID == activityID {
			return &we.ActivityHistory[i], nil
		}
	}
	return nil, ErrActivityExecutionNotFound
}

func (we *WorkflowExecution) mergeOutputToContext(outputData map[string]interface{}) {
	for key, value := range outputData {
		we.Context[key] = value
	}
}

func generateExecutionID() string {
	return "WE-" + time.Now().Format("20060102150405")
}

func generateActivityExecutionID() string {
	return "AE-" + time.Now().Format("20060102150405")
}
