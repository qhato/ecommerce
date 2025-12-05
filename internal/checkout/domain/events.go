package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// DomainEvent is the base interface for all domain events
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

// CheckoutInitiatedEvent is emitted when a checkout session is created
type CheckoutInitiatedEvent struct {
	SessionID       string
	OrderID         int64
	CustomerID      *string
	Email           string
	IsGuestCheckout bool
	OccurredOn      time.Time
}

func (e CheckoutInitiatedEvent) EventType() string {
	return "checkout.initiated"
}

func (e CheckoutInitiatedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutCustomerInfoAddedEvent is emitted when customer info is added
type CheckoutCustomerInfoAddedEvent struct {
	SessionID  string
	OrderID    int64
	CustomerID *string
	Email      string
	OccurredOn time.Time
}

func (e CheckoutCustomerInfoAddedEvent) EventType() string {
	return "checkout.customer_info.added"
}

func (e CheckoutCustomerInfoAddedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutShippingAddressAddedEvent is emitted when shipping address is added
type CheckoutShippingAddressAddedEvent struct {
	SessionID string
	OrderID   int64
	AddressID int64
	OccurredOn time.Time
}

func (e CheckoutShippingAddressAddedEvent) EventType() string {
	return "checkout.shipping_address.added"
}

func (e CheckoutShippingAddressAddedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutShippingMethodSelectedEvent is emitted when shipping method is selected
type CheckoutShippingMethodSelectedEvent struct {
	SessionID        string
	OrderID          int64
	ShippingMethodID string
	ShippingCost     decimal.Decimal
	OccurredOn       time.Time
}

func (e CheckoutShippingMethodSelectedEvent) EventType() string {
	return "checkout.shipping_method.selected"
}

func (e CheckoutShippingMethodSelectedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutBillingAddressAddedEvent is emitted when billing address is added
type CheckoutBillingAddressAddedEvent struct {
	SessionID  string
	OrderID    int64
	AddressID  int64
	SameAsShipping bool
	OccurredOn time.Time
}

func (e CheckoutBillingAddressAddedEvent) EventType() string {
	return "checkout.billing_address.added"
}

func (e CheckoutBillingAddressAddedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutPaymentMethodAddedEvent is emitted when payment method is added
type CheckoutPaymentMethodAddedEvent struct {
	SessionID       string
	OrderID         int64
	PaymentMethodID int64
	OccurredOn      time.Time
}

func (e CheckoutPaymentMethodAddedEvent) EventType() string {
	return "checkout.payment_method.added"
}

func (e CheckoutPaymentMethodAddedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutCouponAppliedEvent is emitted when a coupon is applied
type CheckoutCouponAppliedEvent struct {
	SessionID    string
	OrderID      int64
	CouponCode   string
	DiscountAmount decimal.Decimal
	OccurredOn   time.Time
}

func (e CheckoutCouponAppliedEvent) EventType() string {
	return "checkout.coupon.applied"
}

func (e CheckoutCouponAppliedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutCouponRemovedEvent is emitted when a coupon is removed
type CheckoutCouponRemovedEvent struct {
	SessionID  string
	OrderID    int64
	CouponCode string
	OccurredOn time.Time
}

func (e CheckoutCouponRemovedEvent) EventType() string {
	return "checkout.coupon.removed"
}

func (e CheckoutCouponRemovedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutTotalsUpdatedEvent is emitted when totals are recalculated
type CheckoutTotalsUpdatedEvent struct {
	SessionID      string
	OrderID        int64
	Subtotal       decimal.Decimal
	ShippingCost   decimal.Decimal
	TaxAmount      decimal.Decimal
	DiscountAmount decimal.Decimal
	TotalAmount    decimal.Decimal
	OccurredOn     time.Time
}

func (e CheckoutTotalsUpdatedEvent) EventType() string {
	return "checkout.totals.updated"
}

func (e CheckoutTotalsUpdatedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutReadyForSubmissionEvent is emitted when checkout is ready to submit
type CheckoutReadyForSubmissionEvent struct {
	SessionID   string
	OrderID     int64
	TotalAmount decimal.Decimal
	OccurredOn  time.Time
}

func (e CheckoutReadyForSubmissionEvent) EventType() string {
	return "checkout.ready_for_submission"
}

func (e CheckoutReadyForSubmissionEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutSubmittedEvent is emitted when checkout is submitted
type CheckoutSubmittedEvent struct {
	SessionID   string
	OrderID     int64
	CustomerID  *string
	Email       string
	TotalAmount decimal.Decimal
	OccurredOn  time.Time
}

func (e CheckoutSubmittedEvent) EventType() string {
	return "checkout.submitted"
}

func (e CheckoutSubmittedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutConfirmedEvent is emitted when checkout is confirmed after payment
type CheckoutConfirmedEvent struct {
	SessionID   string
	OrderID     int64
	CustomerID  *string
	Email       string
	TotalAmount decimal.Decimal
	OccurredOn  time.Time
}

func (e CheckoutConfirmedEvent) EventType() string {
	return "checkout.confirmed"
}

func (e CheckoutConfirmedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutCancelledEvent is emitted when checkout is cancelled
type CheckoutCancelledEvent struct {
	SessionID  string
	OrderID    int64
	Reason     string
	OccurredOn time.Time
}

func (e CheckoutCancelledEvent) EventType() string {
	return "checkout.cancelled"
}

func (e CheckoutCancelledEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// CheckoutExpiredEvent is emitted when checkout session expires
type CheckoutExpiredEvent struct {
	SessionID  string
	OrderID    int64
	ExpiresAt  time.Time
	OccurredOn time.Time
}

func (e CheckoutExpiredEvent) EventType() string {
	return "checkout.expired"
}

func (e CheckoutExpiredEvent) OccurredAt() time.Time {
	return e.OccurredOn
}
