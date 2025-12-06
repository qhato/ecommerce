package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/cms/domain"
)

// ContentDTO represents content for API responses
type ContentDTO struct {
	ID              int64                  `json:"id"`
	Title           string                 `json:"title"`
	Slug            string                 `json:"slug"`
	Type            string                 `json:"type"`
	Status          string                 `json:"status"`
	Body            string                 `json:"body"`
	Excerpt         string                 `json:"excerpt"`
	FeaturedImage   string                 `json:"featured_image,omitempty"`
	MetaTitle       string                 `json:"meta_title,omitempty"`
	MetaDescription string                 `json:"meta_description,omitempty"`
	MetaKeywords    string                 `json:"meta_keywords,omitempty"`
	Template        string                 `json:"template,omitempty"`
	AuthorID        int64                  `json:"author_id"`
	ParentID        *int64                 `json:"parent_id,omitempty"`
	SortOrder       int                    `json:"sort_order"`
	ViewCount       int                    `json:"view_count"`
	Locale          string                 `json:"locale"`
	PublishedAt     *time.Time             `json:"published_at,omitempty"`
	ScheduledFor    *time.Time             `json:"scheduled_for,omitempty"`
	ExpiresAt       *time.Time             `json:"expires_at,omitempty"`
	CustomFields    map[string]interface{} `json:"custom_fields,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Children        []ContentDTO           `json:"children,omitempty"`
}

// ContentVersionDTO represents a content version for API responses
type ContentVersionDTO struct {
	ID            int64     `json:"id"`
	ContentID     int64     `json:"content_id"`
	VersionNumber int       `json:"version_number"`
	Body          string    `json:"body"`
	CreatedBy     int64     `json:"created_by"`
	Comment       string    `json:"comment"`
	CreatedAt     time.Time `json:"created_at"`
}

// MediaDTO represents media for API responses
type MediaDTO struct {
	ID         int64     `json:"id"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	URL        string    `json:"url"`
	MimeType   string    `json:"mime_type"`
	FileSize   int64     `json:"file_size"`
	Title      string    `json:"title"`
	AltText    string    `json:"alt_text"`
	Caption    string    `json:"caption"`
	UploadedBy int64     `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ToContentDTO converts domain Content to ContentDTO
func ToContentDTO(content *domain.Content) *ContentDTO {
	dto := &ContentDTO{
		ID:              content.ID,
		Title:           content.Title,
		Slug:            content.Slug,
		Type:            string(content.Type),
		Status:          string(content.Status),
		Body:            content.Body,
		Excerpt:         content.Excerpt,
		FeaturedImage:   content.FeaturedImage,
		MetaTitle:       content.MetaTitle,
		MetaDescription: content.MetaDescription,
		MetaKeywords:    content.MetaKeywords,
		Template:        content.Template,
		AuthorID:        content.AuthorID,
		ParentID:        content.ParentID,
		SortOrder:       content.SortOrder,
		ViewCount:       content.ViewCount,
		Locale:          content.Locale,
		PublishedAt:     content.PublishedAt,
		ScheduledFor:    content.ScheduledFor,
		ExpiresAt:       content.ExpiresAt,
		CustomFields:    content.CustomFields,
		CreatedAt:       content.CreatedAt,
		UpdatedAt:       content.UpdatedAt,
	}

	// Convert children
	if len(content.Children) > 0 {
		dto.Children = make([]ContentDTO, len(content.Children))
		for i, child := range content.Children {
			dto.Children[i] = *ToContentDTO(&child)
		}
	}

	return dto
}

// ToContentVersionDTO converts domain ContentVersion to ContentVersionDTO
func ToContentVersionDTO(version *domain.ContentVersion) *ContentVersionDTO {
	return &ContentVersionDTO{
		ID:            version.ID,
		ContentID:     version.ContentID,
		VersionNumber: version.VersionNumber,
		Body:          version.Body,
		CreatedBy:     version.CreatedBy,
		Comment:       version.Comment,
		CreatedAt:     version.CreatedAt,
	}
}

// ToMediaDTO converts domain Media to MediaDTO
func ToMediaDTO(media *domain.Media) *MediaDTO {
	return &MediaDTO{
		ID:         media.ID,
		FileName:   media.FileName,
		FilePath:   media.FilePath,
		URL:        media.GetURL(),
		MimeType:   media.MimeType,
		FileSize:   media.FileSize,
		Title:      media.Title,
		AltText:    media.AltText,
		Caption:    media.Caption,
		UploadedBy: media.UploadedBy,
		CreatedAt:  media.CreatedAt,
		UpdatedAt:  media.UpdatedAt,
	}
}
