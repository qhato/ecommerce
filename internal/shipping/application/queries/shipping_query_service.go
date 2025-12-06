package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/shipping/domain"
	"github.com/shopspring/decimal"
)

type ShippingRate struct {
	MethodID      int64
	Carrier       string
	Name          string
	ServiceCode   string
	EstimatedDays int
	Cost          decimal.Decimal
	DiscountedCost *decimal.Decimal
}

type ShippingQueryService struct {
	carrierRepo domain.CarrierConfigRepository
	methodRepo  domain.ShippingMethodRepository
	bandRepo    domain.ShippingBandRepository
	ruleRepo    domain.ShippingRuleRepository
}

func NewShippingQueryService(
	carrierRepo domain.CarrierConfigRepository,
	methodRepo domain.ShippingMethodRepository,
	bandRepo domain.ShippingBandRepository,
	ruleRepo domain.ShippingRuleRepository,
) *ShippingQueryService {
	return &ShippingQueryService{
		carrierRepo: carrierRepo,
		methodRepo:  methodRepo,
		bandRepo:    bandRepo,
		ruleRepo:    ruleRepo,
	}
}

// Carrier Queries

func (s *ShippingQueryService) GetCarrierConfig(ctx context.Context, query GetCarrierConfigQuery) (*domain.CarrierConfig, error) {
	config, err := s.carrierRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get carrier config: %w", err)
	}
	return config, nil
}

func (s *ShippingQueryService) GetCarrierConfigByCarrier(ctx context.Context, query GetCarrierConfigByCarrierQuery) (*domain.CarrierConfig, error) {
	config, err := s.carrierRepo.FindByCarrier(ctx, domain.ShippingCarrier(query.Carrier))
	if err != nil {
		return nil, fmt.Errorf("failed to get carrier config: %w", err)
	}
	return config, nil
}

func (s *ShippingQueryService) GetAllCarrierConfigs(ctx context.Context, query GetAllCarrierConfigsQuery) ([]*domain.CarrierConfig, error) {
	configs, err := s.carrierRepo.FindAll(ctx, query.EnabledOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to get carrier configs: %w", err)
	}
	return configs, nil
}

// Shipping Method Queries

func (s *ShippingQueryService) GetShippingMethod(ctx context.Context, query GetShippingMethodQuery) (*domain.ShippingMethod, error) {
	method, err := s.methodRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping method: %w", err)
	}
	return method, nil
}

func (s *ShippingQueryService) GetShippingMethodsByCarrier(ctx context.Context, query GetShippingMethodsByCarrierQuery) ([]*domain.ShippingMethod, error) {
	methods, err := s.methodRepo.FindByCarrier(ctx, domain.ShippingCarrier(query.Carrier))
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping methods: %w", err)
	}
	return methods, nil
}

func (s *ShippingQueryService) GetAllEnabledShippingMethods(ctx context.Context, query GetAllEnabledShippingMethodsQuery) ([]*domain.ShippingMethod, error) {
	methods, err := s.methodRepo.FindAllEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get enabled shipping methods: %w", err)
	}
	return methods, nil
}

// Shipping Band Queries

func (s *ShippingQueryService) GetShippingBandsByMethod(ctx context.Context, query GetShippingBandsByMethodQuery) ([]*domain.ShippingBand, error) {
	bands, err := s.bandRepo.FindByMethodID(ctx, query.MethodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping bands: %w", err)
	}
	return bands, nil
}

// Shipping Rule Queries

func (s *ShippingQueryService) GetShippingRule(ctx context.Context, query GetShippingRuleQuery) (*domain.ShippingRule, error) {
	rule, err := s.ruleRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping rule: %w", err)
	}
	return rule, nil
}

func (s *ShippingQueryService) GetAllEnabledShippingRules(ctx context.Context, query GetAllEnabledShippingRulesQuery) ([]*domain.ShippingRule, error) {
	rules, err := s.ruleRepo.FindAllEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get enabled shipping rules: %w", err)
	}
	return rules, nil
}

// Shipping Calculation Queries

func (s *ShippingQueryService) CalculateShippingRates(ctx context.Context, query CalculateShippingRatesQuery) ([]*ShippingRate, error) {
	var methods []*domain.ShippingMethod
	var err error

	// If specific method requested, get just that one
	if query.MethodID != nil {
		method, err := s.methodRepo.FindByID(ctx, *query.MethodID)
		if err != nil {
			return nil, fmt.Errorf("failed to get shipping method: %w", err)
		}
		if method != nil {
			methods = []*domain.ShippingMethod{method}
		}
	} else {
		// Get all enabled methods
		methods, err = s.methodRepo.FindAllEnabled(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get enabled shipping methods: %w", err)
		}
	}

	// Get applicable rules
	rules, err := s.ruleRepo.FindAllEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping rules: %w", err)
	}

	// Calculate rates for each method
	rates := make([]*ShippingRate, 0, len(methods))
	for _, method := range methods {
		if !method.IsEnabled {
			continue
		}

		// Calculate base cost
		cost := method.CalculateRate(query.Weight, query.OrderTotal, query.Quantity)

		rate := &ShippingRate{
			MethodID:      method.ID,
			Carrier:       string(method.Carrier),
			Name:          method.Name,
			ServiceCode:   method.ServiceCode,
			EstimatedDays: method.EstimatedDays,
			Cost:          cost,
		}

		// Apply rules
		discountedCost := cost
		for _, rule := range rules {
			if rule.AppliesTo(query.OrderTotal, query.Country, query.ZipCode) {
				discount := rule.CalculateDiscount(discountedCost)
				discountedCost = discountedCost.Sub(discount)
			}
		}

		if discountedCost.LessThan(cost) {
			rate.DiscountedCost = &discountedCost
		}

		rates = append(rates, rate)
	}

	return rates, nil
}

func (s *ShippingQueryService) GetAvailableShippingMethods(ctx context.Context, query GetAvailableShippingMethodsQuery) ([]*domain.ShippingMethod, error) {
	// Get all enabled methods
	methods, err := s.methodRepo.FindAllEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get enabled shipping methods: %w", err)
	}

	// Get restriction rules
	rules, err := s.ruleRepo.FindAllEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping rules: %w", err)
	}

	// Filter methods by country/zip restrictions
	available := make([]*domain.ShippingMethod, 0)
	for _, method := range methods {
		allowed := true
		for _, rule := range rules {
			if rule.RuleType == domain.RuleTypeRestrictionCountry || rule.RuleType == domain.RuleTypeRestrictionZip {
				if !rule.AppliesTo(decimal.Zero, query.Country, query.ZipCode) {
					allowed = false
					break
				}
			}
		}
		if allowed {
			available = append(available, method)
		}
	}

	return available, nil
}
