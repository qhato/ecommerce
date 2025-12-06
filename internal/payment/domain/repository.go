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

// PaymentTokenRepository defines the interface for payment token persistence
type PaymentTokenRepository interface {
	Create(ctx context.Context, token *PaymentToken) error
	Update(ctx context.Context, token *PaymentToken) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*PaymentToken, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*PaymentToken, error)
	FindDefaultByCustomerID(ctx context.Context, customerID string) (*PaymentToken, error)
	FindActiveByCustomerID(ctx context.Context, customerID string) ([]*PaymentToken, error)
}

// WebhookEventRepository defines the interface for webhook event persistence
type WebhookEventRepository interface {
	Create(ctx context.Context, event *WebhookEvent) error
	Update(ctx context.Context, event *WebhookEvent) error
	FindByID(ctx context.Context, id string) (*WebhookEvent, error)
	FindByEventID(ctx context.Context, gatewayName, eventID string) (*WebhookEvent, error)
	FindPending(ctx context.Context, limit int) ([]*WebhookEvent, error)
	FindByStatus(ctx context.Context, status WebhookStatus, limit int) ([]*WebhookEvent, error)
}

// GatewayConfigRepository defines the interface for gateway configuration persistence
type GatewayConfigRepository interface {
	Create(ctx context.Context, config *GatewayConfig) error
	Update(ctx context.Context, config *GatewayConfig) error
	FindByName(ctx context.Context, gatewayName string) (*GatewayConfig, error)
	FindAllEnabled(ctx context.Context) ([]*GatewayConfig, error)
	FindAll(ctx context.Context) ([]*GatewayConfig, error)
}
