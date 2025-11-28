package domain

// ProductOption represents a configurable option for a product (e.g., Size, Color)
type ProductOption struct {
	ID                 int64
	AttributeName      string
	DisplayOrder       int
	ErrorCode          string
	ErrorMessage       string
	Label              string
	LongDescription    string
	Name               string
	ValidationStrategy string
	ValidationType     string
	Required           bool
	OptionType         string
	UseInSKUGeneration bool
	ValidationString   string
}

// NewProductOption creates a new product option
func NewProductOption(name, label, attributeName string) *ProductOption {
	return &ProductOption{
		Name:          name,
		Label:         label,
		AttributeName: attributeName,
		Required:      false,
	}
}
