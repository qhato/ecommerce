package domain

import (
	"time"
)

// WorkflowType represents the type of workflow
type WorkflowType string

const (
	WorkflowTypeCheckout    WorkflowType = "CHECKOUT"
	WorkflowTypeOrderFulfillment WorkflowType = "ORDER_FULFILLMENT"
	WorkflowTypePaymentProcessing WorkflowType = "PAYMENT_PROCESSING"
	WorkflowTypeReturnProcess WorkflowType = "RETURN_PROCESS"
	WorkflowTypeCustom WorkflowType = "CUSTOM"
)

// WorkflowStatus represents the status of a workflow execution
type WorkflowStatus string

const (
	WorkflowStatusPending    WorkflowStatus = "PENDING"
	WorkflowStatusRunning    WorkflowStatus = "RUNNING"
	WorkflowStatusCompleted  WorkflowStatus = "COMPLETED"
	WorkflowStatusFailed     WorkflowStatus = "FAILED"
	WorkflowStatusCancelled  WorkflowStatus = "CANCELLED"
	WorkflowStatusSuspended  WorkflowStatus = "SUSPENDED"
)

// Workflow represents a workflow definition
// Business Logic: Define reusable business process flows
type Workflow struct {
	ID          string
	Name        string
	Description string
	Type        WorkflowType
	Version     string
	IsActive    bool
	Activities  []Activity
	Transitions []Transition
	StartActivityID string
	EndActivityIDs []string
	Metadata    map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Activity represents a single step/task in a workflow
type Activity struct {
	ID          string
	Name        string
	Description string
	Type        ActivityType
	Config      ActivityConfig
	Rollback    *RollbackConfig
	Timeout     *time.Duration
	RetryPolicy *RetryPolicy
	IsAsync     bool
	Order       int
}

// ActivityType represents the type of activity
type ActivityType string

const (
	ActivityTypeTask       ActivityType = "TASK"        // Execute a task
	ActivityTypeDecision   ActivityType = "DECISION"    // Branch based on condition
	ActivityTypeParallel   ActivityType = "PARALLEL"    // Execute multiple activities in parallel
	ActivityTypeWait       ActivityType = "WAIT"        // Wait for external event
	ActivityTypeSubWorkflow ActivityType = "SUB_WORKFLOW" // Execute another workflow
	ActivityTypeScript     ActivityType = "SCRIPT"      // Execute custom script/code
)

// ActivityConfig holds activity-specific configuration
type ActivityConfig struct {
	Handler     string                 // Handler function/service name
	InputMapping map[string]string      // Map workflow data to activity inputs
	OutputMapping map[string]string     // Map activity outputs to workflow data
	Parameters  map[string]interface{} // Static parameters
	Conditions  []Condition            // For decision activities
}

// Condition represents a conditional expression
type Condition struct {
	Expression string // Expression to evaluate (e.g., "order.total > 100")
	NextActivityID string
}

// RollbackConfig defines how to rollback an activity
type RollbackConfig struct {
	Handler    string
	Parameters map[string]interface{}
}

// RetryPolicy defines retry behavior for failed activities
type RetryPolicy struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
}

// Transition represents a state transition between activities
type Transition struct {
	FromActivityID string
	ToActivityID   string
	Condition      string // Optional condition expression
	Priority       int    // For ordering multiple transitions
}

// NewWorkflow creates a new workflow definition
func NewWorkflow(name, workflowType string) (*Workflow, error) {
	if name == "" {
		return nil, ErrWorkflowNameRequired
	}

	now := time.Now()
	return &Workflow{
		ID:          generateWorkflowID(),
		Name:        name,
		Type:        WorkflowType(workflowType),
		Version:     "1.0.0",
		IsActive:    true,
		Activities:  make([]Activity, 0),
		Transitions: make([]Transition, 0),
		EndActivityIDs: make([]string, 0),
		Metadata:    make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// AddActivity adds an activity to the workflow
func (w *Workflow) AddActivity(activity Activity) error {
	// Validate activity
	if activity.ID == "" {
		return ErrActivityIDRequired
	}
	if activity.Name == "" {
		return ErrActivityNameRequired
	}

	// Check for duplicate IDs
	for _, a := range w.Activities {
		if a.ID == activity.ID {
			return ErrActivityIDDuplicate
		}
	}

	w.Activities = append(w.Activities, activity)
	w.UpdatedAt = time.Now()
	return nil
}

// AddTransition adds a transition between activities
func (w *Workflow) AddTransition(transition Transition) error {
	// Validate activities exist
	if !w.activityExists(transition.FromActivityID) {
		return ErrActivityNotFound
	}
	if !w.activityExists(transition.ToActivityID) {
		return ErrActivityNotFound
	}

	w.Transitions = append(w.Transitions, transition)
	w.UpdatedAt = time.Now()
	return nil
}

// SetStartActivity sets the starting activity
func (w *Workflow) SetStartActivity(activityID string) error {
	if !w.activityExists(activityID) {
		return ErrActivityNotFound
	}
	w.StartActivityID = activityID
	w.UpdatedAt = time.Now()
	return nil
}

// AddEndActivity adds an end activity
func (w *Workflow) AddEndActivity(activityID string) error {
	if !w.activityExists(activityID) {
		return ErrActivityNotFound
	}
	w.EndActivityIDs = append(w.EndActivityIDs, activityID)
	w.UpdatedAt = time.Now()
	return nil
}

// GetActivity retrieves an activity by ID
func (w *Workflow) GetActivity(activityID string) (*Activity, error) {
	for _, a := range w.Activities {
		if a.ID == activityID {
			return &a, nil
		}
	}
	return nil, ErrActivityNotFound
}

// GetNextActivities gets the next activities based on current activity
func (w *Workflow) GetNextActivities(currentActivityID string) []Activity {
	nextActivities := make([]Activity, 0)

	for _, t := range w.Transitions {
		if t.FromActivityID == currentActivityID {
			activity, err := w.GetActivity(t.ToActivityID)
			if err == nil {
				nextActivities = append(nextActivities, *activity)
			}
		}
	}

	return nextActivities
}

// IsEndActivity checks if an activity is an end activity
func (w *Workflow) IsEndActivity(activityID string) bool {
	for _, endID := range w.EndActivityIDs {
		if endID == activityID {
			return true
		}
	}
	return false
}

// Validate validates the workflow definition
func (w *Workflow) Validate() error {
	if w.StartActivityID == "" {
		return ErrStartActivityNotSet
	}

	if len(w.EndActivityIDs) == 0 {
		return ErrEndActivityNotSet
	}

	if len(w.Activities) == 0 {
		return ErrNoActivitiesDefined
	}

	// Validate start activity exists
	if !w.activityExists(w.StartActivityID) {
		return ErrStartActivityNotFound
	}

	// Validate end activities exist
	for _, endID := range w.EndActivityIDs {
		if !w.activityExists(endID) {
			return ErrEndActivityNotFound
		}
	}

	return nil
}

// Activate activates the workflow
func (w *Workflow) Activate() {
	w.IsActive = true
	w.UpdatedAt = time.Now()
}

// Deactivate deactivates the workflow
func (w *Workflow) Deactivate() {
	w.IsActive = false
	w.UpdatedAt = time.Now()
}

// Private helper methods

func (w *Workflow) activityExists(activityID string) bool {
	for _, a := range w.Activities {
		if a.ID == activityID {
			return true
		}
	}
	return false
}

func generateWorkflowID() string {
	return "WF-" + time.Now().Format("20060102150405")
}
