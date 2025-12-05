package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/tax/domain"
)

// TaxCommandHandler handles tax-related commands
type TaxCommandHandler struct {
	jurisdictionRepo domain.TaxJurisdictionRepository
	rateRepo         domain.TaxRateRepository
	exemptionRepo    domain.TaxExemptionRepository
}

// NewTaxCommandHandler creates a new command handler
func NewTaxCommandHandler(
	jurisdictionRepo domain.TaxJurisdictionRepository,
	rateRepo domain.TaxRateRepository,
	exemptionRepo domain.TaxExemptionRepository,
) *TaxCommandHandler {
	return &TaxCommandHandler{
		jurisdictionRepo: jurisdictionRepo,
		rateRepo:         rateRepo,
		exemptionRepo:    exemptionRepo,
	}
}

// HandleCreateTaxJurisdiction handles creating a new tax jurisdiction
func (h *TaxCommandHandler) HandleCreateTaxJurisdiction(ctx context.Context, cmd CreateTaxJurisdictionCommand) (*domain.TaxJurisdiction, error) {
	// Check if code already exists
	exists, err := h.jurisdictionRepo.ExistsByCode(ctx, cmd.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check jurisdiction existence: %w", err)
	}
	if exists {
		return nil, domain.ErrJurisdictionAlreadyExists
	}

	// Parse jurisdiction type
	jurisdictionType := domain.TaxJurisdictionType(cmd.JurisdictionType)

	// Create jurisdiction
	jurisdiction, err := domain.NewTaxJurisdiction(cmd.Code, cmd.Name, jurisdictionType, cmd.Country)
	if err != nil {
		return nil, err
	}

	// Set location details
	jurisdiction.SetLocation(cmd.StateProvince, cmd.County, cmd.City, cmd.PostalCode)

	// Set parent if provided
	if cmd.ParentID != nil {
		jurisdiction.SetParent(*cmd.ParentID)
	}

	// Set priority
	jurisdiction.Priority = cmd.Priority

	// Save
	if err := h.jurisdictionRepo.Create(ctx, jurisdiction); err != nil {
		return nil, fmt.Errorf("failed to create jurisdiction: %w", err)
	}

	return jurisdiction, nil
}

// HandleUpdateTaxJurisdiction handles updating an existing tax jurisdiction
func (h *TaxCommandHandler) HandleUpdateTaxJurisdiction(ctx context.Context, cmd UpdateTaxJurisdictionCommand) (*domain.TaxJurisdiction, error) {
	// Find existing jurisdiction
	jurisdiction, err := h.jurisdictionRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find jurisdiction: %w", err)
	}
	if jurisdiction == nil {
		return nil, domain.ErrJurisdictionNotFound
	}

	// Update fields
	jurisdiction.Name = cmd.Name
	jurisdiction.SetLocation(cmd.StateProvince, cmd.County, cmd.City, cmd.PostalCode)
	jurisdiction.Priority = cmd.Priority

	// Set parent if provided
	if cmd.ParentID != nil {
		jurisdiction.SetParent(*cmd.ParentID)
	}

	// Update active status
	if cmd.IsActive {
		jurisdiction.Activate()
	} else {
		jurisdiction.Deactivate()
	}

	// Save
	if err := h.jurisdictionRepo.Update(ctx, jurisdiction); err != nil {
		return nil, fmt.Errorf("failed to update jurisdiction: %w", err)
	}

	return jurisdiction, nil
}

// HandleDeleteTaxJurisdiction handles deleting a tax jurisdiction
func (h *TaxCommandHandler) HandleDeleteTaxJurisdiction(ctx context.Context, cmd DeleteTaxJurisdictionCommand) error {
	// Check if jurisdiction exists
	jurisdiction, err := h.jurisdictionRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find jurisdiction: %w", err)
	}
	if jurisdiction == nil {
		return domain.ErrJurisdictionNotFound
	}

	// Delete
	if err := h.jurisdictionRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete jurisdiction: %w", err)
	}

	return nil
}

