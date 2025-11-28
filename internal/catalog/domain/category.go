package domain

import (
	"time"
)

// Category represents a product category
type Category struct {
	ID                       int64
	Name                     string
	Description              string
	LongDescription          string
	ActiveStartDate          *time.Time
	ActiveEndDate            *time.Time
	Archived                 bool
	DisplayTemplate          string
	ExternalID               string
	FulfillmentType          string
	InventoryType            string
	MetaDescription          string
	MetaTitle                string
	OverrideGeneratedURL     bool
	ProductDescPattern       string
	ProductTitlePattern      string
	RootDisplayOrder         float64
	TaxCode                  string
	URL                      string
	URLKey                   string
	DefaultParentCategoryID  *int64
	ParentCategories         []Category
	Attributes               []CategoryAttribute
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

// CategoryAttribute represents a custom attribute of a category
type CategoryAttribute struct {
	ID         int64
	Name       string
	Value      string
	CategoryID int64
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
		Attributes:       make([]CategoryAttribute, 0),
		ParentCategories: make([]Category, 0),
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

// SetParentCategory sets the parent category
func (c *Category) SetParentCategory(parentID int64) {
	c.DefaultParentCategoryID = &parentID
	c.UpdatedAt = time.Now()
}

// RemoveDefaultParentCategory removes the default parent category relationship
func (c *Category) RemoveDefaultParentCategory() {
	c.DefaultParentCategoryID = nil
	c.UpdatedAt = time.Now()
}

// AddParentCategory adds a parent category to the list
func (c *Category) AddParentCategory(category Category) {
	c.ParentCategories = append(c.ParentCategories, category)
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

// AddAttribute adds a custom attribute to the category
func (c *Category) AddAttribute(name, value string) {
	c.Attributes = append(c.Attributes, CategoryAttribute{
		Name:       name,
		Value:      value,
		CategoryID: c.ID,
	})
	c.UpdatedAt = time.Now()
}

// UpdateAttribute updates an existing attribute or adds it if not found
func (c *Category) UpdateAttribute(name, value string) {
	for i, attr := range c.Attributes {
		if attr.Name == name {
			c.Attributes[i].Value = value
			c.UpdatedAt = time.Now()
			return
		}
	}
	c.AddAttribute(name, value)
}

// GetAttribute retrieves an attribute value by name
func (c *Category) GetAttribute(name string) (string, bool) {
	for _, attr := range c.Attributes {
		if attr.Name == name {
			return attr.Value, true
		}
	}
	return "", false
}

// RemoveAttribute removes an attribute by name
func (c *Category) RemoveAttribute(name string) {
	for i, attr := range c.Attributes {
		if attr.Name == name {
			c.Attributes = append(c.Attributes[:i], c.Attributes[i+1:]...)
			c.UpdatedAt = time.Now()
			return
		}
	}
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
