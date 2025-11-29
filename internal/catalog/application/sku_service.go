package application

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// SkuService defines the application service for SKU-related operations.
type SkuService interface {
	// CreateSku creates a new SKU.
	CreateSku(ctx context.Context, cmd *CreateSkuCommand) (*SkuDTO, error)

	// GetSkuByID retrieves a SKU by its ID.
	GetSkuByID(ctx context.Context, id int64) (*SkuDTO, error)

	// UpdateSku updates an existing SKU.
	UpdateSku(ctx context.Context, cmd *UpdateSkuCommand) (*SkuDTO, error)

	// SetSkuAvailability updates the available flag for a SKU.
	SetSkuAvailability(ctx context.Context, skuID int64, available bool) error

	// UpdateSkuPricing updates the pricing information for a SKU.
	UpdateSkuPricing(ctx context.Context, skuID int64, retailPrice, salePrice float64) error

	// AddSkuAttribute adds a custom attribute to a SKU.
	AddSkuAttribute(ctx context.Context, skuID int64, name, value string) (*SkuAttributeDTO, error)

	// UpdateSkuAttribute updates an existing SKU attribute.
	UpdateSkuAttribute(ctx context.Context, skuAttributeID int64, name, value string) (*SkuAttributeDTO, error)

	// RemoveSkuAttribute removes a SKU attribute by its ID.
	RemoveSkuAttribute(ctx context.Context, skuAttributeID int64) error

	// AddSkuProductOptionValue adds a product option value association to a SKU.
	AddSkuProductOptionValue(ctx context.Context, skuID, productOptionValueID int64) (*SkuProductOptionValueXrefDTO, error)

	// RemoveSkuProductOptionValue removes a product option value association from a SKU.
	RemoveSkuProductOptionValue(ctx context.Context, skuID, productOptionValueID int64) error
}

