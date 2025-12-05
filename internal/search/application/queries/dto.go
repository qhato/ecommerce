package queries

import "time"

// SearchResultDTO represents search results
type SearchResultDTO struct {
	Products    []*ProductSearchDTO  `json:"products"`
	Total       int                  `json:"total"`
	Facets      map[string]*FacetDTO `json:"facets"`
	Page        int                  `json:"page"`
	PageSize    int                  `json:"page_size"`
	TotalPages  int                  `json:"total_pages"`
	Query       string               `json:"query"`
	RedirectURL string               `json:"redirect_url,omitempty"`
}

// ProductSearchDTO represents a product in search results
type ProductSearchDTO struct {
	ProductID   int64   `json:"product_id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SalePrice   *float64 `json:"sale_price,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
	IsAvailable bool    `json:"is_available"`
	Rating      float64 `json:"rating,omitempty"`
	ReviewCount int     `json:"review_count,omitempty"`
	Score       float64 `json:"score"`
}

// FacetDTO represents a facet
type FacetDTO struct {
	Name   string           `json:"name"`
	Values []*FacetValueDTO `json:"values"`
}

// FacetValueDTO represents a facet value
type FacetValueDTO struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// SynonymDTO represents a search synonym
type SynonymDTO struct {
	ID       int64    `json:"id"`
	Term     string   `json:"term"`
	Synonyms []string `json:"synonyms"`
	IsActive bool     `json:"is_active"`
}

// RedirectDTO represents a search redirect
type RedirectDTO struct {
	ID             int64      `json:"id"`
	SearchTerm     string     `json:"search_term"`
	TargetURL      string     `json:"target_url"`
	Priority       int        `json:"priority"`
	IsActive       bool       `json:"is_active"`
	ActivationDate *time.Time `json:"activation_date,omitempty"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
}

// FacetConfigDTO represents a facet configuration
type FacetConfigDTO struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Label            string `json:"label"`
	FieldName        string `json:"field_name"`
	FacetType        string `json:"facet_type"`
	IsActive         bool   `json:"is_active"`
	ShowInResults    bool   `json:"show_in_results"`
	ShowInNavigation bool   `json:"show_in_navigation"`
	Priority         int    `json:"priority"`
	MinDocCount      int    `json:"min_doc_count"`
	MaxValues        int    `json:"max_values"`
}

// IndexingJobDTO represents an indexing job
type IndexingJobDTO struct {
	ID             int64      `json:"id"`
	Type           string     `json:"type"`
	Status         string     `json:"status"`
	EntityType     string     `json:"entity_type"`
	TotalItems     int        `json:"total_items"`
	ProcessedItems int        `json:"processed_items"`
	FailedItems    int        `json:"failed_items"`
	ErrorMessage   string     `json:"error_message,omitempty"`
	Progress       float64    `json:"progress"`
	Duration       string     `json:"duration,omitempty"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
