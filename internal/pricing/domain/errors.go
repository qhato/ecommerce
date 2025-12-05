package domain

import "errors"

// Domain errors for Pricing
var (
	// PriceList errors
	ErrPriceListNameRequired = errors.New("price list name is required")
	ErrPriceListCodeRequired = errors.New("price list code is required")
	ErrCurrencyRequired      = errors.New("currency is required")
	ErrPriceListNotFound     = errors.New("price list not found")
	ErrPriceListCodeExists   = errors.New("price list code already exists")

	// PriceListItem errors
	ErrPriceListIDRequired        = errors.New("price list ID is required")
	ErrSKUIDRequired              = errors.New("SKU ID is required")
	ErrPriceCannotBeNegative      = errors.New("price cannot be negative")
	ErrMinQuantityCannotBeNegative = errors.New("minimum quantity cannot be negative")
	ErrMaxQuantityLessThanMin     = errors.New("maximum quantity cannot be less than minimum quantity")
	ErrPriceListItemNotFound      = errors.New("price list item not found")

	// PricingRule errors
	ErrPricingRuleNameRequired = errors.New("pricing rule name is required")
	ErrPricingRuleNotFound     = errors.New("pricing rule not found")

	// Pricing errors
	ErrNoPriceFound             = errors.New("no price found for SKU")
	ErrInvalidCurrency          = errors.New("invalid currency")
	ErrInvalidQuantity          = errors.New("invalid quantity")
	ErrPricingContextRequired   = errors.New("pricing context is required")
	ErrNoActivePriceList        = errors.New("no active price list found")
)
