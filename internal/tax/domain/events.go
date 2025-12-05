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

// TaxJurisdictionCreatedEvent is emitted when a tax jurisdiction is created
type TaxJurisdictionCreatedEvent struct {
	JurisdictionID   int64
	Code             string
	Name             string
	JurisdictionType TaxJurisdictionType
	Country          string
	StateProvince    *string
	OccurredOn       time.Time
}

func (e TaxJurisdictionCreatedEvent) EventType() string {
	return "tax.jurisdiction.created"
}

func (e TaxJurisdictionCreatedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// TaxJurisdictionUpdatedEvent is emitted when a tax jurisdiction is updated
type TaxJurisdictionUpdatedEvent struct {
	JurisdictionID int64
	Code           string
	Name           string
	IsActive       bool
	OccurredOn     time.Time
}

func (e TaxJurisdictionUpdatedEvent) EventType() string {
	return "tax.jurisdiction.updated"
}

func (e TaxJurisdictionUpdatedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// TaxRateCreatedEvent is emitted when a tax rate is created
type TaxRateCreatedEvent struct {
	TaxRateID      int64
	JurisdictionID int64
	Name           string
	TaxType        TaxRateType
	Rate           decimal.Decimal
	TaxCategory    TaxCategory
	IsCompound     bool
	OccurredOn     time.Time
}

func (e TaxRateCreatedEvent) EventType() string {
	return "tax.rate.created"
}

func (e TaxRateCreatedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// TaxRateUpdatedEvent is emitted when a tax rate is updated
type TaxRateUpdatedEvent struct {
	TaxRateID      int64
	JurisdictionID int64
	Name           string
	Rate           decimal.Decimal
	IsActive       bool
	OccurredOn     time.Time
}

func (e TaxRateUpdatedEvent) EventType() string {
	return "tax.rate.updated"
}

func (e TaxRateUpdatedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// TaxCalculatedEvent is emitted when taxes are calculated for an order
type TaxCalculatedEvent struct {
	OrderID           *int64
	CustomerID        *string
	TotalTax          decimal.Decimal
	Subtotal          decimal.Decimal
	ShippingTax       decimal.Decimal
	ItemCount         int
	JurisdictionsUsed []string
	CalculatedAt      time.Time
	OccurredOn        time.Time
}

func (e TaxCalculatedEvent) EventType() string {
	return "tax.calculated"
}

func (e TaxCalculatedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// TaxExemptionCreatedEvent is emitted when a tax exemption is created
type TaxExemptionCreatedEvent struct {
	ExemptionID          int64
	CustomerID           string
	ExemptionCertificate string
	JurisdictionID       *int64
	TaxCategory          *TaxCategory
	Reason               string
	OccurredOn           time.Time
}

func (e TaxExemptionCreatedEvent) EventType() string {
	return "tax.exemption.created"
}

func (e TaxExemptionCreatedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// TaxExemptionUpdatedEvent is emitted when a tax exemption is updated
type TaxExemptionUpdatedEvent struct {
	ExemptionID int64
	CustomerID  string
	IsActive    bool
	OccurredOn  time.Time
}

func (e TaxExemptionUpdatedEvent) EventType() string {
	return "tax.exemption.updated"
}

func (e TaxExemptionUpdatedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// TaxExemptionAppliedEvent is emitted when a tax exemption is applied to a calculation
type TaxExemptionAppliedEvent struct {
	ExemptionID          int64
	CustomerID           string
	OrderID              *int64
	ExemptionCertificate string
	JurisdictionCode     string
	TaxCategory          TaxCategory
	ExemptAmount         decimal.Decimal
	OccurredOn           time.Time
}

func (e TaxExemptionAppliedEvent) EventType() string {
	return "tax.exemption.applied"
}

func (e TaxExemptionAppliedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}
