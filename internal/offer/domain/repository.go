package domain

import (
	"context"
)

// OfferRepository provides an interface for managing Offers in the catalog.
type OfferRepository interface {
	// Save stores a new offer or updates an existing one.
	Save(ctx context.Context, offer *Offer) error

	// FindByID retrieves an offer by its unique identifier.
	FindByID(ctx context.Context, id int64) (*Offer, error)

	// FindActiveOffers retrieves all currently active offers.
	FindActiveOffers(ctx context.Context) ([]*Offer, error)

	// Delete removes an offer by its unique identifier.
	Delete(ctx context.Context, id int64) error
}

// OfferCodeRepository provides an interface for managing OfferCodes.
type OfferCodeRepository interface {
	// Save stores a new offer code or updates an existing one.
	Save(ctx context.Context, offerCode *OfferCode) error

	// FindByID retrieves an offer code by its unique identifier.
	FindByID(ctx context.Context, id int64) (*OfferCode, error)

	// FindByCode retrieves an offer code by its code string.
	FindByCode(ctx context.Context, code string) (*OfferCode, error)

	// FindByOfferID retrieves all offer codes associated with a given offer ID.
	FindByOfferID(ctx context.Context, offerID int64) ([]*OfferCode, error)

	// Delete removes an offer code by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteByOfferID removes all offer codes associated with a given offer ID.
	DeleteByOfferID(ctx context.Context, offerID int64) error
}

// OfferItemCriteriaRepository provides an interface for managing OfferItemCriteria.
type OfferItemCriteriaRepository interface {
	// Save stores a new offer item criteria or updates an existing one.
	Save(ctx context.Context, criteria *OfferItemCriteria) error

	// FindByID retrieves an offer item criteria by its unique identifier.
	FindByID(ctx context.Context, id int64) (*OfferItemCriteria, error)

	// FindAll retrieves all offer item criteria.
	FindAll(ctx context.Context) ([]*OfferItemCriteria, error)

	// Delete removes an offer item criteria by its unique identifier.
	Delete(ctx context.Context, id int64) error
}

// OfferRuleRepository provides an interface for managing OfferRules.
type OfferRuleRepository interface {
	// Save stores a new offer rule or updates an existing one.
	Save(ctx context.Context, rule *OfferRule) error

	// FindByID retrieves an offer rule by its unique identifier.
	FindByID(ctx context.Context, id int64) (*OfferRule, error)

	// FindAll retrieves all offer rules.
	FindAll(ctx context.Context) ([]*OfferRule, error)

	// Delete removes an offer rule by its unique identifier.
	Delete(ctx context.Context, id int64) error
}

// OfferPriceDataRepository provides an interface for managing OfferPriceData.
type OfferPriceDataRepository interface {
	// Save stores new offer price data or updates an existing one.
	Save(ctx context.Context, priceData *OfferPriceData) error

	// FindByID retrieves offer price data by its unique identifier.
	FindByID(ctx context.Context, id int64) (*OfferPriceData, error)

	// FindByOfferID retrieves all offer price data associated with a given offer ID.
	FindByOfferID(ctx context.Context, offerID int64) ([]*OfferPriceData, error)

	// FindActiveByOfferID retrieves all currently active offer price data for a given offer ID. (Based on StartDate/EndDate)
	FindActiveByOfferID(ctx context.Context, offerID int64) ([]*OfferPriceData, error)

	// Delete removes offer price data by its unique identifier.
	Delete(ctx context.Context, id int64) error
}

// QualCritOfferXrefRepository provides an interface for managing QualCritOfferXref (qualifying criteria for offers).
type QualCritOfferXrefRepository interface {
	// Save stores a new qualifying criteria xref or updates an existing one.
	Save(ctx context.Context, xref *QualCritOfferXref) error

	// FindByID retrieves a qualifying criteria xref by its unique identifier.
	FindByID(ctx context.Context, id int64) (*QualCritOfferXref, error)

	// FindByOfferID retrieves all qualifying criteria xrefs for a given offer ID.
	FindByOfferID(ctx context.Context, offerID int64) ([]*QualCritOfferXref, error)

	// FindByOfferItemCriteriaID retrieves all qualifying criteria xrefs for a given offer item criteria ID.
	FindByOfferItemCriteriaID(ctx context.Context, offerItemCriteriaID int64) ([]*QualCritOfferXref, error)

	// Delete removes a qualifying criteria xref by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteByOfferID removes all qualifying criteria xrefs for a given offer ID.
	DeleteByOfferID(ctx context.Context, offerID int64) error

	// DeleteByOfferItemCriteriaID removes all qualifying criteria xrefs for a given offer item criteria ID.
	DeleteByOfferItemCriteriaID(ctx context.Context, offerItemCriteriaID int64) error

	// RemoveQualCritOfferXref removes a specific qualifying criteria xref by offer ID and offer item criteria ID.
	RemoveQualCritOfferXref(ctx context.Context, offerID, offerItemCriteriaID int64) error
}

// TarCritOfferXrefRepository provides an interface for managing TarCritOfferXref (target criteria for offers).
type TarCritOfferXrefRepository interface {
	// Save stores a new target criteria xref or updates an existing one.
	Save(ctx context.Context, xref *TarCritOfferXref) error

	// FindByID retrieves a target criteria xref by its unique identifier.
	FindByID(ctx context.Context, id int64) (*TarCritOfferXref, error)

	// FindByOfferID retrieves all target criteria xrefs for a given offer ID.
	FindByOfferID(ctx context.Context, offerID int64) ([]*TarCritOfferXref, error)

	// FindByOfferItemCriteriaID retrieves all target criteria xrefs for a given offer item criteria ID.
	FindByOfferItemCriteriaID(ctx context.Context, offerItemCriteriaID int64) ([]*TarCritOfferXref, error)

	// Delete removes a target criteria xref by its unique identifier.
	Delete(ctx context.Context, id int64) error

	// DeleteByOfferID removes all target criteria xrefs for a given offer ID.
	DeleteByOfferID(ctx context.Context, offerID int64) error

	// DeleteByOfferItemCriteriaID removes all target criteria xrefs for a given offer item criteria ID.
	DeleteByOfferItemCriteriaID(ctx context.Context, offerItemCriteriaID int64) error

	// RemoveTarCritOfferXref removes a specific target criteria xref by offer ID and offer item criteria ID.
	RemoveTarCritOfferXref(ctx context.Context, offerID, offerItemCriteriaID int64) error
}
