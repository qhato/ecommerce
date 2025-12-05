package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// CheckoutState represents the current state of the checkout process
type CheckoutState string

const (
	CheckoutStateInitiated           CheckoutState = "INITIATED"
	CheckoutStateCustomerInfoAdded   CheckoutState = "CUSTOMER_INFO_ADDED"
	CheckoutStateShippingInfoAdded   CheckoutState = "SHIPPING_INFO_ADDED"
	CheckoutStateShippingSelected    CheckoutState = "SHIPPING_SELECTED"
	CheckoutStateBillingInfoAdded    CheckoutState = "BILLING_INFO_ADDED"
	CheckoutStatePaymentInfoAdded    CheckoutState = "PAYMENT_INFO_ADDED"
	CheckoutStateReadyForSubmission  CheckoutState = "READY_FOR_SUBMISSION"
	CheckoutStateSubmitted           CheckoutState = "SUBMITTED"
	CheckoutStateConfirmed           CheckoutState = "CONFIRMED"
	CheckoutStateCancelled           CheckoutState = "CANCELLED"
	CheckoutStateExpired             CheckoutState = "EXPIRED"
)

// CheckoutSession represents a checkout session
// Business Logic: Track checkout progress and state
type CheckoutSession struct {
	ID                   string
	OrderID              int64
	CustomerID           *string
	Email                string
	IsGuestCheckout      bool
	State                CheckoutState
	CurrentStep          int
	CompletedSteps       []string
	ShippingAddressID    *int64
	BillingAddressID     *int64
	ShippingMethodID     *string
	PaymentMethodID      *int64
	Subtotal             decimal.Decimal
	ShippingCost         decimal.Decimal
	TaxAmount            decimal.Decimal
	DiscountAmount       decimal.Decimal
	TotalAmount          decimal.Decimal
	CouponCodes          []string
	CustomerNotes        string
	SessionData          map[string]interface{} // Additional session data
	ExpiresAt            time.Time
	LastActivityAt       time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
	SubmittedAt          *time.Time
	ConfirmedAt          *time.Time
}

