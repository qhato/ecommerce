package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/rule/domain"
)

// RuleCommandHandler handles rule-related commands
type RuleCommandHandler struct {
	ruleRepo domain.RuleRepository
	engine   *domain.RuleEngine
}

// NewRuleCommandHandler creates a new rule command handler
func NewRuleCommandHandler(
	ruleRepo domain.RuleRepository,
	engine *domain.RuleEngine,
) *RuleCommandHandler {
	return &RuleCommandHandler{
		ruleRepo: ruleRepo,
		engine:   engine,
	}
}

// HandleCreateRule handles creating a new rule
func (h *RuleCommandHandler) HandleCreateRule(ctx context.Context, cmd CreateRuleCommand) (*domain.Rule, error) {
	// Check if rule name already exists
	exists, err := h.ruleRepo.ExistsByName(ctx, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check rule name: %w", err)
	}
	if exists {
		return nil, domain.ErrRuleNameTaken
	}

	// Convert conditions
	conditions := make([]domain.Condition, len(cmd.Conditions))
	for i, condDTO := range cmd.Conditions {
		conditions[i] = domain.Condition{
			Field:    condDTO.Field,
			Operator: domain.Operator(condDTO.Operator),
			Value:    condDTO.Value,
		}
	}

	// Convert actions
	actions := make([]domain.Action, len(cmd.Actions))
	for i, actDTO := range cmd.Actions {
		actions[i] = domain.Action{
			Type:       domain.ActionType(actDTO.Type),
			Parameters: actDTO.Parameters,
		}
	}

	rule, err := domain.NewRule(
		cmd.Name,
		cmd.Description,
		domain.RuleType(cmd.Type),
		cmd.Priority,
		conditions,
		actions,
	)
	if err != nil {
		return nil, err
	}

	rule.StartDate = cmd.StartDate
	rule.EndDate = cmd.EndDate
	rule.Context = cmd.Context

	if err := h.ruleRepo.Create(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to create rule: %w", err)
	}

	return rule, nil
}

// HandleUpdateRule handles updating a rule
func (h *RuleCommandHandler) HandleUpdateRule(ctx context.Context, cmd UpdateRuleCommand) (*domain.Rule, error) {
	rule, err := h.ruleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find rule: %w", err)
	}
	if rule == nil {
		return nil, domain.ErrRuleNotFound
	}

	// Convert conditions
	conditions := make([]domain.Condition, len(cmd.Conditions))
	for i, condDTO := range cmd.Conditions {
		conditions[i] = domain.Condition{
			Field:    condDTO.Field,
			Operator: domain.Operator(condDTO.Operator),
			Value:    condDTO.Value,
		}
	}

	// Convert actions
	actions := make([]domain.Action, len(cmd.Actions))
	for i, actDTO := range cmd.Actions {
		actions[i] = domain.Action{
			Type:       domain.ActionType(actDTO.Type),
			Parameters: actDTO.Parameters,
		}
	}

	if err := rule.Update(cmd.Name, cmd.Description, cmd.Priority, conditions, actions); err != nil {
		return nil, err
	}

	rule.StartDate = cmd.StartDate
	rule.EndDate = cmd.EndDate
	rule.Context = cmd.Context

	if err := h.ruleRepo.Update(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to update rule: %w", err)
	}

	return rule, nil
}

// HandleActivateRule handles activating a rule
func (h *RuleCommandHandler) HandleActivateRule(ctx context.Context, cmd ActivateRuleCommand) (*domain.Rule, error) {
	rule, err := h.ruleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find rule: %w", err)
	}
	if rule == nil {
		return nil, domain.ErrRuleNotFound
	}

	rule.Activate()

	if err := h.ruleRepo.Update(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to activate rule: %w", err)
	}

	return rule, nil
}

// HandleDeactivateRule handles deactivating a rule
func (h *RuleCommandHandler) HandleDeactivateRule(ctx context.Context, cmd DeactivateRuleCommand) (*domain.Rule, error) {
	rule, err := h.ruleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find rule: %w", err)
	}
	if rule == nil {
		return nil, domain.ErrRuleNotFound
	}

	rule.Deactivate()

	if err := h.ruleRepo.Update(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to deactivate rule: %w", err)
	}

	return rule, nil
}

// HandleDeleteRule handles deleting a rule
func (h *RuleCommandHandler) HandleDeleteRule(ctx context.Context, cmd DeleteRuleCommand) error {
	rule, err := h.ruleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find rule: %w", err)
	}
	if rule == nil {
		return domain.ErrRuleNotFound
	}

	if err := h.ruleRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete rule: %w", err)
	}

	return nil
}

// HandleEvaluateRule handles evaluating rules
func (h *RuleCommandHandler) HandleEvaluateRule(ctx context.Context, cmd EvaluateRuleCommand) ([]domain.Action, error) {
	rules, err := h.ruleRepo.FindByType(ctx, domain.RuleType(cmd.Type), true)
	if err != nil {
		return nil, fmt.Errorf("failed to find rules: %w", err)
	}

	return h.engine.EvaluateRules(rules, cmd.Context), nil
}
