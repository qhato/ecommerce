package domain

import "time"

// ProductAttribute represents a custom attribute of a product
type ProductAttribute struct {
	ID        int64
	Name      string
	Value     string
	ProductID int64
	// Broadleaf does not explicitly have CreatedAt/UpdatedAt for attributes,
	// but it's good practice for auditing in our Go domain.
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewProductAttribute creates a new ProductAttribute
func NewProductAttribute(productID int64, name, value string) (*ProductAttribute, error) {
	if productID == 0 {
		return nil, NewDomainError("ProductID cannot be zero for ProductAttribute")
	}
	if name == "" {
		return nil, NewDomainError("Name cannot be empty for ProductAttribute")
	}

	now := time.Now()
	return &ProductAttribute{
		ProductID: productID,
		Name:      name,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateValue updates the value of the product attribute
func (pa *ProductAttribute) UpdateValue(value string) {
	pa.Value = value
	pa.UpdatedAt = time.Now()
}
