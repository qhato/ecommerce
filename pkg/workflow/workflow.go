package workflow

import (
	"context"
	"fmt"
	"sort"

	"github.com/qhato/ecommerce/pkg/logger"
)

// ProcessWorkflow represents a sequence of activities to be executed
type ProcessWorkflow interface {
	// Execute runs all activities in order
	Execute(ctx context.Context, seedData interface{}) (ProcessContext, error)

	// AddActivity adds an activity to the workflow
	AddActivity(activity ProcessActivity)

	// GetActivities returns all activities in execution order
	GetActivities() []ProcessActivity

	// GetName returns the workflow name
	GetName() string
}

// DefaultWorkflow is the default implementation of Workflow
type DefaultWorkflow struct {
	name       string
	activities []Activity
	log        logger.Logger
}

// NewWorkflow creates a new Workflow
func NewWorkflow(name string, log logger.Logger) Workflow {
	return &DefaultWorkflow{
		name:       name,
		activities: make([]Activity, 0),
		log:        log,
	}
}

// AddActivity adds an activity to the workflow
func (w *DefaultWorkflow) AddActivity(activity Activity) {
	w.activities = append(w.activities, activity)
	// Sort activities by order
	sort.Slice(w.activities, func(i, j int) bool {
		return w.activities[i].GetOrder() < w.activities[j].GetOrder()
	})
}

// GetActivities returns all activities
func (w *DefaultWorkflow) GetActivities() []Activity {
	return w.activities
}

// GetName returns the workflow name
func (w *DefaultWorkflow) GetName() string {
	return w.name
}

// Execute runs all activities in order
func (w *DefaultWorkflow) Execute(ctx context.Context, seedData interface{}) (ProcessContext, error) {
	processCtx := NewProcessContext(ctx, seedData)

	w.log.Info(fmt.Sprintf("Starting workflow: %s with %d activities", w.name, len(w.activities)))

	executedActivities := make([]Activity, 0)

	// Execute activities in order
	for i, activity := range w.activities {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			w.log.Warn(fmt.Sprintf("Workflow %s cancelled by context", w.name))
			return processCtx, ctx.Err()
		default:
		}

		// Check if process should stop
		if processCtx.IsStopped() {
			w.log.Info(fmt.Sprintf("Workflow %s stopped at activity %d", w.name, i))
			break
		}

		// Check if activity should execute
		if !activity.ShouldExecute(processCtx) {
			w.log.Debug(fmt.Sprintf("Skipping activity %s (order: %d)", activity.GetBeanName(), activity.GetOrder()))
			continue
		}

		w.log.Debug(fmt.Sprintf("Executing activity %s (order: %d)", activity.GetBeanName(), activity.GetOrder()))

		// Execute activity
		err := activity.Execute(processCtx)
		if err != nil {
			w.log.Error(fmt.Sprintf("Activity %s failed: %v", activity.GetBeanName(), err))
			processCtx.SetError(err)

			// Rollback executed activities
			if rollbackErr := w.rollback(processCtx, executedActivities); rollbackErr != nil {
				w.log.Error(fmt.Sprintf("Rollback failed: %v", rollbackErr))
				return processCtx, fmt.Errorf("activity failed: %w, rollback failed: %v", err, rollbackErr)
			}

			return processCtx, fmt.Errorf("workflow %s failed at activity %s: %w", w.name, activity.GetBeanName(), err)
		}

		executedActivities = append(executedActivities, activity)
		w.log.Debug(fmt.Sprintf("Activity %s completed successfully", activity.GetBeanName()))
	}

	w.log.Info(fmt.Sprintf("Workflow %s completed successfully", w.name))
	return processCtx, nil
}

// rollback reverts changes made by executed activities
func (w *DefaultWorkflow) rollback(ctx ProcessContext, executedActivities []Activity) error {
	w.log.Info(fmt.Sprintf("Rolling back %d activities", len(executedActivities)))

	// Rollback in reverse order
	for i := len(executedActivities) - 1; i >= 0; i-- {
		activity := executedActivities[i]

		// Check if activity supports rollback
		if rollbackActivity, ok := activity.(RollbackHandler); ok {
			w.log.Debug(fmt.Sprintf("Rolling back activity %s", activity.GetBeanName()))

			if err := rollbackActivity.RollbackState(ctx); err != nil {
				w.log.Error(fmt.Sprintf("Failed to rollback activity %s: %v", activity.GetBeanName(), err))
				return fmt.Errorf("rollback failed for activity %s: %w", activity.GetBeanName(), err)
			}

			w.log.Debug(fmt.Sprintf("Activity %s rolled back successfully", activity.GetBeanName()))
		}
	}

	return nil
}

// WorkflowRegistry manages multiple workflows
type WorkflowRegistry struct {
	workflows map[string]Workflow
	log       logger.Logger
}

// NewWorkflowRegistry creates a new WorkflowRegistry
func NewWorkflowRegistry(log logger.Logger) *WorkflowRegistry {
	return &WorkflowRegistry{
		workflows: make(map[string]Workflow),
		log:       log,
	}
}

// Register adds a workflow to the registry
func (r *WorkflowRegistry) Register(name string, workflow Workflow) {
	r.workflows[name] = workflow
	r.log.Info(fmt.Sprintf("Registered workflow: %s", name))
}

// Get retrieves a workflow by name
func (r *WorkflowRegistry) Get(name string) (Workflow, error) {
	workflow, ok := r.workflows[name]
	if !ok {
		return nil, fmt.Errorf("workflow not found: %s", name)
	}
	return workflow, nil
}

// Execute runs a workflow by name
func (r *WorkflowRegistry) Execute(ctx context.Context, name string, seedData interface{}) (ProcessContext, error) {
	workflow, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	return workflow.Execute(ctx, seedData)
}
