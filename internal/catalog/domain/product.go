package domain

import (
	"time"
)

// Product represents a product in the catalog
type Product struct {
	ID                          int64
	Archived                    bool
	CanSellWithoutOptions       bool   // From blc_product.can_sell_without_options
	CanonicalURL                string
	DisplayTemplate             string
	EnableDefaultSKUInInventory bool   // From blc_product.enable_default_sku_in_inventory
	Manufacture                 string
	MetaDescription             string
	MetaTitle                   string
	Model                       string
	OverrideGeneratedURL        bool
	URL                         string
	URLKey                      string
	DefaultCategoryID           *int64 // From blc_product.default_category_id
	DefaultSkuID                *int64 // From blc_product.default_sku_id
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
}

// NewProduct creates a new product
func NewProduct(
	manufacture, model, url, urlKey string,
	canSellWithoutOptions, enableDefaultSKUInInventory bool,
) *Product {
	now := time.Now()
	return &Product{
		Manufacture:                 manufacture,
		Model:                       model,
		URL:                         url,
		URLKey:                      urlKey,
		CanSellWithoutOptions:       canSellWithoutOptions,
		EnableDefaultSKUInInventory: enableDefaultSKUInInventory,
		Archived:                    false,
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}
}

// Archive marks the product as archived
func (p *Product) Archive() {
	p.Archived = true
	p.UpdatedAt = time.Now()
}

// Unarchive marks the product as active
func (p *Product) Unarchive() {
	p.Archived = false
	p.UpdatedAt = time.Now()
}

// SetDefaultSKU sets the default SKU
func (p *Product) SetDefaultSKU(skuID int64) {
	p.DefaultSkuID = &skuID
	p.UpdatedAt = time.Now()
}

// IsArchived checks if the product is archived
func (p *Product) IsArchived() bool {
	return p.Archived
}

// UpdateMetadata updates SEO metadata
func (p *Product) UpdateMetadata(title, description string) {
	p.MetaTitle = title
	p.MetaDescription = description
	p.UpdatedAt = time.Now()
}

// UpdateURLs updates URL and URL key
func (p *Product) UpdateURLs(url, urlKey string, overrideGenerated bool) {
	p.URL = url
	p.URLKey = urlKey
	p.OverrideGeneratedURL = overrideGenerated
	p.UpdatedAt = time.Now()
}
