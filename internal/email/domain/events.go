package domain

import "time"

// EmailEvent represents a base email event
type EmailEvent struct {
	EmailID   int64
	Type      string
	Timestamp time.Time
	Data      map[string]interface{}
}

// EmailQueuedEvent is published when an email is queued
type EmailQueuedEvent struct {
	EmailID   int64
	Type      EmailType
	Priority  EmailPriority
	To        []string
	Subject   string
	Timestamp time.Time
}

// EmailSentEvent is published when an email is successfully sent
type EmailSentEvent struct {
	EmailID   int64
	Type      EmailType
	To        []string
	SentAt    time.Time
	Timestamp time.Time
}

// EmailFailedEvent is published when an email fails to send
type EmailFailedEvent struct {
	EmailID      int64
	Type         EmailType
	To           []string
	ErrorMessage string
	RetryCount   int
	CanRetry     bool
	Timestamp    time.Time
}

// EmailCancelledEvent is published when an email is cancelled
type EmailCancelledEvent struct {
	EmailID   int64
	Type      EmailType
	Timestamp time.Time
}

// EmailScheduledEvent is published when an email is scheduled
type EmailScheduledEvent struct {
	EmailID     int64
	Type        EmailType
	ScheduledAt time.Time
	Timestamp   time.Time
}

// NewEmailQueuedEvent creates a new EmailQueuedEvent
func NewEmailQueuedEvent(email *Email) *EmailQueuedEvent {
	return &EmailQueuedEvent{
		EmailID:   email.ID,
		Type:      email.Type,
		Priority:  email.Priority,
		To:        email.To,
		Subject:   email.Subject,
		Timestamp: time.Now(),
	}
}

// NewEmailSentEvent creates a new EmailSentEvent
func NewEmailSentEvent(email *Email) *EmailSentEvent {
	return &EmailSentEvent{
		EmailID:   email.ID,
		Type:      email.Type,
		To:        email.To,
		SentAt:    *email.SentAt,
		Timestamp: time.Now(),
	}
}

// NewEmailFailedEvent creates a new EmailFailedEvent
func NewEmailFailedEvent(email *Email) *EmailFailedEvent {
	return &EmailFailedEvent{
		EmailID:      email.ID,
		Type:         email.Type,
		To:           email.To,
		ErrorMessage: email.ErrorMessage,
		RetryCount:   email.RetryCount,
		CanRetry:     email.CanRetry(),
		Timestamp:    time.Now(),
	}
}

// NewEmailCancelledEvent creates a new EmailCancelledEvent
func NewEmailCancelledEvent(email *Email) *EmailCancelledEvent {
	return &EmailCancelledEvent{
		EmailID:   email.ID,
		Type:      email.Type,
		Timestamp: time.Now(),
	}
}

// NewEmailScheduledEvent creates a new EmailScheduledEvent
func NewEmailScheduledEvent(email *Email) *EmailScheduledEvent {
	return &EmailScheduledEvent{
		EmailID:     email.ID,
		Type:        email.Type,
		ScheduledAt: *email.ScheduledAt,
		Timestamp:   time.Now(),
	}
}
