package domain

import (
	"errors"
	"time"
)

// CustomerSegment represents a customer segmentation for targeted promotions
type CustomerSegment struct {
	ID          int64
	Name        string
	Description string
	Rules       []SegmentRule
	IsActive    bool
	CustomerCount int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SegmentRule represents a single rule in a segment
type SegmentRule struct {
	ID        int64
	SegmentID int64
	Field     string      // total_spent, order_count, last_order_date, avg_order_value, customer_since, etc.
	Operator  string      // >, <, >=, <=, =, IN, BETWEEN
	Value     interface{} // Can be number, string, date, array
	LogicOp   string      // AND, OR (for combining with next rule)
}

// SegmentMembership represents a customer's membership in a segment
type SegmentMembership struct {
	ID         int64
	SegmentID  int64
	CustomerID int64
	AddedAt    time.Time
	LastChecked time.Time
}

// NewCustomerSegment creates a new customer segment
func NewCustomerSegment(name, description string, rules []SegmentRule) (*CustomerSegment, error) {
	if name == "" {
		return nil, errors.New("segment name is required")
	}
	if len(rules) == 0 {
		return nil, errors.New("at least one rule is required")
	}

	now := time.Now()
	return &CustomerSegment{
		Name:          name,
		Description:   description,
		Rules:         rules,
		IsActive:      true,
		CustomerCount: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// Activate activates the segment
func (s *CustomerSegment) Activate() {
	s.IsActive = true
	s.UpdatedAt = time.Now()
}

// Deactivate deactivates the segment
func (s *CustomerSegment) Deactivate() {
	s.IsActive = false
	s.UpdatedAt = time.Now()
}

// UpdateRules updates the segment rules
func (s *CustomerSegment) UpdateRules(rules []SegmentRule) error {
	if len(rules) == 0 {
		return errors.New("at least one rule is required")
	}
	s.Rules = rules
	s.UpdatedAt = time.Now()
	return nil
}

// CustomerData represents customer data for segment evaluation
type CustomerData struct {
	CustomerID       int64
	TotalSpent       float64
	OrderCount       int
	AvgOrderValue    float64
	LastOrderDate    *time.Time
	CustomerSince    time.Time
	FavoriteCategory *string
	CustomerTier     string
	LifetimePoints   int64
}

// EvaluateCustomer evaluates if a customer matches this segment
func (s *CustomerSegment) EvaluateCustomer(customer CustomerData) bool {
	if len(s.Rules) == 0 {
		return false
	}

	// Simple evaluation - all rules must match (AND logic)
	// In a real implementation, you would support complex logic with AND/OR
	for _, rule := range s.Rules {
		if !evaluateRule(rule, customer) {
			return false
		}
	}
	return true
}

// evaluateRule evaluates a single rule against customer data
func evaluateRule(rule SegmentRule, customer CustomerData) bool {
	switch rule.Field {
	case "total_spent":
		return compareFloat(customer.TotalSpent, rule.Operator, rule.Value)
	case "order_count":
		return compareInt(customer.OrderCount, rule.Operator, rule.Value)
	case "avg_order_value":
		return compareFloat(customer.AvgOrderValue, rule.Operator, rule.Value)
	case "customer_tier":
		return compareString(customer.CustomerTier, rule.Operator, rule.Value)
	default:
		return false
	}
}

func compareFloat(actual float64, operator string, expected interface{}) bool {
	expectedVal, ok := expected.(float64)
	if !ok {
		return false
	}

	switch operator {
	case ">":
		return actual > expectedVal
	case ">=":
		return actual >= expectedVal
	case "<":
		return actual < expectedVal
	case "<=":
		return actual <= expectedVal
	case "=":
		return actual == expectedVal
	default:
		return false
	}
}

func compareInt(actual int, operator string, expected interface{}) bool {
	var expectedVal int
	switch v := expected.(type) {
	case int:
		expectedVal = v
	case float64:
		expectedVal = int(v)
	default:
		return false
	}

	switch operator {
	case ">":
		return actual > expectedVal
	case ">=":
		return actual >= expectedVal
	case "<":
		return actual < expectedVal
	case "<=":
		return actual <= expectedVal
	case "=":
		return actual == expectedVal
	default:
		return false
	}
}

func compareString(actual, operator string, expected interface{}) bool {
	expectedVal, ok := expected.(string)
	if !ok {
		return false
	}

	switch operator {
	case "=":
		return actual == expectedVal
	case "IN":
		// For IN operator, expected should be a slice
		return false
	default:
		return false
	}
}
