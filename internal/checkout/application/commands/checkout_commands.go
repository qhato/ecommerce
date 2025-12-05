package commands

import "github.com/shopspring/decimal"

// InitiateCheckoutCommand represents a command to start checkout
type InitiateCheckoutCommand struct {
	OrderID    int64
	CustomerID *string
	Email      string
	IsGuest    bool
}

// AddCustomerInfoCommand represents a command to add customer information
type AddCustomerInfoCommand struct {
	SessionID  string
	CustomerID *string
	Email      string
	FirstName  string
	LastName   string
	Phone      string
}

// AddShippingAddressCommand represents a command to add shipping address
type AddShippingAddressCommand struct {
	SessionID    string
	AddressLine1 string
	AddressLine2 string
	City         string
	StateProvince string
	PostalCode   string
	Country      string
	Phone        string
	FirstName    string
	LastName     string
}

// SelectShippingMethodCommand represents a command to select shipping method
type SelectShippingMethodCommand struct {
	SessionID        string
	ShippingMethodID string
}

// AddBillingAddressCommand represents a command to add billing address
type AddBillingAddressCommand struct {
	SessionID      string
	SameAsShipping bool
	AddressLine1   string
	AddressLine2   string
	City           string
	StateProvince  string
	PostalCode     string
	Country        string
	Phone          string
	FirstName      string
	LastName       string
}

// AddPaymentMethodCommand represents a command to add payment method
type AddPaymentMethodCommand struct {
	SessionID       string
	PaymentType     string // CREDIT_CARD, DEBIT_CARD, PAYPAL, etc.
	CardNumber      string
	CardHolderName  string
	ExpiryMonth     int
	ExpiryYear      int
	CVV             string
	SaveCard        bool
}

// ApplyCouponCommand represents a command to apply a coupon
type ApplyCouponCommand struct {
	SessionID  string
	CouponCode string
}

// RemoveCouponCommand represents a command to remove a coupon
type RemoveCouponCommand struct {
	SessionID  string
	CouponCode string
}

// RecalculateTotalsCommand represents a command to recalculate totals
type RecalculateTotalsCommand struct {
	SessionID string
}

// SubmitCheckoutCommand represents a command to submit the checkout
type SubmitCheckoutCommand struct {
	SessionID string
}

// ConfirmCheckoutCommand represents a command to confirm after payment
type ConfirmCheckoutCommand struct {
	SessionID     string
	PaymentID     string
	TransactionID string
}

// CancelCheckoutCommand represents a command to cancel checkout
type CancelCheckoutCommand struct {
	SessionID string
	Reason    string
}

// ExtendSessionCommand represents a command to extend session expiration
type ExtendSessionCommand struct {
	SessionID string
	Hours     int
}

// UpdateShippingCostCommand represents a command to update shipping cost
type UpdateShippingCostCommand struct {
	SessionID    string
	ShippingCost decimal.Decimal
}

// UpdateTaxCommand represents a command to update tax amount
type UpdateTaxCommand struct {
	SessionID string
	TaxAmount decimal.Decimal
}

// UpdateDiscountCommand represents a command to update discount
type UpdateDiscountCommand struct {
	SessionID      string
	DiscountAmount decimal.Decimal
}
