package domain

import "errors"

var (
	// Content errors
	ErrContentNotFound         = errors.New("content not found")
	ErrContentTitleRequired    = errors.New("content title is required")
	ErrContentSlugRequired     = errors.New("content slug is required")
	ErrContentAlreadyPublished = errors.New("content is already published")

	// Media errors
	ErrMediaNotFound         = errors.New("media not found")
	ErrMediaFileNameRequired = errors.New("media file name is required")
	ErrMediaFilePathRequired = errors.New("media file path is required")

	// Version errors
	ErrVersionNotFound = errors.New("version not found")
)
