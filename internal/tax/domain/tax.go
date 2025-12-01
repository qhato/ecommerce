package domain

import "time"

// TaxDetail represents a tax detail entry, often linked to a tax rate configuration.
type TaxDetail struct {
	ID               int64
	Amount           float64 // From blc_tax_detail.amount (numeric(19,5))
	TaxCountry       string  // From blc_tax_detail.tax_country
	JurisdictionName string  // From blc_tax_detail.jurisdiction_name
	Rate             float64 // From blc_tax_detail.rate (numeric(19,5))
	TaxRegion        string  // From blc_tax_detail.tax_region
	TaxName          string  // From blc_tax_detail.tax_name
	Type             string  // From blc_tax_detail.type
	CurrencyCode     string  // From blc_tax_detail.currency_code
	ModuleConfigID   *int64  // From blc_tax_detail.module_config_id (int8 NULL)
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewTaxDetail creates a new TaxDetail.
func NewTaxDetail(
	amount float64,
	taxCountry, jurisdictionName string,
	rate float64,
	taxRegion, taxName, taxType, currencyCode string,
) (*TaxDetail, error) {
	if amount < 0 {
		return nil, NewDomainError("Tax amount cannot be negative")
	}
	if rate < 0 {
		return nil, NewDomainError("Tax rate cannot be negative")
	}
	if taxCountry == "" || jurisdictionName == "" || taxName == "" || taxType == "" || currencyCode == "" {
		return nil, NewDomainError("TaxCountry, JurisdictionName, TaxName, Type, and CurrencyCode cannot be empty")
	}
	now := time.Now()
	return &TaxDetail{
		Amount:           amount,
		TaxCountry:       taxCountry,
		JurisdictionName: jurisdictionName,
		Rate:             rate,
		TaxRegion:        taxRegion,
		TaxName:          taxName,
		Type:             taxType,
		CurrencyCode:     currencyCode,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// UpdateDetails updates the tax detail fields.
func (td *TaxDetail) UpdateDetails(amount float64, taxCountry, jurisdictionName, taxRegion, taxName, taxType, currencyCode string, rate float64) {
	td.Amount = amount
	td.TaxCountry = taxCountry
	td.JurisdictionName = jurisdictionName
	td.Rate = rate
	td.TaxRegion = taxRegion
	td.TaxName = taxName
	td.Type = taxType
	td.CurrencyCode = currencyCode
	td.UpdatedAt = time.Now()
}

// SetModuleConfigID sets the associated module configuration ID.
func (td *TaxDetail) SetModuleConfigID(moduleConfigID int64) {
	td.ModuleConfigID = &moduleConfigID
	td.UpdatedAt = time.Now()
}

// DomainError represents a business rule validation error within the domain.
type DomainError struct {
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new DomainError.
func NewDomainError(message string) error {
	return &DomainError{Message: message}
}
