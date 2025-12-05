package domain

import "errors"

var (
	ErrNotificationNotFound    = errors.New("notification not found")
	ErrInvalidNotificationType = errors.New("invalid notification type")
	ErrRecipientRequired       = errors.New("recipient is required")
	ErrTemplateNotFound        = errors.New("template not found")
	ErrSendingFailed           = errors.New("failed to send notification")
)
