package workflow

// ProcessActivity represents a single unit of work in a workflow
type ProcessActivity interface {
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

// BaseProcessActivity provides default implementations for ProcessActivity interface
type BaseProcessActivity struct {
	order    int
	beanName string
}

// NewBaseProcessActivity creates a new BaseProcessActivity
func NewBaseProcessActivity(order int, beanName string) *BaseProcessActivity {
	return &BaseProcessActivity{
		order:    order,
		beanName: beanName,
	}
}

func (a *BaseProcessActivity) ShouldExecute(ctx ProcessContext) bool {
	return true
}

func (a *BaseProcessActivity) GetOrder() int {
	return a.order
}

func (a *BaseProcessActivity) GetBeanName() string {
	return a.beanName
}

// ActivityWithRollback combines ProcessActivity and RollbackHandler
type ActivityWithRollback interface {
	ProcessActivity
	RollbackHandler
}