// HandleCreateTaxRate handles creating a new tax rate
func (h *TaxCommandHandler) HandleCreateTaxRate(ctx context.Context, cmd CreateTaxRateCommand) (*domain.TaxRate, error) {
	// Verify jurisdiction exists
	jurisdiction, err := h.jurisdictionRepo.FindByID(ctx, cmd.JurisdictionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find jurisdiction: %w", err)
	}
	if jurisdiction == nil {
		return nil, domain.ErrJurisdictionNotFound
	}

	// Parse tax type and category
	taxType := domain.TaxRateType(cmd.TaxType)
	taxCategory := domain.TaxCategory(cmd.TaxCategory)

	// Create tax rate
	rate, err := domain.NewTaxRate(cmd.JurisdictionID, cmd.Name, taxType, cmd.Rate, taxCategory)
	if err != nil {
		return nil, err
	}

	// Set additional properties
	rate.IsCompound = cmd.IsCompound
	rate.IsShippingTaxable = cmd.IsShippingTaxable
	rate.MinThreshold = cmd.MinThreshold
	rate.MaxThreshold = cmd.MaxThreshold
	rate.Priority = cmd.Priority
	rate.StartDate = cmd.StartDate
	rate.EndDate = cmd.EndDate

	// Save
	if err := h.rateRepo.Create(ctx, rate); err != nil {
		return nil, fmt.Errorf("failed to create tax rate: %w", err)
	}

	return rate, nil
}

// HandleUpdateTaxRate handles updating an existing tax rate
func (h *TaxCommandHandler) HandleUpdateTaxRate(ctx context.Context, cmd UpdateTaxRateCommand) (*domain.TaxRate, error) {
	// Find existing rate
	rate, err := h.rateRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tax rate: %w", err)
	}
	if rate == nil {
		return nil, domain.ErrTaxRateNotFound
	}

	// Update fields
	rate.Name = cmd.Name
	rate.Rate = cmd.Rate
	rate.IsCompound = cmd.IsCompound
	rate.IsShippingTaxable = cmd.IsShippingTaxable
	rate.MinThreshold = cmd.MinThreshold
	rate.MaxThreshold = cmd.MaxThreshold
	rate.Priority = cmd.Priority
	rate.StartDate = cmd.StartDate
	rate.EndDate = cmd.EndDate

	// Update active status
	if cmd.IsActive {
		rate.Activate()
	} else {
		rate.Deactivate()
	}

	// Save
	if err := h.rateRepo.Update(ctx, rate); err != nil {
		return nil, fmt.Errorf("failed to update tax rate: %w", err)
	}

	return rate, nil
}

// HandleDeleteTaxRate handles deleting a tax rate
func (h *TaxCommandHandler) HandleDeleteTaxRate(ctx context.Context, cmd DeleteTaxRateCommand) error {
	// Check if rate exists
	rate, err := h.rateRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find tax rate: %w", err)
	}
	if rate == nil {
		return domain.ErrTaxRateNotFound
	}

	// Delete
	if err := h.rateRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete tax rate: %w", err)
	}

	return nil
}

// HandleBulkCreateTaxRates handles creating multiple tax rates at once
func (h *TaxCommandHandler) HandleBulkCreateTaxRates(ctx context.Context, cmd BulkCreateTaxRatesCommand) ([]*domain.TaxRate, error) {
	rates := make([]*domain.TaxRate, 0, len(cmd.Rates))

	for _, rateCmd := range cmd.Rates {
		// Verify jurisdiction exists
		jurisdiction, err := h.jurisdictionRepo.FindByID(ctx, rateCmd.JurisdictionID)
		if err != nil {
			return nil, fmt.Errorf("failed to find jurisdiction: %w", err)
		}
		if jurisdiction == nil {
			return nil, domain.ErrJurisdictionNotFound
		}

		// Parse tax type and category
		taxType := domain.TaxRateType(rateCmd.TaxType)
		taxCategory := domain.TaxCategory(rateCmd.TaxCategory)

		// Create tax rate
		rate, err := domain.NewTaxRate(rateCmd.JurisdictionID, rateCmd.Name, taxType, rateCmd.Rate, taxCategory)
		if err != nil {
			return nil, err
		}

		// Set additional properties
		rate.IsCompound = rateCmd.IsCompound
		rate.IsShippingTaxable = rateCmd.IsShippingTaxable
		rate.MinThreshold = rateCmd.MinThreshold
		rate.MaxThreshold = rateCmd.MaxThreshold
		rate.Priority = rateCmd.Priority
		rate.StartDate = rateCmd.StartDate
		rate.EndDate = rateCmd.EndDate

		rates = append(rates, rate)
	}

	// Bulk create
	if err := h.rateRepo.BulkCreate(ctx, rates); err != nil {
		return nil, fmt.Errorf("failed to bulk create tax rates: %w", err)
	}

	return rates, nil
}

