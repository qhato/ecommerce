package domain

import "context"

// RuleRepository defines the interface for rule persistence
type RuleRepository interface {
	Create(ctx context.Context, rule *Rule) error
	Update(ctx context.Context, rule *Rule) error
	FindByID(ctx context.Context, id int64) (*Rule, error)
	FindByType(ctx context.Context, ruleType RuleType, status RuleStatus) ([]*Rule, error)
	FindActive(ctx context.Context) ([]*Rule, error)
	FindAll(ctx context.Context, status RuleStatus) ([]*Rule, error)
	Delete(ctx context.Context, id int64) error
	GetConditions(ctx context.Context, ruleID int64) ([]Condition, error)
	GetActions(ctx context.Context, ruleID int64) ([]Action, error)
}
