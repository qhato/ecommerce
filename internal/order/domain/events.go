package domain

import (
	"time"

	"github.com/qhato/ecommerce/pkg/event"
)

const (
	EventOrderCreated   = "order.created"
	EventOrderSubmitted = "order.submitted"
	EventOrderCancelled = "order.cancelled"
	EventOrderShipped   = "order.shipped"
)

type OrderCreatedEvent struct {
	event.BaseEvent
	OrderID      int64  `json:"order_id"`
	OrderNumber  string `json:"order_number"`
	CustomerID   int64  `json:"customer_id"`
	Total        float64 `json:"total"`
	CurrencyCode string `json:"currency_code"`
}

func NewOrderCreatedEvent(orderID int64, orderNumber string, customerID int64, total float64, currencyCode string) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		BaseEvent:    event.BaseEvent{EventType: EventOrderCreated, Timestamp: time.Now()},
		OrderID:      orderID,
		OrderNumber:  orderNumber,
		CustomerID:   customerID,
		Total:        total,
		CurrencyCode: currencyCode,
	}
}

func (e *OrderCreatedEvent) Type() string {
	return e.EventType
}

type OrderSubmittedEvent struct {
	event.BaseEvent
	OrderID     int64   `json:"order_id"`
	OrderNumber string  `json:"order_number"`
	CustomerID  int64   `json:"customer_id"`
	Total       float64 `json:"total"`
}

func (e *OrderSubmittedEvent) Type() string {
	return e.EventType
}

type OrderCancelledEvent struct {
	event.BaseEvent
	OrderID     int64  `json:"order_id"`
	OrderNumber string `json:"order_number"`
	CustomerID  int64  `json:"customer_id"`
}

func (e *OrderCancelledEvent) Type() string {
	return e.EventType
}

type OrderShippedEvent struct {
	event.BaseEvent
	OrderID      int64  `json:"order_id"`
	OrderNumber  string `json:"order_number"`
	CustomerID   int64  `json:"customer_id"`
	TrackingInfo string `json:"tracking_info"`
}

func (e *OrderShippedEvent) Type() string {
	return e.EventType
}
