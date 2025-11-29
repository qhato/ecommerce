package domain

import (
	"time"
)

// QualCritOfferXref represents a cross-reference between an Offer and OfferItemCriteria
// for qualifying items (i.e., conditions that must be met for the offer to apply)
type QualCritOfferXref struct {
	ID                  int64
	OfferID             int64
	OfferItemCriteriaID int64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NewQualCritOfferXref creates a new QualCritOfferXref
func NewQualCritOfferXref(offerID, offerItemCriteriaID int64) (*QualCritOfferXref, error) {
	if offerID == 0 {
		return nil, NewDomainError("OfferID cannot be zero for QualCritOfferXref")
	}
	if offerItemCriteriaID == 0 {
		return nil, NewDomainError("OfferItemCriteriaID cannot be zero for QualCritOfferXref")
	}

	now := time.Now()
	return &QualCritOfferXref{
		OfferID:             offerID,
		OfferItemCriteriaID: offerItemCriteriaID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}
