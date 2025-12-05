package domain

import (
	"context"

	"github.com/shopspring/decimal"
)

// TaxJurisdictionRepository defines the interface for tax jurisdiction persistence
type TaxJurisdictionRepository interface {
	// Create creates a new tax jurisdiction
	Create(ctx context.Context, jurisdiction *TaxJurisdiction) error

	// Update updates an existing tax jurisdiction
	Update(ctx context.Context, jurisdiction *TaxJurisdiction) error

	// FindByID finds a jurisdiction by ID
	FindByID(ctx context.Context, id int64) (*TaxJurisdiction, error)

	// FindByCode finds a jurisdiction by code
	FindByCode(ctx context.Context, code string) (*TaxJurisdiction, error)

	// FindByLocation finds all jurisdictions matching a location
	FindByLocation(ctx context.Context, country, stateProvince, county, city, postalCode string) ([]*TaxJurisdiction, error)

	// FindAll finds all jurisdictions with optional filters
	FindAll(ctx context.Context, activeOnly bool) ([]*TaxJurisdiction, error)

	// FindByCountry finds all jurisdictions in a country
	FindByCountry(ctx context.Context, country string, activeOnly bool) ([]*TaxJurisdiction, error)

	// FindChildren finds all child jurisdictions of a parent
	FindChildren(ctx context.Context, parentID int64) ([]*TaxJurisdiction, error)

	// Delete deletes a jurisdiction
	Delete(ctx context.Context, id int64) error

	// ExistsByCode checks if a jurisdiction exists with the given code
	ExistsByCode(ctx context.Context, code string) (bool, error)
}

// TaxRateRepository defines the interface for tax rate persistence
type TaxRateRepository interface {
	// Create creates a new tax rate
	Create(ctx context.Context, rate *TaxRate) error

	// Update updates an existing tax rate
	Update(ctx context.Context, rate *TaxRate) error

	// FindByID finds a tax rate by ID
	FindByID(ctx context.Context, id int64) (*TaxRate, error)

	// FindByJurisdiction finds all tax rates for a jurisdiction
	FindByJurisdiction(ctx context.Context, jurisdictionID int64, activeOnly bool) ([]*TaxRate, error)

	// FindByJurisdictionAndCategory finds tax rates for a jurisdiction and category
	FindByJurisdictionAndCategory(ctx context.Context, jurisdictionID int64, category TaxCategory, activeOnly bool) ([]*TaxRate, error)

	// FindApplicableRates finds all applicable tax rates for a calculation
	FindApplicableRates(ctx context.Context, jurisdictionIDs []int64, category TaxCategory, activeOnly bool) ([]*TaxRate, error)

	// FindAll finds all tax rates with optional filters
	FindAll(ctx context.Context, activeOnly bool) ([]*TaxRate, error)

	// Delete deletes a tax rate
	Delete(ctx context.Context, id int64) error

	// BulkCreate creates multiple tax rates in a transaction
	BulkCreate(ctx context.Context, rates []*TaxRate) error
}

// TaxExemptionRepository defines the interface for tax exemption persistence
type TaxExemptionRepository interface {
	// Create creates a new tax exemption
	Create(ctx context.Context, exemption *TaxExemption) error

	// Update updates an existing tax exemption
	Update(ctx context.Context, exemption *TaxExemption) error

	// FindByID finds an exemption by ID
	FindByID(ctx context.Context, id int64) (*TaxExemption, error)

	// FindByCustomerID finds all exemptions for a customer
	FindByCustomerID(ctx context.Context, customerID string, activeOnly bool) ([]*TaxExemption, error)

	// FindActiveExemptions finds all currently active exemptions for a customer
	FindActiveExemptions(ctx context.Context, customerID string) ([]*TaxExemption, error)

	// FindByCustomerAndJurisdiction finds exemptions for a customer in a jurisdiction
	FindByCustomerAndJurisdiction(ctx context.Context, customerID string, jurisdictionID int64, activeOnly bool) ([]*TaxExemption, error)

	// FindByCertificate finds an exemption by certificate number
	FindByCertificate(ctx context.Context, certificate string) (*TaxExemption, error)

	// FindAll finds all exemptions with optional filters
	FindAll(ctx context.Context, activeOnly bool) ([]*TaxExemption, error)

	// Delete deletes an exemption
	Delete(ctx context.Context, id int64) error

	// ExistsByCertificate checks if an exemption exists with the given certificate
	ExistsByCertificate(ctx context.Context, certificate string) (bool, error)
}

// TaxCalculatorService defines the interface for tax calculation
type TaxCalculatorService interface {
	// Calculate calculates taxes for a request
	Calculate(ctx context.Context, request *TaxCalculationRequest) (*TaxCalculationResult, error)

	// CalculateWithExemptions calculates taxes applying customer exemptions
	CalculateWithExemptions(ctx context.Context, request *TaxCalculationRequest, exemptions []*TaxExemption) (*TaxCalculationResult, error)

	// EstimateTax provides a quick tax estimate without full calculation
	EstimateTax(ctx context.Context, address Address, subtotal decimal.Decimal) (decimal.Decimal, error)

	// ValidateAddress validates if an address has applicable tax jurisdictions
	ValidateAddress(ctx context.Context, address Address) (bool, error)
}
