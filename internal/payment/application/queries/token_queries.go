package queries

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/payment/domain"
)

type PaymentTokenDTO struct {
	ID          string     `json:"id"`
	CustomerID  string     `json:"customer_id"`
	TokenType   string     `json:"token_type"`
	GatewayName string     `json:"gateway_name"`
	Last4Digits *string    `json:"last_4_digits,omitempty"`
	CardBrand   *string    `json:"card_brand,omitempty"`
	ExpiryMonth *int       `json:"expiry_month,omitempty"`
	ExpiryYear  *int       `json:"expiry_year,omitempty"`
	IsDefault   bool       `json:"is_default"`
	IsActive    bool       `json:"is_active"`
	IsExpired   bool       `json:"is_expired"`
	CreatedAt   time.Time  `json:"created_at"`
}

type PaymentTokenQueryService struct {
	tokenRepo domain.PaymentTokenRepository
}

func NewPaymentTokenQueryService(tokenRepo domain.PaymentTokenRepository) *PaymentTokenQueryService {
	return &PaymentTokenQueryService{
		tokenRepo: tokenRepo,
	}
}

func (s *PaymentTokenQueryService) GetToken(ctx context.Context, id string) (*PaymentTokenDTO, error) {
	token, err := s.tokenRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find token: %w", err)
	}
	if token == nil {
		return nil, fmt.Errorf("token not found")
	}

	return toPaymentTokenDTO(token), nil
}

func (s *PaymentTokenQueryService) GetCustomerTokens(ctx context.Context, customerID string) ([]*PaymentTokenDTO, error) {
	tokens, err := s.tokenRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find customer tokens: %w", err)
	}

	dtos := make([]*PaymentTokenDTO, len(tokens))
	for i, token := range tokens {
		dtos[i] = toPaymentTokenDTO(token)
	}

	return dtos, nil
}

func (s *PaymentTokenQueryService) GetCustomerActiveTokens(ctx context.Context, customerID string) ([]*PaymentTokenDTO, error) {
	tokens, err := s.tokenRepo.FindActiveByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find customer active tokens: %w", err)
	}

	dtos := make([]*PaymentTokenDTO, len(tokens))
	for i, token := range tokens {
		dtos[i] = toPaymentTokenDTO(token)
	}

	return dtos, nil
}

func (s *PaymentTokenQueryService) GetDefaultToken(ctx context.Context, customerID string) (*PaymentTokenDTO, error) {
	token, err := s.tokenRepo.FindDefaultByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find default token: %w", err)
	}
	if token == nil {
		return nil, fmt.Errorf("no default token found")
	}

	return toPaymentTokenDTO(token), nil
}

func toPaymentTokenDTO(token *domain.PaymentToken) *PaymentTokenDTO {
	return &PaymentTokenDTO{
		ID:          token.ID,
		CustomerID:  token.CustomerID,
		TokenType:   string(token.TokenType),
		GatewayName: token.GatewayName,
		Last4Digits: token.Last4Digits,
		CardBrand:   token.CardBrand,
		ExpiryMonth: token.ExpiryMonth,
		ExpiryYear:  token.ExpiryYear,
		IsDefault:   token.IsDefault,
		IsActive:    token.IsActive,
		IsExpired:   token.IsExpired(),
		CreatedAt:   token.CreatedAt,
	}
}
