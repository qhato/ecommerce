package commands

import (
	"context"

	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
)

// CreateSKUCommand represents a command to create a SKU
type CreateSKUCommand struct {
	Name             string            `json:"name" validate:"required"`
	Description      string            `json:"description,omitempty"`
	LongDescription  string            `json:"long_description,omitempty"`
	UPC              string            `json:"upc,omitempty"`
	CurrencyCode     string            `json:"currency_code" validate:"required"`
	RetailPrice      float64           `json:"retail_price" validate:"required,min=0"`
	SalePrice        float64           `json:"sale_price,omitempty"`
	Cost             float64           `json:"cost,omitempty"`
	Available        bool              `json:"available"`
	Discountable     bool              `json:"discountable"`
	Taxable          bool              `json:"taxable"`
	TaxCode          string            `json:"tax_code,omitempty"`
	DefaultProductID *int64            `json:"default_product_id,omitempty"`
	Attributes       map[string]string `json:"attributes,omitempty"`
}

// UpdateSKUCommand represents a command to update a SKU
type UpdateSKUCommand struct {
	ID              int64             `json:"id" validate:"required"`
	Name            string            `json:"name,omitempty"`
	Description     string            `json:"description,omitempty"`
	LongDescription string            `json:"long_description,omitempty"`
	UPC             string            `json:"upc,omitempty"`
	RetailPrice     *float64          `json:"retail_price,omitempty"`
	SalePrice       *float64          `json:"sale_price,omitempty"`
	Cost            *float64          `json:"cost,omitempty"`
	Available       *bool             `json:"available,omitempty"`
	Discountable    *bool             `json:"discountable,omitempty"`
	Taxable         *bool             `json:"taxable,omitempty"`
	TaxCode         string            `json:"tax_code,omitempty"`
	Attributes      map[string]string `json:"attributes,omitempty"`
}

// UpdateSKUPricingCommand represents a command to update SKU pricing
type UpdateSKUPricingCommand struct {
	ID          int64   `json:"id" validate:"required"`
	RetailPrice float64 `json:"retail_price" validate:"required,min=0"`
	SalePrice   float64 `json:"sale_price,omitempty"`
}

// UpdateSKUAvailabilityCommand represents a command to update SKU availability
type UpdateSKUAvailabilityCommand struct {
	ID        int64 `json:"id" validate:"required"`
	Available bool  `json:"available"`
}

// DeleteSKUCommand represents a command to delete a SKU
type DeleteSKUCommand struct {
	ID int64 `json:"id" validate:"required"`
}

// SKUCommandHandler handles SKU commands
type SKUCommandHandler struct {
	repo      domain.SKURepository
	attrRepo  domain.SKUAttributeRepository
	eventBus  event.Bus
	validator *validator.Validator
	logger    *logger.Logger
}

// NewSKUCommandHandler creates a new SKU command handler
func NewSKUCommandHandler(
	repo domain.SKURepository,
	attrRepo domain.SKUAttributeRepository,
	eventBus event.Bus,
	validator *validator.Validator,
	logger *logger.Logger,
) *SKUCommandHandler {
	return &SKUCommandHandler{
		repo:      repo,
		attrRepo:  attrRepo,
		eventBus:  eventBus,
		validator: validator,
		logger:    logger,
	}
}

// HandleCreateSKU handles the create SKU command
func (h *SKUCommandHandler) HandleCreateSKU(ctx context.Context, cmd *CreateSKUCommand) (int64, error) {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return 0, errors.ValidationError("invalid create SKU command").WithInternal(err)
	}

	// Create SKU entity
	sku := domain.NewSKU(
		cmd.Name,
		cmd.Description,
		cmd.UPC,
		cmd.CurrencyCode,
		cmd.Cost,
		cmd.RetailPrice,
		cmd.SalePrice,
	)

	// Set optional fields
	sku.LongDescription = cmd.LongDescription
	sku.Available = cmd.Available
	sku.Discountable = cmd.Discountable
	sku.Taxable = cmd.Taxable
	sku.TaxCode = cmd.TaxCode
	sku.DefaultProductID = cmd.DefaultProductID

	// Save to repository
	if err := h.repo.Create(ctx, sku); err != nil {
		h.logger.WithError(err).Error("failed to create SKU")
		return 0, errors.InternalWrap(err, "failed to create SKU")
	}

	// Add attributes
	if cmd.Attributes != nil {
		for name, value := range cmd.Attributes {
			attr, err := domain.NewSKUAttribute(sku.ID, name, value)
			if err != nil {
				return 0, err
			}
			if err := h.attrRepo.Save(ctx, attr); err != nil {
				return 0, errors.InternalWrap(err, "failed to save SKU attribute")
			}
		}
	}

	// Publish domain event
	event := domain.NewSKUCreatedEvent(sku.ID, sku.DefaultProductID, sku.Name, sku.RetailPrice)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		h.logger.WithError(err).Error("failed to publish SKU created event")
	}

	h.logger.WithField("sku_id", sku.ID).Info("SKU created")
	return sku.ID, nil
}

