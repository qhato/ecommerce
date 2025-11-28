package domain

import (
	"time"

	"github.com/qhato/ecommerce/pkg/event"
)

const (
	EventPaymentCreated    = "payment.created"
	EventPaymentAuthorized = "payment.authorized"
	EventPaymentCaptured   = "payment.captured"
	EventPaymentCompleted  = "payment.completed"
	EventPaymentFailed     = "payment.failed"
	EventPaymentRefunded   = "payment.refunded"
)

type PaymentCreatedEvent struct {
	event.BaseEvent
	PaymentID     int64         `json:"payment_id"`
	OrderID       int64         `json:"order_id"`
	CustomerID    int64         `json:"customer_id"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Amount        float64       `json:"amount"`
	CurrencyCode  string        `json:"currency_code"`
}

func NewPaymentCreatedEvent(paymentID, orderID, customerID int64, paymentMethod PaymentMethod, amount float64, currencyCode string) *PaymentCreatedEvent {
	return &PaymentCreatedEvent{
		BaseEvent:     event.BaseEvent{EventType: EventPaymentCreated, Timestamp: time.Now()},
		PaymentID:     paymentID,
		OrderID:       orderID,
		CustomerID:    customerID,
		PaymentMethod: paymentMethod,
		Amount:        amount,
		CurrencyCode:  currencyCode,
	}
}

func (e *PaymentCreatedEvent) Type() string {
	return e.EventType
}

type PaymentAuthorizedEvent struct {
	event.BaseEvent
	PaymentID         int64  `json:"payment_id"`
	OrderID           int64  `json:"order_id"`
	TransactionID     string `json:"transaction_id"`
	AuthorizationCode string `json:"authorization_code"`
	Amount            float64 `json:"amount"`
}

func (e *PaymentAuthorizedEvent) Type() string {
	return e.EventType
}

type PaymentCapturedEvent struct {
	event.BaseEvent
	PaymentID     int64  `json:"payment_id"`
	OrderID       int64  `json:"order_id"`
	TransactionID string `json:"transaction_id"`
	Amount        float64 `json:"amount"`
}

func (e *PaymentCapturedEvent) Type() string {
	return e.EventType
}

type PaymentCompletedEvent struct {
	event.BaseEvent
	PaymentID     int64  `json:"payment_id"`
	OrderID       int64  `json:"order_id"`
	TransactionID string `json:"transaction_id"`
	Amount        float64 `json:"amount"`
}

func (e *PaymentCompletedEvent) Type() string {
	return e.EventType
}

type PaymentFailedEvent struct {
	event.BaseEvent
	PaymentID     int64  `json:"payment_id"`
	OrderID       int64  `json:"order_id"`
	FailureReason string `json:"failure_reason"`
	Amount        float64 `json:"amount"`
}

func (e *PaymentFailedEvent) Type() string {
	return e.EventType
}

type PaymentRefundedEvent struct {
	event.BaseEvent
	PaymentID    int64   `json:"payment_id"`
	OrderID      int64   `json:"order_id"`
	RefundAmount float64 `json:"refund_amount"`
	TotalRefunded float64 `json:"total_refunded"`
}

func (e *PaymentRefundedEvent) Type() string {
	return e.EventType
}
