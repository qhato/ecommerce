package domain

import "errors"

var (
	ErrMediaNameRequired           = errors.New("media name is required")
	ErrMediaMimeTypeRequired       = errors.New("media mime type is required")
	ErrMediaFilePathRequired       = errors.New("media file path is required")
	ErrMediaNotFound               = errors.New("media not found")
	ErrCannotActivateDeletedMedia  = errors.New("cannot activate deleted media")
	ErrInvalidMediaType            = errors.New("invalid media type")
	ErrFileTooLarge                = errors.New("file size exceeds maximum allowed")
)