// SkuDTO represents a SKU data transfer object.
type SkuDTO struct {
	ID                      int64
	Name                    string
	Description             string
	LongDescription         string
	ActiveStartDate         *time.Time
	ActiveEndDate           *time.Time
	Available               bool
	Cost                    float64
	RetailPrice             float64
	SalePrice               float64
	Taxable                 bool
	TaxCode                 string
	UPC                     string
	URLKey                  string
	Weight                  float64
	WeightUnitOfMeasure     string
	CurrencyCode            string
	DefaultProductID        *int64
	AdditionalProductID     *int64
	ContainerShape          string
	Depth                   float64
	DimensionUnitOfMeasure  string
	Girth                   float64
	Height                  float64
	ContainerSize           string
	Width                   float64
	Discountable            bool
	DisplayTemplate         string
	ExternalID              string
	FulfillmentType         string
	InventoryType           string
	IsMachineSortable       bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// SkuAttributeDTO represents a SKU attribute data transfer object.
type SkuAttributeDTO struct {
	ID    int64
	Name  string
	Value string
	SKUID int64
}

// SkuProductOptionValueXrefDTO represents a SKU product option value cross-reference data transfer object.
type SkuProductOptionValueXrefDTO struct {
	ID                  int64
	SKUID               int64
	ProductOptionValueID int64
}

// CreateSkuCommand is a command to create a new SKU.
type CreateSkuCommand struct {
	Name                    string
	Description             string
	LongDescription         string
	UPC                     string
	CurrencyCode            string
	Cost                    float64
	RetailPrice             float64
	SalePrice               float64
	Taxable                 bool
	TaxCode                 string
	Available               bool
	ActiveStartDate         *time.Time
	ActiveEndDate           *time.Time
	DefaultProductID        *int64
	AdditionalProductID     *int64
	ContainerShape          string
	Depth                   float64
	DimensionUnitOfMeasure  string
	Girth                   float64
	Height                  float64
	ContainerSize           string
	Width                   float64
	Discountable            bool
	DisplayTemplate         string
	ExternalID              string
	FulfillmentType         string
	InventoryType           string
	IsMachineSortable       bool
	URLKey                  string
	Weight                  float64
	WeightUnitOfMeasure     string
}

// UpdateSkuCommand is a command to update an existing SKU.
type UpdateSkuCommand struct {
	ID                      int64
	Name                    *string
	Description             *string
	LongDescription         *string
	UPC                     *string
	CurrencyCode            *string
	Cost                    *float64
	RetailPrice             *float64
	SalePrice               *float64
	Taxable                 *bool
	TaxCode                 *string
	Available               *bool
	ActiveStartDate         *time.Time
	ActiveEndDate           *time.Time
	DefaultProductID        *int64
	AdditionalProductID     *int64
	ContainerShape          *string
	Depth                   *float64
	DimensionUnitOfMeasure  *string
	Girth                   *float64
	Height                  *float64
	ContainerSize           *string
	Width                   *float64
	Discountable            *bool
	DisplayTemplate         *string
	ExternalID              *string
	FulfillmentType         *string
	InventoryType           *string
	IsMachineSortable       *bool
	URLKey                  *string
	Weight                  *float64
	WeightUnitOfMeasure     *string
}

type skuService struct {
	skuRepo                     domain.SKURepository
	skuAttributeRepo            domain.SKUAttributeRepository
	skuProductOptionValueXrefRepo domain.SkuProductOptionValueXrefRepository
}

// NewSkuService creates a new instance of SkuService.
func NewSkuService(
	skuRepo domain.SKURepository,
	skuAttributeRepo domain.SKUAttributeRepository,
	skuProductOptionValueXrefRepo domain.SkuProductOptionValueXrefRepository,
) SkuService {
	return &skuService{
		skuRepo:                     skuRepo,
		skuAttributeRepo:            skuAttributeRepo,
		skuProductOptionValueXrefRepo: skuProductOptionValueXrefRepo,
	}
}

func (s *skuService) CreateSku(ctx context.Context, cmd *CreateSkuCommand) (*SkuDTO, error) {
	sku := domain.NewSKU(
		cmd.Name, cmd.Description, cmd.UPC, cmd.CurrencyCode,
		cmd.Cost, cmd.RetailPrice, cmd.SalePrice,
	)

	sku.LongDescription = cmd.LongDescription
	sku.ActiveStartDate = cmd.ActiveStartDate
	sku.ActiveEndDate = cmd.ActiveEndDate
	sku.Available = cmd.Available
	sku.Taxable = cmd.Taxable
	sku.TaxCode = cmd.TaxCode
	sku.DefaultProductID = cmd.DefaultProductID
	sku.AdditionalProductID = cmd.AdditionalProductID
	sku.ContainerShape = cmd.ContainerShape
	sku.Depth = cmd.Depth
	sku.DimensionUnitOfMeasure = cmd.DimensionUnitOfMeasure
	sku.Girth = cmd.Girth
	sku.Height = cmd.Height
	sku.ContainerSize = cmd.ContainerSize
	sku.Width = cmd.Width
	sku.Discountable = cmd.Discountable
	sku.DisplayTemplate = cmd.DisplayTemplate
	sku.ExternalID = cmd.ExternalID
	sku.FulfillmentType = cmd.FulfillmentType
	sku.InventoryType = cmd.InventoryType
	sku.IsMachineSortable = cmd.IsMachineSortable
	sku.URLKey = cmd.URLKey
	sku.Weight = cmd.Weight
	sku.WeightUnitOfMeasure = cmd.WeightUnitOfMeasure

	err := s.skuRepo.Save(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to save SKU: %w", err)
	}

	return toSkuDTO(sku), nil
}

func (s *skuService) GetSkuByID(ctx context.Context, id int64) (*SkuDTO, error) {
	sku, err := s.skuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find SKU by ID: %w", err)
	}
	if sku == nil {
		return nil, fmt.Errorf("SKU with ID %d not found", id)
	}
	return toSkuDTO(sku), nil
}

