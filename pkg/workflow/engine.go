package workflow

import (
	"context"
	"fmt"
	"time"
)

// Activity represents a single step in a workflow
type Activity interface {
	Execute(ctx context.Context, input interface{}) (output interface{}, err error)
	Compensate(ctx context.Context, input interface{}) error
	Name() string
}

// Status represents the status of a workflow execution
type Status string

const (
	StatusPending    Status = "PENDING"
	StatusRunning    Status = "RUNNING"
	StatusCompleted  Status = "COMPLETED"
	StatusFailed     Status = "FAILED"
	StatusCompensating Status = "COMPENSATING"
	StatusCompensated  Status = "COMPENSATED"
)

// ExecutionContext contains workflow execution state
type ExecutionContext struct {
	WorkflowID   string
	ExecutionID  string
	Status       Status
	StartTime    time.Time
	EndTime      *time.Time
	Input        interface{}
	Output       interface{}
	Error        error
	Activities   []ActivityExecution
	Metadata     map[string]interface{}
}

// ActivityExecution represents the execution of a single activity
type ActivityExecution struct {
	Name      string
	Status    Status
	StartTime time.Time
	EndTime   *time.Time
	Input     interface{}
	Output    interface{}
	Error     error
	Attempts  int
}

// Workflow defines a workflow with activities
type Workflow struct {
	ID          string
	Name        string
	Description string
	Activities  []Activity
	Options     WorkflowOptions
}

// WorkflowOptions contains workflow execution options
type WorkflowOptions struct {
	MaxRetries       int
	RetryDelay       time.Duration
	Timeout          time.Duration
	CompensateOnFail bool
}

// DefaultWorkflowOptions returns default workflow options
func DefaultWorkflowOptions() WorkflowOptions {
	return WorkflowOptions{
		MaxRetries:       3,
		RetryDelay:       1 * time.Second,
		Timeout:          5 * time.Minute,
		CompensateOnFail: true,
	}
}

// Engine is the workflow execution engine
type Engine struct {
	workflows map[string]*Workflow
	logger    Logger
	metrics   MetricsRecorder
	tracer    Tracer
}

// Logger interface for workflow logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	WithContext(ctx context.Context) Logger
}

// MetricsRecorder interface for workflow metrics
type MetricsRecorder interface {
	RecordWorkflowExecution(workflowName string, duration time.Duration, status Status)
	RecordActivityExecution(workflowName, activityName string, duration time.Duration, status Status)
	IncrementWorkflowCounter(workflowName string, status Status)
	IncrementActivityCounter(workflowName, activityName string, status Status)
}

// Tracer interface for distributed tracing
type Tracer interface {
	StartSpan(ctx context.Context, name string) (context.Context, Span)
}

// Span interface for tracing spans
type Span interface {
	End()
	SetAttribute(key string, value interface{})
	RecordError(err error)
}

// NewEngine creates a new workflow engine
func NewEngine(logger Logger, metrics MetricsRecorder, tracer Tracer) *Engine {
	return &Engine{
		workflows: make(map[string]*Workflow),
		logger:    logger,
		metrics:   metrics,
		tracer:    tracer,
	}
}

// RegisterWorkflow registers a workflow with the engine
func (e *Engine) RegisterWorkflow(workflow *Workflow) error {
	if workflow.ID == "" {
		return fmt.Errorf("workflow ID cannot be empty")
	}
	if workflow.Name == "" {
		return fmt.Errorf("workflow name cannot be empty")
	}
	if len(workflow.Activities) == 0 {
		return fmt.Errorf("workflow must have at least one activity")
	}

	e.workflows[workflow.ID] = workflow
	e.logger.Info("Workflow registered",
		"workflow_id", workflow.ID,
		"workflow_name", workflow.Name,
		"activities_count", len(workflow.Activities),
	)
	return nil
}

