package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/notification/domain"
)

type NotificationDTO struct {
	ID             int64                  `json:"id"`
	Type           string                 `json:"type"`
	Status         string                 `json:"status"`
	Priority       string                 `json:"priority"`
	RecipientID    string                 `json:"recipient_id"`
	RecipientEmail string                 `json:"recipient_email,omitempty"`
	RecipientPhone string                 `json:"recipient_phone,omitempty"`
	Subject        string                 `json:"subject"`
	Body           string                 `json:"body"`
	TemplateID     *int64                 `json:"template_id,omitempty"`
	TemplateData   map[string]interface{} `json:"template_data,omitempty"`
	ScheduledFor   *time.Time             `json:"scheduled_for,omitempty"`
	SentAt         *time.Time             `json:"sent_at,omitempty"`
	FailedAt       *time.Time             `json:"failed_at,omitempty"`
	ErrorMsg       *string                `json:"error_msg,omitempty"`
	RetryCount     int                    `json:"retry_count"`
	MaxRetries     int                    `json:"max_retries"`
	CanRetry       bool                   `json:"can_retry"`
	IsScheduled    bool                   `json:"is_scheduled"`
	IsReadyToSend  bool                   `json:"is_ready_to_send"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

func ToNotificationDTO(n *domain.Notification) *NotificationDTO {
	return &NotificationDTO{
		ID:             n.ID,
		Type:           string(n.Type),
		Status:         string(n.Status),
		Priority:       string(n.Priority),
		RecipientID:    n.RecipientID,
		RecipientEmail: n.RecipientEmail,
		RecipientPhone: n.RecipientPhone,
		Subject:        n.Subject,
		Body:           n.Body,
		TemplateID:     n.TemplateID,
		TemplateData:   n.TemplateData,
		ScheduledFor:   n.ScheduledFor,
		SentAt:         n.SentAt,
		FailedAt:       n.FailedAt,
		ErrorMsg:       n.ErrorMsg,
		RetryCount:     n.RetryCount,
		MaxRetries:     n.MaxRetries,
		CanRetry:       n.CanRetry(),
		IsScheduled:    n.IsScheduled(),
		IsReadyToSend:  n.IsReadyToSend(),
		CreatedAt:      n.CreatedAt,
		UpdatedAt:      n.UpdatedAt,
	}
}
