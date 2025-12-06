package commands

// CreateMediaCommand creates a new media
type CreateMediaCommand struct {
	Name        string
	Title       string
	Description string
	MimeType    string
	FilePath    string
	FileSize    int64
	UploadedBy  string
	EntityType  *string
	EntityID    *string
	Tags        []string
}

// UpdateMediaCommand updates media information
type UpdateMediaCommand struct {
	ID          string
	Title       string
	Description string
	Tags        []string
}

// AttachMediaCommand attaches media to an entity
type AttachMediaCommand struct {
	ID         string
	EntityType string
	EntityID   string
}

// SetMediaURLCommand sets the public URL
type SetMediaURLCommand struct {
	ID           string
	URL          string
	ThumbnailURL string
}

// SetMediaDimensionsCommand sets dimensions
type SetMediaDimensionsCommand struct {
	ID     string
	Width  int
	Height int
}

// ActivateMediaCommand activates media
type ActivateMediaCommand struct {
	ID string
}

// ArchiveMediaCommand archives media
type ArchiveMediaCommand struct {
	ID string
}

// DeleteMediaCommand deletes media
type DeleteMediaCommand struct {
	ID string
}
