package domain

import "time"

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeEmail NotificationType = "EMAIL"
	NotificationTypeSMS   NotificationType = "SMS"
	NotificationTypePush  NotificationType = "PUSH"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "PENDING"
	NotificationStatusSending   NotificationStatus = "SENDING"
	NotificationStatusSent      NotificationStatus = "SENT"
	NotificationStatusFailed    NotificationStatus = "FAILED"
	NotificationStatusCancelled NotificationStatus = "CANCELLED"
)

// NotificationPriority represents the priority of a notification
type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "LOW"
	NotificationPriorityNormal NotificationPriority = "NORMAL"
	NotificationPriorityHigh   NotificationPriority = "HIGH"
	NotificationPriorityUrgent NotificationPriority = "URGENT"
)

// Notification represents a notification to be sent
type Notification struct {
	ID           int64
	Type         NotificationType
	Status       NotificationStatus
	Priority     NotificationPriority
	RecipientID  string
	RecipientEmail string
	RecipientPhone string
	Subject      string
	Body         string
	TemplateID   *int64
	TemplateData map[string]interface{}
	ScheduledFor *time.Time
	SentAt       *time.Time
	FailedAt     *time.Time
	ErrorMsg     *string
	RetryCount   int
	MaxRetries   int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewNotification creates a new notification
func NewNotification(notifType NotificationType, recipientID, subject, body string) (*Notification, error) {
	now := time.Now()
	return &Notification{
		Type:         notifType,
		Status:       NotificationStatusPending,
		Priority:     NotificationPriorityNormal,
		RecipientID:  recipientID,
		Subject:      subject,
		Body:         body,
		TemplateData: make(map[string]interface{}),
		MaxRetries:   3,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// MarkAsSending marks the notification as being sent
func (n *Notification) MarkAsSending() {
	n.Status = NotificationStatusSending
	n.UpdatedAt = time.Now()
}

// MarkAsSent marks the notification as successfully sent
func (n *Notification) MarkAsSent() {
	now := time.Now()
	n.Status = NotificationStatusSent
	n.SentAt = &now
	n.UpdatedAt = now
}

// MarkAsFailed marks the notification as failed
func (n *Notification) MarkAsFailed(errorMsg string) {
	now := time.Now()
	n.Status = NotificationStatusFailed
	n.ErrorMsg = &errorMsg
	n.FailedAt = &now
	n.RetryCount++
	n.UpdatedAt = now
}

// Cancel cancels the notification
func (n *Notification) Cancel() {
	n.Status = NotificationStatusCancelled
	n.UpdatedAt = time.Now()
}

// CanRetry checks if the notification can be retried
func (n *Notification) CanRetry() bool {
	return n.Status == NotificationStatusFailed && n.RetryCount < n.MaxRetries
}

// IsScheduled checks if the notification is scheduled for future delivery
func (n *Notification) IsScheduled() bool {
	return n.ScheduledFor != nil && n.ScheduledFor.After(time.Now())
}

// IsReadyToSend checks if the notification is ready to be sent
func (n *Notification) IsReadyToSend() bool {
	if n.Status != NotificationStatusPending {
		return false
	}
	if n.ScheduledFor != nil {
		return time.Now().After(*n.ScheduledFor)
	}
	return true
}
