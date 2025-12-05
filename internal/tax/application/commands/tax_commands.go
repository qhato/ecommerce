package commands

import (
	"time"

	"github.com/shopspring/decimal"
)

// CreateTaxJurisdictionCommand represents a command to create a tax jurisdiction
type CreateTaxJurisdictionCommand struct {
	Code             string
	Name             string
	JurisdictionType string // FEDERAL, STATE, COUNTY, CITY, DISTRICT
	Country          string
	StateProvince    *string
	County           *string
	City             *string
	PostalCode       *string
	ParentID         *int64
	Priority         int
}

// UpdateTaxJurisdictionCommand represents a command to update a tax jurisdiction
type UpdateTaxJurisdictionCommand struct {
	ID               int64
	Name             string
	StateProvince    *string
	County           *string
	City             *string
	PostalCode       *string
	ParentID         *int64
	Priority         int
	IsActive         bool
}

// CreateTaxRateCommand represents a command to create a tax rate
type CreateTaxRateCommand struct {
	JurisdictionID    int64
	Name              string
	TaxType           string // PERCENTAGE, FLAT, COMPOUND
	Rate              decimal.Decimal
	TaxCategory       string // GENERAL, FOOD, CLOTHING, DIGITAL, SHIPPING, SERVICE, EXEMPT
	IsCompound        bool
	IsShippingTaxable bool
	MinThreshold      *decimal.Decimal
	MaxThreshold      *decimal.Decimal
	Priority          int
	StartDate         *time.Time
	EndDate           *time.Time
}

// UpdateTaxRateCommand represents a command to update a tax rate
type UpdateTaxRateCommand struct {
	ID                int64
	Name              string
	Rate              decimal.Decimal
	IsCompound        bool
	IsShippingTaxable bool
	MinThreshold      *decimal.Decimal
	MaxThreshold      *decimal.Decimal
	Priority          int
	IsActive          bool
	StartDate         *time.Time
	EndDate           *time.Time
}

// CreateTaxExemptionCommand represents a command to create a tax exemption
type CreateTaxExemptionCommand struct {
	CustomerID           string
	ExemptionCertificate string
	JurisdictionID       *int64
	TaxCategory          *string
	Reason               string
	StartDate            *time.Time
	EndDate              *time.Time
}

// UpdateTaxExemptionCommand represents a command to update a tax exemption
type UpdateTaxExemptionCommand struct {
	ID               int64
	JurisdictionID   *int64
	TaxCategory      *string
	Reason           string
	IsActive         bool
	StartDate        *time.Time
	EndDate          *time.Time
}

// DeleteTaxJurisdictionCommand represents a command to delete a tax jurisdiction
type DeleteTaxJurisdictionCommand struct {
	ID int64
}

// DeleteTaxRateCommand represents a command to delete a tax rate
type DeleteTaxRateCommand struct {
	ID int64
}

// DeleteTaxExemptionCommand represents a command to delete a tax exemption
type DeleteTaxExemptionCommand struct {
	ID int64
}

// BulkCreateTaxRatesCommand represents a command to create multiple tax rates at once
type BulkCreateTaxRatesCommand struct {
	Rates []CreateTaxRateCommand
}
