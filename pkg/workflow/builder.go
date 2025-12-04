package workflow

import (
	"fmt"
)

// BaseActivity provides common functionality for activities
type BaseActivity struct {
	name        string
	description string
}

// NewBaseActivity creates a new base activity
func NewBaseActivity(name, description string) BaseActivity {
	return BaseActivity{
		name:        name,
		description: description,
	}
}

// Name returns the activity name
func (a *BaseActivity) Name() string {
	return a.name
}

// Description returns the activity description
func (a *BaseActivity) Description() string {
	return a.description
}

// ActivityBuilder helps build workflows
type ActivityBuilder struct {
	activities []Activity
}

// NewActivityBuilder creates a new activity builder
func NewActivityBuilder() *ActivityBuilder {
	return &ActivityBuilder{
		activities: make([]Activity, 0),
	}
}

// Add adds an activity to the builder
func (b *ActivityBuilder) Add(activity Activity) *ActivityBuilder {
	b.activities = append(b.activities, activity)
	return b
}

// Build returns the list of activities
func (b *ActivityBuilder) Build() []Activity {
	return b.activities
}

// WorkflowBuilder helps build workflows
type WorkflowBuilder struct {
	id          string
	name        string
	description string
	activities  []Activity
	options     WorkflowOptions
}

// NewWorkflowBuilder creates a new workflow builder
func NewWorkflowBuilder(id, name string) *WorkflowBuilder {
	return &WorkflowBuilder{
		id:         id,
		name:       name,
		activities: make([]Activity, 0),
		options:    DefaultWorkflowOptions(),
	}
}

// Description sets the workflow description
func (b *WorkflowBuilder) Description(desc string) *WorkflowBuilder {
	b.description = desc
	return b
}

// AddActivity adds an activity to the workflow
func (b *WorkflowBuilder) AddActivity(activity Activity) *WorkflowBuilder {
	b.activities = append(b.activities, activity)
	return b
}

// AddActivities adds multiple activities to the workflow
func (b *WorkflowBuilder) AddActivities(activities ...Activity) *WorkflowBuilder {
	b.activities = append(b.activities, activities...)
	return b
}

// WithOptions sets workflow options
func (b *WorkflowBuilder) WithOptions(options WorkflowOptions) *WorkflowBuilder {
	b.options = options
	return b
}

// MaxRetries sets max retry attempts
func (b *WorkflowBuilder) MaxRetries(retries int) *WorkflowBuilder {
	b.options.MaxRetries = retries
	return b
}

// CompensateOnFail sets compensation flag
func (b *WorkflowBuilder) CompensateOnFail(compensate bool) *WorkflowBuilder {
	b.options.CompensateOnFail = compensate
	return b
}

// Build builds the workflow
func (b *WorkflowBuilder) Build() (*Workflow, error) {
	if b.id == "" {
		return nil, fmt.Errorf("workflow ID is required")
	}
	if b.name == "" {
		return nil, fmt.Errorf("workflow name is required")
	}
	if len(b.activities) == 0 {
		return nil, fmt.Errorf("workflow must have at least one activity")
	}

	return &Workflow{
		ID:          b.id,
		Name:        b.name,
		Description: b.description,
		Activities:  b.activities,
		Options:     b.options,
	}, nil
}

// ConditionalActivity wraps an activity with a condition
type ConditionalActivity struct {
	BaseActivity
	activity  Activity
	condition func(interface{}) bool
}

// NewConditionalActivity creates a conditional activity
func NewConditionalActivity(name string, activity Activity, condition func(interface{}) bool) *ConditionalActivity {
	return &ConditionalActivity{
		BaseActivity: NewBaseActivity(name, "Conditional: "+activity.Name()),
		activity:     activity,
		condition:    condition,
	}
}

// Execute executes the activity if condition is met
func (a *ConditionalActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	if a.condition(input) {
		return a.activity.Execute(ctx, input)
	}
	return input, nil // Skip activity
}

// Compensate compensates the activity if condition was met
func (a *ConditionalActivity) Compensate(ctx context.Context, input interface{}) error {
	if a.condition(input) {
		return a.activity.Compensate(ctx, input)
	}
	return nil
}

// ParallelActivity executes multiple activities in parallel
type ParallelActivity struct {
	BaseActivity
	activities []Activity
}

// NewParallelActivity creates a parallel activity
func NewParallelActivity(name string, activities ...Activity) *ParallelActivity {
	return &ParallelActivity{
		BaseActivity: NewBaseActivity(name, "Parallel execution"),
		activities:   activities,
	}
}

// Execute executes all activities in parallel
func (a *ParallelActivity) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	type result struct {
		output interface{}
		err    error
		index  int
	}

	results := make(chan result, len(a.activities))
	
	for i, activity := range a.activities {
		go func(idx int, act Activity) {
			output, err := act.Execute(ctx, input)
			results <- result{output: output, err: err, index: idx}
		}(i, activity)
	}

	outputs := make([]interface{}, len(a.activities))
	var firstError error

	for i := 0; i < len(a.activities); i++ {
		res := <-results
		outputs[res.index] = res.output
		if res.err != nil && firstError == nil {
			firstError = res.err
		}
	}

	if firstError != nil {
		return nil, firstError
	}

	return outputs, nil
}

// Compensate compensates all activities in parallel
func (a *ParallelActivity) Compensate(ctx context.Context, input interface{}) error {
	type result struct {
		err   error
		index int
	}

	results := make(chan result, len(a.activities))
	
	for i, activity := range a.activities {
		go func(idx int, act Activity) {
			err := act.Compensate(ctx, input)
			results <- result{err: err, index: idx}
		}(i, activity)
	}

	var firstError error
	for i := 0; i < len(a.activities); i++ {
		res := <-results
		if res.err != nil && firstError == nil {
			firstError = res.err
		}
	}

	return firstError
}
