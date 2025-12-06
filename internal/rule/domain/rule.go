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
	Priority    int
	IsActive    bool
	Conditions  []Condition
	Actions     []Action
	StartDate   *time.Time
	EndDate     *time.Time
	Context     map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Operator represents a comparison operator
type Operator string

const (
	OperatorEquals      Operator = "EQUALS"
	OperatorNotEquals   Operator = "NOT_EQUALS"
	OperatorGreaterThan Operator = "GREATER_THAN"
	OperatorLessThan    Operator = "LESS_THAN"
	OperatorContains    Operator = "CONTAINS"
)

// ActionType represents an action type
type ActionType string

const (
	ActionTypeDiscount ActionType = "APPLY_DISCOUNT"
	ActionTypeSetPrice ActionType = "SET_PRICE"
	ActionTypeSendEmail ActionType = "SEND_EMAIL"
)

// Condition represents a rule condition
type Condition struct {
	Field    string   // e.g., "order.total", "customer.type"
	Operator Operator // e.g., EQUALS, GREATER_THAN
	Value    string
}

// Action represents a rule action
type Action struct {
	Type       ActionType             // e.g., APPLY_DISCOUNT, SET_PRICE
	Parameters map[string]interface{}
}

// NewRule creates a new rule
func NewRule(name, description string, ruleType RuleType, priority int, conditions []Condition, actions []Action) (*Rule, error) {
	if name == "" {
		return nil, ErrRuleNameRequired
	}

	now := time.Now()
	return &Rule{
		Name:        name,
		Description: description,
		Type:        ruleType,
		Priority:    priority,
		IsActive:    true,
		Conditions:  conditions,
		Actions:     actions,
		Context:     make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update updates rule information
func (r *Rule) Update(name, description string, priority int, conditions []Condition, actions []Action) error {
	if name == "" {
		return ErrRuleNameRequired
	}
	r.Name = name
	r.Description = description
	r.Priority = priority
	r.Conditions = conditions
	r.Actions = actions
	r.UpdatedAt = time.Now()
	return nil
}

// Activate activates the rule
func (r *Rule) Activate() {
	r.IsActive = true
	r.UpdatedAt = time.Now()
}

// Deactivate deactivates the rule
func (r *Rule) Deactivate() {
	r.IsActive = false
	r.UpdatedAt = time.Now()
}

// IsActiveNow checks if the rule is currently active based on dates
func (r *Rule) IsActiveNow() bool {
	if !r.IsActive {
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
