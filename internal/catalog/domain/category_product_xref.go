package domain

import "time"

// CategoryProductXref represents the cross-reference between a Category and a Product
type CategoryProductXref struct {
	ID               int64
	CategoryID       int64
	ProductID        int64
	DefaultReference bool    // From blc_category_product_xref.default_reference
	DisplayOrder     float64 // From blc_category_product_xref.display_order
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewCategoryProductXref creates a new CategoryProductXref
func NewCategoryProductXref(categoryID, productID int64) (*CategoryProductXref, error) {
	if categoryID == 0 {
		return nil, NewDomainError("CategoryID cannot be zero for CategoryProductXref")
	}
	if productID == 0 {
		return nil, NewDomainError("ProductID cannot be zero for CategoryProductXref")
	}

	now := time.Now()
	return &CategoryProductXref{
		CategoryID:       categoryID,
		ProductID:        productID,
		DefaultReference: false, // Default value
		DisplayOrder:     0.0,   // Default value
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// SetDefaultReference sets the default reference flag
func (cpx *CategoryProductXref) SetDefaultReference(isDefault bool) {
	cpx.DefaultReference = isDefault
	cpx.UpdatedAt = time.Now()
}

// SetDisplayOrder sets the display order
func (cpx *CategoryProductXref) SetDisplayOrder(order float64) {
	cpx.DisplayOrder = order
	cpx.UpdatedAt = time.Now()
}
