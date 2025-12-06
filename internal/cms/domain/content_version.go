package domain

import "time"

// ContentVersion represents a version of content
type ContentVersion struct {
	ID            int64
	ContentID     int64
	VersionNumber int
	Body          string
	CreatedBy     int64
	Comment       string
	CreatedAt     time.Time
}

// NewContentVersion creates a new content version
func NewContentVersion(contentID int64, body string, createdBy int64, comment string) *ContentVersion {
	return &ContentVersion{
		ContentID:     contentID,
		VersionNumber: 1,
		Body:          body,
		CreatedBy:     createdBy,
		Comment:       comment,
		CreatedAt:     time.Now(),
	}
}
