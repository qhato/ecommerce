package domain

import "errors"

var (
	// ErrEmailNotFound is returned when an email is not found
	ErrEmailNotFound = errors.New("email not found")

	// ErrInvalidEmail is returned when email validation fails
	ErrInvalidEmail = errors.New("invalid email")

	// ErrEmailAlreadySent is returned when attempting to modify a sent email
	ErrEmailAlreadySent = errors.New("email already sent")

	// ErrEmailCancelled is returned when attempting to send a cancelled email
	ErrEmailCancelled = errors.New("email cancelled")

	// ErrInvalidRecipient is returned when a recipient address is invalid
	ErrInvalidRecipient = errors.New("invalid recipient address")

	// ErrTemplateNotFound is returned when an email template is not found
	ErrTemplateNotFound = errors.New("email template not found")

	// ErrTemplateRenderFailed is returned when template rendering fails
	ErrTemplateRenderFailed = errors.New("failed to render email template")

	// ErrSMTPConnectionFailed is returned when SMTP connection fails
	ErrSMTPConnectionFailed = errors.New("SMTP connection failed")

	// ErrEmailSendFailed is returned when email sending fails
	ErrEmailSendFailed = errors.New("failed to send email")

	// ErrQueueFull is returned when the email queue is full
	ErrQueueFull = errors.New("email queue is full")

	// ErrMaxRetriesExceeded is returned when max retries are exceeded
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")

	// ErrInvalidScheduleTime is returned when schedule time is invalid
	ErrInvalidScheduleTime = errors.New("invalid schedule time")
)
