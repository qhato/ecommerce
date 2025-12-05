package domain

import "errors"

var (
	ErrContentNotFound          = errors.New("content not found")
	ErrContentTitleRequired     = errors.New("content title is required")
	ErrContentSlugRequired      = errors.New("content slug is required")
	ErrContentSlugTaken         = errors.New("content slug is already taken")
	ErrContentAlreadyPublished  = errors.New("content is already published")
	ErrInvalidContentType       = errors.New("invalid content type")
)
