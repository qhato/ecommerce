package domain

import (
	"time"
)

// OfferPriceData represents price data associated with an offer for specific items or conditions
type OfferPriceData struct {
	ID              int64
	OfferID         int64      // Foreign key to the Offer it's associated with
	Amount          float64    // From blc_offer_price_data.amount (numeric(19,5))
	DiscountType    string     // From blc_offer_price_data.discount_type
	IdentifierType  string     // From blc_offer_price_data.identifier_type
	IdentifierValue string     // From blc_offer_price_data.identifier_value
	Quantity        int        // From blc_offer_price_data.quantity (int4)
	StartDate       *time.Time // From blc_offer_price_data.start_date
	EndDate         *time.Time // From blc_offer_price_data.end_date
	Archived        bool       // From blc_offer_price_data.archived (bpchar(1) 'Y'/'N')
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewOfferPriceData creates a new OfferPriceData
func NewOfferPriceData(offerID int64, amount float64, discountType, identifierType, identifierValue string, quantity int) (*OfferPriceData, error) {
	if offerID == 0 {
		return nil, NewDomainError("OfferID cannot be zero for OfferPriceData")
	}
	if amount <= 0 {
		return nil, NewDomainError("Amount must be greater than zero for OfferPriceData")
	}
	if discountType == "" {
		return nil, NewDomainError("DiscountType cannot be empty for OfferPriceData")
	}

	now := time.Now()
	return &OfferPriceData{
		OfferID:         offerID,
		Amount:          amount,
		DiscountType:    discountType,
		IdentifierType:  identifierType,
		IdentifierValue: identifierValue,
		Quantity:        quantity,
		Archived:        false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// SetValidityPeriod sets the start and end dates for the price data's validity
func (opd *OfferPriceData) SetValidityPeriod(startDate, endDate *time.Time) {
	opd.StartDate = startDate
	opd.EndDate = endDate
	opd.UpdatedAt = time.Now()
}

// UpdateData updates the amount, discount type, identifier, and quantity.
func (opd *OfferPriceData) UpdateData(amount float64, discountType, identifierType, identifierValue string, quantity int) {
	opd.Amount = amount
	opd.DiscountType = discountType
	opd.IdentifierType = identifierType
	opd.IdentifierValue = identifierValue
	opd.Quantity = quantity
	opd.UpdatedAt = time.Now()
}

// Archive marks the offer price data as archived
func (opd *OfferPriceData) Archive() {
	opd.Archived = true
	opd.UpdatedAt = time.Now()
}

// Unarchive marks the offer price data as active
func (opd *OfferPriceData) Unarchive() {
	opd.Archived = false
	opd.UpdatedAt = time.Now()
}
