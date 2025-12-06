package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/shipping/domain"
)

type ShippingCommandHandler struct {
	carrierRepo domain.CarrierConfigRepository
	methodRepo  domain.ShippingMethodRepository
	bandRepo    domain.ShippingBandRepository
	ruleRepo    domain.ShippingRuleRepository
}

func NewShippingCommandHandler(
	carrierRepo domain.CarrierConfigRepository,
	methodRepo domain.ShippingMethodRepository,
	bandRepo domain.ShippingBandRepository,
	ruleRepo domain.ShippingRuleRepository,
) *ShippingCommandHandler {
	return &ShippingCommandHandler{
		carrierRepo: carrierRepo,
		methodRepo:  methodRepo,
		bandRepo:    bandRepo,
		ruleRepo:    ruleRepo,
	}
}

// Carrier Commands

func (h *ShippingCommandHandler) HandleCreateCarrierConfig(ctx context.Context, cmd CreateCarrierConfigCommand) (*domain.CarrierConfig, error) {
	config := domain.NewCarrierConfig(domain.ShippingCarrier(cmd.Carrier), cmd.Name)
	config.IsEnabled = cmd.IsEnabled
	config.Priority = cmd.Priority
	config.APIKey = cmd.APIKey
	config.APISecret = cmd.APISecret
	config.AccountID = cmd.AccountID
	config.Config = cmd.Config

	if err := h.carrierRepo.Create(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to create carrier config: %w", err)
	}

	return config, nil
}

func (h *ShippingCommandHandler) HandleUpdateCarrierConfig(ctx context.Context, cmd UpdateCarrierConfigCommand) (*domain.CarrierConfig, error) {
	config, err := h.carrierRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, fmt.Errorf("carrier config not found")
	}

	config.Name = cmd.Name
	config.IsEnabled = cmd.IsEnabled
	config.Priority = cmd.Priority
	config.APIKey = cmd.APIKey
	config.APISecret = cmd.APISecret
	config.AccountID = cmd.AccountID
	config.Config = cmd.Config
	config.UpdatedAt = time.Now()

	if err := h.carrierRepo.Update(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to update carrier config: %w", err)
	}

	return config, nil
}

// Shipping Method Commands

func (h *ShippingCommandHandler) HandleCreateShippingMethod(ctx context.Context, cmd CreateShippingMethodCommand) (*domain.ShippingMethod, error) {
	method := domain.NewShippingMethod(
		domain.ShippingCarrier(cmd.Carrier),
		cmd.Name,
		cmd.ServiceCode,
		domain.PricingType(cmd.PricingType),
	)
	method.Description = cmd.Description
	method.EstimatedDays = cmd.EstimatedDays
	method.FlatRate = cmd.FlatRate
	method.IsEnabled = cmd.IsEnabled

	if err := h.methodRepo.Create(ctx, method); err != nil {
		return nil, fmt.Errorf("failed to create shipping method: %w", err)
	}

	return method, nil
}

func (h *ShippingCommandHandler) HandleUpdateShippingMethod(ctx context.Context, cmd UpdateShippingMethodCommand) (*domain.ShippingMethod, error) {
	method, err := h.methodRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if method == nil {
		return nil, fmt.Errorf("shipping method not found")
	}

	method.Name = cmd.Name
	method.Description = cmd.Description
	method.ServiceCode = cmd.ServiceCode
	method.EstimatedDays = cmd.EstimatedDays
	method.PricingType = domain.PricingType(cmd.PricingType)
	method.FlatRate = cmd.FlatRate
	method.IsEnabled = cmd.IsEnabled
	method.UpdatedAt = time.Now()

	if err := h.methodRepo.Update(ctx, method); err != nil {
		return nil, fmt.Errorf("failed to update shipping method: %w", err)
	}

	return method, nil
}

func (h *ShippingCommandHandler) HandleDeleteShippingMethod(ctx context.Context, cmd DeleteShippingMethodCommand) error {
	// Delete associated bands first
	if err := h.bandRepo.DeleteByMethodID(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete shipping bands: %w", err)
	}

	if err := h.methodRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete shipping method: %w", err)
	}

	return nil
}

// Shipping Band Commands

func (h *ShippingCommandHandler) HandleCreateShippingBand(ctx context.Context, cmd CreateShippingBandCommand) (*domain.ShippingBand, error) {
	band := &domain.ShippingBand{
		MethodID:      cmd.MethodID,
		BandType:      domain.BandType(cmd.BandType),
		MinValue:      cmd.MinValue,
		MaxValue:      cmd.MaxValue,
		Price:         cmd.Price,
		PercentCharge: cmd.PercentCharge,
		CreatedAt:     time.Now(),
	}

	if err := h.bandRepo.Create(ctx, band); err != nil {
		return nil, fmt.Errorf("failed to create shipping band: %w", err)
	}

	return band, nil
}

func (h *ShippingCommandHandler) HandleDeleteShippingBand(ctx context.Context, cmd DeleteShippingBandCommand) error {
	if err := h.bandRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete shipping band: %w", err)
	}
	return nil
}

func (h *ShippingCommandHandler) HandleDeleteBandsByMethod(ctx context.Context, cmd DeleteBandsByMethodCommand) error {
	if err := h.bandRepo.DeleteByMethodID(ctx, cmd.MethodID); err != nil {
		return fmt.Errorf("failed to delete shipping bands: %w", err)
	}
	return nil
}

// Shipping Rule Commands

func (h *ShippingCommandHandler) HandleCreateShippingRule(ctx context.Context, cmd CreateShippingRuleCommand) (*domain.ShippingRule, error) {
	rule := domain.NewShippingRule(cmd.Name, domain.RuleType(cmd.RuleType))
	rule.Description = cmd.Description
	rule.IsEnabled = cmd.IsEnabled
	rule.Priority = cmd.Priority
	rule.MinOrderValue = cmd.MinOrderValue
	rule.Countries = cmd.Countries
	rule.ExcludedZips = cmd.ExcludedZips
	rule.DiscountType = domain.DiscountType(cmd.DiscountType)
	rule.DiscountValue = cmd.DiscountValue

	if err := h.ruleRepo.Create(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to create shipping rule: %w", err)
	}

	return rule, nil
}

func (h *ShippingCommandHandler) HandleUpdateShippingRule(ctx context.Context, cmd UpdateShippingRuleCommand) (*domain.ShippingRule, error) {
	rule, err := h.ruleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, fmt.Errorf("shipping rule not found")
	}

	rule.Name = cmd.Name
	rule.Description = cmd.Description
	rule.RuleType = domain.RuleType(cmd.RuleType)
	rule.IsEnabled = cmd.IsEnabled
	rule.Priority = cmd.Priority
	rule.MinOrderValue = cmd.MinOrderValue
	rule.Countries = cmd.Countries
	rule.ExcludedZips = cmd.ExcludedZips
	rule.DiscountType = domain.DiscountType(cmd.DiscountType)
	rule.DiscountValue = cmd.DiscountValue
	rule.UpdatedAt = time.Now()

	if err := h.ruleRepo.Update(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to update shipping rule: %w", err)
	}

	return rule, nil
}

func (h *ShippingCommandHandler) HandleDeleteShippingRule(ctx context.Context, cmd DeleteShippingRuleCommand) error {
	if err := h.ruleRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete shipping rule: %w", err)
	}
	return nil
}
