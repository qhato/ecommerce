package domain

import "context"

// EmailRepository defines the interface for email persistence
type EmailRepository interface {
	// Create creates a new email
	Create(ctx context.Context, email *Email) error

	// Update updates an existing email
	Update(ctx context.Context, email *Email) error

	// FindByID finds an email by ID
	FindByID(ctx context.Context, id int64) (*Email, error)

	// FindPendingEmails finds emails that are pending to be sent
	FindPendingEmails(ctx context.Context, limit int) ([]*Email, error)

	// FindByStatus finds emails by status
	FindByStatus(ctx context.Context, status EmailStatus, offset, limit int) ([]*Email, error)

	// FindByType finds emails by type
	FindByType(ctx context.Context, emailType EmailType, offset, limit int) ([]*Email, error)

	// FindByOrderID finds emails associated with an order
	FindByOrderID(ctx context.Context, orderID int64) ([]*Email, error)

	// FindByCustomerID finds emails associated with a customer
	FindByCustomerID(ctx context.Context, customerID int64, offset, limit int) ([]*Email, error)

	// FindScheduledEmails finds emails scheduled to be sent
	FindScheduledEmails(ctx context.Context) ([]*Email, error)

	// FindFailedEmails finds emails that failed and can be retried
	FindFailedEmailsForRetry(ctx context.Context, limit int) ([]*Email, error)

	// Delete deletes an email by ID
	Delete(ctx context.Context, id int64) error

	// Count counts total emails
	Count(ctx context.Context) (int64, error)

	// CountByStatus counts emails by status
	CountByStatus(ctx context.Context, status EmailStatus) (int64, error)
}
