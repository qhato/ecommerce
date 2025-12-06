package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/payment/domain"
)

type GatewayConfigDTO struct {
	GatewayName string            `json:"gateway_name"`
	Enabled     bool              `json:"enabled"`
	Priority    int               `json:"priority"`
	Environment string            `json:"environment"`
	Config      map[string]string `json:"config,omitempty"`
}

type GatewayQueryService struct {
	gatewayConfigRepo domain.GatewayConfigRepository
}

func NewGatewayQueryService(gatewayConfigRepo domain.GatewayConfigRepository) *GatewayQueryService {
	return &GatewayQueryService{
		gatewayConfigRepo: gatewayConfigRepo,
	}
}

func (s *GatewayQueryService) GetGatewayConfig(ctx context.Context, gatewayName string) (*GatewayConfigDTO, error) {
	config, err := s.gatewayConfigRepo.FindByName(ctx, gatewayName)
	if err != nil {
		return nil, fmt.Errorf("failed to find gateway config: %w", err)
	}
	if config == nil {
		return nil, fmt.Errorf("gateway config not found")
	}

	return toGatewayConfigDTO(config), nil
}

func (s *GatewayQueryService) GetAllGatewayConfigs(ctx context.Context) ([]*GatewayConfigDTO, error) {
	configs, err := s.gatewayConfigRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find gateway configs: %w", err)
	}

	dtos := make([]*GatewayConfigDTO, len(configs))
	for i, config := range configs {
		dtos[i] = toGatewayConfigDTO(config)
	}

	return dtos, nil
}

func (s *GatewayQueryService) GetEnabledGatewayConfigs(ctx context.Context) ([]*GatewayConfigDTO, error) {
	configs, err := s.gatewayConfigRepo.FindAllEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find enabled gateway configs: %w", err)
	}

	dtos := make([]*GatewayConfigDTO, len(configs))
	for i, config := range configs {
		dtos[i] = toGatewayConfigDTO(config)
	}

	return dtos, nil
}

func toGatewayConfigDTO(config *domain.GatewayConfig) *GatewayConfigDTO {
	return &GatewayConfigDTO{
		GatewayName: config.GatewayName,
		Enabled:     config.Enabled,
		Priority:    config.Priority,
		Environment: config.Environment,
		Config:      config.Config,
	}
}
