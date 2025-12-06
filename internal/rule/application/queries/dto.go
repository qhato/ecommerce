package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/rule/domain"
)

// RuleDTO represents a rule for API responses
type RuleDTO struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Priority    int                    `json:"priority"`
	IsActive    bool                   `json:"is_active"`
	Conditions  []ConditionDTO         `json:"conditions"`
	Actions     []ActionDTO            `json:"actions"`
	StartDate   *time.Time             `json:"start_date,omitempty"`
	EndDate     *time.Time             `json:"end_date,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ConditionDTO represents a rule condition for API responses
type ConditionDTO struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// ActionDTO represents a rule action for API responses
type ActionDTO struct {
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
}

// ToRuleDTO converts domain Rule to RuleDTO
func ToRuleDTO(rule *domain.Rule) *RuleDTO {
	conditions := make([]ConditionDTO, len(rule.Conditions))
	for i, cond := range rule.Conditions {
		conditions[i] = ConditionDTO{
			Field:    cond.Field,
			Operator: string(cond.Operator),
			Value:    cond.Value,
		}
	}

	actions := make([]ActionDTO, len(rule.Actions))
	for i, action := range rule.Actions {
		actions[i] = ActionDTO{
			Type:       string(action.Type),
			Parameters: action.Parameters,
		}
	}

	return &RuleDTO{
		ID:          rule.ID,
		Name:        rule.Name,
		Description: rule.Description,
		Type:        string(rule.Type),
		Priority:    rule.Priority,
		IsActive:    rule.IsActive,
		Conditions:  conditions,
		Actions:     actions,
		StartDate:   rule.StartDate,
		EndDate:     rule.EndDate,
		Context:     rule.Context,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}
}
