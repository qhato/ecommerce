package application

import (
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// ProductDTO represents a product data transfer object
type ProductDTO struct {
	ID                    int64              `json:"id"`
	Archived              bool               `json:"archived"`
	CanSellWithoutOptions bool               `json:"can_sell_without_options"`
	CanonicalURL          string             `json:"canonical_url,omitempty"`
	DisplayTemplate       string             `json:"display_template,omitempty"`
	EnableDefaultSKU      bool               `json:"enable_default_sku"`
	Manufacture           string             `json:"manufacture"`
	MetaDescription       string             `json:"meta_description,omitempty"`
	MetaTitle             string             `json:"meta_title,omitempty"`
	Model                 string             `json:"model"`
	OverrideGeneratedURL  bool               `json:"override_generated_url"`
	URL                   string             `json:"url"`
	URLKey                string             `json:"url_key"`
	DefaultCategoryID     *int64             `json:"default_category_id,omitempty"`
	DefaultSKUID          *int64             `json:"default_sku_id,omitempty"`
	Attributes            map[string]string  `json:"attributes,omitempty"`
	CreatedAt             time.Time          `json:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at"`
}

// CategoryDTO represents a category data transfer object
type CategoryDTO struct {
	ID                      int64              `json:"id"`
	Name                    string             `json:"name"`
	Description             string             `json:"description,omitempty"`
	LongDescription         string             `json:"long_description,omitempty"`
	ActiveStartDate         *time.Time         `json:"active_start_date,omitempty"`
	ActiveEndDate           *time.Time         `json:"active_end_date,omitempty"`
	Archived                bool               `json:"archived"`
	DisplayTemplate         string             `json:"display_template,omitempty"`
	ExternalID              string             `json:"external_id,omitempty"`
	FulfillmentType         string             `json:"fulfillment_type,omitempty"`
	InventoryType           string             `json:"inventory_type,omitempty"`
	MetaDescription         string             `json:"meta_description,omitempty"`
	MetaTitle               string             `json:"meta_title,omitempty"`
	OverrideGeneratedURL    bool               `json:"override_generated_url"`
	ProductDescPattern      string             `json:"product_desc_pattern,omitempty"`
	ProductTitlePattern     string             `json:"product_title_pattern,omitempty"`
	RootDisplayOrder        float64            `json:"root_display_order"`
	TaxCode                 string             `json:"tax_code,omitempty"`
	URL                     string             `json:"url"`
	URLKey                  string             `json:"url_key"`
	DefaultParentCategoryID *int64             `json:"default_parent_category_id,omitempty"`
	Attributes              map[string]string  `json:"attributes,omitempty"`
	IsActive                bool               `json:"is_active"`
	CreatedAt               time.Time          `json:"created_at"`
	UpdatedAt               time.Time          `json:"updated_at"`
}

// SKUDTO represents a SKU data transfer object
type SKUDTO struct {
	ID                    int64              `json:"id"`
	Name                  string             `json:"name"`
	Description           string             `json:"description,omitempty"`
	LongDescription       string             `json:"long_description,omitempty"`
	ActiveStartDate       *time.Time         `json:"active_start_date,omitempty"`
	ActiveEndDate         *time.Time         `json:"active_end_date,omitempty"`
	Available             bool               `json:"available"`
	Cost                  float64            `json:"cost,omitempty"`
	ContainerShape        string             `json:"container_shape,omitempty"`
	Depth                 float64            `json:"depth,omitempty"`
	DimensionUnitOfMeasure string            `json:"dimension_unit_of_measure,omitempty"`
	Girth                 float64            `json:"girth,omitempty"`
	Height                float64            `json:"height,omitempty"`
	ContainerSize         string             `json:"container_size,omitempty"`
	Width                 float64            `json:"width,omitempty"`
	Discountable          bool               `json:"discountable"`
	DisplayTemplate       string             `json:"display_template,omitempty"`
	ExternalID            string             `json:"external_id,omitempty"`
	FulfillmentType       string             `json:"fulfillment_type,omitempty"`
	InventoryType         string             `json:"inventory_type,omitempty"`
	IsMachineSortable     bool               `json:"is_machine_sortable"`
	OverrideGeneratedURL  bool               `json:"override_generated_url"`
	Price                 float64            `json:"price"`
	RetailPrice           float64            `json:"retail_price"`
	SalePrice             float64            `json:"sale_price,omitempty"`
	EffectivePrice        float64            `json:"effective_price"`
	Taxable               bool               `json:"taxable"`
	TaxCode               string             `json:"tax_code,omitempty"`
	UPC                   string             `json:"upc,omitempty"`
	URLKey                string             `json:"url_key,omitempty"`
	Weight                float64            `json:"weight,omitempty"`
	WeightUnitOfMeasure   string             `json:"weight_unit_of_measure,omitempty"`
	CurrencyCode          string             `json:"currency_code"`
	DefaultProductID      *int64             `json:"default_product_id,omitempty"`
	AdditionalProductID   *int64             `json:"additional_product_id,omitempty"`
	Attributes            map[string]string  `json:"attributes,omitempty"`
	IsActive              bool               `json:"is_active"`
	CreatedAt             time.Time          `json:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at"`
}

// ToProductDTO converts a domain Product to ProductDTO
func ToProductDTO(product *domain.Product) *ProductDTO {
	attributes := make(map[string]string)
	for _, attr := range product.Attributes {
		attributes[attr.Name] = attr.Value
	}

	return &ProductDTO{
		ID:                    product.ID,
		Archived:              product.Archived,
		CanSellWithoutOptions: product.CanSellWithoutOptions,
		CanonicalURL:          product.CanonicalURL,
		DisplayTemplate:       product.DisplayTemplate,
		EnableDefaultSKU:      product.EnableDefaultSKU,
		Manufacture:           product.Manufacture,
		MetaDescription:       product.MetaDescription,
		MetaTitle:             product.MetaTitle,
		Model:                 product.Model,
		OverrideGeneratedURL:  product.OverrideGeneratedURL,
		URL:                   product.URL,
		URLKey:                product.URLKey,
		DefaultCategoryID:     product.DefaultCategoryID,
		DefaultSKUID:          product.DefaultSKUID,
		Attributes:            attributes,
		CreatedAt:             product.CreatedAt,
		UpdatedAt:             product.UpdatedAt,
	}
}

// ToCategoryDTO converts a domain Category to CategoryDTO
func ToCategoryDTO(category *domain.Category) *CategoryDTO {
	attributes := make(map[string]string)
	for _, attr := range category.Attributes {
		attributes[attr.Name] = attr.Value
	}

	return &CategoryDTO{
		ID:                      category.ID,
		Name:                    category.Name,
		Description:             category.Description,
		LongDescription:         category.LongDescription,
		ActiveStartDate:         category.ActiveStartDate,
		ActiveEndDate:           category.ActiveEndDate,
		Archived:                category.Archived,
		DisplayTemplate:         category.DisplayTemplate,
		ExternalID:              category.ExternalID,
		FulfillmentType:         category.FulfillmentType,
		InventoryType:           category.InventoryType,
		MetaDescription:         category.MetaDescription,
		MetaTitle:               category.MetaTitle,
		OverrideGeneratedURL:    category.OverrideGeneratedURL,
		ProductDescPattern:      category.ProductDescPattern,
		ProductTitlePattern:     category.ProductTitlePattern,
		RootDisplayOrder:        category.RootDisplayOrder,
		TaxCode:                 category.TaxCode,
		URL:                     category.URL,
		URLKey:                  category.URLKey,
		DefaultParentCategoryID: category.DefaultParentCategoryID,
		Attributes:              attributes,
		IsActive:                category.IsActive(),
		CreatedAt:               category.CreatedAt,
		UpdatedAt:               category.UpdatedAt,
	}
}

// ToSKUDTO converts a domain SKU to SKUDTO
func ToSKUDTO(sku *domain.SKU) *SKUDTO {
	attributes := make(map[string]string)
	for _, attr := range sku.Attributes {
		attributes[attr.Name] = attr.Value
	}

	return &SKUDTO{
		ID:                    sku.ID,
		Name:                  sku.Name,
		Description:           sku.Description,
		LongDescription:       sku.LongDescription,
		ActiveStartDate:       sku.ActiveStartDate,
		ActiveEndDate:         sku.ActiveEndDate,
		Available:             sku.Available,
		Cost:                  sku.Cost,
		ContainerShape:        sku.ContainerShape,
		Depth:                 sku.Depth,
		DimensionUnitOfMeasure: sku.DimensionUnitOfMeasure,
		Girth:                 sku.Girth,
		Height:                sku.Height,
		ContainerSize:         sku.ContainerSize,
		Width:                 sku.Width,
		Discountable:          sku.Discountable,
		DisplayTemplate:       sku.DisplayTemplate,
		ExternalID:            sku.ExternalID,
		FulfillmentType:       sku.FulfillmentType,
		InventoryType:         sku.InventoryType,
		IsMachineSortable:     sku.IsMachineSortable,
		OverrideGeneratedURL:  sku.OverrideGeneratedURL,
		Price:                 sku.Price,
		RetailPrice:           sku.RetailPrice,
		SalePrice:             sku.SalePrice,
		EffectivePrice:        sku.GetEffectivePrice(),
		Taxable:               sku.Taxable,
		TaxCode:               sku.TaxCode,
		UPC:                   sku.UPC,
		URLKey:                sku.URLKey,
		Weight:                sku.Weight,
		WeightUnitOfMeasure:   sku.WeightUnitOfMeasure,
		CurrencyCode:          sku.CurrencyCode,
		DefaultProductID:      sku.DefaultProductID,
		AdditionalProductID:   sku.AdditionalProductID,
		Attributes:            attributes,
		IsActive:              sku.IsActive(),
		CreatedAt:             sku.CreatedAt,
		UpdatedAt:             sku.UpdatedAt,
	}
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int64       `json:"total_pages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, page, pageSize int, totalItems int64) *PaginatedResponse {
	totalPages := totalItems / int64(pageSize)
	if totalItems%int64(pageSize) > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
