package domain

import (
	"time"
)

// OfferCode represents a promotional code associated with an offer
type OfferCode struct {
	ID           int64
	OfferID      int64      // Foreign key to the Offer it's associated with
	Code         string     // The actual promotional code string
	MaxUses      *int       // From blc_offer_code.max_uses
	Uses         int        // From blc_offer_code.uses
	EmailAddress *string    // From blc_offer_code.email_address
	StartDate    *time.Time // From blc_offer_code.start_date
	EndDate      *time.Time // From blc_offer_code.end_date
	Archived     bool       // From blc_offer_code.archived (bpchar(1) 'Y'/'N')
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewOfferCode creates a new OfferCode
func NewOfferCode(offerID int64, code string) (*OfferCode, error) {
	if offerID == 0 {
		return nil, NewDomainError("OfferID cannot be zero for OfferCode")
	}
	if code == "" {
		return nil, NewDomainError("Code cannot be empty for OfferCode")
	}

	now := time.Now()
	return &OfferCode{
		OfferID:      offerID,
		Code:         code,
		Uses:         0,
		Archived:     false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// SetMaxUses sets the maximum number of uses for this specific offer code
func (oc *OfferCode) SetMaxUses(maxUses int) {
	oc.MaxUses = &maxUses
	oc.UpdatedAt = time.Now()
}

// IncrementUses increments the number of times this offer code has been used
func (oc *OfferCode) IncrementUses() {
	oc.Uses++
	oc.UpdatedAt = time.Now()
}

// SetEmailAddress sets the email address associated with this offer code
func (oc *OfferCode) SetEmailAddress(email string) {
	oc.EmailAddress = &email
	oc.UpdatedAt = time.Now()
}

// SetValidityPeriod sets the start and end dates for the offer code's validity
func (oc *OfferCode) SetValidityPeriod(startDate, endDate *time.Time) {
	oc.StartDate = startDate
	oc.EndDate = endDate
	oc.UpdatedAt = time.Now()
}

// IsActive checks if the offer code is currently active and usable
func (oc *OfferCode) IsActive() bool {
	if oc.Archived {
		return false
	}
	// Check if max uses exceeded (if set)
	if oc.MaxUses != nil && oc.Uses >= *oc.MaxUses {
		return false
	}
	// Check validity period
	now := time.Now()
	if oc.StartDate != nil && now.Before(*oc.StartDate) {
		return false
	}
	if oc.EndDate != nil && now.After(*oc.EndDate) {
		return false
	}
	return true
}
