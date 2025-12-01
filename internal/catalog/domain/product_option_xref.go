package domain

import "time"

// ProductOptionXref represents the cross-reference between a Product and a ProductOption
type ProductOptionXref struct {
	ID              int64
	ProductID       int64
	ProductOptionID int64
	// Broadleaf does not explicitly have CreatedAt/UpdatedAt for xrefs,
	// but it's good practice for auditing in our Go domain.
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewProductOptionXref creates a new ProductOptionXref
func NewProductOptionXref(productID, productOptionID int64) (*ProductOptionXref, error) {
	if productID == 0 {
		return nil, NewDomainError("ProductID cannot be zero for ProductOptionXref")
	}
	if productOptionID == 0 {
		return nil, NewDomainError("ProductOptionID cannot be zero for ProductOptionXref")
	}

	now := time.Now()
	return &ProductOptionXref{
		ProductID:       productID,
		ProductOptionID: productOptionID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}
