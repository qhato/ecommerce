package commands

import "time"

// Content Commands
type CreateContentCommand struct {
	Title           string
	Slug            string
	Type            string
	Body            string
	Excerpt         string
	FeaturedImage   string
	MetaTitle       string
	MetaDescription string
	MetaKeywords    string
	Template        string
	AuthorID        int64
	ParentID        *int64
	SortOrder       int
	Locale          string
	CustomFields    map[string]interface{}
}

type UpdateContentCommand struct {
	ID              int64
	Title           string
	Slug            string
	Body            string
	Excerpt         string
	FeaturedImage   string
	MetaTitle       string
	MetaDescription string
	MetaKeywords    string
	Template        string
	CustomFields    map[string]interface{}
}

type PublishContentCommand struct {
	ID          int64
	PublishedAt *time.Time
}

type UnpublishContentCommand struct {
	ID int64
}

type ArchiveContentCommand struct {
	ID int64
}

type DeleteContentCommand struct {
	ID int64
}

// Media Commands
type CreateMediaCommand struct {
	FileName   string
	FilePath   string
	MimeType   string
	FileSize   int64
	Title      string
	AltText    string
	Caption    string
	UploadedBy int64
}

type UpdateMediaCommand struct {
	ID      int64
	Title   string
	AltText string
	Caption string
}

type DeleteMediaCommand struct {
	ID int64
}

// Content Version Commands
type CreateContentVersionCommand struct {
	ContentID int64
	Body      string
	CreatedBy int64
	Comment   string
}
