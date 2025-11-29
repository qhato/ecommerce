package domain

import (
	"time"
)

// OfferItemCriteria represents criteria for matching order items for an offer
type OfferItemCriteria struct {
	ID                 int64
	Quantity           int    // From blc_offer_item_criteria.quantity
	OrderItemMatchRule string // From blc_offer_item_criteria.order_item_match_rule (text)
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// NewOfferItemCriteria creates a new OfferItemCriteria
func NewOfferItemCriteria(quantity int, orderItemMatchRule string) (*OfferItemCriteria, error) {
	if quantity <= 0 {
		return nil, NewDomainError("Quantity must be greater than zero for OfferItemCriteria")
	}
	if orderItemMatchRule == "" {
		return nil, NewDomainError("OrderItemMatchRule cannot be empty for OfferItemCriteria")
	}

	now := time.Now()
	return &OfferItemCriteria{
		Quantity:           quantity,
		OrderItemMatchRule: orderItemMatchRule,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

// UpdateCriteria updates the quantity and match rule
func (oic *OfferItemCriteria) UpdateCriteria(quantity int, orderItemMatchRule string) {
	oic.Quantity = quantity
	oic.OrderItemMatchRule = orderItemMatchRule
	oic.UpdatedAt = time.Now()
}
