package domain

import "time"

// OrderItemAdjustment represents an adjustment (e.g., discount) applied to an individual order item
type OrderItemAdjustment struct {
	ID                 int64
	OrderItemID        int64
	OfferID            int64   // Reference to the applied offer
	AdjustmentReason   string  // From blc_order_item_adjustment.adjustment_reason
	AdjustmentValue    float64 // From blc_order_item_adjustment.adjustment_value
	AppliedToSalePrice bool    // From blc_order_item_adjustment.applied_to_sale_price
	CreatedAt          time.Time
}

// NewOrderItemAdjustment creates a new OrderItemAdjustment
func NewOrderItemAdjustment(
	orderItemID, offerID int64,
	adjustmentReason string,
	adjustmentValue float64,
	appliedToSalePrice bool,
) (*OrderItemAdjustment, error) {
	if orderItemID == 0 {
		return nil, NewDomainError("OrderItemID cannot be zero for OrderItemAdjustment")
	}
	if offerID == 0 {
		return nil, NewDomainError("OfferID cannot be zero for OrderItemAdjustment")
	}
	if adjustmentReason == "" {
		return nil, NewDomainError("AdjustmentReason cannot be empty for OrderItemAdjustment")
	}
	if adjustmentValue == 0.0 {
		return nil, NewDomainError("AdjustmentValue cannot be zero for OrderItemAdjustment")
	}

	now := time.Now()
	return &OrderItemAdjustment{
		OrderItemID:        orderItemID,
		OfferID:            offerID,
		AdjustmentReason:   adjustmentReason,
		AdjustmentValue:    adjustmentValue,
		AppliedToSalePrice: appliedToSalePrice,
		CreatedAt:          now,
	}, nil
}
