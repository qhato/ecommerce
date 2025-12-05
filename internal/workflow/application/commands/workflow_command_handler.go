package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/workflow/domain"
)

// WorkflowCommandHandler handles workflow-related commands
type WorkflowCommandHandler struct {
	workflowRepo   domain.WorkflowRepository
	executionRepo  domain.WorkflowExecutionRepository
	engineService  WorkflowEngineService
}

// WorkflowEngineService defines the interface for workflow engine operations
type WorkflowEngineService interface {
	ExecuteActivity(ctx context.Context, execution *domain.WorkflowExecution, activity *domain.Activity) error
	EvaluateCondition(ctx context.Context, execution *domain.WorkflowExecution, condition string) (bool, error)
	ProcessTransition(ctx context.Context, execution *domain.WorkflowExecution, workflow *domain.Workflow) error
}

// NewWorkflowCommandHandler creates a new workflow command handler
func NewWorkflowCommandHandler(
	workflowRepo domain.WorkflowRepository,
	executionRepo domain.WorkflowExecutionRepository,
	engineService WorkflowEngineService,
) *WorkflowCommandHandler {
	return &WorkflowCommandHandler{
		workflowRepo:  workflowRepo,
		executionRepo: executionRepo,
		engineService: engineService,
	}
}

// HandleCreateWorkflow handles creating a new workflow definition
func (h *WorkflowCommandHandler) HandleCreateWorkflow(ctx context.Context, cmd CreateWorkflowCommand) (*domain.Workflow, error) {
	// Check if workflow with the same name already exists
	exists, err := h.workflowRepo.ExistsByName(ctx, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check workflow existence: %w", err)
	}
	if exists {
		return nil, domain.ErrWorkflowAlreadyExists
	}

	// Create workflow
	workflow, err := domain.NewWorkflow(cmd.Name, cmd.Type)
	if err != nil {
		return nil, err
	}

	workflow.Description = cmd.Description
	workflow.Metadata = cmd.Metadata

	// Add activities
	for _, activityDTO := range cmd.Activities {
		activity := activityDTO.ToActivity()
		if err := workflow.AddActivity(activity); err != nil {
			return nil, fmt.Errorf("failed to add activity %s: %w", activity.ID, err)
		}
	}

	// Add transitions
	for _, transitionDTO := range cmd.Transitions {
		transition := transitionDTO.ToTransition()
		if err := workflow.AddTransition(transition); err != nil {
			return nil, fmt.Errorf("failed to add transition: %w", err)
		}
	}

	// Set start activity
	if cmd.StartActivityID != "" {
		if err := workflow.SetStartActivity(cmd.StartActivityID); err != nil {
			return nil, fmt.Errorf("failed to set start activity: %w", err)
		}
	}

	// Add end activities
	for _, endActivityID := range cmd.EndActivityIDs {
		if err := workflow.AddEndActivity(endActivityID); err != nil {
			return nil, fmt.Errorf("failed to add end activity %s: %w", endActivityID, err)
		}
	}

	// Validate workflow
	if err := workflow.Validate(); err != nil {
		return nil, fmt.Errorf("workflow validation failed: %w", err)
	}

	// Save to repository
	if err := h.workflowRepo.Create(ctx, workflow); err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	return workflow, nil
}

// HandleUpdateWorkflow handles updating an existing workflow definition
func (h *WorkflowCommandHandler) HandleUpdateWorkflow(ctx context.Context, cmd UpdateWorkflowCommand) (*domain.Workflow, error) {
	// Find existing workflow
	workflow, err := h.workflowRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find workflow: %w", err)
	}
	if workflow == nil {
		return nil, domain.ErrWorkflowNotFound
	}

	// Update basic fields
	workflow.Name = cmd.Name
	workflow.Description = cmd.Description
	workflow.Metadata = cmd.Metadata

	// Replace activities
	workflow.Activities = make([]domain.Activity, 0)
	for _, activityDTO := range cmd.Activities {
		activity := activityDTO.ToActivity()
		if err := workflow.AddActivity(activity); err != nil {
			return nil, fmt.Errorf("failed to add activity %s: %w", activity.ID, err)
		}
	}

	// Replace transitions
	workflow.Transitions = make([]domain.Transition, 0)
	for _, transitionDTO := range cmd.Transitions {
		transition := transitionDTO.ToTransition()
		if err := workflow.AddTransition(transition); err != nil {
			return nil, fmt.Errorf("failed to add transition: %w", err)
		}
	}

	// Update start activity
	if cmd.StartActivityID != "" {
		if err := workflow.SetStartActivity(cmd.StartActivityID); err != nil {
			return nil, fmt.Errorf("failed to set start activity: %w", err)
		}
	}

	// Update end activities
	workflow.EndActivityIDs = make([]string, 0)
	for _, endActivityID := range cmd.EndActivityIDs {
		if err := workflow.AddEndActivity(endActivityID); err != nil {
			return nil, fmt.Errorf("failed to add end activity %s: %w", endActivityID, err)
		}
	}

	// Validate workflow
	if err := workflow.Validate(); err != nil {
		return nil, fmt.Errorf("workflow validation failed: %w", err)
	}

	// Save to repository
	if err := h.workflowRepo.Update(ctx, workflow); err != nil {
		return nil, fmt.Errorf("failed to update workflow: %w", err)
	}

	return workflow, nil
}

