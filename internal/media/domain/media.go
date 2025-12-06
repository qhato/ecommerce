package domain

import (
	"time"

	"github.com/google/uuid"
)

// MediaType represents the type of media
type MediaType string

const (
	MediaTypeImage    MediaType = "IMAGE"
	MediaTypeVideo    MediaType = "VIDEO"
	MediaTypeDocument MediaType = "DOCUMENT"
	MediaTypeAudio    MediaType = "AUDIO"
	MediaTypeOther    MediaType = "OTHER"
)

// MediaStatus represents the status of media
type MediaStatus string

const (
	MediaStatusPending   MediaStatus = "PENDING"
	MediaStatusUploading MediaStatus = "UPLOADING"
	MediaStatusActive    MediaStatus = "ACTIVE"
	MediaStatusArchived  MediaStatus = "ARCHIVED"
	MediaStatusDeleted   MediaStatus = "DELETED"
)

// Media represents a media file (image, video, document, etc.)
type Media struct {
	ID           string
	Name         string
	Title        string
	Description  string
	MediaType    MediaType
	Status       MediaStatus
	MimeType     string
	FileSize     int64
	FilePath     string
	URL          string
	ThumbnailURL string
	Width        *int
	Height       *int
	Duration     *int // for videos/audio in seconds
	Tags         []string
	Metadata     map[string]interface{}
	UploadedBy   string
	EntityType   *string // Product, Category, etc.
	EntityID     *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewMedia creates a new media entity
func NewMedia(name, mimeType, filePath, uploadedBy string, fileSize int64) (*Media, error) {
	if name == "" {
		return nil, ErrMediaNameRequired
	}
	if mimeType == "" {
		return nil, ErrMediaMimeTypeRequired
	}
	if filePath == "" {
		return nil, ErrMediaFilePathRequired
	}

	mediaType := DetectMediaType(mimeType)

	now := time.Now()
	return &Media{
		ID:         uuid.New().String(),
		Name:       name,
		MimeType:   mimeType,
		MediaType:  mediaType,
		Status:     MediaStatusPending,
		FilePath:   filePath,
		FileSize:   fileSize,
		UploadedBy: uploadedBy,
		Tags:       make([]string, 0),
		Metadata:   make(map[string]interface{}),
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// DetectMediaType detects media type from mime type
func DetectMediaType(mimeType string) MediaType {
	switch {
	case len(mimeType) >= 6 && mimeType[:6] == "image/":
		return MediaTypeImage
	case len(mimeType) >= 6 && mimeType[:6] == "video/":
		return MediaTypeVideo
	case len(mimeType) >= 6 && mimeType[:6] == "audio/":
		return MediaTypeAudio
	case mimeType == "application/pdf" ||
		 mimeType == "application/msword" ||
		 mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return MediaTypeDocument
	default:
		return MediaTypeOther
	}
}

// UpdateInfo updates media information
func (m *Media) UpdateInfo(title, description string) {
	m.Title = title
	m.Description = description
	m.UpdatedAt = time.Now()
}

// SetURL sets the public URL for the media
func (m *Media) SetURL(url string) {
	m.URL = url
	m.UpdatedAt = time.Now()
}

// SetThumbnail sets the thumbnail URL
func (m *Media) SetThumbnail(thumbnailURL string) {
	m.ThumbnailURL = thumbnailURL
	m.UpdatedAt = time.Now()
}

// SetDimensions sets width and height for images/videos
func (m *Media) SetDimensions(width, height int) {
	m.Width = &width
	m.Height = &height
	m.UpdatedAt = time.Now()
}

// SetDuration sets duration for video/audio
func (m *Media) SetDuration(duration int) {
	m.Duration = &duration
	m.UpdatedAt = time.Now()
}

// AddTag adds a tag to the media
func (m *Media) AddTag(tag string) {
	if tag == "" {
		return
	}
	for _, t := range m.Tags {
		if t == tag {
			return
		}
	}
	m.Tags = append(m.Tags, tag)
	m.UpdatedAt = time.Now()
}

// RemoveTag removes a tag from the media
func (m *Media) RemoveTag(tag string) {
	for i, t := range m.Tags {
		if t == tag {
			m.Tags = append(m.Tags[:i], m.Tags[i+1:]...)
			m.UpdatedAt = time.Now()
			return
		}
	}
}

// AttachToEntity attaches media to an entity
func (m *Media) AttachToEntity(entityType, entityID string) {
	m.EntityType = &entityType
	m.EntityID = &entityID
	m.UpdatedAt = time.Now()
}

// Activate activates the media
func (m *Media) Activate() error {
	if m.Status == MediaStatusDeleted {
		return ErrCannotActivateDeletedMedia
	}
	m.Status = MediaStatusActive
	m.UpdatedAt = time.Now()
	return nil
}

// Archive archives the media
func (m *Media) Archive() {
	m.Status = MediaStatusArchived
	m.UpdatedAt = time.Now()
}

// Delete soft deletes the media
func (m *Media) Delete() {
	m.Status = MediaStatusDeleted
	m.UpdatedAt = time.Now()
}

// IsImage checks if media is an image
func (m *Media) IsImage() bool {
	return m.MediaType == MediaTypeImage
}

// IsVideo checks if media is a video
func (m *Media) IsVideo() bool {
	return m.MediaType == MediaTypeVideo
}

// IsDocument checks if media is a document
func (m *Media) IsDocument() bool {
	return m.MediaType == MediaTypeDocument
}
