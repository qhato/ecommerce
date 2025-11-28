package domain

import (
	"context"
)

// OrderRepository defines the interface for order persistence
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id int64) (*Order, error)
	FindByOrderNumber(ctx context.Context, orderNumber string) (*Order, error)
	FindByCustomerID(ctx context.Context, customerID int64, filter *OrderFilter) ([]*Order, int64, error)
	FindAll(ctx context.Context, filter *OrderFilter) ([]*Order, int64, error)
}

// OrderFilter represents filtering options for orders
type OrderFilter struct {
	Page       int
	PageSize   int
	Status     OrderStatus
	CustomerID int64
	SortBy     string
	SortOrder  string
}
