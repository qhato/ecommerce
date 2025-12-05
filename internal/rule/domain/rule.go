package domain

import "time"

// RuleType represents the type of rule
type RuleType string

const (
	RuleTypePrice     RuleType = "PRICE"
	RuleTypePromotion RuleType = "PROMOTION"
	RuleTypeInventory RuleType = "INVENTORY"
	RuleTypeTax       RuleType = "TAX"
	RuleTypeShipping  RuleType = "SHIPPING"
	RuleTypeCustom    RuleType = "CUSTOM"
)

// RuleStatus represents the status of a rule
type RuleStatus string

const (
	RuleStatusActive   RuleStatus = "ACTIVE"
	RuleStatusInactive RuleStatus = "INACTIVE"
	RuleStatusExpired  RuleStatus = "EXPIRED"
)

// Rule represents a business rule
type Rule struct {
	ID          int64
	Name        string
	Description string
	Type        RuleType
	Status      RuleStatus
	Priority    int
	Conditions  []Condition
	Actions     []Action
	StartDate   *time.Time
	EndDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Condition represents a rule condition
type Condition struct {
	ID         int64
	RuleID     int64
	Field      string // e.g., "order.total", "customer.type", "product.category"
	Operator   string // e.g., "EQUALS", "GREATER_THAN", "CONTAINS"
	Value      string
	LogicOperator string // "AND", "OR"
	SortOrder  int
}

// Action represents a rule action
type Action struct {
	ID         int64
	RuleID     int64
	ActionType string // e.g., "APPLY_DISCOUNT", "SET_PRICE", "SEND_EMAIL"
	Parameters map[string]interface{}
	SortOrder  int
}

// NewRule creates a new rule
func NewRule(name string, ruleType RuleType, priority int) (*Rule, error) {
	if name == "" {
		return nil, ErrRuleNameRequired
	}

	now := time.Now()
	return &Rule{
		Name:       name,
		Type:       ruleType,
		Status:     RuleStatusActive,
		Priority:   priority,
		Conditions: make([]Condition, 0),
		Actions:    make([]Action, 0),
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// Activate activates the rule
func (r *Rule) Activate() {
	r.Status = RuleStatusActive
	r.UpdatedAt = time.Now()
}

// Deactivate deactivates the rule
func (r *Rule) Deactivate() {
	r.Status = RuleStatusInactive
	r.UpdatedAt = time.Now()
}

// IsActive checks if the rule is currently active
func (r *Rule) IsActive() bool {
	if r.Status != RuleStatusActive {
		return false
	}

	now := time.Now()

	if r.StartDate != nil && now.Before(*r.StartDate) {
		return false
	}

	if r.EndDate != nil && now.After(*r.EndDate) {
		return false
	}

	return true
}

// AddCondition adds a condition to the rule
func (r *Rule) AddCondition(condition Condition) {
	r.Conditions = append(r.Conditions, condition)
	r.UpdatedAt = time.Now()
}

// AddAction adds an action to the rule
func (r *Rule) AddAction(action Action) {
	r.Actions = append(r.Actions, action)
	r.UpdatedAt = time.Now()
}
