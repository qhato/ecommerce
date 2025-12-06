package queries

import "github.com/shopspring/decimal"

// Carrier Queries
type GetCarrierConfigQuery struct {
	ID int64
}

type GetCarrierConfigByCarrierQuery struct {
	Carrier string
}

type GetAllCarrierConfigsQuery struct {
	EnabledOnly bool
}

// Shipping Method Queries
type GetShippingMethodQuery struct {
	ID int64
}

type GetShippingMethodsByCarrierQuery struct {
	Carrier string
}

type GetAllEnabledShippingMethodsQuery struct{}

// Shipping Band Queries
type GetShippingBandsByMethodQuery struct {
	MethodID int64
}

// Shipping Rule Queries
type GetShippingRuleQuery struct {
	ID int64
}

type GetAllEnabledShippingRulesQuery struct{}

// Shipping Calculation Queries
type CalculateShippingRatesQuery struct {
	Weight     decimal.Decimal
	OrderTotal decimal.Decimal
	Quantity   int
	Country    string
	ZipCode    string
	MethodID   *int64 // If nil, return all available methods
}

type GetAvailableShippingMethodsQuery struct {
	Country string
	ZipCode string
}
