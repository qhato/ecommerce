package domain

import "time"

// ProductOption represents a configurable option for a product (e.g., Size, Color)
type ProductOption struct {
	ID                     int64
	AttributeName          string
	DisplayOrder           int
	ErrorCode              string
	ErrorMessage           string
	Label                  string
	LongDescription        string
	Name                   string
	ValidationStrategyType string // Changed from ValidationStrategy to ValidationStrategyType
	ValidationType         string
	Required               bool
	OptionType             string
	UseInSKUGeneration     bool
	ValidationString       string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// NewProductOption creates a new product option
func NewProductOption(name, label, attributeName string) *ProductOption {
	now := time.Now()
	return &ProductOption{
		Name:          name,
		Label:         label,
		AttributeName: attributeName,
		Required:      false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// SetDisplayOrder sets the display order
func (po *ProductOption) SetDisplayOrder(order int) {
	po.DisplayOrder = order
	po.UpdatedAt = time.Now()
}

// UpdateValidation sets validation rules
func (po *ProductOption) UpdateValidation(strategyType, validationType, validationString, errorCode, errorMessage string, required bool) {
	po.ValidationStrategyType = strategyType
	po.ValidationType = validationType
	po.ValidationString = validationString
	po.ErrorCode = errorCode
	po.ErrorMessage = errorMessage
	po.Required = required
	po.UpdatedAt = time.Now()
}

// SetSKUGenerationFlag sets whether this option is used in SKU generation
func (po *ProductOption) SetSKUGenerationFlag(useInSKU bool) {
	po.UseInSKUGeneration = useInSKU
	po.UpdatedAt = time.Now()
}

// UpdateDescription updates description and long description
func (po *ProductOption) UpdateDescription(longDescription string) {
	po.LongDescription = longDescription
	po.UpdatedAt = time.Now()
}

// UpdateNameAndLabel updates the name and label of the product option
func (po *ProductOption) UpdateNameAndLabel(name, label string) {
	po.Name = name
	po.Label = label
	po.UpdatedAt = time.Now()
}
