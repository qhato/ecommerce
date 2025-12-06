package domain

import "context"

type ContentRepository interface {
	Create(ctx context.Context, content *Content) error
	Update(ctx context.Context, content *Content) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Content, error)
	FindBySlug(ctx context.Context, slug, locale string) (*Content, error)
	FindByType(ctx context.Context, contentType ContentType, locale string, publishedOnly bool) ([]*Content, error)
	FindAll(ctx context.Context, locale string, publishedOnly bool) ([]*Content, error)
	FindHierarchy(ctx context.Context, contentType ContentType, locale string) ([]*Content, error)
	FindChildren(ctx context.Context, parentID int64) ([]*Content, error)
	Search(ctx context.Context, query, locale string, publishedOnly bool) ([]*Content, error)
	IncrementViewCount(ctx context.Context, id int64) error
}

type MediaRepository interface {
	Create(ctx context.Context, media *Media) error
	Update(ctx context.Context, media *Media) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Media, error)
	FindAll(ctx context.Context, mimeType string, limit int) ([]*Media, error)
	FindByUploader(ctx context.Context, uploaderID int64, limit int) ([]*Media, error)
}

type ContentVersionRepository interface {
	Create(ctx context.Context, version *ContentVersion) error
	FindByID(ctx context.Context, id int64) (*ContentVersion, error)
	FindByContentID(ctx context.Context, contentID int64) ([]*ContentVersion, error)
}
