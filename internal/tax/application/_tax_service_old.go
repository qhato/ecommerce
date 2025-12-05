package application

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/tax/domain"
)

// TaxService defines the application service for tax-related operations.
type TaxService interface {
	// CreateTaxDetail creates a new tax detail record.
	CreateTaxDetail(ctx context.Context, cmd *CreateTaxDetailCommand) (*TaxDetailDTO, error)

	// GetTaxDetailByID retrieves a tax detail record by its ID.
	GetTaxDetailByID(ctx context.Context, id int64) (*TaxDetailDTO, error)

	// UpdateTaxDetail updates an existing tax detail record.
	UpdateTaxDetail(ctx context.Context, cmd *UpdateTaxDetailCommand) (*TaxDetailDTO, error)

	// DeleteTaxDetail deletes a tax detail record by its ID.
	DeleteTaxDetail(ctx context.Context, id int64) error

	// FindApplicableTaxDetails retrieves tax details applicable to a given country, region, and type.
	FindApplicableTaxDetails(ctx context.Context, taxCountry, taxRegion, taxType string) ([]*TaxDetailDTO, error)

	// CalculateTaxForItem calculates the tax amount for a given item price, category, and order details.
	CalculateTaxForItem(ctx context.Context, orderID int64, itemTotalPrice float64, itemTaxCategory string) (float64, error)
}

// TaxDetailDTO represents a tax detail data transfer object.
type TaxDetailDTO struct {
	ID               int64
	Amount           float64
	TaxCountry       string
	JurisdictionName string
	Rate             float64
	TaxRegion        string
	TaxName          string
	Type             string
	CurrencyCode     string
	ModuleConfigID   *int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// CreateTaxDetailCommand is a command to create a new tax detail.
type CreateTaxDetailCommand struct {
	Amount           float64
	TaxCountry       string
	JurisdictionName string
	Rate             float64
	TaxRegion        string
	TaxName          string
	Type             string
	CurrencyCode     string
	ModuleConfigID   *int64
}

// UpdateTaxDetailCommand is a command to update an existing tax detail.
type UpdateTaxDetailCommand struct {
	ID               int64
	Amount           *float64
	TaxCountry       *string
	JurisdictionName *string
	Rate             *float64
	TaxRegion        *string
	TaxName          *string
	Type             *string
	CurrencyCode     *string
	ModuleConfigID   *int64
}

type taxService struct {
	taxDetailRepo domain.TaxDetailRepository
}

// NewTaxService creates a new instance of TaxService.
func NewTaxService(taxDetailRepo domain.TaxDetailRepository) TaxService {
	return &taxService{
		taxDetailRepo: taxDetailRepo,
	}
}

func (s *taxService) CreateTaxDetail(ctx context.Context, cmd *CreateTaxDetailCommand) (*TaxDetailDTO, error) {
	taxDetail, err := domain.NewTaxDetail(
		cmd.Amount,
		cmd.TaxCountry,
		cmd.JurisdictionName,
		cmd.Rate,
		cmd.TaxRegion,
		cmd.TaxName,
		cmd.Type,
		cmd.CurrencyCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tax detail domain entity: %w", err)
	}

	if cmd.ModuleConfigID != nil {
		taxDetail.SetModuleConfigID(*cmd.ModuleConfigID)
	}

	err = s.taxDetailRepo.Save(ctx, taxDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to save tax detail: %w", err)
	}

	return toTaxDetailDTO(taxDetail), nil
}

func (s *taxService) GetTaxDetailByID(ctx context.Context, id int64) (*TaxDetailDTO, error) {
	taxDetail, err := s.taxDetailRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tax detail by ID: %w", err)
	}
	if taxDetail == nil {
		return nil, fmt.Errorf("tax detail with ID %d not found", id)
	}
	return toTaxDetailDTO(taxDetail), nil
}

