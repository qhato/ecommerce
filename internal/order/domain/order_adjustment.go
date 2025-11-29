package domain

import "time"

// OrderAdjustment represents an adjustment (e.g., discount) applied to the entire order
type OrderAdjustment struct {
	ID               int64
	OrderID          int64
	OfferID          int64   // Reference to the applied offer
	AdjustmentReason string  // From blc_order_adjustment.adjustment_reason
	AdjustmentValue  float64 // From blc_order_adjustment.adjustment_value
	IsFutureCredit   bool    // From blc_order_adjustment.is_future_credit
	CreatedAt        time.Time
}

// NewOrderAdjustment creates a new OrderAdjustment
func NewOrderAdjustment(
	orderID, offerID int64,
	adjustmentReason string,
	adjustmentValue float64,
	isFutureCredit bool,
) (*OrderAdjustment, error) {
	if orderID == 0 {
		return nil, NewDomainError("OrderID cannot be zero for OrderAdjustment")
	}
	if offerID == 0 {
		return nil, NewDomainError("OfferID cannot be zero for OrderAdjustment")
	}
	if adjustmentReason == "" {
		return nil, NewDomainError("AdjustmentReason cannot be empty for OrderAdjustment")
	}
	if adjustmentValue == 0.0 {
		return nil, NewDomainError("AdjustmentValue cannot be zero for OrderAdjustment")
	}

	now := time.Now()
	return &OrderAdjustment{
		OrderID:          orderID,
		OfferID:          offerID,
		AdjustmentReason: adjustmentReason,
		AdjustmentValue:  adjustmentValue,
		IsFutureCredit:   isFutureCredit,
		CreatedAt:        now,
	}, nil
}
