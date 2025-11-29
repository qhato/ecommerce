package domain

import "time"

// Category represents a product category
type Category struct {
	ID                      int64
	Name                    string
	Description             string
	LongDescription         string
	ActiveStartDate         *time.Time
	ActiveEndDate           *time.Time
	Archived                bool // From blc_category.archived (bpchar(1) 'Y'/'N')
	DisplayTemplate         string
	ExternalID              string
	FulfillmentType         string
	InventoryType           string
	MetaDescription         string // From blc_category.meta_desc
	MetaTitle               string // From blc_category.meta_title
	OverrideGeneratedURL    bool
	ProductDescPattern      string  // From blc_category.product_desc_pattern_override
	ProductTitlePattern     string  // From blc_category.product_title_pattern_override
	RootDisplayOrder        float64 // From blc_category.root_display_order
	TaxCode                 string
	URL                     string
	URLKey                  string
	DefaultParentCategoryID *int64 // From blc_category.default_parent_category_id
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// NewCategory creates a new category
func NewCategory(name, description, url, urlKey string) *Category {
	now := time.Now()
	return &Category{
		Name:        name,
		Description: description,
		URL:         url,
		URLKey:      urlKey,
		Archived:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Archive marks the category as archived
func (c *Category) Archive() {
	c.Archived = true
	c.UpdatedAt = time.Now()
}

// Unarchive marks the category as active
func (c *Category) Unarchive() {
	c.Archived = false
	c.UpdatedAt = time.Now()
}

// SetParentCategory sets the default parent category
func (c *Category) SetParentCategory(parentID int64) {
	c.DefaultParentCategoryID = &parentID
	c.UpdatedAt = time.Now()
}

// RemoveDefaultParentCategory removes the default parent category relationship
func (c *Category) RemoveDefaultParentCategory() {
	c.DefaultParentCategoryID = nil
	c.UpdatedAt = time.Now()
}

// SetActiveDate sets the active date range
func (c *Category) SetActiveDate(startDate, endDate *time.Time) {
	c.ActiveStartDate = startDate
	c.ActiveEndDate = endDate
	c.UpdatedAt = time.Now()
}

// IsActive checks if the category is currently active
func (c *Category) IsActive() bool {
	if c.Archived {
		return false
	}

	now := time.Now()
	if c.ActiveStartDate != nil && now.Before(*c.ActiveStartDate) {
		return false
	}
	if c.ActiveEndDate != nil && now.After(*c.ActiveEndDate) {
		return false
	}

	return true
}

// UpdateMetadata updates SEO metadata
func (c *Category) UpdateMetadata(title, description string) {
	c.MetaTitle = title
	c.MetaDescription = description
	c.UpdatedAt = time.Now()
}

// UpdateURLs updates URL and URL key
func (c *Category) UpdateURLs(url, urlKey string, overrideGenerated bool) {
	c.URL = url
	c.URLKey = urlKey
	c.OverrideGeneratedURL = overrideGenerated
	c.UpdatedAt = time.Now()
}

// UpdateDescription updates description and long description
func (c *Category) UpdateDescription(description, longDescription string) {
	c.Description = description
	c.LongDescription = longDescription
	c.UpdatedAt = time.Now()
}

// SetDisplayOrder sets the root display order
func (c *Category) SetDisplayOrder(order float64) {
	c.RootDisplayOrder = order
	c.UpdatedAt = time.Now()
}
