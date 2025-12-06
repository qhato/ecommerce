package domain

import (
	"fmt"
	"strconv"
	"strings"
)

// RuleEngine evaluates rules and executes actions
type RuleEngine struct{}

// NewRuleEngine creates a new rule engine
func NewRuleEngine() *RuleEngine {
	return &RuleEngine{}
}

// EvaluateRules evaluates a list of rules against a context and returns matching actions
func (e *RuleEngine) EvaluateRules(rules []*Rule, context map[string]interface{}) []Action {
	var actions []Action

	for _, rule := range rules {
		if !rule.IsActiveNow() {
			continue
		}

		if e.evaluateConditions(rule.Conditions, context) {
			actions = append(actions, rule.Actions...)
		}
	}

	return actions
}

// evaluateConditions evaluates all conditions for a rule
func (e *RuleEngine) evaluateConditions(conditions []Condition, context map[string]interface{}) bool {
	if len(conditions) == 0 {
		return true
	}

	for _, condition := range conditions {
		if !e.evaluateCondition(condition, context) {
			return false
		}
	}

	return true
}

// evaluateCondition evaluates a single condition
func (e *RuleEngine) evaluateCondition(condition Condition, context map[string]interface{}) bool {
	value, ok := e.getValueFromContext(condition.Field, context)
	if !ok {
		return false
	}

	switch condition.Operator {
	case OperatorEquals:
		return fmt.Sprintf("%v", value) == condition.Value
	case OperatorNotEquals:
		return fmt.Sprintf("%v", value) != condition.Value
	case OperatorGreaterThan:
		return e.compareNumbers(value, condition.Value, ">")
	case OperatorLessThan:
		return e.compareNumbers(value, condition.Value, "<")
	case OperatorContains:
		return strings.Contains(fmt.Sprintf("%v", value), condition.Value)
	default:
		return false
	}
}

// getValueFromContext retrieves a value from context using dot notation (e.g., "order.total")
func (e *RuleEngine) getValueFromContext(field string, context map[string]interface{}) (interface{}, bool) {
	parts := strings.Split(field, ".")
	current := context

	for i, part := range parts {
		value, ok := current[part]
		if !ok {
			return nil, false
		}

		if i == len(parts)-1 {
			return value, true
		}

		if nested, ok := value.(map[string]interface{}); ok {
			current = nested
		} else {
			return nil, false
		}
	}

	return nil, false
}

// compareNumbers compares two values as numbers
func (e *RuleEngine) compareNumbers(value interface{}, compareValue string, operator string) bool {
	var num1 float64
	var num2 float64

	switch v := value.(type) {
	case int:
		num1 = float64(v)
	case int64:
		num1 = float64(v)
	case float64:
		num1 = v
	case string:
		var err error
		num1, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return false
		}
	default:
		return false
	}

	var err error
	num2, err = strconv.ParseFloat(compareValue, 64)
	if err != nil {
		return false
	}

	switch operator {
	case ">":
		return num1 > num2
	case "<":
		return num1 < num2
	default:
		return false
	}
}
