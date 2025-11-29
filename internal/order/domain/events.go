package domain

import (
	"time"
)

// OrderCreatedEvent is published when a new order is successfully created.
type OrderCreatedEvent struct {
	OrderID      int64
	CustomerID   int64
	OrderNumber  string
	CreationTime time.Time
}

// OrderItemAddedEvent is published when an item is added to an existing order.
type OrderItemAddedEvent struct {
	OrderID      int64
	OrderItemID  int64
	SKUID        int64
	Quantity     int
	AddedTime    time.Time
}

// OrderStatusUpdatedEvent is published when an order's status changes.
type OrderStatusUpdatedEvent struct {
	OrderID     int64
	OldStatus   OrderStatus
	NewStatus   OrderStatus
	UpdateTime  time.Time
}

// OrderSubmittedEvent is published when an order is submitted for processing.
type OrderSubmittedEvent struct {
	OrderID      int64
	OrderNumber  string
	SubmitTime   time.Time
	Total        float64
	CurrencyCode string
}

// OrderCancelledEvent is published when an order is cancelled.
type OrderCancelledEvent struct {
	OrderID      int64
	OrderNumber  string
	CancelTime   time.Time
	Reason       string // Optional: reason for cancellation
}

// OrderAdjustedEvent is published when an adjustment (discount/promotion) is applied to an order.
type OrderAdjustedEvent struct {
	OrderID     int64
	AdjustmentID int64
	OfferID     int64
	Amount      float64
	Description string
	AdjustmentTime time.Time
}

// FulfillmentGroupAddedEvent is published when a new fulfillment group is added to an order.
type FulfillmentGroupAddedEvent struct {
	OrderID          int64
	FulfillmentGroupID int64
	Type             string
	AddressID        int64
	ItemIDs          []int64
	CreationTime     time.Time
}