package queries

import "time"

// EmailDTO represents an email data transfer object
type EmailDTO struct {
	ID              int64      `json:"id"`
	Type            string     `json:"type"`
	Status          string     `json:"status"`
	Priority        int        `json:"priority"`
	From            string     `json:"from"`
	To              []string   `json:"to"`
	CC              []string   `json:"cc,omitempty"`
	BCC             []string   `json:"bcc,omitempty"`
	ReplyTo         string     `json:"reply_to,omitempty"`
	Subject         string     `json:"subject"`
	TemplateName    string     `json:"template_name,omitempty"`
	MaxRetries      int        `json:"max_retries"`
	RetryCount      int        `json:"retry_count"`
	ErrorMessage    string     `json:"error_message,omitempty"`
	OrderID         *int64     `json:"order_id,omitempty"`
	CustomerID      *int64     `json:"customer_id,omitempty"`
	HasAttachments  bool       `json:"has_attachments"`
	AttachmentCount int        `json:"attachment_count"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ScheduledAt     *time.Time `json:"scheduled_at,omitempty"`
	SentAt          *time.Time `json:"sent_at,omitempty"`
	FailedAt        *time.Time `json:"failed_at,omitempty"`
}

// EmailStatsDTO represents email statistics
type EmailStatsDTO struct {
	Total   int64 `json:"total"`
	Pending int64 `json:"pending"`
	Queued  int64 `json:"queued"`
	Sent    int64 `json:"sent"`
	Failed  int64 `json:"failed"`
}
