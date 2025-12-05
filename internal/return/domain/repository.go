package domain

import "context"

// ReturnRepository defines the interface for return persistence
type ReturnRepository interface {
	Create(ctx context.Context, returnReq *ReturnRequest) error
	Update(ctx context.Context, returnReq *ReturnRequest) error
	FindByID(ctx context.Context, id int64) (*ReturnRequest, error)
	FindByRMA(ctx context.Context, rma string) (*ReturnRequest, error)
	FindByOrderID(ctx context.Context, orderID int64) ([]*ReturnRequest, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*ReturnRequest, error)
	FindByStatus(ctx context.Context, status ReturnStatus, limit int) ([]*ReturnRequest, error)
	Delete(ctx context.Context, id int64) error
	GetReturnItems(ctx context.Context, returnID int64) ([]ReturnItem, error)
}
