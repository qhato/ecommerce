package commands

import "github.com/shopspring/decimal"

// Carrier Commands
type CreateCarrierConfigCommand struct {
	Carrier   string
	Name      string
	IsEnabled bool
	Priority  int
	APIKey    string
	APISecret string
	AccountID string
	Config    map[string]string
}

type UpdateCarrierConfigCommand struct {
	ID        int64
	Name      string
	IsEnabled bool
	Priority  int
	APIKey    string
	APISecret string
	AccountID string
	Config    map[string]string
}

// Shipping Method Commands
type CreateShippingMethodCommand struct {
	Carrier       string
	Name          string
	Description   string
	ServiceCode   string
	EstimatedDays int
	PricingType   string
	FlatRate      decimal.Decimal
	IsEnabled     bool
}

type UpdateShippingMethodCommand struct {
	ID            int64
	Name          string
	Description   string
	ServiceCode   string
	EstimatedDays int
	PricingType   string
	FlatRate      decimal.Decimal
	IsEnabled     bool
}

type DeleteShippingMethodCommand struct {
	ID int64
}

// Shipping Band Commands
type CreateShippingBandCommand struct {
	MethodID      int64
	BandType      string
	MinValue      decimal.Decimal
	MaxValue      *decimal.Decimal
	Price         decimal.Decimal
	PercentCharge *decimal.Decimal
}

type DeleteShippingBandCommand struct {
	ID int64
}

type DeleteBandsByMethodCommand struct {
	MethodID int64
}

// Shipping Rule Commands
type CreateShippingRuleCommand struct {
	Name          string
	Description   string
	RuleType      string
	IsEnabled     bool
	Priority      int
	MinOrderValue *decimal.Decimal
	Countries     []string
	ExcludedZips  []string
	DiscountType  string
	DiscountValue decimal.Decimal
}

type UpdateShippingRuleCommand struct {
	ID            int64
	Name          string
	Description   string
	RuleType      string
	IsEnabled     bool
	Priority      int
	MinOrderValue *decimal.Decimal
	Countries     []string
	ExcludedZips  []string
	DiscountType  string
	DiscountValue decimal.Decimal
}

type DeleteShippingRuleCommand struct {
	ID int64
}

// Shipping Calculation Command
type CalculateShippingCommand struct {
	Weight      decimal.Decimal
	OrderTotal  decimal.Decimal
	Quantity    int
	Country     string
	ZipCode     string
	MethodID    *int64 // If nil, return all available methods
}
