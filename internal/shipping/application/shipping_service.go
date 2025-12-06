package application

import (
	"context"

	"github.com/qhato/ecommerce/internal/shipping/domain"
	"github.com/shopspring/decimal"
)

type ShippingService struct {
	methodRepo domain.ShippingMethodRepository
	bandRepo   domain.ShippingBandRepository
	ruleRepo   domain.ShippingRuleRepository
	carrierRepo domain.CarrierConfigRepository
}

func NewShippingService(
	methodRepo domain.ShippingMethodRepository,
	bandRepo domain.ShippingBandRepository,
	ruleRepo domain.ShippingRuleRepository,
	carrierRepo domain.CarrierConfigRepository,
) *ShippingService {
	return &ShippingService{
		methodRepo:  methodRepo,
		bandRepo:    bandRepo,
		ruleRepo:    ruleRepo,
		carrierRepo: carrierRepo,
	}
}

type ShippingRateRequest struct {
	Weight     decimal.Decimal
	OrderTotal decimal.Decimal
	Quantity   int
	Country    string
	ZipCode    string
}

type ShippingRateResponse struct {
	MethodID      int64           `json:"method_id"`
	MethodName    string          `json:"method_name"`
	Carrier       string          `json:"carrier"`
	BaseRate      decimal.Decimal `json:"base_rate"`
	Discount      decimal.Decimal `json:"discount"`
	FinalRate     decimal.Decimal `json:"final_rate"`
	EstimatedDays int             `json:"estimated_days"`
}

func (s *ShippingService) CalculateRates(ctx context.Context, req ShippingRateRequest) ([]*ShippingRateResponse, error) {
	methods, err := s.methodRepo.FindAllEnabled(ctx)
	if err != nil {
		return nil, err
	}

	rules, err := s.ruleRepo.FindAllEnabled(ctx)
	if err != nil {
		return nil, err
	}

	var rates []*ShippingRateResponse
	for _, method := range methods {
		bands, _ := s.bandRepo.FindByMethodID(ctx, method.ID)
		method.Bands = make([]domain.ShippingBand, len(bands))
		for i, band := range bands {
			method.Bands[i] = *band
		}

		baseRate := method.CalculateRate(req.Weight, req.OrderTotal, req.Quantity)
		discount := decimal.Zero

		for _, rule := range rules {
			if rule.AppliesTo(req.OrderTotal, req.Country, req.ZipCode) {
				ruleDiscount := rule.CalculateDiscount(baseRate)
				if ruleDiscount.GreaterThan(discount) {
					discount = ruleDiscount
				}
			}
		}

		finalRate := baseRate.Sub(discount)
		if finalRate.LessThan(decimal.Zero) {
			finalRate = decimal.Zero
		}

		rates = append(rates, &ShippingRateResponse{
			MethodID:      method.ID,
			MethodName:    method.Name,
			Carrier:       string(method.Carrier),
			BaseRate:      baseRate,
			Discount:      discount,
			FinalRate:     finalRate,
			EstimatedDays: method.EstimatedDays,
		})
	}

	return rates, nil
}
