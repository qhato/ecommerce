package domain

import "time"

// SKU represents a Stock Keeping Unit
type SKU struct {
	ID                     int64
	Name                   string
	Description            string
	LongDescription        string
	ActiveStartDate        *time.Time
	ActiveEndDate          *time.Time
	Available              bool // From blc_sku.available_flag (bpchar(1) 'Y'/'N')
	Cost                   float64
	ContainerShape         string
	Depth                  float64
	DimensionUnitOfMeasure string
	Girth                  float64
	Height                 float64
	ContainerSize          string
	Width                  float64
	Discountable           bool // From blc_sku.discountable_flag (bpchar(1) 'Y'/'N')
	DisplayTemplate        string
	ExternalID             string
	FulfillmentType        string
	InventoryType          string
	IsMachineSortable      bool
	RetailPrice            float64
	SalePrice              float64
	Taxable                bool // From blc_sku.taxable_flag (bpchar(1) 'Y'/'N')
	TaxCode                string
	UPC                    string
	URLKey                 string
	Weight                 float64
	WeightUnitOfMeasure    string
	CurrencyCode           string
	DefaultProductID       *int64
	AdditionalProductID    *int64 // From blc_sku.addl_product_id
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// NewSKU creates a new SKU
func NewSKU(
	name, description, upc, currencyCode string,
	cost, retailPrice, salePrice float64,
) *SKU {
	now := time.Now()
	return &SKU{
		Name:         name,
		Description:  description,
		UPC:          upc,
		CurrencyCode: currencyCode,
		Cost:         cost,
		RetailPrice:  retailPrice,
		SalePrice:    salePrice,
		Available:    true,
		Discountable: true,
		Taxable:      true,
		CreatedAt:    now,
		UpdatedAt:    now,
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
func (s *SKU) UpdatePricing(retailPrice, salePrice float64) {
	s.RetailPrice = retailPrice
	s.SalePrice = salePrice
	s.UpdatedAt = time.Now()
}

// SetDimensions sets the physical dimensions
func (s *SKU) SetDimensions(height, width, depth, girth float64, containerShape, dimensionUnit, containerSize string) {
	s.Height = height
	s.Width = width
	s.Depth = depth
	s.Girth = girth
	s.ContainerShape = containerShape
	s.DimensionUnitOfMeasure = dimensionUnit
	s.ContainerSize = containerSize
	s.UpdatedAt = time.Now()
}

// SetWeight sets the weight information
func (s *SKU) SetWeight(weight float64, unit string) {
	s.Weight = weight
	s.WeightUnitOfMeasure = unit
	s.UpdatedAt = time.Now()
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
