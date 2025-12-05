package domain

import "time"

// WorkflowEventType represents the type of workflow event
type WorkflowEventType string

const (
	// Workflow Definition Events
	EventWorkflowCreated     WorkflowEventType = "workflow.created"
	EventWorkflowUpdated     WorkflowEventType = "workflow.updated"
	EventWorkflowActivated   WorkflowEventType = "workflow.activated"
	EventWorkflowDeactivated WorkflowEventType = "workflow.deactivated"
	EventWorkflowDeleted     WorkflowEventType = "workflow.deleted"

	// Workflow Execution Events
	EventWorkflowStarted    WorkflowEventType = "workflow.execution.started"
	EventWorkflowCompleted  WorkflowEventType = "workflow.execution.completed"
	EventWorkflowFailed     WorkflowEventType = "workflow.execution.failed"
	EventWorkflowSuspended  WorkflowEventType = "workflow.execution.suspended"
	EventWorkflowResumed    WorkflowEventType = "workflow.execution.resumed"
	EventWorkflowCancelled  WorkflowEventType = "workflow.execution.cancelled"

	// Activity Execution Events
	EventActivityStarted    WorkflowEventType = "workflow.activity.started"
	EventActivityCompleted  WorkflowEventType = "workflow.activity.completed"
	EventActivityFailed     WorkflowEventType = "workflow.activity.failed"
	EventActivityRetrying   WorkflowEventType = "workflow.activity.retrying"
	EventActivitySkipped    WorkflowEventType = "workflow.activity.skipped"

	// Transition Events
	EventTransitionExecuted WorkflowEventType = "workflow.transition.executed"
)

// WorkflowEvent is the base event for all workflow-related events
type WorkflowEvent struct {
	EventType   WorkflowEventType
	EventID     string
	WorkflowID  string
	ExecutionID string
	OccurredAt  time.Time
	Data        interface{}
}

// WorkflowCreatedEvent is published when a new workflow definition is created
type WorkflowCreatedEvent struct {
	WorkflowEvent
	WorkflowName    string
	WorkflowType    WorkflowType
	WorkflowVersion string
	CreatedBy       string
}

// WorkflowUpdatedEvent is published when a workflow definition is updated
type WorkflowUpdatedEvent struct {
	WorkflowEvent
	WorkflowName    string
	WorkflowVersion string
	UpdatedBy       string
	Changes         map[string]interface{}
}

// WorkflowActivatedEvent is published when a workflow is activated
type WorkflowActivatedEvent struct {
	WorkflowEvent
	WorkflowName    string
	ActivatedBy     string
}

// WorkflowDeactivatedEvent is published when a workflow is deactivated
type WorkflowDeactivatedEvent struct {
	WorkflowEvent
	WorkflowName    string
	DeactivatedBy   string
}

// WorkflowDeletedEvent is published when a workflow definition is deleted
type WorkflowDeletedEvent struct {
	WorkflowEvent
	WorkflowName    string
	DeletedBy       string
}

// WorkflowStartedEvent is published when a workflow execution starts
type WorkflowStartedEvent struct {
	WorkflowEvent
	WorkflowName        string
	WorkflowVersion     string
	StartActivityID     string
	StartedBy           string
	InputData           map[string]interface{}
	EntityType          string
	EntityID            string
}

// WorkflowCompletedEvent is published when a workflow execution completes successfully
type WorkflowCompletedEvent struct {
	WorkflowEvent
	WorkflowName        string
	Duration            time.Duration
	OutputData          map[string]interface{}
	TotalActivities     int
	CompletedActivities int
}

// WorkflowFailedEvent is published when a workflow execution fails
type WorkflowFailedEvent struct {
	WorkflowEvent
	WorkflowName        string
	CurrentActivityID   string
	ErrorMessage        string
	Duration            time.Duration
	RetryCount          int
}

// WorkflowSuspendedEvent is published when a workflow execution is suspended
type WorkflowSuspendedEvent struct {
	WorkflowEvent
	WorkflowName        string
	CurrentActivityID   string
	SuspendReason       string
}

// WorkflowResumedEvent is published when a workflow execution is resumed
type WorkflowResumedEvent struct {
	WorkflowEvent
	WorkflowName        string
	CurrentActivityID   string
	ResumedBy           string
}

// WorkflowCancelledEvent is published when a workflow execution is cancelled
type WorkflowCancelledEvent struct {
	WorkflowEvent
	WorkflowName        string
	CurrentActivityID   string
	CancelReason        string
	CancelledBy         string
}

// ActivityStartedEvent is published when an activity execution starts
type ActivityStartedEvent struct {
	WorkflowEvent
	ActivityID          string
	ActivityName        string
	ActivityType        ActivityType
	InputData           map[string]interface{}
}

// ActivityCompletedEvent is published when an activity execution completes
type ActivityCompletedEvent struct {
	WorkflowEvent
	ActivityID          string
	ActivityName        string
	ActivityType        ActivityType
	Duration            time.Duration
	OutputData          map[string]interface{}
}

// ActivityFailedEvent is published when an activity execution fails
type ActivityFailedEvent struct {
	WorkflowEvent
	ActivityID          string
	ActivityName        string
	ActivityType        ActivityType
	ErrorMessage        string
	Duration            time.Duration
	RetryCount          int
}

// ActivityRetryingEvent is published when an activity is retrying
type ActivityRetryingEvent struct {
	WorkflowEvent
	ActivityID          string
	ActivityName        string
	RetryCount          int
	MaxRetries          int
	NextRetryAt         time.Time
}

