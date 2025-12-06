package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/media/domain"
)

type MediaDTO struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	MediaType    string                 `json:"media_type"`
	Status       string                 `json:"status"`
	MimeType     string                 `json:"mime_type"`
	FileSize     int64                  `json:"file_size"`
	FilePath     string                 `json:"file_path"`
	URL          string                 `json:"url"`
	ThumbnailURL string                 `json:"thumbnail_url,omitempty"`
	Width        *int                   `json:"width,omitempty"`
	Height       *int                   `json:"height,omitempty"`
	Duration     *int                   `json:"duration,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	UploadedBy   string                 `json:"uploaded_by"`
	EntityType   *string                `json:"entity_type,omitempty"`
	EntityID     *string                `json:"entity_id,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

func ToMediaDTO(media *domain.Media) *MediaDTO {
	return &MediaDTO{
		ID:           media.ID,
		Name:         media.Name,
		Title:        media.Title,
		Description:  media.Description,
		MediaType:    string(media.MediaType),
		Status:       string(media.Status),
		MimeType:     media.MimeType,
		FileSize:     media.FileSize,
		FilePath:     media.FilePath,
		URL:          media.URL,
		ThumbnailURL: media.ThumbnailURL,
		Width:        media.Width,
		Height:       media.Height,
		Duration:     media.Duration,
		Tags:         media.Tags,
		Metadata:     media.Metadata,
		UploadedBy:   media.UploadedBy,
		EntityType:   media.EntityType,
		EntityID:     media.EntityID,
		CreatedAt:    media.CreatedAt,
		UpdatedAt:    media.UpdatedAt,
	}
}