func (s *taxService) UpdateTaxDetail(ctx context.Context, cmd *UpdateTaxDetailCommand) (*TaxDetailDTO, error) {
	taxDetail, err := s.taxDetailRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tax detail by ID for update: %w", err)
	}
	if taxDetail == nil {
		return nil, fmt.Errorf("tax detail with ID %d not found for update", cmd.ID)
	}

	amount := taxDetail.Amount
	if cmd.Amount != nil {
		amount = *cmd.Amount
	}
	taxCountry := taxDetail.TaxCountry
	if cmd.TaxCountry != nil {
		taxCountry = *cmd.TaxCountry
	}
	jurisdictionName := taxDetail.JurisdictionName
	if cmd.JurisdictionName != nil {
		jurisdictionName = *cmd.JurisdictionName
	}
	rate := taxDetail.Rate
	if cmd.Rate != nil {
		rate = *cmd.Rate
	}
	taxRegion := taxDetail.TaxRegion
	if cmd.TaxRegion != nil {
		taxRegion = *cmd.TaxRegion
	}
	taxName := taxDetail.TaxName
	if cmd.TaxName != nil {
		taxName = *cmd.TaxName
	}
	taxType := taxDetail.Type
	if cmd.Type != nil {
		taxType = *cmd.Type
	}
	currencyCode := taxDetail.CurrencyCode
	if cmd.CurrencyCode != nil {
		currencyCode = *cmd.CurrencyCode
	}

	taxDetail.UpdateDetails(amount, taxCountry, jurisdictionName, taxRegion, taxName, taxType, currencyCode, rate)

	if cmd.ModuleConfigID != nil {
		taxDetail.SetModuleConfigID(*cmd.ModuleConfigID)
	}

	err = s.taxDetailRepo.Save(ctx, taxDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to update tax detail: %w", err)
	}

	return toTaxDetailDTO(taxDetail), nil
}

func (s *taxService) DeleteTaxDetail(ctx context.Context, id int64) error {
	err := s.taxDetailRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete tax detail: %w", err)
	}
	return nil
}

func (s *taxService) FindApplicableTaxDetails(ctx context.Context, taxCountry, taxRegion, taxType string) ([]*TaxDetailDTO, error) {
	details, err := s.taxDetailRepo.FindApplicableTaxDetails(ctx, taxCountry, taxRegion, taxType)
	if err != nil {
		return nil, fmt.Errorf("failed to find applicable tax details: %w", err)
	}

	detailsDTO := make([]*TaxDetailDTO, len(details))
	for i, detail := range details {
		detailsDTO[i] = toTaxDetailDTO(detail)
	}
	return detailsDTO, nil
}

// CalculateTaxForItem calculates the tax amount for a given item price, category, and order details.
// This is a simplified example; a real tax calculation could involve complex rules, multiple tax details,
// and external tax providers.
func (s *taxService) CalculateTaxForItem(ctx context.Context, orderID int64, itemTotalPrice float64, itemTaxCategory string) (float64, error) {
	// For now, let's assume a simplified scenario where we fetch a default tax detail
	// based on some hardcoded criteria or a simple lookup.
	// In a real system, you'd likely derive country/region from the order's shipping address
	// and use the itemTaxCategory.

	// Placeholder values for demonstration
	defaultTaxCountry := "US"
	defaultTaxRegion := "CA"
	defaultTaxType := "SALES_TAX"

	applicableDetails, err := s.FindApplicableTaxDetails(ctx, defaultTaxCountry, defaultTaxRegion, defaultTaxType)
	if err != nil {
		return 0, fmt.Errorf("failed to find applicable tax details for item calculation: %w", err)
	}

	totalTaxRate := 0.0
	for _, detail := range applicableDetails {
		totalTaxRate += detail.Rate
	}

	return itemTotalPrice * totalTaxRate, nil
}

func toTaxDetailDTO(taxDetail *domain.TaxDetail) *TaxDetailDTO {
	return &TaxDetailDTO{
		ID:               taxDetail.ID,
		Amount:           taxDetail.Amount,
		TaxCountry:       taxDetail.TaxCountry,
		JurisdictionName: taxDetail.JurisdictionName,
		Rate:             taxDetail.Rate,
		TaxRegion:        taxDetail.TaxRegion,
		TaxName:          taxDetail.TaxName,
		Type:             taxDetail.Type,
		CurrencyCode:     taxDetail.CurrencyCode,
		ModuleConfigID:   taxDetail.ModuleConfigID,
		CreatedAt:        taxDetail.CreatedAt,
		UpdatedAt:        taxDetail.UpdatedAt,
	}
}