// HandleUpdateSKU handles the update SKU command
func (h *SKUCommandHandler) HandleUpdateSKU(ctx context.Context, cmd *UpdateSKUCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.ValidationError("invalid update SKU command").WithInternal(err)
	}

	// Find existing SKU
	sku, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.InternalWrap(err, "SKU not found")
	}

	// Update fields if provided
	if cmd.Name != "" && cmd.Name != sku.Name {
		sku.Name = cmd.Name
	}
	if cmd.Description != "" || cmd.LongDescription != "" {
		sku.UpdateDescription(cmd.Description, cmd.LongDescription)
	}
	if cmd.UPC != "" && cmd.UPC != sku.UPC {
		sku.UPC = cmd.UPC
	}
	if cmd.Available != nil {
		sku.SetAvailability(*cmd.Available)
	}
	if cmd.Discountable != nil {
		sku.SetDiscountable(*cmd.Discountable)
	}
	if cmd.Taxable != nil {
		sku.SetTaxable(*cmd.Taxable, cmd.TaxCode)
	}
	if cmd.RetailPrice != nil || cmd.SalePrice != nil {
		retailPrice := sku.RetailPrice
		if cmd.RetailPrice != nil {
			retailPrice = *cmd.RetailPrice
		}
		salePrice := sku.SalePrice
		if cmd.SalePrice != nil {
			salePrice = *cmd.SalePrice
		}
		sku.UpdatePricing(retailPrice, salePrice)
	}

	// Update attributes
	if cmd.Attributes != nil {
		for name, value := range cmd.Attributes {
			attr, err := domain.NewSKUAttribute(sku.ID, name, value)
			if err != nil {
				return err
			}
			if err := h.attrRepo.Save(ctx, attr); err != nil {
				return errors.InternalWrap(err, "failed to save SKU attribute")
			}
		}
	}

	// Save to repository
	if err := h.repo.Update(ctx, sku); err != nil {
		h.logger.WithField("sku_id", cmd.ID).WithError(err).Error("failed to update SKU")
		return errors.InternalWrap(err, "failed to update SKU")
	}

	h.logger.WithField("sku_id", sku.ID).Info("SKU updated")
	return nil
}

// HandleUpdateSKUPricing handles the update SKU pricing command
func (h *SKUCommandHandler) HandleUpdateSKUPricing(ctx context.Context, cmd *UpdateSKUPricingCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.ValidationError("invalid update SKU pricing command").WithInternal(err)
	}

	// Find existing SKU
	sku, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.InternalWrap(err, "SKU not found")
	}

	// Track old price for event
	oldPrice := sku.RetailPrice

	// Update pricing
	sku.UpdatePricing(cmd.RetailPrice, cmd.SalePrice)

	// Save to repository
	if err := h.repo.Update(ctx, sku); err != nil {
		h.logger.WithField("sku_id", cmd.ID).WithError(err).Error("failed to update SKU pricing")
		return errors.InternalWrap(err, "failed to update SKU pricing")
	}

	// Publish price changed event if price actually changed
	if oldPrice != cmd.RetailPrice {
		event := domain.NewSKUPriceChangedEvent(sku.ID, oldPrice, cmd.RetailPrice)
		if err := h.eventBus.Publish(ctx, event); err != nil {
			h.logger.WithError(err).Error("failed to publish SKU price changed event")
		}
	}

	h.logger.WithField("sku_id", sku.ID).Info("SKU pricing updated")
	return nil
}

// HandleUpdateSKUAvailability handles the update SKU availability command
func (h *SKUCommandHandler) HandleUpdateSKUAvailability(ctx context.Context, cmd *UpdateSKUAvailabilityCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.ValidationError("invalid update SKU availability command").WithInternal(err)
	}

	// Update availability directly
	if err := h.repo.UpdateAvailability(ctx, cmd.ID, cmd.Available); err != nil {
		h.logger.WithField("sku_id", cmd.ID).WithError(err).Error("failed to update SKU availability")
		return errors.InternalWrap(err, "failed to update SKU availability")
	}

	// Publish availability changed event
	event := domain.NewSKUAvailabilityChangedEvent(cmd.ID, cmd.Available)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		h.logger.WithError(err).Error("failed to publish SKU availability changed event")
	}

	h.logger.WithField("sku_id", cmd.ID).WithField("available", cmd.Available).Info("SKU availability updated")
	return nil
}

// HandleDeleteSKU handles the delete SKU command
func (h *SKUCommandHandler) HandleDeleteSKU(ctx context.Context, cmd *DeleteSKUCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.ValidationError("invalid delete SKU command").WithInternal(err)
	}

	// Check if SKU exists
	_, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.InternalWrap(err, "SKU not found")
	}

	// Delete SKU
	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		h.logger.WithField("sku_id", cmd.ID).WithError(err).Error("failed to delete SKU")
		return errors.InternalWrap(err, "failed to delete SKU")
	}

	h.logger.WithField("sku_id", cmd.ID).Info("SKU deleted")
	return nil
}
