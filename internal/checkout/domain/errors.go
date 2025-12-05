package domain

import "errors"

// Checkout Session Errors
var (
	ErrOrderIDRequired            = errors.New("order ID is required")
	ErrEmailRequired              = errors.New("email is required")
	ErrCheckoutSessionNotFound    = errors.New("checkout session not found")
	ErrCheckoutSessionExpired     = errors.New("checkout session has expired")
	ErrCheckoutSessionAlreadyExists = errors.New("checkout session already exists for this order")
)

// Address Errors
var (
	ErrAddressIDRequired         = errors.New("address ID is required")
	ErrShippingAddressRequired   = errors.New("shipping address is required")
	ErrBillingAddressRequired    = errors.New("billing address is required")
	ErrInvalidAddress            = errors.New("invalid address")
)

// Shipping Errors
var (
	ErrShippingMethodRequired    = errors.New("shipping method is required")
	ErrShippingMethodNotFound    = errors.New("shipping method not found")
	ErrShippingMethodUnavailable = errors.New("shipping method is not available")
	ErrShippingCostInvalid       = errors.New("shipping cost is invalid")
)

// Payment Errors
var (
	ErrPaymentMethodRequired     = errors.New("payment method is required")
	ErrPaymentMethodNotFound     = errors.New("payment method not found")
	ErrPaymentMethodInvalid      = errors.New("payment method is invalid")
	ErrPaymentMethodExpired      = errors.New("payment method has expired")
)

// Checkout State Errors
var (
	ErrCheckoutNotReady          = errors.New("checkout is not ready for submission")
	ErrCheckoutNotSubmitted      = errors.New("checkout has not been submitted")
	ErrCheckoutAlreadySubmitted  = errors.New("checkout has already been submitted")
	ErrCheckoutAlreadyConfirmed  = errors.New("checkout has already been confirmed")
	ErrCheckoutCancelled         = errors.New("checkout has been cancelled")
	ErrInvalidCheckoutState      = errors.New("invalid checkout state")
)

// Coupon Errors
var (
	ErrCouponCodeRequired        = errors.New("coupon code is required")
	ErrCouponAlreadyApplied      = errors.New("coupon has already been applied")
	ErrCouponNotFound            = errors.New("coupon not found")
	ErrCouponInvalid             = errors.New("coupon is invalid")
	ErrCouponExpired             = errors.New("coupon has expired")
)

// Validation Errors
var (
	ErrInvalidTotalAmount        = errors.New("invalid total amount")
	ErrCartIsEmpty               = errors.New("cart is empty")
	ErrItemOutOfStock            = errors.New("item is out of stock")
	ErrInsufficientStock         = errors.New("insufficient stock")
	ErrInvalidQuantity           = errors.New("invalid quantity")
)

// Step Errors
var (
	ErrStepNotFound              = errors.New("checkout step not found")
	ErrStepNotCompleted          = errors.New("checkout step is not completed")
	ErrStepAlreadyCompleted      = errors.New("checkout step is already completed")
	ErrPreviousStepNotCompleted  = errors.New("previous step must be completed first")
)

// Repository Errors
var (
	ErrRepositoryOperation       = errors.New("repository operation failed")
	ErrTransactionFailed         = errors.New("database transaction failed")
	ErrConcurrentUpdate          = errors.New("concurrent update detected")
)
