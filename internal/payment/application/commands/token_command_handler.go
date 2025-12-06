package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/payment/domain"
)

type PaymentTokenCommandHandler struct {
	tokenRepo domain.PaymentTokenRepository
}

func NewPaymentTokenCommandHandler(tokenRepo domain.PaymentTokenRepository) *PaymentTokenCommandHandler {
	return &PaymentTokenCommandHandler{
		tokenRepo: tokenRepo,
	}
}

func (h *PaymentTokenCommandHandler) HandleCreatePaymentToken(ctx context.Context, cmd CreatePaymentTokenCommand) (*domain.PaymentToken, error) {
	tokenType := domain.TokenType(cmd.TokenType)
	token := domain.NewPaymentToken(cmd.CustomerID, cmd.Token, cmd.GatewayName, tokenType)

	token.Last4Digits = cmd.Last4Digits
	token.CardBrand = cmd.CardBrand
	token.ExpiryMonth = cmd.ExpiryMonth
	token.ExpiryYear = cmd.ExpiryYear

	if cmd.IsDefault {
		// Unset any existing default token
		existingDefault, err := h.tokenRepo.FindDefaultByCustomerID(ctx, cmd.CustomerID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing default token: %w", err)
		}
		if existingDefault != nil {
			existingDefault.IsDefault = false
			if err := h.tokenRepo.Update(ctx, existingDefault); err != nil {
				return nil, fmt.Errorf("failed to update existing default token: %w", err)
			}
		}
		token.SetAsDefault()
	}

	if err := h.tokenRepo.Create(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to create payment token: %w", err)
	}

	return token, nil
}

func (h *PaymentTokenCommandHandler) HandleSetDefaultToken(ctx context.Context, cmd SetDefaultTokenCommand) (*domain.PaymentToken, error) {
	token, err := h.tokenRepo.FindByID(ctx, cmd.TokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to find token: %w", err)
	}
	if token == nil {
		return nil, fmt.Errorf("token not found")
	}

	if token.CustomerID != cmd.CustomerID {
		return nil, fmt.Errorf("token does not belong to customer")
	}

	// Unset any existing default token
	existingDefault, err := h.tokenRepo.FindDefaultByCustomerID(ctx, cmd.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing default token: %w", err)
	}
	if existingDefault != nil && existingDefault.ID != cmd.TokenID {
		existingDefault.IsDefault = false
		if err := h.tokenRepo.Update(ctx, existingDefault); err != nil {
			return nil, fmt.Errorf("failed to update existing default token: %w", err)
		}
	}

	token.SetAsDefault()
	if err := h.tokenRepo.Update(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to update token: %w", err)
	}

	return token, nil
}

func (h *PaymentTokenCommandHandler) HandleDeactivateToken(ctx context.Context, cmd DeactivateTokenCommand) error {
	token, err := h.tokenRepo.FindByID(ctx, cmd.TokenID)
	if err != nil {
		return fmt.Errorf("failed to find token: %w", err)
	}
	if token == nil {
		return fmt.Errorf("token not found")
	}

	token.Deactivate()
	if err := h.tokenRepo.Update(ctx, token); err != nil {
		return fmt.Errorf("failed to deactivate token: %w", err)
	}

	return nil
}

func (h *PaymentTokenCommandHandler) HandleDeleteToken(ctx context.Context, cmd DeleteTokenCommand) error {
	if err := h.tokenRepo.Delete(ctx, cmd.TokenID); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}
