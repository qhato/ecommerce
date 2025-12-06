package domain

import "context"

// ReviewRepository defines the interface for review persistence
type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	Update(ctx context.Context, review *Review) error
	FindByID(ctx context.Context, id string) (*Review, error)
	FindByProductID(ctx context.Context, productID string, status *ReviewStatus, limit, offset int) ([]*Review, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*Review, error)
	FindByStatus(ctx context.Context, status ReviewStatus) ([]*Review, error)
	GetAverageRating(ctx context.Context, productID string) (float64, error)
	GetRatingDistribution(ctx context.Context, productID string) (map[int]int, error)
	CountByProductID(ctx context.Context, productID string) (int64, error)
	ExistsByCustomerAndProduct(ctx context.Context, customerID, productID string) (bool, error)
	Delete(ctx context.Context, id string) error
}
