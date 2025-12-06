package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

type ProductBundleCommandHandler struct {
	bundleRepo     domain.ProductBundleRepository
	bundleItemRepo domain.ProductBundleItemRepository
}

func NewProductBundleCommandHandler(
	bundleRepo domain.ProductBundleRepository,
	bundleItemRepo domain.ProductBundleItemRepository,
) *ProductBundleCommandHandler {
	return &ProductBundleCommandHandler{
		bundleRepo:     bundleRepo,
		bundleItemRepo: bundleItemRepo,
	}
}

func (h *ProductBundleCommandHandler) HandleCreateProductBundle(ctx context.Context, cmd CreateProductBundleCommand) (*domain.ProductBundle, error) {
	bundle := domain.NewProductBundle(cmd.Name, cmd.Description, cmd.BundlePrice)

	if err := h.bundleRepo.Create(ctx, bundle); err != nil {
		return nil, fmt.Errorf("failed to create bundle: %w", err)
	}

	for _, itemInput := range cmd.Items {
		if err := bundle.AddItem(itemInput.ProductID, itemInput.SKUID, itemInput.Quantity, itemInput.SortOrder); err != nil {
			return nil, err
		}
		item := &bundle.Items[len(bundle.Items)-1]
		item.BundleID = bundle.ID
		if err := h.bundleItemRepo.Create(ctx, item); err != nil {
			return nil, fmt.Errorf("failed to create bundle item: %w", err)
		}
	}

	return bundle, nil
}

func (h *ProductBundleCommandHandler) HandleUpdateProductBundle(ctx context.Context, cmd UpdateProductBundleCommand) (*domain.ProductBundle, error) {
	bundle, err := h.bundleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find bundle: %w", err)
	}
	if bundle == nil {
		return nil, domain.ErrBundleNotFound
	}

	bundle.Name = cmd.Name
	bundle.Description = cmd.Description
	if err := bundle.UpdatePrice(cmd.BundlePrice); err != nil {
		return nil, err
	}

	if err := h.bundleRepo.Update(ctx, bundle); err != nil {
		return nil, fmt.Errorf("failed to update bundle: %w", err)
	}

	return bundle, nil
}

func (h *ProductBundleCommandHandler) HandleActivateBundle(ctx context.Context, cmd ActivateBundleCommand) error {
	bundle, err := h.bundleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find bundle: %w", err)
	}
	if bundle == nil {
		return domain.ErrBundleNotFound
	}

	bundle.Activate()
	return h.bundleRepo.Update(ctx, bundle)
}

func (h *ProductBundleCommandHandler) HandleDeactivateBundle(ctx context.Context, cmd DeactivateBundleCommand) error {
	bundle, err := h.bundleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find bundle: %w", err)
	}
	if bundle == nil {
		return domain.ErrBundleNotFound
	}

	bundle.Deactivate()
	return h.bundleRepo.Update(ctx, bundle)
}

func (h *ProductBundleCommandHandler) HandleDeleteBundle(ctx context.Context, cmd DeleteBundleCommand) error {
	if err := h.bundleItemRepo.DeleteByBundleID(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete bundle items: %w", err)
	}
	return h.bundleRepo.Delete(ctx, cmd.ID)
}
