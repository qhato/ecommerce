package domain

import "errors"

// Workflow Definition Errors
var (
	ErrWorkflowNameRequired      = errors.New("workflow name is required")
	ErrWorkflowNotFound          = errors.New("workflow not found")
	ErrWorkflowAlreadyExists     = errors.New("workflow already exists")
	ErrWorkflowInactive          = errors.New("workflow is inactive")
	ErrWorkflowInvalidVersion    = errors.New("invalid workflow version")
)

// Activity Errors
var (
	ErrActivityIDRequired        = errors.New("activity ID is required")
	ErrActivityNameRequired      = errors.New("activity name is required")
	ErrActivityIDDuplicate       = errors.New("activity ID already exists")
	ErrActivityNotFound          = errors.New("activity not found")
	ErrActivityConfigInvalid     = errors.New("activity configuration is invalid")
	ErrActivityHandlerNotFound   = errors.New("activity handler not found")
	ErrActivityTimeout           = errors.New("activity execution timed out")
	ErrActivityExecutionNotFound = errors.New("activity execution not found")
)

// Transition Errors
var (
	ErrTransitionInvalid         = errors.New("invalid transition")
	ErrNoValidTransition         = errors.New("no valid transition found")
)

// Workflow Structure Errors
var (
	ErrStartActivityNotSet       = errors.New("start activity not set")
	ErrStartActivityNotFound     = errors.New("start activity not found")
	ErrEndActivityNotSet         = errors.New("end activity not set")
	ErrEndActivityNotFound       = errors.New("end activity not found")
	ErrNoActivitiesDefined       = errors.New("no activities defined in workflow")
	ErrCircularDependency        = errors.New("circular dependency detected in workflow")
)

// Workflow Execution Errors
var (
	ErrWorkflowExecutionNotFound    = errors.New("workflow execution not found")
	ErrWorkflowNotRunning           = errors.New("workflow is not in running state")
	ErrWorkflowNotInPendingState    = errors.New("workflow is not in pending state")
	ErrWorkflowNotSuspended         = errors.New("workflow is not suspended")
	ErrWorkflowAlreadyFinished      = errors.New("workflow has already finished")
	ErrWorkflowAlreadyCompleted     = errors.New("workflow is already completed")
	ErrWorkflowAlreadyCancelled     = errors.New("workflow is already cancelled")
	ErrWorkflowAlreadyRunning       = errors.New("workflow is already running")
)

// Workflow Context Errors
var (
	ErrContextKeyNotFound        = errors.New("context key not found")
	ErrContextValueInvalid       = errors.New("context value is invalid")
)

// Retry and Rollback Errors
var (
	ErrMaxRetriesExceeded        = errors.New("maximum retry attempts exceeded")
	ErrRollbackFailed            = errors.New("rollback operation failed")
	ErrRollbackNotDefined        = errors.New("rollback not defined for activity")
)

// Condition and Decision Errors
var (
	ErrConditionEvaluationFailed = errors.New("condition evaluation failed")
	ErrNoMatchingCondition       = errors.New("no matching condition found")
	ErrInvalidExpression         = errors.New("invalid expression")
)

// Input/Output Mapping Errors
var (
	ErrInputMappingFailed        = errors.New("input mapping failed")
	ErrOutputMappingFailed       = errors.New("output mapping failed")
	ErrRequiredInputMissing      = errors.New("required input parameter missing")
)

// Repository Errors
var (
	ErrWorkflowRepositoryFailure      = errors.New("workflow repository operation failed")
	ErrWorkflowExecutionRepositoryFailure = errors.New("workflow execution repository operation failed")
)
