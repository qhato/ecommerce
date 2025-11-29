package domain

import "time"

// SkuProductOptionValueXref represents the cross-reference between a SKU and a ProductOptionValue
type SkuProductOptionValueXref struct {
	ID                  int64
	SKUID               int64
	ProductOptionValueID int64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NewSkuProductOptionValueXref creates a new SkuProductOptionValueXref
func NewSkuProductOptionValueXref(skuID, productOptionValueID int64) (*SkuProductOptionValueXref, error) {
	if skuID == 0 {
		return nil, NewDomainError("SKUID cannot be zero for SkuProductOptionValueXref")
	}
	if productOptionValueID == 0 {
		return nil, NewDomainError("ProductOptionValueID cannot be zero for SkuProductOptionValueXref")
	}

	now := time.Now()
	return &SkuProductOptionValueXref{
		SKUID:               skuID,
		ProductOptionValueID: productOptionValueID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}