// ActivitySkippedEvent is published when an activity is skipped
type ActivitySkippedEvent struct {
	WorkflowEvent
	ActivityID          string
	ActivityName        string
	SkipReason          string
}

// TransitionExecutedEvent is published when a transition is executed
type TransitionExecutedEvent struct {
	WorkflowEvent
	FromActivityID      string
	ToActivityID        string
	Condition           string
	ConditionResult     bool
}

// NewWorkflowEvent creates a new base workflow event
func NewWorkflowEvent(eventType WorkflowEventType, workflowID, executionID string) WorkflowEvent {
	return WorkflowEvent{
		EventType:   eventType,
		EventID:     generateEventID(),
		WorkflowID:  workflowID,
		ExecutionID: executionID,
		OccurredAt:  time.Now(),
	}
}

// NewWorkflowCreatedEvent creates a workflow created event
func NewWorkflowCreatedEvent(workflowID, workflowName string, workflowType WorkflowType, version, createdBy string) *WorkflowCreatedEvent {
	return &WorkflowCreatedEvent{
		WorkflowEvent:   NewWorkflowEvent(EventWorkflowCreated, workflowID, ""),
		WorkflowName:    workflowName,
		WorkflowType:    workflowType,
		WorkflowVersion: version,
		CreatedBy:       createdBy,
	}
}

// NewWorkflowStartedEvent creates a workflow started event
func NewWorkflowStartedEvent(workflowID, executionID, workflowName, version, startActivityID, startedBy string, inputData map[string]interface{}) *WorkflowStartedEvent {
	return &WorkflowStartedEvent{
		WorkflowEvent:   NewWorkflowEvent(EventWorkflowStarted, workflowID, executionID),
		WorkflowName:    workflowName,
		WorkflowVersion: version,
		StartActivityID: startActivityID,
		StartedBy:       startedBy,
		InputData:       inputData,
	}
}

// NewWorkflowCompletedEvent creates a workflow completed event
func NewWorkflowCompletedEvent(workflowID, executionID, workflowName string, duration time.Duration, outputData map[string]interface{}, totalActivities, completedActivities int) *WorkflowCompletedEvent {
	return &WorkflowCompletedEvent{
		WorkflowEvent:       NewWorkflowEvent(EventWorkflowCompleted, workflowID, executionID),
		WorkflowName:        workflowName,
		Duration:            duration,
		OutputData:          outputData,
		TotalActivities:     totalActivities,
		CompletedActivities: completedActivities,
	}
}

// NewWorkflowFailedEvent creates a workflow failed event
func NewWorkflowFailedEvent(workflowID, executionID, workflowName, currentActivityID, errorMsg string, duration time.Duration, retryCount int) *WorkflowFailedEvent {
	return &WorkflowFailedEvent{
		WorkflowEvent:     NewWorkflowEvent(EventWorkflowFailed, workflowID, executionID),
		WorkflowName:      workflowName,
		CurrentActivityID: currentActivityID,
		ErrorMessage:      errorMsg,
		Duration:          duration,
		RetryCount:        retryCount,
	}
}

// NewActivityStartedEvent creates an activity started event
func NewActivityStartedEvent(workflowID, executionID, activityID, activityName string, activityType ActivityType, inputData map[string]interface{}) *ActivityStartedEvent {
	return &ActivityStartedEvent{
		WorkflowEvent: NewWorkflowEvent(EventActivityStarted, workflowID, executionID),
		ActivityID:    activityID,
		ActivityName:  activityName,
		ActivityType:  activityType,
		InputData:     inputData,
	}
}

// NewActivityCompletedEvent creates an activity completed event
func NewActivityCompletedEvent(workflowID, executionID, activityID, activityName string, activityType ActivityType, duration time.Duration, outputData map[string]interface{}) *ActivityCompletedEvent {
	return &ActivityCompletedEvent{
		WorkflowEvent: NewWorkflowEvent(EventActivityCompleted, workflowID, executionID),
		ActivityID:    activityID,
		ActivityName:  activityName,
		ActivityType:  activityType,
		Duration:      duration,
		OutputData:    outputData,
	}
}

// NewActivityFailedEvent creates an activity failed event
func NewActivityFailedEvent(workflowID, executionID, activityID, activityName string, activityType ActivityType, errorMsg string, duration time.Duration, retryCount int) *ActivityFailedEvent {
	return &ActivityFailedEvent{
		WorkflowEvent: NewWorkflowEvent(EventActivityFailed, workflowID, executionID),
		ActivityID:    activityID,
		ActivityName:  activityName,
		ActivityType:  activityType,
		ErrorMessage:  errorMsg,
		Duration:      duration,
		RetryCount:    retryCount,
	}
}

// NewTransitionExecutedEvent creates a transition executed event
func NewTransitionExecutedEvent(workflowID, executionID, fromActivityID, toActivityID, condition string, conditionResult bool) *TransitionExecutedEvent {
	return &TransitionExecutedEvent{
		WorkflowEvent:   NewWorkflowEvent(EventTransitionExecuted, workflowID, executionID),
		FromActivityID:  fromActivityID,
		ToActivityID:    toActivityID,
		Condition:       condition,
		ConditionResult: conditionResult,
	}
}

func generateEventID() string {
	return "EVT-" + time.Now().Format("20060102150405")
}