func (s *skuService) UpdateSku(ctx context.Context, cmd *UpdateSkuCommand) (*SkuDTO, error) {
	sku, err := s.skuRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find SKU by ID for update: %w", err)
	}
	if sku == nil {
		return nil, fmt.Errorf("SKU with ID %d not found for update", cmd.ID)
	}

	if cmd.Name != nil {
		sku.Name = *cmd.Name
	}
	if cmd.Description != nil {
		sku.Description = *cmd.Description
	}
	if cmd.LongDescription != nil {
		sku.LongDescription = *cmd.LongDescription
	}
	if cmd.UPC != nil {
		sku.UPC = *cmd.UPC
	}
	if cmd.CurrencyCode != nil {
		sku.CurrencyCode = *cmd.CurrencyCode
	}
	if cmd.Cost != nil {
		sku.Cost = *cmd.Cost
	}
	if cmd.RetailPrice != nil && cmd.SalePrice != nil {
		sku.UpdatePricing(*cmd.RetailPrice, *cmd.SalePrice)
	}
	if cmd.Taxable != nil && cmd.TaxCode != nil {
		sku.SetTaxable(*cmd.Taxable, *cmd.TaxCode)
	} else if cmd.Taxable != nil { // If only taxable flag is updated
		sku.SetTaxable(*cmd.Taxable, sku.TaxCode)
	}
	if cmd.Available != nil {
		sku.SetAvailability(*cmd.Available)
	}
	if cmd.ActiveStartDate != nil || cmd.ActiveEndDate != nil {
		sku.SetActiveDate(cmd.ActiveStartDate, cmd.ActiveEndDate)
	}
	if cmd.DefaultProductID != nil {
		sku.DefaultProductID = cmd.DefaultProductID
	}
	if cmd.AdditionalProductID != nil {
		sku.AdditionalProductID = cmd.AdditionalProductID
	}
	if cmd.ContainerShape != nil {
		sku.ContainerShape = *cmd.ContainerShape
	}
	if cmd.Depth != nil {
		sku.Depth = *cmd.Depth
	}
	if cmd.DimensionUnitOfMeasure != nil {
		sku.DimensionUnitOfMeasure = *cmd.DimensionUnitOfMeasure
	}
	if cmd.Girth != nil {
		sku.Girth = *cmd.Girth
	}
	if cmd.Height != nil {
		sku.Height = *cmd.Height
	}
	if cmd.ContainerSize != nil {
		sku.ContainerSize = *cmd.ContainerSize
	}
	if cmd.Width != nil {
		sku.Width = *cmd.Width
	}
	if cmd.Discountable != nil {
		sku.SetDiscountable(*cmd.Discountable)
	}
	if cmd.DisplayTemplate != nil {
		sku.DisplayTemplate = *cmd.DisplayTemplate
	}
	if cmd.ExternalID != nil {
		sku.ExternalID = *cmd.ExternalID
	}
	if cmd.FulfillmentType != nil {
		sku.FulfillmentType = *cmd.FulfillmentType
	}
	if cmd.InventoryType != nil {
		sku.InventoryType = *cmd.InventoryType
	}
	if cmd.IsMachineSortable != nil {
		sku.IsMachineSortable = *cmd.IsMachineSortable
	}
	if cmd.URLKey != nil {
		sku.URLKey = *cmd.URLKey
	}
	if cmd.Weight != nil {
		sku.SetWeight(*cmd.Weight, sku.WeightUnitOfMeasure) // Assuming WeightUnitOfMeasure doesn't change here
	}
	if cmd.WeightUnitOfMeasure != nil {
		sku.WeightUnitOfMeasure = *cmd.WeightUnitOfMeasure
	}

	err = s.skuRepo.Save(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to update SKU: %w", err)
	}

	return toSkuDTO(sku), nil
}

func (s *skuService) SetSkuAvailability(ctx context.Context, skuID int64, available bool) error {
	sku, err := s.skuRepo.FindByID(ctx, skuID)
	if err != nil {
		return fmt.Errorf("failed to find SKU by ID: %w", err)
	}
	if sku == nil {
		return fmt.Errorf("SKU with ID %d not found", skuID)
	}

	sku.SetAvailability(available)
	err = s.skuRepo.Save(ctx, sku)
	if err != nil {
		return fmt.Errorf("failed to set SKU availability: %w", err)
	}
	return nil
}

func (s *skuService) UpdateSkuPricing(ctx context.Context, skuID int64, retailPrice, salePrice float64) error {
	sku, err := s.skuRepo.FindByID(ctx, skuID)
	if err != nil {
		return fmt.Errorf("failed to find SKU by ID: %w", err)
	}
	if sku == nil {
		return fmt.Errorf("SKU with ID %d not found", skuID)
	}

	sku.UpdatePricing(retailPrice, salePrice)
	err = s.skuRepo.Save(ctx, sku)
	if err != nil {
		return fmt.Errorf("failed to update SKU pricing: %w", err)
	}
	return nil
}

func (s *skuService) AddSkuAttribute(ctx context.Context, skuID int64, name, value string) (*SkuAttributeDTO, error) {
	attribute, err := domain.NewSKUAttribute(skuID, name, value)
	if err != nil {
		return nil, fmt.Errorf("failed to create new SKU attribute: %w", err)
	}

	err = s.skuAttributeRepo.Save(ctx, attribute)
	if err != nil {
		return nil, fmt.Errorf("failed to save SKU attribute: %w", err)
	}

	return toSkuAttributeDTO(attribute), nil
}

