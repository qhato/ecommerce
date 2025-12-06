package domain

import "time"

// ContentType represents the type of content
type ContentType string

const (
	ContentTypePage       ContentType = "PAGE"
	ContentTypeArticle    ContentType = "ARTICLE"
	ContentTypeBanner     ContentType = "BANNER"
	ContentTypeBlock      ContentType = "BLOCK"
	ContentTypeWidget     ContentType = "WIDGET"
)

// ContentStatus represents the publication status
type ContentStatus string

const (
	ContentStatusDraft     ContentStatus = "DRAFT"
	ContentStatusReview    ContentStatus = "REVIEW"
	ContentStatusPublished ContentStatus = "PUBLISHED"
	ContentStatusArchived  ContentStatus = "ARCHIVED"
)

// Content represents a CMS content item
type Content struct {
	ID              int64
	Title           string
	Slug            string
	Type            ContentType
	Status          ContentStatus
	Body            string
	Excerpt         string
	FeaturedImage   string
	MetaTitle       string
	MetaDescription string
	MetaKeywords    string
	Template        string
	AuthorID        int64
	PublishedAt     *time.Time
	ScheduledFor    *time.Time
	ExpiresAt       *time.Time
	Version         int
	ParentID        *int64
	SortOrder       int
	ViewCount       int
	Locale          string
	IsActive        bool
	CustomFields    map[string]interface{}
	Children        []Content
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewContent creates a new content item
func NewContent(title, slug string, contentType ContentType, body string, authorID int64, locale string) (*Content, error) {
	if title == "" {
		return nil, ErrContentTitleRequired
	}
	if slug == "" {
		return nil, ErrContentSlugRequired
	}
	if locale == "" {
		locale = "en"
	}

	now := time.Now()
	return &Content{
		Title:        title,
		Slug:         slug,
		Type:         contentType,
		Status:       ContentStatusDraft,
		Body:         body,
		AuthorID:     authorID,
		Version:      1,
		Locale:       locale,
		IsActive:     true,
		CustomFields: make(map[string]interface{}),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Publish publishes the content
func (c *Content) Publish(publishedAt time.Time) error {
	if c.Status == ContentStatusPublished {
		return ErrContentAlreadyPublished
	}
	c.Status = ContentStatusPublished
	c.PublishedAt = &publishedAt
	c.UpdatedAt = time.Now()
	return nil
}

// Unpublish unpublishes the content
func (c *Content) Unpublish() {
	c.Status = ContentStatusDraft
	c.UpdatedAt = time.Now()
}

// Archive archives the content
func (c *Content) Archive() {
	c.Status = ContentStatusArchived
	c.IsActive = false
	c.UpdatedAt = time.Now()
}

// Update updates content fields
func (c *Content) Update(title, slug, body string) error {
	if title == "" {
		return ErrContentTitleRequired
	}
	if slug == "" {
		return ErrContentSlugRequired
	}
	c.Title = title
	c.Slug = slug
	c.Body = body
	c.UpdatedAt = time.Now()
	return nil
}

// UpdateContent updates content fields
func (c *Content) UpdateContent(title, body, metaTitle, metaDescription string) {
	c.Title = title
	c.Body = body
	c.MetaTitle = metaTitle
	c.MetaDescription = metaDescription
	c.UpdatedAt = time.Now()
}
