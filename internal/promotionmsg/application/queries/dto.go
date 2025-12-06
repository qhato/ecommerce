package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/promotionmsg/domain"
)

type PromotionMessageDTO struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Priority    int                    `json:"priority"`
	Status      string                 `json:"status"`
	Message     string                 `json:"message"`
	Description string                 `json:"description"`
	Rules       []MessageRuleDTO       `json:"rules,omitempty"`
	Triggers    []MessageTriggerDTO    `json:"triggers,omitempty"`
	Placements  []string               `json:"placements,omitempty"`
	StartDate   *time.Time             `json:"start_date,omitempty"`
	EndDate     *time.Time             `json:"end_date,omitempty"`
	MaxViews    *int                   `json:"max_views,omitempty"`
	ViewCount   int                    `json:"view_count"`
	ClickCount  int                    `json:"click_count"`
	IsActive    bool                   `json:"is_active"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type MessageRuleDTO struct {
	Field    string                 `json:"field"`
	Operator string                 `json:"operator"`
	Value    string                 `json:"value"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type MessageTriggerDTO struct {
	Event      string                 `json:"event"`
	Conditions []MessageRuleDTO       `json:"conditions,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

func ToPromotionMessageDTO(m *domain.PromotionMessage) *PromotionMessageDTO {
	rules := make([]MessageRuleDTO, len(m.Rules))
	for i, r := range m.Rules {
		rules[i] = MessageRuleDTO{
			Field:    r.Field,
			Operator: r.Operator,
			Value:    r.Value,
			Metadata: r.Metadata,
		}
	}

	triggers := make([]MessageTriggerDTO, len(m.Triggers))
	for i, t := range m.Triggers {
		conditions := make([]MessageRuleDTO, len(t.Conditions))
		for j, c := range t.Conditions {
			conditions[j] = MessageRuleDTO{
				Field:    c.Field,
				Operator: c.Operator,
				Value:    c.Value,
				Metadata: c.Metadata,
			}
		}
		triggers[i] = MessageTriggerDTO{
			Event:      t.Event,
			Conditions: conditions,
			Metadata:   t.Metadata,
		}
	}

	return &PromotionMessageDTO{
		ID:          m.ID,
		Name:        m.Name,
		Type:        string(m.Type),
		Priority:    m.Priority,
		Status:      string(m.Status),
		Message:     m.Message,
		Description: m.Description,
		Rules:       rules,
		Triggers:    triggers,
		Placements:  m.Placements,
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		MaxViews:    m.MaxViews,
		ViewCount:   m.ViewCount,
		ClickCount:  m.ClickCount,
		IsActive:    m.IsActive(),
		Metadata:    m.Metadata,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
