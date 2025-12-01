package commands

import (
	"context"

	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
)

// CreateProductCommand represents a command to create a product
type CreateProductCommand struct {
	Manufacture           string            `json:"manufacture" validate:"required"`
	Model                 string            `json:"model" validate:"required"`
	URL                   string            `json:"url" validate:"required,url"`
	URLKey                string            `json:"url_key" validate:"required"`
	CanSellWithoutOptions bool              `json:"can_sell_without_options"`
	EnableDefaultSKU      bool              `json:"enable_default_sku"`
	CanonicalURL          string            `json:"canonical_url,omitempty" validate:"omitempty,url"`
	DisplayTemplate       string            `json:"display_template,omitempty"`
	MetaDescription       string            `json:"meta_description,omitempty"`
	MetaTitle             string            `json:"meta_title,omitempty"`
	OverrideGeneratedURL  bool              `json:"override_generated_url"`
	DefaultCategoryID     *int64            `json:"default_category_id,omitempty"`
	Attributes            map[string]string `json:"attributes,omitempty"`
}

// UpdateProductCommand represents a command to update a product
type UpdateProductCommand struct {
	ID                    int64             `json:"id" validate:"required"`
	Manufacture           string            `json:"manufacture,omitempty"`
	Model                 string            `json:"model,omitempty"`
	URL                   string            `json:"url,omitempty" validate:"omitempty,url"`
	URLKey                string            `json:"url_key,omitempty"`
	CanSellWithoutOptions *bool             `json:"can_sell_without_options,omitempty"`
	EnableDefaultSKU      *bool             `json:"enable_default_sku,omitempty"`
	CanonicalURL          string            `json:"canonical_url,omitempty" validate:"omitempty,url"`
	DisplayTemplate       string            `json:"display_template,omitempty"`
	MetaDescription       string            `json:"meta_description,omitempty"`
	MetaTitle             string            `json:"meta_title,omitempty"`
	OverrideGeneratedURL  *bool             `json:"override_generated_url,omitempty"`
	DefaultCategoryID     *int64            `json:"default_category_id,omitempty"`
	DefaultSKUID          *int64            `json:"default_sku_id,omitempty"`
	Attributes            map[string]string `json:"attributes,omitempty"`
}

// DeleteProductCommand represents a command to delete a product
type DeleteProductCommand struct {
	ID int64 `json:"id" validate:"required"`
}

// ArchiveProductCommand represents a command to archive a product
type ArchiveProductCommand struct {
	ID int64 `json:"id" validate:"required"`
}

// ProductCommandHandler handles product commands
type ProductCommandHandler struct {
	repo      domain.ProductRepository
	attrRepo  domain.ProductAttributeRepository
	eventBus  event.Bus
	validator *validator.Validator
	logger    *logger.Logger
}

// NewProductCommandHandler creates a new product command handler
func NewProductCommandHandler(
	repo domain.ProductRepository,
	attrRepo domain.ProductAttributeRepository,
	eventBus event.Bus,
	validator *validator.Validator,
	logger *logger.Logger,
) *ProductCommandHandler {
	return &ProductCommandHandler{
		repo:      repo,
		attrRepo:  attrRepo,
		eventBus:  eventBus,
		validator: validator,
		logger:    logger,
	}
}

// HandleCreateProduct handles the create product command
func (h *ProductCommandHandler) HandleCreateProduct(ctx context.Context, cmd *CreateProductCommand) (int64, error) {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return 0, errors.ValidationError("invalid create product command").WithInternal(err)
	}

	// Create product entity
	product := domain.NewProduct(
		cmd.Manufacture,
		cmd.Model,
		cmd.URL,
		cmd.URLKey,
		cmd.CanSellWithoutOptions,
		cmd.EnableDefaultSKU,
	)

	// Set optional fields
	product.CanonicalURL = cmd.CanonicalURL
	product.DisplayTemplate = cmd.DisplayTemplate
	product.MetaDescription = cmd.MetaDescription
	product.MetaTitle = cmd.MetaTitle
	product.OverrideGeneratedURL = cmd.OverrideGeneratedURL
	if cmd.DefaultCategoryID != nil {
		product.SetDefaultCategory(*cmd.DefaultCategoryID)
	}

	// Save to repository
	if err := h.repo.Create(ctx, product); err != nil {
		h.logger.WithError(err).Error("failed to create product")
		return 0, errors.InternalWrap(err, "failed to create product")
	}

	// Add attributes
	if cmd.Attributes != nil {
		for name, value := range cmd.Attributes {
			attr, err := domain.NewProductAttribute(product.ID, name, value)
			if err != nil {
				return 0, err
			}
			if err := h.attrRepo.Save(ctx, attr); err != nil {
				return 0, errors.InternalWrap(err, "failed to save product attribute")
			}
		}
	}

	// Publish domain event
	event := domain.NewProductCreatedEvent(product.ID, product.Model, product.Manufacture)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		h.logger.WithError(err).Error("failed to publish product created event")
	}

	h.logger.WithField("product_id", product.ID).Info("product created")
	return product.ID, nil
}