// HandleCreateTaxExemption handles creating a new tax exemption
func (h *TaxCommandHandler) HandleCreateTaxExemption(ctx context.Context, cmd CreateTaxExemptionCommand) (*domain.TaxExemption, error) {
	// Check if certificate already exists
	exists, err := h.exemptionRepo.ExistsByCertificate(ctx, cmd.ExemptionCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to check exemption certificate: %w", err)
	}
	if exists {
		return nil, domain.ErrExemptionAlreadyExists
	}

	// If jurisdiction is specified, verify it exists
	if cmd.JurisdictionID != nil {
		jurisdiction, err := h.jurisdictionRepo.FindByID(ctx, *cmd.JurisdictionID)
		if err != nil {
			return nil, fmt.Errorf("failed to find jurisdiction: %w", err)
		}
		if jurisdiction == nil {
			return nil, domain.ErrJurisdictionNotFound
		}
	}

	// Create exemption
	exemption, err := domain.NewTaxExemption(cmd.CustomerID, cmd.ExemptionCertificate, cmd.Reason)
	if err != nil {
		return nil, err
	}

	// Set additional fields
	exemption.JurisdictionID = cmd.JurisdictionID
	if cmd.TaxCategory != nil {
		category := domain.TaxCategory(*cmd.TaxCategory)
		exemption.TaxCategory = &category
	}
	exemption.StartDate = cmd.StartDate
	exemption.EndDate = cmd.EndDate

	// Validate date range
	if exemption.StartDate != nil && exemption.EndDate != nil && exemption.EndDate.Before(*exemption.StartDate) {
		return nil, domain.ErrInvalidExemptionDateRange
	}

	// Save
	if err := h.exemptionRepo.Create(ctx, exemption); err != nil {
		return nil, fmt.Errorf("failed to create tax exemption: %w", err)
	}

	return exemption, nil
}

// HandleUpdateTaxExemption handles updating an existing tax exemption
func (h *TaxCommandHandler) HandleUpdateTaxExemption(ctx context.Context, cmd UpdateTaxExemptionCommand) (*domain.TaxExemption, error) {
	// Find existing exemption
	exemption, err := h.exemptionRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tax exemption: %w", err)
	}
	if exemption == nil {
		return nil, domain.ErrExemptionNotFound
	}

	// If jurisdiction is specified, verify it exists
	if cmd.JurisdictionID != nil {
		jurisdiction, err := h.jurisdictionRepo.FindByID(ctx, *cmd.JurisdictionID)
		if err != nil {
			return nil, fmt.Errorf("failed to find jurisdiction: %w", err)
		}
		if jurisdiction == nil {
			return nil, domain.ErrJurisdictionNotFound
		}
	}

	// Update fields
	exemption.JurisdictionID = cmd.JurisdictionID
	if cmd.TaxCategory != nil {
		category := domain.TaxCategory(*cmd.TaxCategory)
		exemption.TaxCategory = &category
	}
	exemption.Reason = cmd.Reason
	exemption.IsActive = cmd.IsActive
	exemption.StartDate = cmd.StartDate
	exemption.EndDate = cmd.EndDate

	// Validate date range
	if exemption.StartDate != nil && exemption.EndDate != nil && exemption.EndDate.Before(*exemption.StartDate) {
		return nil, domain.ErrInvalidExemptionDateRange
	}

	// Save
	if err := h.exemptionRepo.Update(ctx, exemption); err != nil {
		return nil, fmt.Errorf("failed to update tax exemption: %w", err)
	}

	return exemption, nil
}

// HandleDeleteTaxExemption handles deleting a tax exemption
func (h *TaxCommandHandler) HandleDeleteTaxExemption(ctx context.Context, cmd DeleteTaxExemptionCommand) error {
	// Check if exemption exists
	exemption, err := h.exemptionRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find tax exemption: %w", err)
	}
	if exemption == nil {
		return domain.ErrExemptionNotFound
	}

	// Delete
	if err := h.exemptionRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete tax exemption: %w", err)
	}

	return nil
}
