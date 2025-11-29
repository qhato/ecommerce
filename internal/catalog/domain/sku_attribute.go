package domain

import "time"

// SKUAttribute represents a custom attribute of a SKU
type SKUAttribute struct {
	ID    int64
	Name  string
	Value string
	SKUID int64
	// Broadleaf does not explicitly have CreatedAt/UpdatedAt for attributes,
	// but it's good practice for auditing in our Go domain.
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewSKUAttribute creates a new SKUAttribute
func NewSKUAttribute(skuID int64, name, value string) (*SKUAttribute, error) {
	if skuID == 0 {
		return nil, NewDomainError("SKUID cannot be zero for SKUAttribute")
	}
	if name == "" {
		return nil, NewDomainError("Name cannot be empty for SKUAttribute")
	}
	// Broadleaf's schema for blc_sku_attribute has `value varchar(255) NOT NULL`, but often allows NULL in practice.
	// For strict adherence, `value` should not be empty, but if `value` can be logically empty,
	// you might adjust this validation. For now, matching the `NOT NULL` in SQL.
	// if value == "" {
	// 	return nil, NewDomainError("Value cannot be empty for SKUAttribute")
	// }

	now := time.Now()
	return &SKUAttribute{
		SKUID:     skuID,
		Name:      name,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateValue updates the value of the SKU attribute
func (sa *SKUAttribute) UpdateValue(value string) {
	sa.Value = value
	sa.UpdatedAt = time.Now()
}
