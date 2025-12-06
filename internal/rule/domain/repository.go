package domain

import "context"

// RuleRepository defines the interface for rule persistence
type RuleRepository interface {
	Create(ctx context.Context, rule *Rule) error
	Update(ctx context.Context, rule *Rule) error
	FindByID(ctx context.Context, id int64) (*Rule, error)
	FindByType(ctx context.Context, ruleType RuleType, activeOnly bool) ([]*Rule, error)
	FindAll(ctx context.Context, activeOnly bool) ([]*Rule, error)
	Delete(ctx context.Context, id int64) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}
