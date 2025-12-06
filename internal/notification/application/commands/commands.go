package commands

import "time"

// CreateNotificationCommand creates a new notification
type CreateNotificationCommand struct {
	Type           string                 `json:"type"`
	RecipientID    string                 `json:"recipient_id"`
	RecipientEmail string                 `json:"recipient_email,omitempty"`
	RecipientPhone string                 `json:"recipient_phone,omitempty"`
	Subject        string                 `json:"subject"`
	Body           string                 `json:"body"`
	Priority       string                 `json:"priority,omitempty"`
	TemplateID     *int64                 `json:"template_id,omitempty"`
	TemplateData   map[string]interface{} `json:"template_data,omitempty"`
	ScheduledFor   *time.Time             `json:"scheduled_for,omitempty"`
}

// SendNotificationCommand sends a pending notification
type SendNotificationCommand struct {
	ID int64 `json:"id"`
}

// MarkAsSentCommand marks a notification as sent
type MarkAsSentCommand struct {
	ID int64 `json:"id"`
}

// MarkAsFailedCommand marks a notification as failed
type MarkAsFailedCommand struct {
	ID       int64  `json:"id"`
	ErrorMsg string `json:"error_msg"`
}

// RetryNotificationCommand retries a failed notification
type RetryNotificationCommand struct {
	ID int64 `json:"id"`
}

// CancelNotificationCommand cancels a notification
type CancelNotificationCommand struct {
	ID int64 `json:"id"`
}

// DeleteNotificationCommand deletes a notification
type DeleteNotificationCommand struct {
	ID int64 `json:"id"`
}
