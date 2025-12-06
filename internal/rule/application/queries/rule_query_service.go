package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/rule/domain"
)

// RuleQueryService handles rule-related queries
type RuleQueryService struct {
	ruleRepo domain.RuleRepository
}

// NewRuleQueryService creates a new rule query service
func NewRuleQueryService(ruleRepo domain.RuleRepository) *RuleQueryService {
	return &RuleQueryService{
		ruleRepo: ruleRepo,
	}
}

// GetRule retrieves a rule by ID
func (s *RuleQueryService) GetRule(ctx context.Context, id int64) (*RuleDTO, error) {
	rule, err := s.ruleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find rule: %w", err)
	}
	if rule == nil {
		return nil, domain.ErrRuleNotFound
	}

	return ToRuleDTO(rule), nil
}

// GetAllRules retrieves all rules
func (s *RuleQueryService) GetAllRules(ctx context.Context, activeOnly bool) ([]*RuleDTO, error) {
	rules, err := s.ruleRepo.FindAll(ctx, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find rules: %w", err)
	}

	dtos := make([]*RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = ToRuleDTO(rule)
	}

	return dtos, nil
}

// GetRulesByType retrieves rules by type
func (s *RuleQueryService) GetRulesByType(ctx context.Context, ruleType string, activeOnly bool) ([]*RuleDTO, error) {
	rules, err := s.ruleRepo.FindByType(ctx, domain.RuleType(ruleType), activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find rules: %w", err)
	}

	dtos := make([]*RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = ToRuleDTO(rule)
	}

	return dtos, nil
}
