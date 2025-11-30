package commands

import (
	"context"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
)

// CreateCategoryCommand represents a command to create a category
type CreateCategoryCommand struct {
	Name                    string            `json:"name" validate:"required"`
	Description             string            `json:"description,omitempty"`
	LongDescription         string            `json:"long_description,omitempty"`
	URL                     string            `json:"url" validate:"required,url"`
	URLKey                  string            `json:"url_key" validate:"required"`
	ActiveStartDate         *time.Time        `json:"active_start_date,omitempty"`
	ActiveEndDate           *time.Time        `json:"active_end_date,omitempty"`
	DisplayTemplate         string            `json:"display_template,omitempty"`
	ExternalID              string            `json:"external_id,omitempty"`
	FulfillmentType         string            `json:"fulfillment_type,omitempty"`
	InventoryType           string            `json:"inventory_type,omitempty"`
	MetaDescription         string            `json:"meta_description,omitempty"`
	MetaTitle               string            `json:"meta_title,omitempty"`
	OverrideGeneratedURL    bool              `json:"override_generated_url"`
	RootDisplayOrder        float64           `json:"root_display_order"`
	TaxCode                 string            `json:"tax_code,omitempty"`
	DefaultParentCategoryID *int64            `json:"default_parent_category_id,omitempty"`
	Attributes              map[string]string `json:"attributes,omitempty"`
}

// UpdateCategoryCommand represents a command to update a category
type UpdateCategoryCommand struct {
	ID                      int64             `json:"id" validate:"required"`
	Name                    string            `json:"name,omitempty"`
	Description             string            `json:"description,omitempty"`
	LongDescription         string            `json:"long_description,omitempty"`
	URL                     string            `json:"url,omitempty" validate:"omitempty,url"`
	URLKey                  string            `json:"url_key,omitempty"`
	ActiveStartDate         *time.Time        `json:"active_start_date,omitempty"`
	ActiveEndDate           *time.Time        `json:"active_end_date,omitempty"`
	DisplayTemplate         string            `json:"display_template,omitempty"`
	MetaDescription         string            `json:"meta_description,omitempty"`
	MetaTitle               string            `json:"meta_title,omitempty"`
	OverrideGeneratedURL    *bool             `json:"override_generated_url,omitempty"`
	RootDisplayOrder        *float64          `json:"root_display_order,omitempty"`
	TaxCode                 string            `json:"tax_code,omitempty"`
	DefaultParentCategoryID *int64            `json:"default_parent_category_id,omitempty"`
	Attributes              map[string]string `json:"attributes,omitempty"`
}

// DeleteCategoryCommand represents a command to delete a category
type DeleteCategoryCommand struct {
	ID int64 `json:"id" validate:"required"`
}

// CategoryCommandHandler handles category commands
type CategoryCommandHandler struct {
	repo      domain.CategoryRepository
	attrRepo  domain.CategoryAttributeRepository
	eventBus  event.Bus
	validator *validator.Validator
	logger    *logger.Logger
}

// NewCategoryCommandHandler creates a new category command handler
func NewCategoryCommandHandler(
	repo domain.CategoryRepository,
	attrRepo domain.CategoryAttributeRepository,
	eventBus event.Bus,
	validator *validator.Validator,
	logger *logger.Logger,
) *CategoryCommandHandler {
	return &CategoryCommandHandler{
		repo:      repo,
		attrRepo:  attrRepo,
		eventBus:  eventBus,
		validator: validator,
		logger:    logger,
	}
}

// HandleCreateCategory handles the create category command
func (h *CategoryCommandHandler) HandleCreateCategory(ctx context.Context, cmd *CreateCategoryCommand) (int64, error) {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return 0, errors.ValidationError("invalid create category command").WithInternal(err)
	}

	// Create category entity
	category := domain.NewCategory(
		cmd.Name,
		cmd.Description,
		cmd.URL,
		cmd.URLKey,
	)

	// Set optional fields
	category.LongDescription = cmd.LongDescription
	category.DisplayTemplate = cmd.DisplayTemplate
	category.ExternalID = cmd.ExternalID
	category.FulfillmentType = cmd.FulfillmentType
	category.InventoryType = cmd.InventoryType
	category.MetaDescription = cmd.MetaDescription
	category.MetaTitle = cmd.MetaTitle
	category.OverrideGeneratedURL = cmd.OverrideGeneratedURL
	category.RootDisplayOrder = cmd.RootDisplayOrder
	category.TaxCode = cmd.TaxCode
	category.DefaultParentCategoryID = cmd.DefaultParentCategoryID

	// Set active dates
	if cmd.ActiveStartDate != nil || cmd.ActiveEndDate != nil {
		category.SetActiveDate(cmd.ActiveStartDate, cmd.ActiveEndDate)
	}

	// Save to repository
	if err := h.repo.Create(ctx, category); err != nil {
		h.logger.WithError(err).Error("failed to create category")
		return 0, errors.InternalWrap(err, "failed to create category")
	}

	// Add attributes
	if cmd.Attributes != nil {
		for name, value := range cmd.Attributes {
			attr, err := domain.NewCategoryAttribute(category.ID, name, value)
			if err != nil {
				return 0, err
			}
			if err := h.attrRepo.Save(ctx, attr); err != nil {
				return 0, errors.InternalWrap(err, "failed to save category attribute")
			}
		}
	}

	// Publish domain event
	event := domain.NewCategoryCreatedEvent(category.ID, category.Name, category.DefaultParentCategoryID)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		h.logger.WithError(err).Error("failed to publish category created event")
	}

	h.logger.WithField("category_id", category.ID).Info("category created")
	return category.ID, nil
}

