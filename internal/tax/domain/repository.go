package domain

import (
	"context"
)

// TaxDetailRepository provides an interface for managing TaxDetails.
type TaxDetailRepository interface {
	// Save stores a new tax detail or updates an existing one.
	Save(ctx context.Context, taxDetail *TaxDetail) error

	// FindByID retrieves a tax detail by its unique identifier.
	FindByID(ctx context.Context, id int64) (*TaxDetail, error)

	// FindApplicableTaxDetails retrieves tax details applicable to a given country, region, and type.
	// This will need to be refined based on how Broadleaf manages tax applicability.
	FindApplicableTaxDetails(ctx context.Context, taxCountry, taxRegion, taxType string) ([]*TaxDetail, error)

	// Delete removes a tax detail by its unique identifier.
	Delete(ctx context.Context, id int64) error
}