package rules

import (
	"fmt"
	"sync"
)

// RuleEngine manages and executes business rules
type RuleEngine struct {
	rules    map[string]Rule
	ruleSets map[string]*RuleSet
	mu       sync.RWMutex
}

// NewRuleEngine creates a new RuleEngine
func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		rules:    make(map[string]Rule),
		ruleSets: make(map[string]*RuleSet),
	}
}

// RegisterRule registers a new rule
func (re *RuleEngine) RegisterRule(rule Rule) error {
	re.mu.Lock()
	defer re.mu.Unlock()

	if _, exists := re.rules[rule.GetName()]; exists {
		return fmt.Errorf("rule %s already registered", rule.GetName())
	}

	re.rules[rule.GetName()] = rule
	return nil
}

// RegisterRuleSet registers a new rule set
func (re *RuleEngine) RegisterRuleSet(ruleSet *RuleSet) error {
	re.mu.Lock()
	defer re.mu.Unlock()

	if _, exists := re.ruleSets[ruleSet.GetName()]; exists {
		return fmt.Errorf("rule set %s already registered", ruleSet.GetName())
	}

	re.ruleSets[ruleSet.GetName()] = ruleSet
	return nil
}

// EvaluateRule evaluates a single rule by name
func (re *RuleEngine) EvaluateRule(ruleName string, env map[string]interface{}) (bool, error) {
	re.mu.RLock()
	defer re.mu.RUnlock()

	rule, exists := re.rules[ruleName]
	if !exists {
		return false, fmt.Errorf("rule not found: %s", ruleName)
	}

	return rule.Evaluate(env)
}

// EvaluateRuleSet evaluates a rule set by name
func (re *RuleEngine) EvaluateRuleSet(ruleSetName string, env map[string]interface{}) (bool, error) {
	re.mu.RLock()
	defer re.mu.RUnlock()

	ruleSet, exists := re.ruleSets[ruleSetName]
	if !exists {
		return false, fmt.Errorf("rule set not found: %s", ruleSetName)
	}

	return ruleSet.Evaluate(env)
}

// GetRule retrieves a rule by name
func (re *RuleEngine) GetRule(name string) (Rule, bool) {
	re.mu.RLock()
	defer re.mu.RUnlock()

	rule, exists := re.rules[name]
	return rule, exists
}

// GetRuleSet retrieves a rule set by name
func (re *RuleEngine) GetRuleSet(name string) (*RuleSet, bool) {
	re.mu.RLock()
	defer re.mu.RUnlock()

	ruleSet, exists := re.ruleSets[name]
	return ruleSet, exists
}

// ListRules returns all registered rule names
func (re *RuleEngine) ListRules() []string {
	re.mu.RLock()
	defer re.mu.RUnlock()

	names := make([]string, 0, len(re.rules))
	for name := range re.rules {
		names = append(names, name)
	}
	return names
}

// ListRuleSets returns all registered rule set names
func (re *RuleEngine) ListRuleSets() []string {
	re.mu.RLock()
	defer re.mu.RUnlock()

	names := make([]string, 0, len(re.ruleSets))
	for name := range re.ruleSets {
		names = append(names, name)
	}
	return names
}

// Common helper functions for building rule environments

// BuildOrderEnv builds an environment for order-related rules
func BuildOrderEnv(order interface{}) map[string]interface{} {
	return map[string]interface{}{
		"order": order,
	}
}

// BuildCustomerEnv builds an environment for customer-related rules
func BuildCustomerEnv(customer interface{}) map[string]interface{} {
	return map[string]interface{}{
		"customer": customer,
	}
}

// BuildOfferEnv builds an environment for offer-related rules
func BuildOfferEnv(order, customer, item interface{}) map[string]interface{} {
	env := make(map[string]interface{})
	if order != nil {
		env["order"] = order
	}
	if customer != nil {
		env["customer"] = customer
	}
	if item != nil {
		env["item"] = item
	}
	return env
}

// BuildTaxEnv builds an environment for tax-related rules
func BuildTaxEnv(address, order interface{}) map[string]interface{} {
	return map[string]interface{}{
		"address": address,
		"order":   order,
	}
}
