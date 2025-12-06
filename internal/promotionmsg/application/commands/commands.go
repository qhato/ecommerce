package commands

import "time"

// CreatePromotionMessageCommand creates a new promotion message
type CreatePromotionMessageCommand struct {
	Name        string                 `json:"name"`
	Message     string                 `json:"message"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Priority    int                    `json:"priority"`
	Rules       []MessageRuleCmd       `json:"rules,omitempty"`
	Triggers    []MessageTriggerCmd    `json:"triggers,omitempty"`
	Placements  []string               `json:"placements,omitempty"`
	StartDate   *time.Time             `json:"start_date,omitempty"`
	EndDate     *time.Time             `json:"end_date,omitempty"`
	MaxViews    *int                   `json:"max_views,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdatePromotionMessageCommand updates a promotion message
type UpdatePromotionMessageCommand struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Message     string                 `json:"message"`
	Description string                 `json:"description"`
	Priority    int                    `json:"priority"`
	Rules       []MessageRuleCmd       `json:"rules,omitempty"`
	Triggers    []MessageTriggerCmd    `json:"triggers,omitempty"`
	Placements  []string               `json:"placements,omitempty"`
	StartDate   *time.Time             `json:"start_date,omitempty"`
	EndDate     *time.Time             `json:"end_date,omitempty"`
	MaxViews    *int                   `json:"max_views,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ActivateMessageCommand activates a message
type ActivateMessageCommand struct {
	ID int64 `json:"id"`
}

// DeactivateMessageCommand deactivates a message
type DeactivateMessageCommand struct {
	ID int64 `json:"id"`
}

// IncrementViewCommand increments view counter
type IncrementViewCommand struct {
	ID int64 `json:"id"`
}

// IncrementClickCommand increments click counter
type IncrementClickCommand struct {
	ID int64 `json:"id"`
}

// DeleteMessageCommand deletes a message
type DeleteMessageCommand struct {
	ID int64 `json:"id"`
}

// MessageRuleCmd represents a message rule
type MessageRuleCmd struct {
	Field    string                 `json:"field"`
	Operator string                 `json:"operator"`
	Value    string                 `json:"value"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// MessageTriggerCmd represents a message trigger
type MessageTriggerCmd struct {
	Event      string                 `json:"event"`
	Conditions []MessageRuleCmd       `json:"conditions,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}
