package domain

import (
	"time"
)

// Product represents a product in the catalog
type Product struct {
	ID                      int64
	Archived                bool
	CanSellWithoutOptions   bool
	CanonicalURL            string
	DisplayTemplate         string
	EnableDefaultSKU        bool
	Manufacture             string
	MetaDescription         string
	MetaTitle               string
	Model                   string
	OverrideGeneratedURL    bool
	URL                     string
	URLKey                  string
	DefaultCategoryID       *int64
	DefaultSKUID            *int64
	Attributes              []ProductAttribute
	Options                 []ProductOption
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// ProductAttribute represents a custom attribute of a product
type ProductAttribute struct {
	ID        int64
	Name      string
	Value     string
	ProductID int64
}

// NewProduct creates a new product
func NewProduct(
	manufacture, model, url, urlKey string,
	canSellWithoutOptions, enableDefaultSKU bool,
) *Product {
	now := time.Now()
	return &Product{
		Manufacture:           manufacture,
		Model:                 model,
		URL:                   url,
		URLKey:                urlKey,
		CanSellWithoutOptions: canSellWithoutOptions,
		EnableDefaultSKU:      enableDefaultSKU,
		Archived:              false,
		CreatedAt:             now,
		UpdatedAt:             now,
		Attributes:            make([]ProductAttribute, 0),
		Options:               make([]ProductOption, 0),
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

// SetDefaultCategory sets the default category
func (p *Product) SetDefaultCategory(categoryID int64) {
	p.DefaultCategoryID = &categoryID
	p.UpdatedAt = time.Now()
}

// SetDefaultSKU sets the default SKU
func (p *Product) SetDefaultSKU(skuID int64) {
	p.DefaultSKUID = &skuID
	p.UpdatedAt = time.Now()
}

// AddAttribute adds a custom attribute to the product
func (p *Product) AddAttribute(name, value string) {
	p.Attributes = append(p.Attributes, ProductAttribute{
		Name:      name,
		Value:     value,
		ProductID: p.ID,
	})
	p.UpdatedAt = time.Now()
}

// UpdateAttribute updates an existing attribute or adds it if not found
func (p *Product) UpdateAttribute(name, value string) {
	for i, attr := range p.Attributes {
		if attr.Name == name {
			p.Attributes[i].Value = value
			p.UpdatedAt = time.Now()
			return
		}
	}
	p.AddAttribute(name, value)
}

// GetAttribute retrieves an attribute value by name
func (p *Product) GetAttribute(name string) (string, bool) {
	for _, attr := range p.Attributes {
		if attr.Name == name {
			return attr.Value, true
		}
	}
	return "", false
}

// RemoveAttribute removes an attribute by name
func (p *Product) RemoveAttribute(name string) {
	for i, attr := range p.Attributes {
		if attr.Name == name {
			p.Attributes = append(p.Attributes[:i], p.Attributes[i+1:]...)
			p.UpdatedAt = time.Now()
			return
		}
	}
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
