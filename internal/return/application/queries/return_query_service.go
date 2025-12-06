package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/return/domain"
)

type ReturnQueryService struct {
	returnRepo domain.ReturnRepository
}

func NewReturnQueryService(returnRepo domain.ReturnRepository) *ReturnQueryService {
	return &ReturnQueryService{returnRepo: returnRepo}
}

func (s *ReturnQueryService) GetReturn(ctx context.Context, id int64) (*ReturnRequestDTO, error) {
	returnReq, err := s.returnRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find return: %w", err)
	}
	if returnReq == nil {
		return nil, domain.ErrReturnNotFound
	}

	return ToReturnRequestDTO(returnReq), nil
}

func (s *ReturnQueryService) GetReturnByRMA(ctx context.Context, rma string) (*ReturnRequestDTO, error) {
	returnReq, err := s.returnRepo.FindByRMA(ctx, rma)
	if err != nil {
		return nil, fmt.Errorf("failed to find return: %w", err)
	}
	if returnReq == nil {
		return nil, domain.ErrReturnNotFound
	}

	return ToReturnRequestDTO(returnReq), nil
}

func (s *ReturnQueryService) GetReturnsByCustomer(ctx context.Context, customerID string) ([]*ReturnRequestDTO, error) {
	returns, err := s.returnRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find returns: %w", err)
	}

	dtos := make([]*ReturnRequestDTO, len(returns))
	for i, r := range returns {
		dtos[i] = ToReturnRequestDTO(r)
	}

	return dtos, nil
}

func (s *ReturnQueryService) GetReturnsByStatus(ctx context.Context, status string) ([]*ReturnRequestDTO, error) {
	returns, err := s.returnRepo.FindByStatus(ctx, domain.ReturnStatus(status), 100) // Default limit of 100
	if err != nil {
		return nil, fmt.Errorf("failed to find returns: %w", err)
	}

	dtos := make([]*ReturnRequestDTO, len(returns))
	for i, r := range returns {
		dtos[i] = ToReturnRequestDTO(r)
	}

	return dtos, nil
}
