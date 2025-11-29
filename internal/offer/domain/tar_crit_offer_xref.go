package domain

import "time"

// TarCritOfferXref represents a cross-reference between an Offer and OfferItemCriteria
// for target items (i.e., items that will receive the discount if the offer applies)
type TarCritOfferXref struct {
	ID                  int64
	OfferID             int64
	OfferItemCriteriaID int64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NewTarCritOfferXref creates a new TarCritOfferXref
func NewTarCritOfferXref(offerID, offerItemCriteriaID int64) (*TarCritOfferXref, error) {
	if offerID == 0 {
		return nil, NewDomainError("OfferID cannot be zero for TarCritOfferXref")
	}
	if offerItemCriteriaID == 0 {
		return nil, NewDomainError("OfferItemCriteriaID cannot be zero for TarCritOfferXref")
	}

	now := time.Now()
	return &TarCritOfferXref{
		OfferID:             offerID,
		OfferItemCriteriaID: offerItemCriteriaID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}
