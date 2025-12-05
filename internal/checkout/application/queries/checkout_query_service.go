package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/checkout/domain"
)

// CheckoutQueryService handles checkout queries
type CheckoutQueryService struct {
	sessionRepo        domain.CheckoutSessionRepository
	shippingOptionRepo domain.ShippingOptionRepository
}

// NewCheckoutQueryService creates a new query service
func NewCheckoutQueryService(
	sessionRepo domain.CheckoutSessionRepository,
	shippingOptionRepo domain.ShippingOptionRepository,
) *CheckoutQueryService {
	return &CheckoutQueryService{
		sessionRepo:        sessionRepo,
		shippingOptionRepo: shippingOptionRepo,
	}
}

// GetCheckoutSession retrieves a checkout session by ID
func (s *CheckoutQueryService) GetCheckoutSession(ctx context.Context, sessionID string) (*CheckoutSessionDTO, error) {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find checkout session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	dto := ToCheckoutSessionDTO(session)
	return &dto, nil
}

// GetCheckoutByOrderID retrieves a checkout session by order ID
func (s *CheckoutQueryService) GetCheckoutByOrderID(ctx context.Context, orderID int64) (*CheckoutSessionDTO, error) {
	session, err := s.sessionRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find checkout session: %w", err)
	}
	if session == nil {
		return nil, domain.ErrCheckoutSessionNotFound
	}

	dto := ToCheckoutSessionDTO(session)
	return &dto, nil
}

// GetAvailableShippingOptions retrieves available shipping options
func (s *CheckoutQueryService) GetAvailableShippingOptions(ctx context.Context, country, stateProvince, postalCode string) ([]ShippingOptionDTO, error) {
	options, err := s.shippingOptionRepo.FindAvailableForAddress(ctx, country, stateProvince, postalCode)
	if err != nil {
		return nil, fmt.Errorf("failed to find shipping options: %w", err)
	}

	dtos := make([]ShippingOptionDTO, 0, len(options))
	for _, option := range options {
		dto := ToShippingOptionDTO(option, nil)
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

// GetAllShippingOptions retrieves all shipping options
func (s *CheckoutQueryService) GetAllShippingOptions(ctx context.Context, activeOnly bool) ([]ShippingOptionDTO, error) {
	options, err := s.shippingOptionRepo.FindAll(ctx, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find shipping options: %w", err)
	}

	dtos := make([]ShippingOptionDTO, 0, len(options))
	for _, option := range options {
		dto := ToShippingOptionDTO(option, nil)
		dtos = append(dtos, dto)
	}

	return dtos, nil
}
