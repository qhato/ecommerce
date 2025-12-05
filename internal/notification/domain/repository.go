package domain

import "context"

// NotificationRepository defines the interface for notification persistence
type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	Update(ctx context.Context, notification *Notification) error
	FindByID(ctx context.Context, id int64) (*Notification, error)
	FindByRecipientID(ctx context.Context, recipientID string, limit int) ([]*Notification, error)
	FindByStatus(ctx context.Context, status NotificationStatus, limit int) ([]*Notification, error)
	FindPending(ctx context.Context, limit int) ([]*Notification, error)
	FindScheduled(ctx context.Context, limit int) ([]*Notification, error)
	FindFailed(ctx context.Context, limit int) ([]*Notification, error)
	Delete(ctx context.Context, id int64) error
}
