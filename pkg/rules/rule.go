package rules

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// Rule represents a business rule with a condition
type Rule interface {
	// Evaluate evaluates the rule against the provided environment
	Evaluate(env map[string]interface{}) (bool, error)

	// GetName returns the rule name
	GetName() string

	// GetExpression returns the rule expression
	GetExpression() string

	// GetDescription returns a description of what the rule does
	GetDescription() string
}

// CompiledRule is a pre-compiled rule for better performance
type CompiledRule struct {
	name        string
	expression  string
	description string
	program     *vm.Program
}

// NewRule creates and compiles a new rule
func NewRule(name, expression, description string) (*CompiledRule, error) {
	program, err := expr.Compile(expression, expr.AsBool())
	if err != nil {
		return nil, fmt.Errorf("failed to compile rule %s: %w", name, err)
	}

	return &CompiledRule{
		name:        name,
		expression:  expression,
		description: description,
		program:     program,
	}, nil
}

// Evaluate evaluates the compiled rule
func (r *CompiledRule) Evaluate(env map[string]interface{}) (bool, error) {
	output, err := expr.Run(r.program, env)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate rule %s: %w", r.name, err)
	}

	result, ok := output.(bool)
	if !ok {
		return false, fmt.Errorf("rule %s did not return a boolean value", r.name)
	}

	return result, nil
}

func (r *CompiledRule) GetName() string {
	return r.name
}

func (r *CompiledRule) GetExpression() string {
	return r.expression
}

func (r *CompiledRule) GetDescription() string {
	return r.description
}

// RuleSet represents a collection of rules
type RuleSet struct {
	name  string
	rules []Rule
	mode  RuleSetMode
}

// RuleSetMode determines how rules in a set are evaluated
type RuleSetMode string

const (
	// AllRulesMode requires all rules to pass
	AllRulesMode RuleSetMode = "ALL"
	// AnyRuleMode requires at least one rule to pass
	AnyRuleMode RuleSetMode = "ANY"
	// FirstMatchMode stops at the first matching rule
	FirstMatchMode RuleSetMode = "FIRST_MATCH"
)

// NewRuleSet creates a new RuleSet
func NewRuleSet(name string, mode RuleSetMode) *RuleSet {
	return &RuleSet{
		name:  name,
		rules: make([]Rule, 0),
		mode:  mode,
	}
}

// AddRule adds a rule to the set
func (rs *RuleSet) AddRule(rule Rule) {
	rs.rules = append(rs.rules, rule)
}

// Evaluate evaluates all rules in the set according to the mode
func (rs *RuleSet) Evaluate(env map[string]interface{}) (bool, error) {
	if len(rs.rules) == 0 {
		return true, nil // Empty rule set always passes
	}

	switch rs.mode {
	case AllRulesMode:
		return rs.evaluateAll(env)
	case AnyRuleMode:
		return rs.evaluateAny(env)
	case FirstMatchMode:
		return rs.evaluateFirstMatch(env)
	default:
		return false, fmt.Errorf("unknown rule set mode: %s", rs.mode)
	}
}

func (rs *RuleSet) evaluateAll(env map[string]interface{}) (bool, error) {
	for _, rule := range rs.rules {
		result, err := rule.Evaluate(env)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

func (rs *RuleSet) evaluateAny(env map[string]interface{}) (bool, error) {
	for _, rule := range rs.rules {
		result, err := rule.Evaluate(env)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

func (rs *RuleSet) evaluateFirstMatch(env map[string]interface{}) (bool, error) {
	for _, rule := range rs.rules {
		result, err := rule.Evaluate(env)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

func (rs *RuleSet) GetName() string {
	return rs.name
}

func (rs *RuleSet) GetRules() []Rule {
	return rs.rules
}
