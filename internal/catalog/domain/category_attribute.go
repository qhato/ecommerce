package domain

import "time"

// CategoryAttribute represents a custom attribute of a category
type CategoryAttribute struct {
	ID         int64
	Name       string
	Value      string
	CategoryID int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewCategoryAttribute creates a new CategoryAttribute
func NewCategoryAttribute(categoryID int64, name, value string) (*CategoryAttribute, error) {
	if categoryID == 0 {
		return nil, NewDomainError("CategoryID cannot be zero for CategoryAttribute")
	}
	if name == "" {
		return nil, NewDomainError("Name cannot be empty for CategoryAttribute")
	}

	now := time.Now()
	return &CategoryAttribute{
		CategoryID: categoryID,
		Name:       name,
		Value:      value,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// UpdateValue updates the value of the category attribute
func (ca *CategoryAttribute) UpdateValue(value string) {
	ca.Value = value
	ca.UpdatedAt = time.Now()
}
