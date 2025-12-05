package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/pricing/domain"
)

// PricingCommandHandler handles pricing commands
type PricingCommandHandler struct {
	priceListRepo     domain.PriceListRepository
	priceListItemRepo domain.PriceListItemRepository
	pricingRuleRepo   domain.PricingRuleRepository
}

// NewPricingCommandHandler creates a new PricingCommandHandler
func NewPricingCommandHandler(
	priceListRepo domain.PriceListRepository,
	priceListItemRepo domain.PriceListItemRepository,
	pricingRuleRepo domain.PricingRuleRepository,
) *PricingCommandHandler {
	return &PricingCommandHandler{
		priceListRepo:     priceListRepo,
		priceListItemRepo: priceListItemRepo,
		pricingRuleRepo:   pricingRuleRepo,
	}
}

// HandleCreatePriceList handles creating a new price list
func (h *PricingCommandHandler) HandleCreatePriceList(ctx context.Context, cmd *CreatePriceListCommand) (int64, error) {
	// Check if code already exists
	existing, err := h.priceListRepo.FindByCode(ctx, cmd.Code)
	if err == nil && existing != nil {
		return 0, domain.ErrPriceListCodeExists
	}

	// Create price list
	priceList, err := domain.NewPriceList(
		cmd.Name,
		cmd.Code,
		cmd.PriceListType,
		cmd.Currency,
		cmd.Priority,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create price list: %w", err)
	}

	priceList.Description = cmd.Description
	priceList.SetDateRange(cmd.StartDate, cmd.EndDate)

	// Add customer segments
	for _, segment := range cmd.CustomerSegments {
		priceList.AddCustomerSegment(segment)
	}

	// Save to repository
	err = h.priceListRepo.Save(ctx, priceList)
	if err != nil {
		return 0, fmt.Errorf("failed to save price list: %w", err)
	}

	return priceList.ID, nil
}

// HandleUpdatePriceList handles updating an existing price list
func (h *PricingCommandHandler) HandleUpdatePriceList(ctx context.Context, cmd *UpdatePriceListCommand) error {
	priceList, err := h.priceListRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find price list: %w", err)
	}
	if priceList == nil {
		return domain.ErrPriceListNotFound
	}

	// Apply updates
	if cmd.Name != nil {
		priceList.Name = *cmd.Name
	}
	if cmd.Priority != nil {
		priceList.Priority = *cmd.Priority
	}
	if cmd.IsActive != nil {
		if *cmd.IsActive {
			priceList.Activate()
		} else {
			priceList.Deactivate()
		}
	}
	if cmd.Description != nil {
		priceList.Description = *cmd.Description
	}
	if cmd.StartDate != nil || cmd.EndDate != nil {
		priceList.SetDateRange(cmd.StartDate, cmd.EndDate)
	}

	// Update customer segments if provided
	if cmd.CustomerSegments != nil {
		priceList.CustomerSegments = cmd.CustomerSegments
	}

	// Save to repository
	err = h.priceListRepo.Save(ctx, priceList)
	if err != nil {
		return fmt.Errorf("failed to save price list: %w", err)
	}

	return nil
}

// HandleDeletePriceList handles deleting a price list
func (h *PricingCommandHandler) HandleDeletePriceList(ctx context.Context, id int64) error {
	// Delete all items in the price list first
	err := h.priceListItemRepo.DeleteByPriceListID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete price list items: %w", err)
	}

	// Delete the price list
	err = h.priceListRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete price list: %w", err)
	}

	return nil
}

