package domain

import "context"

// ContentRepository defines the interface for content persistence
type ContentRepository interface {
	Create(ctx context.Context, content *Content) error
	Update(ctx context.Context, content *Content) error
	FindByID(ctx context.Context, id int64) (*Content, error)
	FindBySlug(ctx context.Context, slug string) (*Content, error)
	FindByType(ctx context.Context, contentType ContentType, status ContentStatus) ([]*Content, error)
	FindAll(ctx context.Context, status ContentStatus, limit int) ([]*Content, error)
	Delete(ctx context.Context, id int64) error
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
}