func (s *skuService) UpdateSkuAttribute(ctx context.Context, skuAttributeID int64, name, value string) (*SkuAttributeDTO, error) {
	attribute, err := s.skuAttributeRepo.FindByID(ctx, skuAttributeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find SKU attribute by ID for update: %w", err)
	}
	if attribute == nil {
		return nil, fmt.Errorf("SKU attribute with ID %d not found for update", skuAttributeID)
	}

	attribute.UpdateValue(value)
	err = s.skuAttributeRepo.Save(ctx, attribute)
	if err != nil {
		return nil, fmt.Errorf("failed to update SKU attribute: %w", err)
	}

	return toSkuAttributeDTO(attribute), nil
}

func (s *skuService) RemoveSkuAttribute(ctx context.Context, skuAttributeID int64) error {
	err := s.skuAttributeRepo.Delete(ctx, skuAttributeID)
	if err != nil {
		return fmt.Errorf("failed to remove SKU attribute: %w", err)
	}
	return nil
}

func (s *skuService) AddSkuProductOptionValue(ctx context.Context, skuID, productOptionValueID int64) (*SkuProductOptionValueXrefDTO, error) {
	xref, err := domain.NewSkuProductOptionValueXref(skuID, productOptionValueID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new SKU product option value xref: %w", err)
	}

	err = s.skuProductOptionValueXrefRepo.Save(ctx, xref)
	if err != nil {
		return nil, fmt.Errorf("failed to save SKU product option value xref: %w", err)
	}

	return toSkuProductOptionValueXrefDTO(xref), nil
}

func (s *skuService) RemoveSkuProductOptionValue(ctx context.Context, skuID, productOptionValueID int64) error {
	err := s.skuProductOptionValueXrefRepo.RemoveSkuProductOptionValueXref(ctx, skuID, productOptionValueID)
	if err != nil {
		return fmt.Errorf("failed to remove SKU product option value xref: %w", err)
	}
	return nil
}

func toSkuDTO(sku *domain.SKU) *SkuDTO {
	return &SkuDTO{
		ID:                      sku.ID,
		Name:                    sku.Name,
		Description:             sku.Description,
		LongDescription:         sku.LongDescription,
		ActiveStartDate:         sku.ActiveStartDate,
		ActiveEndDate:           sku.ActiveEndDate,
		Available:               sku.Available,
		Cost:                    sku.Cost,
		RetailPrice:             sku.RetailPrice,
		SalePrice:               sku.SalePrice,
		Taxable:                 sku.Taxable,
		TaxCode:                 sku.TaxCode,
		UPC:                     sku.UPC,
		URLKey:                  sku.URLKey,
		Weight:                  sku.Weight,
		WeightUnitOfMeasure:     sku.WeightUnitOfMeasure,
		CurrencyCode:            sku.CurrencyCode,
		DefaultProductID:        sku.DefaultProductID,
		AdditionalProductID:     sku.AdditionalProductID,
		ContainerShape:          sku.ContainerShape,
		Depth:                   sku.Depth,
		DimensionUnitOfMeasure:  sku.DimensionUnitOfMeasure,
		Girth:                   sku.Girth,
		Height:                  sku.Height,
		ContainerSize:           sku.ContainerSize,
		Width:                   sku.Width,
		Discountable:            sku.Discountable,
		DisplayTemplate:         sku.DisplayTemplate,
		ExternalID:              sku.ExternalID,
		FulfillmentType:         sku.FulfillmentType,
		InventoryType:           sku.InventoryType,
		IsMachineSortable:       sku.IsMachineSortable,
		CreatedAt:               sku.CreatedAt,
		UpdatedAt:               sku.UpdatedAt,
	}
}

func toSkuAttributeDTO(attribute *domain.SKUAttribute) *SkuAttributeDTO {
	return &SkuAttributeDTO{
		ID:    attribute.ID,
		Name:  attribute.Name,
		Value: attribute.Value,
		SKUID: attribute.SKUID,
	}
}

func toSkuProductOptionValueXrefDTO(xref *domain.SkuProductOptionValueXref) *SkuProductOptionValueXrefDTO {
	return &SkuProductOptionValueXrefDTO{
		ID:                  xref.ID,
		SKUID:               xref.SKUID,
		ProductOptionValueID: xref.ProductOptionValueID,
	}
}