// Execute executes a workflow
func (e *Engine) Execute(ctx context.Context, workflowID string, input interface{}) (*ExecutionContext, error) {
	workflow, exists := e.workflows[workflowID]
	if !exists {
		return nil, fmt.Errorf("workflow not found: %s", workflowID)
	}

	// Create execution context
	execCtx := &ExecutionContext{
		WorkflowID:  workflowID,
		ExecutionID: generateExecutionID(),
		Status:      StatusRunning,
		StartTime:   time.Now(),
		Input:       input,
		Activities:  make([]ActivityExecution, 0, len(workflow.Activities)),
		Metadata:    make(map[string]interface{}),
	}

	// Start workflow span
	ctx, span := e.tracer.StartSpan(ctx, "workflow."+workflow.Name)
	defer span.End()

	span.SetAttribute("workflow.id", workflow.ID)
	span.SetAttribute("workflow.execution_id", execCtx.ExecutionID)

	e.logger.WithContext(ctx).Info("Workflow execution started",
		"workflow_id", workflowID,
		"execution_id", execCtx.ExecutionID,
	)

	// Apply timeout
	if workflow.Options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, workflow.Options.Timeout)
		defer cancel()
	}

	// Execute activities in sequence
	currentInput := input
	executedActivities := make([]Activity, 0, len(workflow.Activities))

	for _, activity := range workflow.Activities {
		activityExec := ActivityExecution{
			Name:      activity.Name(),
			Status:    StatusRunning,
			StartTime: time.Now(),
			Input:     currentInput,
			Attempts:  0,
		}

		// Execute activity with retries
		output, err := e.executeActivityWithRetry(ctx, workflow, activity, currentInput)
		
		now := time.Now()
		activityExec.EndTime = &now
		activityExec.Output = output
		activityExec.Error = err

		if err != nil {
			activityExec.Status = StatusFailed
			execCtx.Activities = append(execCtx.Activities, activityExec)

			e.logger.WithContext(ctx).Error("Activity execution failed",
				"workflow_id", workflowID,
				"execution_id", execCtx.ExecutionID,
				"activity", activity.Name(),
				"error", err,
			)

			// Compensate if enabled
			if workflow.Options.CompensateOnFail {
				e.logger.WithContext(ctx).Info("Starting compensation",
					"workflow_id", workflowID,
					"execution_id", execCtx.ExecutionID,
				)
				
				execCtx.Status = StatusCompensating
				if compErr := e.compensate(ctx, executedActivities, input); compErr != nil {
					e.logger.WithContext(ctx).Error("Compensation failed",
						"workflow_id", workflowID,
						"execution_id", execCtx.ExecutionID,
						"error", compErr,
					)
					execCtx.Status = StatusFailed
				} else {
					execCtx.Status = StatusCompensated
				}
			} else {
				execCtx.Status = StatusFailed
			}

			execCtx.Error = err
			now := time.Now()
			execCtx.EndTime = &now
			span.RecordError(err)

			e.metrics.RecordWorkflowExecution(workflow.Name, time.Since(execCtx.StartTime), execCtx.Status)
			e.metrics.IncrementWorkflowCounter(workflow.Name, execCtx.Status)

			return execCtx, err
		}

		activityExec.Status = StatusCompleted
		execCtx.Activities = append(execCtx.Activities, activityExec)
		executedActivities = append(executedActivities, activity)

		// Record activity metrics
		e.metrics.RecordActivityExecution(
			workflow.Name,
			activity.Name(),
			time.Since(activityExec.StartTime),
			activityExec.Status,
		)
		e.metrics.IncrementActivityCounter(workflow.Name, activity.Name(), activityExec.Status)

		// Output becomes input for next activity
		currentInput = output

		e.logger.WithContext(ctx).Debug("Activity completed",
			"workflow_id", workflowID,
			"execution_id", execCtx.ExecutionID,
			"activity", activity.Name(),
		)
	}

	// Workflow completed successfully
	execCtx.Status = StatusCompleted
	execCtx.Output = currentInput
	now := time.Now()
	execCtx.EndTime = &now

	e.logger.WithContext(ctx).Info("Workflow execution completed",
		"workflow_id", workflowID,
		"execution_id", execCtx.ExecutionID,
		"duration", time.Since(execCtx.StartTime),
	)

	e.metrics.RecordWorkflowExecution(workflow.Name, time.Since(execCtx.StartTime), execCtx.Status)
	e.metrics.IncrementWorkflowCounter(workflow.Name, execCtx.Status)

	return execCtx, nil
}

// executeActivityWithRetry executes an activity with retry logic
func (e *Engine) executeActivityWithRetry(ctx context.Context, workflow *Workflow, activity Activity, input interface{}) (interface{}, error) {
	var output interface{}
	var err error

	for attempt := 0; attempt <= workflow.Options.MaxRetries; attempt++ {
		if attempt > 0 {
			e.logger.WithContext(ctx).Warn("Retrying activity",
				"activity", activity.Name(),
				"attempt", attempt,
			)
			time.Sleep(workflow.Options.RetryDelay)
		}

		// Start activity span
		activityCtx, span := e.tracer.StartSpan(ctx, "activity."+activity.Name())
		span.SetAttribute("activity.name", activity.Name())
		span.SetAttribute("activity.attempt", attempt)

		output, err = activity.Execute(activityCtx, input)
		
		if err != nil {
			span.RecordError(err)
			span.End()
			
			// Don't retry context cancellation or timeout
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			
			continue
		}

		span.End()
		return output, nil
	}

	return nil, fmt.Errorf("activity %s failed after %d attempts: %w", activity.Name(), workflow.Options.MaxRetries+1, err)
}

// compensate executes compensation for all executed activities in reverse order
func (e *Engine) compensate(ctx context.Context, activities []Activity, input interface{}) error {
	// Compensate in reverse order
	for i := len(activities) - 1; i >= 0; i-- {
		activity := activities[i]
		
		e.logger.WithContext(ctx).Debug("Compensating activity", "activity", activity.Name())
		
		activityCtx, span := e.tracer.StartSpan(ctx, "compensate."+activity.Name())
		span.SetAttribute("activity.name", activity.Name())
		
		err := activity.Compensate(activityCtx, input)
		
		if err != nil {
			e.logger.WithContext(ctx).Error("Compensation failed for activity",
				"activity", activity.Name(),
				"error", err,
			)
			span.RecordError(err)
			span.End()
			return err
		}
		
		span.End()
	}

	return nil
}

// generateExecutionID generates a unique execution ID
func generateExecutionID() string {
	return fmt.Sprintf("exec-%d", time.Now().UnixNano())
}

// GetWorkflow retrieves a registered workflow
func (e *Engine) GetWorkflow(workflowID string) (*Workflow, bool) {
	workflow, exists := e.workflows[workflowID]
	return workflow, exists
}

// ListWorkflows returns all registered workflows
func (e *Engine) ListWorkflows() []*Workflow {
	workflows := make([]*Workflow, 0, len(e.workflows))
	for _, wf := range e.workflows {
		workflows = append(workflows, wf)
	}
	return workflows
}