// HandleCreatePriceListItem handles creating a new price list item
func (h *PricingCommandHandler) HandleCreatePriceListItem(ctx context.Context, cmd *CreatePriceListItemCommand) (int64, error) {
	// Verify price list exists
	priceList, err := h.priceListRepo.FindByID(ctx, cmd.PriceListID)
	if err != nil {
		return 0, fmt.Errorf("failed to find price list: %w", err)
	}
	if priceList == nil {
		return 0, domain.ErrPriceListNotFound
	}

	// Create price list item
	item, err := domain.NewPriceListItem(
		cmd.PriceListID,
		cmd.SKUID,
		cmd.Price,
		cmd.MinQuantity,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create price list item: %w", err)
	}

	item.ProductID = cmd.ProductID
	if cmd.CompareAtPrice != nil {
		item.SetCompareAtPrice(*cmd.CompareAtPrice)
	}
	if cmd.MaxQuantity != nil {
		err = item.SetQuantityRange(cmd.MinQuantity, cmd.MaxQuantity)
		if err != nil {
			return 0, err
		}
	}
	item.SetDateRange(cmd.StartDate, cmd.EndDate)

	// Save to repository
	err = h.priceListItemRepo.Save(ctx, item)
	if err != nil {
		return 0, fmt.Errorf("failed to save price list item: %w", err)
	}

	return item.ID, nil
}

// HandleUpdatePriceListItem handles updating an existing price list item
func (h *PricingCommandHandler) HandleUpdatePriceListItem(ctx context.Context, cmd *UpdatePriceListItemCommand) error {
	item, err := h.priceListItemRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find price list item: %w", err)
	}
	if item == nil {
		return domain.ErrPriceListItemNotFound
	}

	// Apply updates
	if cmd.Price != nil {
		err = item.UpdatePrice(*cmd.Price)
		if err != nil {
			return err
		}
	}
	if cmd.CompareAtPrice != nil {
		item.SetCompareAtPrice(*cmd.CompareAtPrice)
	}
	if cmd.MinQuantity != nil || cmd.MaxQuantity != nil {
		minQty := item.MinQuantity
		if cmd.MinQuantity != nil {
			minQty = *cmd.MinQuantity
		}
		err = item.SetQuantityRange(minQty, cmd.MaxQuantity)
		if err != nil {
			return err
		}
	}
	if cmd.IsActive != nil {
		if *cmd.IsActive {
			item.Activate()
		} else {
			item.Deactivate()
		}
	}
	if cmd.StartDate != nil || cmd.EndDate != nil {
		item.SetDateRange(cmd.StartDate, cmd.EndDate)
	}

	// Save to repository
	err = h.priceListItemRepo.Save(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to save price list item: %w", err)
	}

	return nil
}

// HandleDeletePriceListItem handles deleting a price list item
func (h *PricingCommandHandler) HandleDeletePriceListItem(ctx context.Context, id int64) error {
	err := h.priceListItemRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete price list item: %w", err)
	}
	return nil
}

// HandleBulkCreatePriceListItems handles bulk creation of price list items
func (h *PricingCommandHandler) HandleBulkCreatePriceListItems(ctx context.Context, cmd *BulkCreatePriceListItemsCommand) error {
	// Verify price list exists
	priceList, err := h.priceListRepo.FindByID(ctx, cmd.PriceListID)
	if err != nil {
		return fmt.Errorf("failed to find price list: %w", err)
	}
	if priceList == nil {
		return domain.ErrPriceListNotFound
	}

	// Create all items
	for _, bulkItem := range cmd.Items {
		item, err := domain.NewPriceListItem(
			cmd.PriceListID,
			bulkItem.SKUID,
			bulkItem.Price,
			bulkItem.MinQuantity,
		)
		if err != nil {
			return fmt.Errorf("failed to create price list item for SKU %s: %w", bulkItem.SKUID, err)
		}

		item.ProductID = bulkItem.ProductID
		if bulkItem.CompareAtPrice != nil {
			item.SetCompareAtPrice(*bulkItem.CompareAtPrice)
		}
		if bulkItem.MaxQuantity != nil {
			err = item.SetQuantityRange(bulkItem.MinQuantity, bulkItem.MaxQuantity)
			if err != nil {
				return fmt.Errorf("failed to set quantity range for SKU %s: %w", bulkItem.SKUID, err)
			}
		}

		err = h.priceListItemRepo.Save(ctx, item)
		if err != nil {
			return fmt.Errorf("failed to save price list item for SKU %s: %w", bulkItem.SKUID, err)
		}
	}

	return nil
}

