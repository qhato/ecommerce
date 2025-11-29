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

// OrderItemRepository defines the interface for order item persistence
type OrderItemRepository interface {
	// Save stores a new order item or updates an existing one.
	Save(ctx context.Context, item *OrderItem) error

	// FindByID retrieves an order item by its unique identifier.
	FindByID(ctx context.Context, id int64) (*OrderItem, error)

	// FindByOrderID retrieves all order items for a given order ID.
	FindByOrderID(ctx context.Context, orderID int64) ([]*OrderItem, error)

	// Delete removes an order item by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteByOrderID removes all order items for a given order ID.
	DeleteByOrderID(ctx context.Context, orderID int64) error
}

// OrderAdjustmentRepository defines the interface for order adjustment persistence
type OrderAdjustmentRepository interface {
	// Save stores a new order adjustment or updates an existing one.
	Save(ctx context.Context, adjustment *OrderAdjustment) error

	// FindByID retrieves an order adjustment by its unique identifier.
	FindByID(ctx context.Context, id int64) (*OrderAdjustment, error)

	// FindByOrderID retrieves all order adjustments for a given order ID.
	FindByOrderID(ctx context.Context, orderID int64) ([]*OrderAdjustment, error)

	// Delete removes an order adjustment by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteByOrderID removes all order adjustments for a given order ID.
	DeleteByOrderID(ctx context.Context, orderID int64) error
}

// OrderItemAdjustmentRepository defines the interface for order item adjustment persistence
type OrderItemAdjustmentRepository interface {
	// Save stores a new order item adjustment or updates an existing one.
	Save(ctx context.Context, adjustment *OrderItemAdjustment) error

	// FindByID retrieves an order item adjustment by its unique identifier.
	FindByID(ctx context.Context, id int64) (*OrderItemAdjustment, error)

	// FindByOrderItemID retrieves all order item adjustments for a given order item ID.
	FindByOrderItemID(ctx context.Context, orderItemID int64) ([]*OrderItemAdjustment, error)

	// Delete removes an order item adjustment by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteByOrderItemID removes all order item adjustments for a given order item ID.
	DeleteByOrderItemID(ctx context.Context, orderItemID int64) error
}

// OrderItemAttributeRepository defines the interface for order item attribute persistence
type OrderItemAttributeRepository interface {
	// Save stores a new order item attribute or updates an existing one.
	Save(ctx context.Context, attribute *OrderItemAttribute) error

	// FindByOrderItemIDAndName retrieves an order item attribute by order item ID and name.
	FindByOrderItemIDAndName(ctx context.Context, orderItemID int64, name string) (*OrderItemAttribute, error)

	// FindByOrderItemID retrieves all order item attributes for a given order item ID.
	FindByOrderItemID(ctx context.Context, orderItemID int64) ([]*OrderItemAttribute, error)

	// Delete removes an order item attribute by order item ID and name.
	Delete(ctx context.Context, orderItemID int64, name string) error

	// DeleteByOrderItemID removes all order item attributes for a given order item ID.
	DeleteByOrderItemID(ctx context.Context, orderItemID int64) error
}

// FulfillmentGroupRepository defines the interface for fulfillment group persistence
type FulfillmentGroupRepository interface {
	// Save stores a new fulfillment group or updates an existing one.
	Save(ctx context.Context, group *FulfillmentGroup) error

	// FindByID retrieves a fulfillment group by its unique identifier.
	FindByID(ctx context.Context, id int64) (*FulfillmentGroup, error)

	// FindByOrderID retrieves all fulfillment groups for a given order ID.
	FindByOrderID(ctx context.Context, orderID int64) ([]*FulfillmentGroup, error)

	// Delete removes a fulfillment group by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteByOrderID removes all fulfillment groups for a given order ID.
	DeleteByOrderID(ctx context.Context, orderID int64) error
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

// OrderItemFilter represents filtering options for order items
type OrderItemFilter struct {
	Page      int
	PageSize  int
	OrderID   int64
	SKUID     int64
	ProductID int64
	SortBy    string
	SortOrder string
}

// FulfillmentGroupFilter represents filtering options for fulfillment groups
type FulfillmentGroupFilter struct {
	Page    int
	PageSize int
	OrderID int64
	Status  string
	SortBy  string
	SortOrder string
}