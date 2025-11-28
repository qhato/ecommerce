package domain

import (
	"time"
)

// SKU represents a Stock Keeping Unit
type SKU struct {
	ID                    int64
	Name                  string
	Description           string
	LongDescription       string
	ActiveStartDate       *time.Time
	ActiveEndDate         *time.Time
	Available             bool
	Cost                  float64
	ContainerShape        string
	Depth                 float64
	DimensionUnitOfMeasure string
	Girth                 float64
	Height                float64
	ContainerSize         string
	Width                 float64
	Discountable          bool
	DisplayTemplate       string
	ExternalID            string
	FulfillmentType       string
	InventoryType         string
	IsMachineSortable     bool
	OverrideGeneratedURL  bool
	Price                 float64
	RetailPrice           float64
	SalePrice             float64
	Taxable               bool
	TaxCode               string
	UPC                   string
	URLKey                string
	Weight                float64
	WeightUnitOfMeasure   string
	CurrencyCode          string
	DefaultProductID      *int64
	AdditionalProductID   *int64
	Attributes            []SKUAttribute
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// SKUAttribute represents a custom attribute of a SKU
type SKUAttribute struct {
	ID    int64
	Name  string
	Value string
	SKUID int64
}

// NewSKU creates a new SKU
func NewSKU(name, description, upc, currencyCode string, price, retailPrice float64) *SKU {
	now := time.Now()
	return &SKU{
		Name:          name,
		Description:   description,
		UPC:           upc,
		CurrencyCode:  currencyCode,
		Price:         price,
		RetailPrice:   retailPrice,
		Available:     true,
		Discountable:  true,
		Taxable:       true,
		CreatedAt:     now,
		UpdatedAt:     now,
		Attributes:    make([]SKUAttribute, 0),
	}
}

// SetAvailability sets the availability status
func (s *SKU) SetAvailability(available bool) {
	s.Available = available
	s.UpdatedAt = time.Now()
}

// SetActiveDate sets the active date range
func (s *SKU) SetActiveDate(startDate, endDate *time.Time) {
	s.ActiveStartDate = startDate
	s.ActiveEndDate = endDate
	s.UpdatedAt = time.Now()
}

// IsActive checks if the SKU is currently active
func (s *SKU) IsActive() bool {
	if !s.Available {
		return false
	}

	now := time.Now()
	if s.ActiveStartDate != nil && now.Before(*s.ActiveStartDate) {
		return false
	}
	if s.ActiveEndDate != nil && now.After(*s.ActiveEndDate) {
		return false
	}

	return true
}

// UpdatePricing updates pricing information
func (s *SKU) UpdatePricing(price, retailPrice, salePrice float64) {
	s.Price = price
	s.RetailPrice = retailPrice
	s.SalePrice = salePrice
	s.UpdatedAt = time.Now()
}

// GetEffectivePrice returns the effective selling price (sale price if set, otherwise regular price)
func (s *SKU) GetEffectivePrice() float64 {
	if s.SalePrice > 0 {
		return s.SalePrice
	}
	return s.Price
}

// SetDimensions sets the physical dimensions
func (s *SKU) SetDimensions(height, width, depth, weight float64, dimUnit, weightUnit string) {
	s.Height = height
	s.Width = width
	s.Depth = depth
	s.Weight = weight
	s.DimensionUnitOfMeasure = dimUnit
	s.WeightUnitOfMeasure = weightUnit
	s.UpdatedAt = time.Now()
}

// AddAttribute adds a custom attribute to the SKU
func (s *SKU) AddAttribute(name, value string) {
	s.Attributes = append(s.Attributes, SKUAttribute{
		Name:  name,
		Value: value,
		SKUID: s.ID,
	})
	s.UpdatedAt = time.Now()
}

// UpdateAttribute updates an existing attribute or adds it if not found
func (s *SKU) UpdateAttribute(name, value string) {
	for i, attr := range s.Attributes {
		if attr.Name == name {
			s.Attributes[i].Value = value
			s.UpdatedAt = time.Now()
			return
		}
	}
	s.AddAttribute(name, value)
}

// GetAttribute retrieves an attribute value by name
func (s *SKU) GetAttribute(name string) (string, bool) {
	for _, attr := range s.Attributes {
		if attr.Name == name {
			return attr.Value, true
		}
	}
	return "", false
}

// RemoveAttribute removes an attribute by name
func (s *SKU) RemoveAttribute(name string) {
	for i, attr := range s.Attributes {
		if attr.Name == name {
			s.Attributes = append(s.Attributes[:i], s.Attributes[i+1:]...)
			s.UpdatedAt = time.Now()
			return
		}
	}
}

// UpdateDescription updates description and long description
func (s *SKU) UpdateDescription(description, longDescription string) {
	s.Description = description
	s.LongDescription = longDescription
	s.UpdatedAt = time.Now()
}

// SetTaxable sets whether the SKU is taxable
func (s *SKU) SetTaxable(taxable bool, taxCode string) {
	s.Taxable = taxable
	s.TaxCode = taxCode
	s.UpdatedAt = time.Now()
}

// SetDiscountable sets whether the SKU can be discounted
func (s *SKU) SetDiscountable(discountable bool) {
	s.Discountable = discountable
	s.UpdatedAt = time.Now()
}
