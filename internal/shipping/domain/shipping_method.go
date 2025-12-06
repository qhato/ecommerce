package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// ShippingMethod represents a shipping method
type ShippingMethod struct {
	ID              int64
	Carrier         ShippingCarrier
	Name            string
	Description     string
	ServiceCode     string
	EstimatedDays   int
	PricingType     PricingType
	FlatRate        decimal.Decimal
	IsEnabled       bool
	Bands           []ShippingBand
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PricingType represents the type of pricing calculation
type PricingType string

const (
	PricingTypeFlat         PricingType = "FLAT"
	PricingTypeBanded       PricingType = "BANDED"
	PricingTypeRealTime     PricingType = "REALTIME"
	PricingTypePercentage   PricingType = "PERCENTAGE"
)

// ShippingBand represents a price band for banded pricing
type ShippingBand struct {
	ID            int64
	MethodID      int64
	BandType      BandType
	MinValue      decimal.Decimal
	MaxValue      *decimal.Decimal
	Price         decimal.Decimal
	PercentCharge *decimal.Decimal
	CreatedAt     time.Time
}

// BandType represents the type of banding
type BandType string

const (
	BandTypeWeight  BandType = "WEIGHT"
	BandTypePrice   BandType = "PRICE"
	BandTypeQuantity BandType = "QUANTITY"
)

// NewShippingMethod creates a new shipping method
func NewShippingMethod(carrier ShippingCarrier, name, serviceCode string, pricingType PricingType) *ShippingMethod {
	now := time.Now()
	return &ShippingMethod{
		Carrier:       carrier,
		Name:          name,
		ServiceCode:   serviceCode,
		PricingType:   pricingType,
		FlatRate:      decimal.Zero,
		IsEnabled:     false,
		Bands:         make([]ShippingBand, 0),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// CalculateRate calculates the shipping rate for a given order
func (m *ShippingMethod) CalculateRate(weight, orderTotal decimal.Decimal, quantity int) decimal.Decimal {
	switch m.PricingType {
	case PricingTypeFlat:
		return m.FlatRate
	case PricingTypeBanded:
		return m.calculateBandedRate(weight, orderTotal, decimal.NewFromInt(int64(quantity)))
	case PricingTypeRealTime:
		return decimal.Zero // Will be calculated by carrier API
	default:
		return m.FlatRate
	}
}

func (m *ShippingMethod) calculateBandedRate(weight, orderTotal, quantity decimal.Decimal) decimal.Decimal {
	for _, band := range m.Bands {
		var value decimal.Decimal
		switch band.BandType {
		case BandTypeWeight:
			value = weight
		case BandTypePrice:
			value = orderTotal
		case BandTypeQuantity:
			value = quantity
		}

		if value.GreaterThanOrEqual(band.MinValue) {
			if band.MaxValue == nil || value.LessThan(*band.MaxValue) {
				if band.PercentCharge != nil {
					return orderTotal.Mul(*band.PercentCharge).Div(decimal.NewFromInt(100))
				}
				return band.Price
			}
		}
	}
	return m.FlatRate
}

// Enable enables the shipping method
func (m *ShippingMethod) Enable() {
	m.IsEnabled = true
	m.UpdatedAt = time.Now()
}

// Disable disables the shipping method
func (m *ShippingMethod) Disable() {
	m.IsEnabled = false
	m.UpdatedAt = time.Now()
}