// HandleActivateWorkflow handles activating a workflow
func (h *WorkflowCommandHandler) HandleActivateWorkflow(ctx context.Context, cmd ActivateWorkflowCommand) (*domain.Workflow, error) {
	workflow, err := h.workflowRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find workflow: %w", err)
	}
	if workflow == nil {
		return nil, domain.ErrWorkflowNotFound
	}

	workflow.Activate()

	if err := h.workflowRepo.Update(ctx, workflow); err != nil {
		return nil, fmt.Errorf("failed to activate workflow: %w", err)
	}

	return workflow, nil
}

// HandleDeactivateWorkflow handles deactivating a workflow
func (h *WorkflowCommandHandler) HandleDeactivateWorkflow(ctx context.Context, cmd DeactivateWorkflowCommand) (*domain.Workflow, error) {
	workflow, err := h.workflowRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find workflow: %w", err)
	}
	if workflow == nil {
		return nil, domain.ErrWorkflowNotFound
	}

	workflow.Deactivate()

	if err := h.workflowRepo.Update(ctx, workflow); err != nil {
		return nil, fmt.Errorf("failed to deactivate workflow: %w", err)
	}

	return workflow, nil
}

// HandleDeleteWorkflow handles deleting a workflow
func (h *WorkflowCommandHandler) HandleDeleteWorkflow(ctx context.Context, cmd DeleteWorkflowCommand) error {
	workflow, err := h.workflowRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find workflow: %w", err)
	}
	if workflow == nil {
		return domain.ErrWorkflowNotFound
	}

	if err := h.workflowRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete workflow: %w", err)
	}

	return nil
}

// HandleStartWorkflowExecution handles starting a new workflow execution
func (h *WorkflowCommandHandler) HandleStartWorkflowExecution(ctx context.Context, cmd StartWorkflowExecutionCommand) (*domain.WorkflowExecution, error) {
	// Find workflow definition
	workflow, err := h.workflowRepo.FindByID(ctx, cmd.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to find workflow: %w", err)
	}
	if workflow == nil {
		return nil, domain.ErrWorkflowNotFound
	}
	if !workflow.IsActive {
		return nil, domain.ErrWorkflowInactive
	}

	// Validate workflow
	if err := workflow.Validate(); err != nil {
		return nil, fmt.Errorf("workflow validation failed: %w", err)
	}

	// Create execution
	execution := domain.NewWorkflowExecution(
		workflow.ID,
		workflow.Version,
		cmd.StartedBy,
		cmd.InputData,
	)

	if cmd.EntityType != "" && cmd.EntityID != "" {
		execution.SetEntityReference(cmd.EntityType, cmd.EntityID)
	}

	// Start execution
	if err := execution.Start(workflow.StartActivityID); err != nil {
		return nil, fmt.Errorf("failed to start execution: %w", err)
	}

	// Save execution
	if err := h.executionRepo.Create(ctx, execution); err != nil {
		return nil, fmt.Errorf("failed to create execution: %w", err)
	}

	// Start first activity (async)
	go func() {
		startActivity, _ := workflow.GetActivity(workflow.StartActivityID)
		if startActivity != nil {
			h.engineService.ExecuteActivity(context.Background(), execution, startActivity)
		}
	}()

	return execution, nil
}

// HandleCompleteActivity handles completing an activity
func (h *WorkflowCommandHandler) HandleCompleteActivity(ctx context.Context, cmd CompleteActivityCommand) (*domain.WorkflowExecution, error) {
	execution, err := h.executionRepo.FindByID(ctx, cmd.ExecutionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find execution: %w", err)
	}
	if execution == nil {
		return nil, domain.ErrWorkflowExecutionNotFound
	}

	// Complete the activity
	if err := execution.CompleteActivity(cmd.ActivityID, cmd.OutputData); err != nil {
		return nil, fmt.Errorf("failed to complete activity: %w", err)
	}

	// Save execution
	if err := h.executionRepo.Update(ctx, execution); err != nil {
		return nil, fmt.Errorf("failed to update execution: %w", err)
	}

	// Get workflow definition
	workflow, err := h.workflowRepo.FindByID(ctx, execution.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to find workflow: %w", err)
	}

	// Check if this is an end activity
	if workflow.IsEndActivity(cmd.ActivityID) {
		// Complete the workflow
		if err := execution.Complete(execution.OutputData); err != nil {
			return nil, fmt.Errorf("failed to complete workflow: %w", err)
		}
		if err := h.executionRepo.Update(ctx, execution); err != nil {
			return nil, fmt.Errorf("failed to update execution: %w", err)
		}
	} else {
		// Process transition to next activity (async)
		go h.engineService.ProcessTransition(context.Background(), execution, workflow)
	}

	return execution, nil
}

