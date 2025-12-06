package domain

import (
	"fmt"
	"time"
)

// Media represents a media file
type Media struct {
	ID         int64
	FileName   string
	FilePath   string
	MimeType   string
	FileSize   int64
	Title      string
	AltText    string
	Caption    string
	UploadedBy int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewMedia creates a new media item
func NewMedia(fileName, filePath, mimeType string, fileSize int64, uploadedBy int64) (*Media, error) {
	if fileName == "" {
		return nil, ErrMediaFileNameRequired
	}
	if filePath == "" {
		return nil, ErrMediaFilePathRequired
	}

	now := time.Now()
	return &Media{
		FileName:   fileName,
		FilePath:   filePath,
		MimeType:   mimeType,
		FileSize:   fileSize,
		UploadedBy: uploadedBy,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// UpdateMetadata updates media metadata
func (m *Media) UpdateMetadata(title, altText, caption string) error {
	m.Title = title
	m.AltText = altText
	m.Caption = caption
	m.UpdatedAt = time.Now()
	return nil
}

// GetURL returns the URL for the media
func (m *Media) GetURL() string {
	return fmt.Sprintf("/media/%s", m.FilePath)
}
