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
		BaseEvent:    event.BaseEvent{Type: EventOrderCreated, OccurredOn: time.Now()},
		OrderID:      orderID,
		OrderNumber:  orderNumber,
		CustomerID:   customerID,
		Total:        total,
		CurrencyCode: currencyCode,
	}
}



type OrderSubmittedEvent struct {
	event.BaseEvent
	OrderID     int64   `json:"order_id"`
	OrderNumber string  `json:"order_number"`
	CustomerID  int64   `json:"customer_id"`
	Total       float64 `json:"total"`
}



type OrderCancelledEvent struct {
	event.BaseEvent
	OrderID     int64  `json:"order_id"`
	OrderNumber string `json:"order_number"`
	CustomerID  int64  `json:"customer_id"`
}



type OrderShippedEvent struct {
	event.BaseEvent
	OrderID      int64  `json:"order_id"`
	OrderNumber  string `json:"order_number"`
	CustomerID   int64  `json:"customer_id"`
	TrackingInfo string `json:"tracking_info"`
}


