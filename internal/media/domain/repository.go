package domain

import "context"

// MediaRepository defines the interface for media persistence
type MediaRepository interface {
	Create(ctx context.Context, media *Media) error
	Update(ctx context.Context, media *Media) error
	FindByID(ctx context.Context, id string) (*Media, error)
	FindByEntityID(ctx context.Context, entityType, entityID string) ([]*Media, error)
	FindByType(ctx context.Context, mediaType MediaType) ([]*Media, error)
	FindByStatus(ctx context.Context, status MediaStatus) ([]*Media, error)
	FindByTags(ctx context.Context, tags []string) ([]*Media, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Media, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}
