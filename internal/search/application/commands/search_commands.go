package commands

import "time"

// IndexProductCommand represents a command to index a product
type IndexProductCommand struct {
	ProductID    int64
	SKU          string
	Name         string
	Description  string
	LongDesc     string
	Price        float64
	SalePrice    *float64
	OnSale       bool
	CategoryID   int64
	CategoryName string
	CategoryPath []string
	Brand        string
	Color        []string
	Size         []string
	Attributes   map[string]string
	IsAvailable  bool
	StockLevel   int
	ImageURL     string
	ThumbnailURL string
	IsActive     bool
	IsFeatured   bool
	Rating       float64
	ReviewCount  int
	Tags         []string
}

// BulkIndexProductsCommand represents a command to index multiple products
type BulkIndexProductsCommand struct {
	Products []IndexProductCommand
}

// DeleteProductFromIndexCommand represents a command to delete a product from index
type DeleteProductFromIndexCommand struct {
	ProductID int64
	SKU       string
}

// ReindexAllProductsCommand represents a command to reindex all products
type ReindexAllProductsCommand struct {
	CreatedBy int64
}

// ReindexCategoryCommand represents a command to reindex products in a category
type ReindexCategoryCommand struct {
	CategoryID int64
	CreatedBy  int64
}

// CreateSynonymCommand represents a command to create a synonym
type CreateSynonymCommand struct {
	Term     string
	Synonyms []string
	IsActive bool
}

// UpdateSynonymCommand represents a command to update a synonym
type UpdateSynonymCommand struct {
	ID       int64
	Term     string
	Synonyms []string
	IsActive bool
}

// DeleteSynonymCommand represents a command to delete a synonym
type DeleteSynonymCommand struct {
	ID int64
}

// CreateRedirectCommand represents a command to create a search redirect
type CreateRedirectCommand struct {
	SearchTerm     string
	TargetURL      string
	Priority       int
	IsActive       bool
	ActivationDate *time.Time
	ExpirationDate *time.Time
}

// UpdateRedirectCommand represents a command to update a search redirect
type UpdateRedirectCommand struct {
	ID             int64
	SearchTerm     string
	TargetURL      string
	Priority       int
	IsActive       bool
	ActivationDate *time.Time
	ExpirationDate *time.Time
}

// DeleteRedirectCommand represents a command to delete a search redirect
type DeleteRedirectCommand struct {
	ID int64
}

// CreateFacetConfigCommand represents a command to create a facet configuration
type CreateFacetConfigCommand struct {
	Name             string
	Label            string
	FieldName        string
	FacetType        string
	IsActive         bool
	ShowInResults    bool
	ShowInNavigation bool
	Priority         int
	MinDocCount      int
	MaxValues        int
}

// UpdateFacetConfigCommand represents a command to update a facet configuration
type UpdateFacetConfigCommand struct {
	ID               int64
	Name             string
	Label            string
	FieldName        string
	FacetType        string
	IsActive         bool
	ShowInResults    bool
	ShowInNavigation bool
	Priority         int
	MinDocCount      int
	MaxValues        int
}

// DeleteFacetConfigCommand represents a command to delete a facet configuration
type DeleteFacetConfigCommand struct {
	ID int64
}
