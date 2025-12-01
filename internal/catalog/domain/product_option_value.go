package domain

import "time"

// ProductOptionValue represents a specific value for a product option (e.g., "Small", "Red")
type ProductOptionValue struct {
	ID              int64
	AttributeValue  string
	DisplayOrder    int
	PriceAdjustment float64
	ProductOptionID int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewProductOptionValue creates a new ProductOptionValue
func NewProductOptionValue(productOptionID int64, attributeValue string, displayOrder int, priceAdjustment float64) (*ProductOptionValue, error) {
	if productOptionID == 0 {
		return nil, NewDomainError("ProductOptionID cannot be zero for ProductOptionValue")
	}
	if attributeValue == "" {
		return nil, NewDomainError("AttributeValue cannot be empty for ProductOptionValue")
	}

	now := time.Now()
	return &ProductOptionValue{
		ProductOptionID: productOptionID,
		AttributeValue:  attributeValue,
		DisplayOrder:    displayOrder,
		PriceAdjustment: priceAdjustment,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// UpdateValue updates the attribute value, display order, and price adjustment
func (pov *ProductOptionValue) UpdateValue(attributeValue string, displayOrder int, priceAdjustment float64) {
	pov.AttributeValue = attributeValue
	pov.DisplayOrder = displayOrder
	pov.PriceAdjustment = priceAdjustment
	pov.UpdatedAt = time.Now()
}