// HandleCreatePricingRule handles creating a new pricing rule
func (h *PricingCommandHandler) HandleCreatePricingRule(ctx context.Context, cmd *CreatePricingRuleCommand) (int64, error) {
	// Create pricing rule
	rule, err := domain.NewPricingRule(cmd.Name, cmd.RuleType, cmd.Priority)
	if err != nil {
		return 0, fmt.Errorf("failed to create pricing rule: %w", err)
	}

	rule.Description = cmd.Description
	rule.ConditionExpression = cmd.ConditionExpression
	rule.SetAction(cmd.ActionType, cmd.ActionValue)

	// Set quantity range
	if cmd.MaxQuantity != nil {
		err = rule.SetQuantityRange(cmd.MinQuantity, cmd.MaxQuantity)
		if err != nil {
			return 0, err
		}
	} else {
		rule.MinQuantity = cmd.MinQuantity
	}

	// Set applicable SKUs
	rule.ApplicableSKUs = cmd.ApplicableSKUs
	rule.ApplicableCategories = cmd.ApplicableCategories

	// Set customer segments
	for _, segment := range cmd.CustomerSegments {
		rule.AddCustomerSegment(segment)
	}

	rule.MinOrderValue = cmd.MinOrderValue
	rule.StartDate = cmd.StartDate
	rule.EndDate = cmd.EndDate

	// Save to repository
	err = h.pricingRuleRepo.Save(ctx, rule)
	if err != nil {
		return 0, fmt.Errorf("failed to save pricing rule: %w", err)
	}

	return rule.ID, nil
}

// HandleUpdatePricingRule handles updating an existing pricing rule
func (h *PricingCommandHandler) HandleUpdatePricingRule(ctx context.Context, cmd *UpdatePricingRuleCommand) error {
	rule, err := h.pricingRuleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find pricing rule: %w", err)
	}
	if rule == nil {
		return domain.ErrPricingRuleNotFound
	}

	// Apply updates
	if cmd.Name != nil {
		rule.Name = *cmd.Name
	}
	if cmd.Description != nil {
		rule.Description = *cmd.Description
	}
	if cmd.Priority != nil {
		rule.Priority = *cmd.Priority
	}
	if cmd.IsActive != nil {
		if *cmd.IsActive {
			rule.Activate()
		} else {
			rule.Deactivate()
		}
	}
	if cmd.ConditionExpression != nil {
		rule.ConditionExpression = *cmd.ConditionExpression
	}
	if cmd.ActionType != nil && cmd.ActionValue != nil {
		rule.SetAction(*cmd.ActionType, *cmd.ActionValue)
	}
	if cmd.MinQuantity != nil || cmd.MaxQuantity != nil {
		minQty := rule.MinQuantity
		if cmd.MinQuantity != nil {
			minQty = *cmd.MinQuantity
		}
		err = rule.SetQuantityRange(minQty, cmd.MaxQuantity)
		if err != nil {
			return err
		}
	}

	// Update collections if provided
	if cmd.ApplicableSKUs != nil {
		rule.ApplicableSKUs = cmd.ApplicableSKUs
	}
	if cmd.ApplicableCategories != nil {
		rule.ApplicableCategories = cmd.ApplicableCategories
	}
	if cmd.CustomerSegments != nil {
		rule.CustomerSegments = cmd.CustomerSegments
	}

	if cmd.MinOrderValue != nil {
		rule.MinOrderValue = cmd.MinOrderValue
	}
	if cmd.StartDate != nil {
		rule.StartDate = cmd.StartDate
	}
	if cmd.EndDate != nil {
		rule.EndDate = cmd.EndDate
	}

	// Save to repository
	err = h.pricingRuleRepo.Save(ctx, rule)
	if err != nil {
		return fmt.Errorf("failed to save pricing rule: %w", err)
	}

	return nil
}

// HandleDeletePricingRule handles deleting a pricing rule
func (h *PricingCommandHandler) HandleDeletePricingRule(ctx context.Context, id int64) error {
	err := h.pricingRuleRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete pricing rule: %w", err)
	}
	return nil
}
