package commands

import "time"

// CreateRuleCommand creates a new rule
type CreateRuleCommand struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Type        string                   `json:"type"`
	Priority    int                      `json:"priority"`
	Conditions  []ConditionDTO           `json:"conditions"`
	Actions     []ActionDTO              `json:"actions"`
	StartDate   *time.Time               `json:"start_date"`
	EndDate     *time.Time               `json:"end_date"`
	Context     map[string]interface{}   `json:"context"`
}

// UpdateRuleCommand updates an existing rule
type UpdateRuleCommand struct {
	ID          int64                    `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Priority    int                      `json:"priority"`
	Conditions  []ConditionDTO           `json:"conditions"`
	Actions     []ActionDTO              `json:"actions"`
	StartDate   *time.Time               `json:"start_date"`
	EndDate     *time.Time               `json:"end_date"`
	Context     map[string]interface{}   `json:"context"`
}

// ActivateRuleCommand activates a rule
type ActivateRuleCommand struct {
	ID int64 `json:"id"`
}

// DeactivateRuleCommand deactivates a rule
type DeactivateRuleCommand struct {
	ID int64 `json:"id"`
}

// DeleteRuleCommand deletes a rule
type DeleteRuleCommand struct {
	ID int64 `json:"id"`
}

// EvaluateRuleCommand evaluates rules against context
type EvaluateRuleCommand struct {
	Type    string                 `json:"type"`
	Context map[string]interface{} `json:"context"`
}

// ConditionDTO represents a rule condition
type ConditionDTO struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// ActionDTO represents a rule action
type ActionDTO struct {
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
}
