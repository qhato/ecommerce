package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// ProductSearchDocument represents a product in the search index
// Business Logic: Estructura optimizada para búsqueda de productos
type ProductSearchDocument struct {
	// Identificación
	ProductID int64  `json:"product_id"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`

	// Contenido para búsqueda
	Description  string   `json:"description"`
	LongDesc     string   `json:"long_description,omitempty"`
	Tags         []string `json:"tags,omitempty"`

	// Pricing
	Price        decimal.Decimal `json:"price"`
	SalePrice    *decimal.Decimal `json:"sale_price,omitempty"`
	OnSale       bool            `json:"on_sale"`

	// Categorización
	CategoryID   int64    `json:"category_id"`
	CategoryName string   `json:"category_name"`
	CategoryPath []string `json:"category_path"` // Jerarquía completa

	// Facets/Filtros
	Brand        string            `json:"brand,omitempty"`
	Color        []string          `json:"color,omitempty"`
	Size         []string          `json:"size,omitempty"`
	Attributes   map[string]string `json:"attributes,omitempty"`

	// Disponibilidad
	IsAvailable  bool   `json:"is_available"`
	StockLevel   int    `json:"stock_level"`

	// Imágenes
	ImageURL     string   `json:"image_url,omitempty"`
	ThumbnailURL string   `json:"thumbnail_url,omitempty"`

	// Metadata
	IsActive     bool      `json:"is_active"`
	IsFeatured   bool      `json:"is_featured"`
	Rating       float64   `json:"rating,omitempty"`
	ReviewCount  int       `json:"review_count,omitempty"`
	ViewCount    int       `json:"view_count,omitempty"`

	// Timestamps
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IndexedAt    time.Time `json:"indexed_at"`
}

// ToSearchDocument converts ProductSearchDocument to generic SearchDocument
func (p *ProductSearchDocument) ToSearchDocument() *SearchDocument {
	facets := make(map[string][]string)

	// Add brand facet
	if p.Brand != "" {
		facets["brand"] = []string{p.Brand}
	}

	// Add color facets
	if len(p.Color) > 0 {
		facets["color"] = p.Color
	}

	// Add size facets
	if len(p.Size) > 0 {
		facets["size"] = p.Size
	}

	// Add category facets
	facets["category"] = []string{p.CategoryName}
	if len(p.CategoryPath) > 0 {
		facets["category_path"] = p.CategoryPath
	}

	// Add availability facet
	if p.IsAvailable {
		facets["availability"] = []string{"in_stock"}
	} else {
		facets["availability"] = []string{"out_of_stock"}
	}

	// Add price range facets
	priceFloat, _ := p.Price.Float64()
	switch {
	case priceFloat < 50:
		facets["price_range"] = []string{"0-50"}
	case priceFloat < 100:
		facets["price_range"] = []string{"50-100"}
	case priceFloat < 200:
		facets["price_range"] = []string{"100-200"}
	case priceFloat < 500:
		facets["price_range"] = []string{"200-500"}
	default:
		facets["price_range"] = []string{"500+"}
	}

	fields := make(map[string]interface{})
	fields["product_id"] = p.ProductID
	fields["sku"] = p.SKU
	fields["price"] = p.Price
	fields["sale_price"] = p.SalePrice
	fields["on_sale"] = p.OnSale
	fields["category_id"] = p.CategoryID
	fields["is_available"] = p.IsAvailable
	fields["stock_level"] = p.StockLevel
	fields["image_url"] = p.ImageURL
	fields["thumbnail_url"] = p.ThumbnailURL
	fields["is_featured"] = p.IsFeatured
	fields["rating"] = p.Rating
	fields["review_count"] = p.ReviewCount
	fields["view_count"] = p.ViewCount

	return &SearchDocument{
		ID:          p.SKU,
		Type:        "product",
		Title:       p.Name,
		Description: p.Description,
		Content:     p.LongDesc,
		Fields:      fields,
		Facets:      facets,
		IndexedAt:   p.IndexedAt,
	}
}

// GetFinalPrice returns the effective price (sale price if on sale, otherwise regular price)
func (p *ProductSearchDocument) GetFinalPrice() decimal.Decimal {
	if p.OnSale && p.SalePrice != nil {
		return *p.SalePrice
	}
	return p.Price
}

// HasDiscount checks if the product has a discount
func (p *ProductSearchDocument) HasDiscount() bool {
	return p.OnSale && p.SalePrice != nil && p.SalePrice.LessThan(p.Price)
}

// GetDiscountPercentage calculates the discount percentage
func (p *ProductSearchDocument) GetDiscountPercentage() float64 {
	if !p.HasDiscount() {
		return 0
	}

	discount := p.Price.Sub(*p.SalePrice)
	percentage := discount.Div(p.Price).Mul(decimal.NewFromInt(100))
	result, _ := percentage.Float64()
	return result
}