// NewCheckoutSession creates a new checkout session
func NewCheckoutSession(orderID int64, email string, isGuest bool) (*CheckoutSession, error) {
	if orderID == 0 {
		return nil, ErrOrderIDRequired
	}
	if email == "" {
		return nil, ErrEmailRequired
	}

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // 24 hour default expiration

	return &CheckoutSession{
		ID:              generateCheckoutID(),
		OrderID:         orderID,
		Email:           email,
		IsGuestCheckout: isGuest,
		State:           CheckoutStateInitiated,
		CurrentStep:     1,
		CompletedSteps:  make([]string, 0),
		Subtotal:        decimal.Zero,
		ShippingCost:    decimal.Zero,
		TaxAmount:       decimal.Zero,
		DiscountAmount:  decimal.Zero,
		TotalAmount:     decimal.Zero,
		CouponCodes:     make([]string, 0),
		SessionData:     make(map[string]interface{}),
		ExpiresAt:       expiresAt,
		LastActivityAt:  now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// SetCustomerInfo sets customer information
func (cs *CheckoutSession) SetCustomerInfo(customerID *string, email string) error {
	if email == "" {
		return ErrEmailRequired
	}

	cs.CustomerID = customerID
	cs.Email = email
	cs.State = CheckoutStateCustomerInfoAdded
	cs.markStepCompleted("customer_info")
	cs.UpdatedAt = time.Now()
	cs.LastActivityAt = time.Now()

	return nil
}

// SetShippingAddress sets the shipping address
func (cs *CheckoutSession) SetShippingAddress(addressID int64) error {
	if addressID == 0 {
		return ErrAddressIDRequired
	}

	cs.ShippingAddressID = &addressID
	cs.State = CheckoutStateShippingInfoAdded
	cs.markStepCompleted("shipping_address")
	cs.UpdatedAt = time.Now()
	cs.LastActivityAt = time.Now()

	return nil
}

// SetShippingMethod sets the shipping method
func (cs *CheckoutSession) SetShippingMethod(methodID string, cost decimal.Decimal) error {
	if methodID == "" {
		return ErrShippingMethodRequired
	}

	cs.ShippingMethodID = &methodID
	cs.ShippingCost = cost
	cs.State = CheckoutStateShippingSelected
	cs.markStepCompleted("shipping_method")
	cs.recalculateTotal()
	cs.UpdatedAt = time.Now()
	cs.LastActivityAt = time.Now()

	return nil
}

// SetBillingAddress sets the billing address
func (cs *CheckoutSession) SetBillingAddress(addressID int64) error {
	if addressID == 0 {
		return ErrAddressIDRequired
	}

	cs.BillingAddressID = &addressID
	cs.State = CheckoutStateBillingInfoAdded
	cs.markStepCompleted("billing_address")
	cs.UpdatedAt = time.Now()
	cs.LastActivityAt = time.Now()

	return nil
}

// UseSameAddressForBilling uses shipping address for billing
func (cs *CheckoutSession) UseSameAddressForBilling() error {
	if cs.ShippingAddressID == nil {
		return ErrShippingAddressRequired
	}

	cs.BillingAddressID = cs.ShippingAddressID
	cs.State = CheckoutStateBillingInfoAdded
	cs.markStepCompleted("billing_address")
	cs.UpdatedAt = time.Now()
	cs.LastActivityAt = time.Now()

	return nil
}

// SetPaymentMethod sets the payment method
func (cs *CheckoutSession) SetPaymentMethod(paymentMethodID int64) error {
	if paymentMethodID == 0 {
		return ErrPaymentMethodRequired
	}

	cs.PaymentMethodID = &paymentMethodID
	cs.State = CheckoutStatePaymentInfoAdded
	cs.markStepCompleted("payment_method")
	cs.UpdatedAt = time.Now()
	cs.LastActivityAt = time.Now()

	return nil
}

// SetPricing sets the pricing information
func (cs *CheckoutSession) SetPricing(subtotal, shippingCost, taxAmount, discountAmount decimal.Decimal) {
	cs.Subtotal = subtotal
	cs.ShippingCost = shippingCost
	cs.TaxAmount = taxAmount
	cs.DiscountAmount = discountAmount
	cs.recalculateTotal()
	cs.UpdatedAt = time.Now()
}

// ApplyCouponCode applies a coupon code
func (cs *CheckoutSession) ApplyCouponCode(code string) error {
	if code == "" {
		return ErrCouponCodeRequired
	}

	// Check if already applied
	for _, c := range cs.CouponCodes {
		if c == code {
			return ErrCouponAlreadyApplied
		}
	}

	cs.CouponCodes = append(cs.CouponCodes, code)
	cs.UpdatedAt = time.Now()
	cs.LastActivityAt = time.Now()

	return nil
}

// RemoveCouponCode removes a coupon code
func (cs *CheckoutSession) RemoveCouponCode(code string) {
	newCodes := make([]string, 0)
	for _, c := range cs.CouponCodes {
		if c != code {
			newCodes = append(newCodes, c)
		}
	}
	cs.CouponCodes = newCodes
	cs.UpdatedAt = time.Now()
}

// MarkReadyForSubmission marks the checkout as ready for submission
func (cs *CheckoutSession) MarkReadyForSubmission() error {
	// Validate all required information is present
	if err := cs.ValidateForSubmission(); err != nil {
		return err
	}

	cs.State = CheckoutStateReadyForSubmission
	cs.markStepCompleted("order_review")
	cs.UpdatedAt = time.Now()
	cs.LastActivityAt = time.Now()

	return nil
}

// Submit submits the checkout
func (cs *CheckoutSession) Submit() error {
	if cs.State != CheckoutStateReadyForSubmission {
		return ErrCheckoutNotReady
	}

	now := time.Now()
	cs.State = CheckoutStateSubmitted
	cs.SubmittedAt = &now
	cs.UpdatedAt = now
	cs.LastActivityAt = now

	return nil
}

// Confirm confirms the checkout after payment
func (cs *CheckoutSession) Confirm() error {
	if cs.State != CheckoutStateSubmitted {
		return ErrCheckoutNotSubmitted
	}

	now := time.Now()
	cs.State = CheckoutStateConfirmed
	cs.ConfirmedAt = &now
	cs.UpdatedAt = now
	cs.LastActivityAt = now

	return nil
}

// Cancel cancels the checkout
func (cs *CheckoutSession) Cancel() {
	cs.State = CheckoutStateCancelled
	cs.UpdatedAt = time.Now()
}

// IsExpired checks if the session is expired
func (cs *CheckoutSession) IsExpired() bool {
	return time.Now().After(cs.ExpiresAt)
}

// ExtendExpiration extends the session expiration
func (cs *CheckoutSession) ExtendExpiration(duration time.Duration) {
	cs.ExpiresAt = time.Now().Add(duration)
	cs.UpdatedAt = time.Now()
}

// ValidateForSubmission validates that all required information is present
func (cs *CheckoutSession) ValidateForSubmission() error {
	if cs.Email == "" {
		return ErrEmailRequired
	}
	if cs.ShippingAddressID == nil {
		return ErrShippingAddressRequired
	}
	if cs.ShippingMethodID == nil {
		return ErrShippingMethodRequired
	}
	if cs.BillingAddressID == nil {
		return ErrBillingAddressRequired
	}
	if cs.PaymentMethodID == nil {
		return ErrPaymentMethodRequired
	}
	if cs.TotalAmount.IsZero() || cs.TotalAmount.IsNegative() {
		return ErrInvalidTotalAmount
	}

	return nil
}

// GetProgress returns the checkout progress as a percentage
func (cs *CheckoutSession) GetProgress() int {
	totalSteps := 6 // customer, shipping address, shipping method, billing, payment, review
	completed := len(cs.CompletedSteps)

	if completed >= totalSteps {
		return 100
	}

	return (completed * 100) / totalSteps
}

// Private helper methods

func (cs *CheckoutSession) markStepCompleted(step string) {
	// Check if already completed
	for _, s := range cs.CompletedSteps {
		if s == step {
			return
		}
	}
	cs.CompletedSteps = append(cs.CompletedSteps, step)
	cs.CurrentStep = len(cs.CompletedSteps) + 1
}

func (cs *CheckoutSession) recalculateTotal() {
	cs.TotalAmount = cs.Subtotal.
		Add(cs.ShippingCost).
		Add(cs.TaxAmount).
		Sub(cs.DiscountAmount)
}

func generateCheckoutID() string {
	// Generate a unique checkout ID (timestamp + random)
	return "CHK-" + time.Now().Format("20060102150405")
}
