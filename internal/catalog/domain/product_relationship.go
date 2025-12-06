package domain

import (
	"time"
)

// ProductRelationshipType represents the type of relationship between products
type ProductRelationshipType string

const (
	RelationshipTypeCrossSell ProductRelationshipType = "CROSS_SELL"
	RelationshipTypeUpSell    ProductRelationshipType = "UP_SELL"
	RelationshipTypeRelated   ProductRelationshipType = "RELATED"
	RelationshipTypeAccessory ProductRelationshipType = "ACCESSORY"
	RelationshipTypeReplacement ProductRelationshipType = "REPLACEMENT"
)

// ProductRelationship represents a relationship between two products
type ProductRelationship struct {
	ID                 int64
	ProductID          int64 // Source product
	RelatedProductID   int64 // Target product
	RelationshipType   ProductRelationshipType
	Sequence           int       // Display order
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// NewProductRelationship creates a new product relationship
func NewProductRelationship(productID, relatedProductID int64, relationshipType ProductRelationshipType) *ProductRelationship {
	now := time.Now()
	return &ProductRelationship{
		ProductID:        productID,
		RelatedProductID: relatedProductID,
		RelationshipType: relationshipType,
		Sequence:         0,
		IsActive:         true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// Activate activates the relationship
func (r *ProductRelationship) Activate() {
	r.IsActive = true
	r.UpdatedAt = time.Now()
}

// Deactivate deactivates the relationship
func (r *ProductRelationship) Deactivate() {
	r.IsActive = false
	r.UpdatedAt = time.Now()
}

// UpdateSequence updates the display sequence
func (r *ProductRelationship) UpdateSequence(sequence int) {
	r.Sequence = sequence
	r.UpdatedAt = time.Now()
}

// IsCrossSell checks if this is a cross-sell relationship
func (r *ProductRelationship) IsCrossSell() bool {
	return r.RelationshipType == RelationshipTypeCrossSell
}

// IsUpSell checks if this is an up-sell relationship
func (r *ProductRelationship) IsUpSell() bool {
	return r.RelationshipType == RelationshipTypeUpSell
}

// IsRelated checks if this is a related product relationship
func (r *ProductRelationship) IsRelated() bool {
	return r.RelationshipType == RelationshipTypeRelated
}
