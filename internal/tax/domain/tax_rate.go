package domain

import "time"

// TaxType represents the type of tax
type TaxType string

const (
	TaxTypeSales    TaxType = "SALES"    // Sales tax
	TaxTypeVAT      TaxType = "VAT"      // Value Added Tax
	TaxTypeGST      TaxType = "GST"      // Goods and Services Tax
	TaxTypeExcise   TaxType = "EXCISE"   // Excise tax
	TaxTypeCustoms  TaxType = "CUSTOMS"  // Customs duty
	TaxTypeProperty TaxType = "PROPERTY" // Property tax
)

// TaxRate represents a tax rate configuration for a jurisdiction
type TaxRate struct {
	ID               int64
	Country          string  // Country code (e.g., "US", "CA", "GB")
	Region           string  // State/Province/Region code (e.g., "CA", "NY", "ON")
	JurisdictionName string  // Full jurisdiction name (e.g., "California", "New York")
	TaxName          string  // Name of the tax (e.g., "State Sales Tax", "VAT")
	TaxType          TaxType // Type of tax
	Rate             float64 // Tax rate as decimal (e.g., 0.0825 for 8.25%)
	Priority         int     // Priority for applying multiple taxes (lower = applied first)
	Active           bool    // Whether this tax rate is currently active
	EffectiveFrom    time.Time
	EffectiveTo      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewTaxRate creates a new TaxRate
func NewTaxRate(
	country, region, jurisdictionName, taxName string,
	taxType TaxType,
	rate float64,
	effectiveFrom time.Time,
) (*TaxRate, error) {
	if country == "" {
		return nil, NewDomainError("Country cannot be empty for TaxRate")
	}
	if jurisdictionName == "" {
		return nil, NewDomainError("JurisdictionName cannot be empty for TaxRate")
	}
	if taxName == "" {
		return nil, NewDomainError("TaxName cannot be empty for TaxRate")
	}
	if rate < 0 || rate > 1 {
		return nil, NewDomainError("Tax rate must be between 0 and 1")
	}

	now := time.Now()
	return &TaxRate{
		Country:          country,
		Region:           region,
		JurisdictionName: jurisdictionName,
		TaxName:          taxName,
		TaxType:          taxType,
		Rate:             rate,
		Priority:         50, // Default priority
		Active:           true,
		EffectiveFrom:    effectiveFrom,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// IsEffective checks if the tax rate is effective for the given date
func (tr *TaxRate) IsEffective(date time.Time) bool {
	if !tr.Active {
		return false
	}

	if date.Before(tr.EffectiveFrom) {
		return false
	}

	if tr.EffectiveTo != nil && date.After(*tr.EffectiveTo) {
		return false
	}

	return true
}

// SetEffectiveTo sets the end date for this tax rate
func (tr *TaxRate) SetEffectiveTo(date time.Time) {
	tr.EffectiveTo = &date
	tr.UpdatedAt = time.Now()
}

// Deactivate marks the tax rate as inactive
func (tr *TaxRate) Deactivate() {
	tr.Active = false
	tr.UpdatedAt = time.Now()
}

// Activate marks the tax rate as active
func (tr *TaxRate) Activate() {
	tr.Active = true
	tr.UpdatedAt = time.Now()
}

// UpdateRate updates the tax rate
func (tr *TaxRate) UpdateRate(rate float64) error {
	if rate < 0 || rate > 1 {
		return NewDomainError("Tax rate must be between 0 and 1")
	}
	tr.Rate = rate
	tr.UpdatedAt = time.Now()
	return nil
}

// SetPriority sets the priority for applying this tax
func (tr *TaxRate) SetPriority(priority int) {
	tr.Priority = priority
	tr.UpdatedAt = time.Now()
}

// TaxJurisdiction represents a taxable jurisdiction
type TaxJurisdiction struct {
	ID               int64
	Country          string
	Region           string
	JurisdictionName string
	PostalCode       *string // Optional postal code for specific areas
	City             *string // Optional city for specific areas
	County           *string // Optional county for specific areas
	Active           bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewTaxJurisdiction creates a new TaxJurisdiction
func NewTaxJurisdiction(country, region, jurisdictionName string) (*TaxJurisdiction, error) {
	if country == "" || jurisdictionName == "" {
		return nil, NewDomainError("Country and JurisdictionName cannot be empty")
	}

	now := time.Now()
	return &TaxJurisdiction{
		Country:          country,
		Region:           region,
		JurisdictionName: jurisdictionName,
		Active:           true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// MatchesAddress checks if this jurisdiction matches the given address
func (tj *TaxJurisdiction) MatchesAddress(country, region, postalCode, city, county string) bool {
	if !tj.Active {
		return false
	}

	// Country must match
	if tj.Country != country {
		return false
	}

	// Region must match if specified
	if tj.Region != "" && tj.Region != region {
		return false
	}

	// Check optional filters
	if tj.PostalCode != nil && *tj.PostalCode != postalCode {
		return false
	}

	if tj.City != nil && *tj.City != city {
		return false
	}

	if tj.County != nil && *tj.County != county {
		return false
	}

	return true
}

// TaxExemption represents a tax exemption for a customer or category
type TaxExemption struct {
	ID             int64
	CustomerID     *string   // Optional: customer-specific exemption
	CategoryID     *string   // Optional: category-specific exemption
	ProductID      *string   // Optional: product-specific exemption
	ExemptionCode  string    // Code identifying the exemption (e.g., "RESALE", "NONPROFIT")
	Reason         string    // Reason for exemption
	DocumentNumber *string   // Certificate or document number
	Country        string    // Country where exemption applies
	Region         *string   // Optional: region where exemption applies
	ExemptTaxTypes []TaxType // Which tax types are exempt
	Active         bool
	EffectiveFrom  time.Time
	EffectiveTo    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewTaxExemption creates a new TaxExemption
func NewTaxExemption(
	exemptionCode, reason, country string,
	exemptTaxTypes []TaxType,
	effectiveFrom time.Time,
) (*TaxExemption, error) {
	if exemptionCode == "" || reason == "" || country == "" {
		return nil, NewDomainError("ExemptionCode, Reason, and Country cannot be empty")
	}

	if len(exemptTaxTypes) == 0 {
		return nil, NewDomainError("At least one tax type must be specified for exemption")
	}

	now := time.Now()
	return &TaxExemption{
		ExemptionCode:  exemptionCode,
		Reason:         reason,
		Country:        country,
		ExemptTaxTypes: exemptTaxTypes,
		Active:         true,
		EffectiveFrom:  effectiveFrom,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// IsEffective checks if the exemption is effective for the given date
func (te *TaxExemption) IsEffective(date time.Time) bool {
	if !te.Active {
		return false
	}

	if date.Before(te.EffectiveFrom) {
		return false
	}

	if te.EffectiveTo != nil && date.After(*te.EffectiveTo) {
		return false
	}

	return true
}

// AppliesToCustomer checks if this exemption applies to a specific customer
func (te *TaxExemption) AppliesToCustomer(customerID string) bool {
	if te.CustomerID == nil {
		return false
	}
	return *te.CustomerID == customerID
}

// AppliesToProduct checks if this exemption applies to a specific product
func (te *TaxExemption) AppliesToProduct(productID string) bool {
	if te.ProductID == nil {
		return false
	}
	return *te.ProductID == productID
}

// AppliesToCategory checks if this exemption applies to a specific category
func (te *TaxExemption) AppliesToCategory(categoryID string) bool {
	if te.CategoryID == nil {
		return false
	}
	return *te.CategoryID == categoryID
}

// ExemptsTaxType checks if this exemption covers a specific tax type
func (te *TaxExemption) ExemptsTaxType(taxType TaxType) bool {
	for _, exemptType := range te.ExemptTaxTypes {
		if exemptType == taxType {
			return true
		}
	}
	return false
}
