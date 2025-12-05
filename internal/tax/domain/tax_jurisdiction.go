package domain

import (
	"time"
)

// TaxJurisdictionType defines the type of tax jurisdiction
type TaxJurisdictionType string

const (
	TaxJurisdictionTypeFederal TaxJurisdictionType = "FEDERAL" // Federal/National level
	TaxJurisdictionTypeState   TaxJurisdictionType = "STATE"   // State/Province level
	TaxJurisdictionTypeCounty  TaxJurisdictionType = "COUNTY"  // County level
	TaxJurisdictionTypeCity    TaxJurisdictionType = "CITY"    // City/Municipal level
	TaxJurisdictionTypeDistrict TaxJurisdictionType = "DISTRICT" // Special district (transit, etc.)
)

// TaxJurisdiction represents a tax jurisdiction (federal, state, county, city)
// Business Logic: Definir jurisdicciones fiscales con tasas y reglas
type TaxJurisdiction struct {
	ID               int64
	Code             string // Unique code (e.g., "US-CA", "US-CA-SF")
	Name             string
	JurisdictionType TaxJurisdictionType
	ParentID         *int64 // Parent jurisdiction (e.g., state for a city)
	Country          string // ISO 3166-1 alpha-2 country code
	StateProvince    *string
	County           *string
	City             *string
	PostalCode       *string // Can be specific postal code or pattern
	IsActive         bool
	Priority         int // Order of application (lower = applied first)
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewTaxJurisdiction creates a new TaxJurisdiction
func NewTaxJurisdiction(code, name string, jurisdictionType TaxJurisdictionType, country string) (*TaxJurisdiction, error) {
	if code == "" {
		return nil, ErrJurisdictionCodeRequired
	}
	if name == "" {
		return nil, ErrJurisdictionNameRequired
	}
	if country == "" {
		return nil, ErrCountryRequired
	}

	now := time.Now()
	return &TaxJurisdiction{
		Code:             code,
		Name:             name,
		JurisdictionType: jurisdictionType,
		Country:          country,
		IsActive:         true,
		Priority:         0,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// SetLocation sets the geographic location details
func (tj *TaxJurisdiction) SetLocation(stateProvince, county, city, postalCode *string) {
	tj.StateProvince = stateProvince
	tj.County = county
	tj.City = city
	tj.PostalCode = postalCode
	tj.UpdatedAt = time.Now()
}

// SetParent sets the parent jurisdiction
func (tj *TaxJurisdiction) SetParent(parentID int64) {
	tj.ParentID = &parentID
	tj.UpdatedAt = time.Now()
}

// Activate activates the jurisdiction
func (tj *TaxJurisdiction) Activate() {
	tj.IsActive = true
	tj.UpdatedAt = time.Now()
}

// Deactivate deactivates the jurisdiction
func (tj *TaxJurisdiction) Deactivate() {
	tj.IsActive = false
	tj.UpdatedAt = time.Now()
}

// MatchesLocation checks if this jurisdiction matches a given location
func (tj *TaxJurisdiction) MatchesLocation(country, stateProvince, county, city, postalCode string) bool {
	// Country must match
	if tj.Country != country {
		return false
	}

	// Check state/province
	if tj.StateProvince != nil && *tj.StateProvince != stateProvince {
		return false
	}

	// Check county
	if tj.County != nil && *tj.County != county {
		return false
	}

	// Check city
	if tj.City != nil && *tj.City != city {
		return false
	}

	// Check postal code (can be pattern or exact match)
	if tj.PostalCode != nil {
		if *tj.PostalCode != postalCode {
			// Could implement pattern matching here (e.g., "9411*")
			return false
		}
	}

	return true
}
