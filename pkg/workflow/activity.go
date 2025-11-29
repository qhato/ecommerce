package workflow

// Activity represents a single unit of work in a workflow
type Activity interface {
	// Execute performs the activity's work
	Execute(ctx ProcessContext) error

	// ShouldExecute determines if this activity should run
	// Returns true by default, can be overridden for conditional execution
	ShouldExecute(ctx ProcessContext) bool

	// GetOrder returns the execution order (lower numbers execute first)
	GetOrder() int

	// GetBeanName returns a unique identifier for this activity
	GetBeanName() string
}

// RollbackHandler defines an interface for handling rollback logic
type RollbackHandler interface {
	// RollbackState reverts changes made by an activity
	RollbackState(ctx ProcessContext) error
}

// BaseActivity provides default implementations for Activity interface
type BaseActivity struct {
	order    int
	beanName string
}

// NewBaseActivity creates a new BaseActivity
func NewBaseActivity(order int, beanName string) *BaseActivity {
	return &BaseActivity{
		order:    order,
		beanName: beanName,
	}
}

func (a *BaseActivity) ShouldExecute(ctx ProcessContext) bool {
	return true
}

func (a *BaseActivity) GetOrder() int {
	return a.order
}

func (a *BaseActivity) GetBeanName() string {
	return a.beanName
}

// ActivityWithRollback combines Activity and RollbackHandler
type ActivityWithRollback interface {
	Activity
	RollbackHandler
}