// HandleFailActivity handles failing an activity
func (h *WorkflowCommandHandler) HandleFailActivity(ctx context.Context, cmd FailActivityCommand) (*domain.WorkflowExecution, error) {
	execution, err := h.executionRepo.FindByID(ctx, cmd.ExecutionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find execution: %w", err)
	}
	if execution == nil {
		return nil, domain.ErrWorkflowExecutionNotFound
	}

	// Fail the activity
	if err := execution.FailActivity(cmd.ActivityID, cmd.ErrorMessage); err != nil {
		return nil, fmt.Errorf("failed to fail activity: %w", err)
	}

	// For now, fail the entire workflow when an activity fails
	// In a more sophisticated implementation, we'd handle retries here
	if err := execution.Fail(cmd.ErrorMessage); err != nil {
		return nil, fmt.Errorf("failed to fail workflow: %w", err)
	}

	// Save execution
	if err := h.executionRepo.Update(ctx, execution); err != nil {
		return nil, fmt.Errorf("failed to update execution: %w", err)
	}

	return execution, nil
}

// HandleSuspendWorkflow handles suspending a workflow execution
func (h *WorkflowCommandHandler) HandleSuspendWorkflow(ctx context.Context, cmd SuspendWorkflowCommand) (*domain.WorkflowExecution, error) {
	execution, err := h.executionRepo.FindByID(ctx, cmd.ExecutionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find execution: %w", err)
	}
	if execution == nil {
		return nil, domain.ErrWorkflowExecutionNotFound
	}

	if err := execution.Suspend(); err != nil {
		return nil, fmt.Errorf("failed to suspend workflow: %w", err)
	}

	if err := h.executionRepo.Update(ctx, execution); err != nil {
		return nil, fmt.Errorf("failed to update execution: %w", err)
	}

	return execution, nil
}

// HandleResumeWorkflow handles resuming a suspended workflow execution
func (h *WorkflowCommandHandler) HandleResumeWorkflow(ctx context.Context, cmd ResumeWorkflowCommand) (*domain.WorkflowExecution, error) {
	execution, err := h.executionRepo.FindByID(ctx, cmd.ExecutionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find execution: %w", err)
	}
	if execution == nil {
		return nil, domain.ErrWorkflowExecutionNotFound
	}

	if err := execution.Resume(); err != nil {
		return nil, fmt.Errorf("failed to resume workflow: %w", err)
	}

	if err := h.executionRepo.Update(ctx, execution); err != nil {
		return nil, fmt.Errorf("failed to update execution: %w", err)
	}

	// Resume workflow processing (async)
	workflow, _ := h.workflowRepo.FindByID(ctx, execution.WorkflowID)
	if workflow != nil {
		go h.engineService.ProcessTransition(context.Background(), execution, workflow)
	}

	return execution, nil
}

// HandleCancelWorkflow handles cancelling a workflow execution
func (h *WorkflowCommandHandler) HandleCancelWorkflow(ctx context.Context, cmd CancelWorkflowCommand) (*domain.WorkflowExecution, error) {
	execution, err := h.executionRepo.FindByID(ctx, cmd.ExecutionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find execution: %w", err)
	}
	if execution == nil {
		return nil, domain.ErrWorkflowExecutionNotFound
	}

	if err := execution.Cancel(); err != nil {
		return nil, fmt.Errorf("failed to cancel workflow: %w", err)
	}

	if err := h.executionRepo.Update(ctx, execution); err != nil {
		return nil, fmt.Errorf("failed to update execution: %w", err)
	}

	return execution, nil
}

// HandleSetWorkflowContext handles setting a value in workflow context
func (h *WorkflowCommandHandler) HandleSetWorkflowContext(ctx context.Context, cmd SetWorkflowContextCommand) (*domain.WorkflowExecution, error) {
	execution, err := h.executionRepo.FindByID(ctx, cmd.ExecutionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find execution: %w", err)
	}
	if execution == nil {
		return nil, domain.ErrWorkflowExecutionNotFound
	}

	execution.SetContext(cmd.Key, cmd.Value)

	if err := h.executionRepo.Update(ctx, execution); err != nil {
		return nil, fmt.Errorf("failed to update execution: %w", err)
	}

	return execution, nil
}