// HandleUpdateProduct handles the update product command
func (h *ProductCommandHandler) HandleUpdateProduct(ctx context.Context, cmd *UpdateProductCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.ValidationError("invalid update product command").WithInternal(err)
	}

	// Find existing product
	product, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.InternalWrap(err, "product not found")
	}

	if product.IsArchived() {
		return errors.Conflict("cannot update archived product")
	}

	// Track changes for event
	changes := make(map[string]interface{})

	// Update fields if provided
	if cmd.Manufacture != "" && cmd.Manufacture != product.Manufacture {
		changes["manufacture"] = cmd.Manufacture
		product.Manufacture = cmd.Manufacture
	}
	if cmd.Model != "" && cmd.Model != product.Model {
		changes["model"] = cmd.Model
		product.Model = cmd.Model
	}
	if cmd.URL != "" && cmd.URL != product.URL {
		product.UpdateURLs(cmd.URL, cmd.URLKey, cmd.OverrideGeneratedURL != nil && *cmd.OverrideGeneratedURL)
		changes["url"] = cmd.URL
	}
	if cmd.MetaTitle != "" || cmd.MetaDescription != "" {
		product.UpdateMetadata(cmd.MetaTitle, cmd.MetaDescription)
		changes["metadata"] = true
	}
	if cmd.DefaultCategoryID != nil {
		product.SetDefaultCategory(*cmd.DefaultCategoryID)
		changes["default_category_id"] = *cmd.DefaultCategoryID
	}
	if cmd.DefaultSKUID != nil {
		product.SetDefaultSKU(*cmd.DefaultSKUID)
		changes["default_sku_id"] = *cmd.DefaultSKUID
	}

	// Update attributes
	if cmd.Attributes != nil {
		for name, value := range cmd.Attributes {
			attr, err := domain.NewProductAttribute(product.ID, name, value)
			if err != nil {
				return err
			}
			if err := h.attrRepo.Save(ctx, attr); err != nil {
				return errors.InternalWrap(err, "failed to save product attribute")
			}
		}
		changes["attributes"] = true
	}

	// Save to repository
	if err := h.repo.Update(ctx, product); err != nil {
		h.logger.WithField("product_id", cmd.ID).WithError(err).Error("failed to update product")
		return errors.InternalWrap(err, "failed to update product")
	}

	// Publish domain event
	if len(changes) > 0 {
		event := domain.NewProductUpdatedEvent(product.ID, changes)
		if err := h.eventBus.Publish(ctx, event); err != nil {
			h.logger.WithError(err).Error("failed to publish product updated event")
		}
	}

	h.logger.WithField("product_id", product.ID).Info("product updated")
	return nil
}

// HandleDeleteProduct handles the delete product command
func (h *ProductCommandHandler) HandleDeleteProduct(ctx context.Context, cmd *DeleteProductCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.ValidationError("invalid delete product command").WithInternal(err)
	}

	// Check if product exists
	_, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.InternalWrap(err, "product not found")
	}

	// Soft delete (archive)
	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		h.logger.WithField("product_id", cmd.ID).WithError(err).Error("failed to delete product")
		return errors.InternalWrap(err, "failed to delete product")
	}

	// Publish domain event
	event := domain.NewProductArchivedEvent(cmd.ID)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		h.logger.WithError(err).Error("failed to publish product archived event")
	}

	h.logger.WithField("product_id", cmd.ID).Info("product deleted (archived)")
	return nil
}

// HandleArchiveProduct handles the archive product command
func (h *ProductCommandHandler) HandleArchiveProduct(ctx context.Context, cmd *ArchiveProductCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.ValidationError("invalid archive product command").WithInternal(err)
	}

	// Find product
	product, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.InternalWrap(err, "product not found")
	}

	if product.IsArchived() {
		return errors.Conflict("product is already archived")
	}

	// Archive product
	product.Archive()

	// Save to repository
	if err := h.repo.Update(ctx, product); err != nil {
		h.logger.WithField("product_id", cmd.ID).WithError(err).Error("failed to archive product")
		return errors.InternalWrap(err, "failed to archive product")
	}

	// Publish domain event
	event := domain.NewProductArchivedEvent(product.ID)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		h.logger.WithError(err).Error("failed to publish product archived event")
	}

	h.logger.WithField("product_id", cmd.ID).Info("product archived")
	return nil
}
