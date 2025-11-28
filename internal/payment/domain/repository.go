package domain

import (
	"context"
)

// PaymentRepository defines the interface for payment persistence
type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	Update(ctx context.Context, payment *Payment) error
	FindByID(ctx context.Context, id int64) (*Payment, error)
	FindByOrderID(ctx context.Context, orderID int64) ([]*Payment, error)
	FindByCustomerID(ctx context.Context, customerID int64, filter *PaymentFilter) ([]*Payment, int64, error)
	FindByTransactionID(ctx context.Context, transactionID string) (*Payment, error)
	FindAll(ctx context.Context, filter *PaymentFilter) ([]*Payment, int64, error)
}

// PaymentFilter represents filtering options for payments
type PaymentFilter struct {
	Page          int
	PageSize      int
	Status        PaymentStatus
	PaymentMethod PaymentMethod
	CustomerID    int64
	OrderID       int64
	SortBy        string
	SortOrder     string
}
