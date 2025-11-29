package domain

import (
	"time"
)

// ContentType represents the type of content
type ContentType string

const (
	ContentTypePage        ContentType = "PAGE"
	ContentTypeArticle     ContentType = "ARTICLE"
	ContentTypeBanner      ContentType = "BANNER"
	ContentTypePromotion   ContentType = "PROMOTION"
	ContentTypeEmailTemplate ContentType = "EMAIL_TEMPLATE"
	ContentTypeWidget      ContentType = "WIDGET"
)

// ContentStatus represents the publication status
type ContentStatus string

const (
	ContentStatusDraft     ContentStatus = "DRAFT"
	ContentStatusPublished ContentStatus = "PUBLISHED"
	ContentStatusArchived  ContentStatus = "ARCHIVED"
	ContentStatusScheduled ContentStatus = "SCHEDULED"
)

// StructuredContent represents a piece of structured content
type StructuredContent struct {
	ID              string
	Type            ContentType
	Title           string
	Slug            string
	Status          ContentStatus
	Fields          map[string]interface{} // Flexible field storage
	Template        *string
	MetaTitle       *string
	MetaDescription *string
	MetaKeywords    *string
	Author          *string
	PublishedAt     *time.Time
	ScheduledAt     *time.Time
	ExpiredAt       *time.Time
	Priority        int
	Tags            []string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ContentField represents a field definition in a content type
type ContentField struct {
	Name        string
	Type        FieldType
	Required    bool
	DefaultValue interface{}
	Validation  *FieldValidation
}

// FieldType represents the type of a field
type FieldType string

const (
	FieldTypeText      FieldType = "TEXT"
	FieldTypeRichText  FieldType = "RICH_TEXT"
	FieldTypeNumber    FieldType = "NUMBER"
	FieldTypeBoolean   FieldType = "BOOLEAN"
	FieldTypeDate      FieldType = "DATE"
	FieldTypeImage     FieldType = "IMAGE"
	FieldTypeVideo     FieldType = "VIDEO"
	FieldTypeFile      FieldType = "FILE"
	FieldTypeURL       FieldType = "URL"
	FieldTypeEmail     FieldType = "EMAIL"
	FieldTypeSelect    FieldType = "SELECT"
	FieldTypeMultiSelect FieldType = "MULTI_SELECT"
	FieldTypeReference FieldType = "REFERENCE"
	FieldTypeJSON      FieldType = "JSON"
)

// FieldValidation represents validation rules for a field
type FieldValidation struct {
	MinLength *int
	MaxLength *int
	Pattern   *string
	Min       *float64
	Max       *float64
	Options   []string
}

// ContentBlock represents a reusable content block
type ContentBlock struct {
	ID        string
	Name      string
	Type      string
	Content   map[string]interface{}
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ContentTemplate represents a content template
type ContentTemplate struct {
	ID          string
	Name        string
	Description string
	ContentType ContentType
	Fields      []*ContentField
	Layout      *string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewStructuredContent creates a new structured content
func NewStructuredContent(contentType ContentType, title, slug string) (*StructuredContent, error) {
	if title == "" || slug == "" {
		return nil, NewDomainError("Title and Slug are required")
	}

	now := time.Now()
	return &StructuredContent{
		Type:      contentType,
		Title:     title,
		Slug:      slug,
		Status:    ContentStatusDraft,
		Fields:    make(map[string]interface{}),
		Tags:      make([]string, 0),
		Priority:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Publish publishes the content
func (c *StructuredContent) Publish() error {
	if c.Status == ContentStatusPublished {
		return NewDomainError("Content is already published")
	}

	now := time.Now()
	c.Status = ContentStatusPublished
	c.PublishedAt = &now
	c.UpdatedAt = now
	return nil
}

// Unpublish unpublishes the content
func (c *StructuredContent) Unpublish() {
	c.Status = ContentStatusDraft
	c.UpdatedAt = time.Now()
}

// Archive archives the content
func (c *StructuredContent) Archive() {
	c.Status = ContentStatusArchived
	c.UpdatedAt = time.Now()
}

// Schedule schedules the content for publication
func (c *StructuredContent) Schedule(publishAt time.Time) error {
	if publishAt.Before(time.Now()) {
		return NewDomainError("Cannot schedule in the past")
	}

	c.Status = ContentStatusScheduled
	c.ScheduledAt = &publishAt
	c.UpdatedAt = time.Now()
	return nil
}

// SetField sets a field value
func (c *StructuredContent) SetField(name string, value interface{}) {
	c.Fields[name] = value
	c.UpdatedAt = time.Now()
}

// GetField gets a field value
func (c *StructuredContent) GetField(name string) (interface{}, bool) {
	value, exists := c.Fields[name]
	return value, exists
}

// AddTag adds a tag
func (c *StructuredContent) AddTag(tag string) {
	c.Tags = append(c.Tags, tag)
	c.UpdatedAt = time.Now()
}

// RemoveTag removes a tag
func (c *StructuredContent) RemoveTag(tag string) {
	for i, t := range c.Tags {
		if t == tag {
			c.Tags = append(c.Tags[:i], c.Tags[i+1:]...)
			break
		}
	}
	c.UpdatedAt = time.Now()
}

// SetMetadata sets SEO metadata
func (c *StructuredContent) SetMetadata(title, description, keywords string) {
	c.MetaTitle = &title
	c.MetaDescription = &description
	c.MetaKeywords = &keywords
	c.UpdatedAt = time.Now()
}

// IsPublished checks if content is published
func (c *StructuredContent) IsPublished() bool {
	if c.Status != ContentStatusPublished {
		return false
	}

	now := time.Now()

	if c.PublishedAt != nil && c.PublishedAt.After(now) {
		return false
	}

	if c.ExpiredAt != nil && c.ExpiredAt.Before(now) {
		return false
	}

	return true
}

// IsScheduled checks if content is scheduled
func (c *StructuredContent) IsScheduled() bool {
	return c.Status == ContentStatusScheduled && c.ScheduledAt != nil
}

// ShouldPublish checks if scheduled content should be published
func (c *StructuredContent) ShouldPublish() bool {
	if !c.IsScheduled() {
		return false
	}

	return time.Now().After(*c.ScheduledAt)
}

// ContentRepository defines repository for structured content
type ContentRepository interface {
	Save(content *StructuredContent) error
	FindByID(id string) (*StructuredContent, error)
	FindBySlug(slug string) (*StructuredContent, error)
	FindByType(contentType ContentType) ([]*StructuredContent, error)
	FindPublished() ([]*StructuredContent, error)
	FindScheduled() ([]*StructuredContent, error)
	Delete(id string) error
}
