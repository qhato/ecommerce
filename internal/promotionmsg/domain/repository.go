package domain

import "context"

// PromotionMessageRepository defines the interface for promotion message persistence
type PromotionMessageRepository interface {
	Create(ctx context.Context, message *PromotionMessage) error
	Update(ctx context.Context, message *PromotionMessage) error
	FindByID(ctx context.Context, id int64) (*PromotionMessage, error)
	FindByType(ctx context.Context, messageType MessageType) ([]*PromotionMessage, error)
	FindByStatus(ctx context.Context, status MessageStatus) ([]*PromotionMessage, error)
	FindActive(ctx context.Context, limit int) ([]*PromotionMessage, error)
	FindByPlacement(ctx context.Context, placement string) ([]*PromotionMessage, error)
	FindByEvent(ctx context.Context, event string) ([]*PromotionMessage, error)
	Delete(ctx context.Context, id int64) error
}
