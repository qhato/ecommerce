package domain

import (
	"errors"
	"time"
)

// PromotionMessage represents a promotional message
type PromotionMessage struct {
	ID          int64
	Name        string
	Type        MessageType
	Priority    int
	Status      MessageStatus
	Message     string
	Description string
	Rules       []MessageRule
	Triggers    []MessageTrigger
	Placements  []string // Where to show: "banner", "popup", "cart", "product", etc.
	StartDate   *time.Time
	EndDate     *time.Time
	MaxViews    *int
	ViewCount   int
	ClickCount  int
	Metadata    map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeBanner      MessageType = "BANNER"
	MessageTypePopup       MessageType = "POPUP"
	MessageTypeNotification MessageType = "NOTIFICATION"
	MessageTypeInline      MessageType = "INLINE"
	MessageTypeTooltip     MessageType = "TOOLTIP"
)

// MessageStatus represents the status of a message
type MessageStatus string

const (
	MessageStatusActive   MessageStatus = "ACTIVE"
	MessageStatusInactive MessageStatus = "INACTIVE"
	MessageStatusExpired  MessageStatus = "EXPIRED"
)

// MessageRule represents a rule for displaying a message
type MessageRule struct {
	Field    string                 `json:"field"`     // "customer_type", "order_total", "product_category", etc.
	Operator string                 `json:"operator"`  // "equals", "greater_than", "less_than", "contains", etc.
	Value    string                 `json:"value"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// MessageTrigger represents when a message should be triggered
type MessageTrigger struct {
	Event      string                 `json:"event"`      // "page_load", "cart_add", "checkout_start", etc.
	Conditions []MessageRule          `json:"conditions,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// NewPromotionMessage creates a new promotion message
func NewPromotionMessage(name, message string, messageType MessageType, priority int) (*PromotionMessage, error) {
	if name == "" {
		return nil, errors.New("message name is required")
	}
	if message == "" {
		return nil, errors.New("message content is required")
	}

	now := time.Now()
	return &PromotionMessage{
		Name:       name,
		Message:    message,
		Type:       messageType,
		Priority:   priority,
		Status:     MessageStatusActive,
		Rules:      make([]MessageRule, 0),
		Triggers:   make([]MessageTrigger, 0),
		Placements: make([]string, 0),
		Metadata:   make(map[string]interface{}),
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// Activate activates the message
func (m *PromotionMessage) Activate() {
	m.Status = MessageStatusActive
	m.UpdatedAt = time.Now()
}

// Deactivate deactivates the message
func (m *PromotionMessage) Deactivate() {
	m.Status = MessageStatusInactive
	m.UpdatedAt = time.Now()
}

// IsActive checks if message is currently active
func (m *PromotionMessage) IsActive() bool {
	if m.Status != MessageStatusActive {
		return false
	}

	now := time.Now()

	if m.StartDate != nil && now.Before(*m.StartDate) {
		return false
	}

	if m.EndDate != nil && now.After(*m.EndDate) {
		return false
	}

	if m.MaxViews != nil && m.ViewCount >= *m.MaxViews {
		return false
	}

	return true
}

// IncrementView increments the view counter
func (m *PromotionMessage) IncrementView() {
	m.ViewCount++
	m.UpdatedAt = time.Now()
}

// IncrementClick increments the click counter
func (m *PromotionMessage) IncrementClick() {
	m.ClickCount++
	m.UpdatedAt = time.Now()
}

// MatchesRules checks if message rules match the given context
func (m *PromotionMessage) MatchesRules(context map[string]interface{}) bool {
	if len(m.Rules) == 0 {
		return true // No rules means always match
	}

	for _, rule := range m.Rules {
		if !evaluateRule(rule, context) {
			return false
		}
	}

	return true
}

// MatchesTrigger checks if message trigger matches the event
func (m *PromotionMessage) MatchesTrigger(event string, context map[string]interface{}) bool {
	if len(m.Triggers) == 0 {
		return true // No triggers means always match
	}

	for _, trigger := range m.Triggers {
		if trigger.Event == event {
			// Check trigger conditions
			allConditionsMet := true
			for _, condition := range trigger.Conditions {
				if !evaluateRule(condition, context) {
					allConditionsMet = false
					break
				}
			}
			if allConditionsMet {
				return true
			}
		}
	}

	return false
}

func evaluateRule(rule MessageRule, context map[string]interface{}) bool {
	value, exists := context[rule.Field]
	if !exists {
		return false
	}

	valueStr := toString(value)

	switch rule.Operator {
	case "equals":
		return valueStr == rule.Value
	case "not_equals":
		return valueStr != rule.Value
	case "contains":
		return contains(valueStr, rule.Value)
	case "greater_than":
		return compareNumbers(valueStr, rule.Value, ">")
	case "less_than":
		return compareNumbers(valueStr, rule.Value, "<")
	case "greater_or_equal":
		return compareNumbers(valueStr, rule.Value, ">=")
	case "less_or_equal":
		return compareNumbers(valueStr, rule.Value, "<=")
	default:
		return false
	}
}

func toString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr))
}

func compareNumbers(a, b string, operator string) bool {
	// Simplified number comparison - in production use proper parsing
	return false
}
