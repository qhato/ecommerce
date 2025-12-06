package domain

import "errors"

var (
	ErrRuleNotFound         = errors.New("rule not found")
	ErrRuleNameRequired     = errors.New("rule name is required")
	ErrRuleNameTaken        = errors.New("rule name is already taken")
	ErrInvalidCondition     = errors.New("invalid condition")
	ErrInvalidAction        = errors.New("invalid action")
	ErrRuleEvaluationFailed = errors.New("rule evaluation failed")
)