// HandleUpdateCategory handles the update category command
func (h *CategoryCommandHandler) HandleUpdateCategory(ctx context.Context, cmd *UpdateCategoryCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.ValidationError("invalid update category command").WithInternal(err)
	}

	// Find existing category
	category, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.InternalWrap(err, "category not found")
	}

	if category.Archived {
		return errors.Conflict("cannot update archived category")
	}

	// Track changes for event
	changes := make(map[string]interface{})

	// Update fields if provided
	if cmd.Name != "" && cmd.Name != category.Name {
		changes["name"] = cmd.Name
		category.Name = cmd.Name
	}
	if cmd.Description != "" || cmd.LongDescription != "" {
		category.UpdateDescription(cmd.Description, cmd.LongDescription)
		changes["description"] = true
	}
	if cmd.URL != "" && cmd.URL != category.URL {
		overrideGenerated := cmd.OverrideGeneratedURL != nil && *cmd.OverrideGeneratedURL
		category.UpdateURLs(cmd.URL, cmd.URLKey, overrideGenerated)
		changes["url"] = cmd.URL
	}
	if cmd.MetaTitle != "" || cmd.MetaDescription != "" {
		category.UpdateMetadata(cmd.MetaTitle, cmd.MetaDescription)
		changes["metadata"] = true
	}
	if cmd.RootDisplayOrder != nil {
		category.SetDisplayOrder(*cmd.RootDisplayOrder)
		changes["display_order"] = *cmd.RootDisplayOrder
	}
	if cmd.DefaultParentCategoryID != nil {
		if *cmd.DefaultParentCategoryID == 0 {
			category.RemoveDefaultParentCategory()
		} else {
			category.SetParentCategory(*cmd.DefaultParentCategoryID)
		}
		changes["parent_category_id"] = *cmd.DefaultParentCategoryID
	}
	if cmd.ActiveStartDate != nil || cmd.ActiveEndDate != nil {
		category.SetActiveDate(cmd.ActiveStartDate, cmd.ActiveEndDate)
		changes["active_dates"] = true
	}

	// Update attributes
	if cmd.Attributes != nil {
		for name, value := range cmd.Attributes {
			attr, err := domain.NewCategoryAttribute(category.ID, name, value)
			if err != nil {
				return err
			}
			if err := h.attrRepo.Save(ctx, attr); err != nil {
				return errors.InternalWrap(err, "failed to save category attribute")
			}
		}
		changes["attributes"] = true
	}

	// Save to repository
	if err := h.repo.Update(ctx, category); err != nil {
		h.logger.WithField("category_id", cmd.ID).WithError(err).Error("failed to update category")
		return errors.InternalWrap(err, "failed to update category")
	}

	// Publish domain event
	if len(changes) > 0 {
		event := domain.NewCategoryUpdatedEvent(category.ID, changes)
		if err := h.eventBus.Publish(ctx, event); err != nil {
			h.logger.WithError(err).Error("failed to publish category updated event")
		}
	}

	h.logger.WithField("category_id", category.ID).Info("category updated")
	return nil
}

// HandleDeleteCategory handles the delete category command
func (h *CategoryCommandHandler) HandleDeleteCategory(ctx context.Context, cmd *DeleteCategoryCommand) error {
	// Validate command
	if err := h.validator.Validate(cmd); err != nil {
		return errors.ValidationError("invalid delete category command").WithInternal(err)
	}

	// Check if category exists
	_, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.InternalWrap(err, "category not found")
	}

	// Soft delete (archive)
	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		h.logger.WithField("category_id", cmd.ID).WithError(err).Error("failed to delete category")
		return errors.InternalWrap(err, "failed to delete category")
	}

	// Publish domain event
	event := domain.NewCategoryUpdatedEvent(cmd.ID, map[string]interface{}{"archived": true})
	if err := h.eventBus.Publish(ctx, event); err != nil {
		h.logger.WithError(err).Error("failed to publish category archived event")
	}

	h.logger.WithField("category_id", cmd.ID).Info("category deleted (archived)")
	return nil
}
