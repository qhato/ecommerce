package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/wishlist/domain"
)

type WishlistCommandHandler struct {
	wishlistRepo     domain.WishlistRepository
	wishlistItemRepo domain.WishlistItemRepository
}

func NewWishlistCommandHandler(
	wishlistRepo domain.WishlistRepository,
	wishlistItemRepo domain.WishlistItemRepository,
) *WishlistCommandHandler {
	return &WishlistCommandHandler{
		wishlistRepo:     wishlistRepo,
		wishlistItemRepo: wishlistItemRepo,
	}
}

func (h *WishlistCommandHandler) HandleCreateWishlist(ctx context.Context, cmd CreateWishlistCommand) (*domain.Wishlist, error) {
	wishlist, err := domain.NewWishlist(cmd.CustomerID, cmd.Name, cmd.IsDefault, cmd.IsPublic)
	if err != nil {
		return nil, fmt.Errorf("failed to create wishlist: %w", err)
	}

	if cmd.IsDefault {
		existingDefault, err := h.wishlistRepo.FindDefaultByCustomerID(ctx, cmd.CustomerID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing default wishlist: %w", err)
		}
		if existingDefault != nil {
			existingDefault.IsDefault = false
			if err := h.wishlistRepo.Update(ctx, existingDefault); err != nil {
				return nil, fmt.Errorf("failed to update existing default wishlist: %w", err)
			}
		}
	}

	if err := h.wishlistRepo.Create(ctx, wishlist); err != nil {
		return nil, fmt.Errorf("failed to save wishlist: %w", err)
	}

	return wishlist, nil
}

func (h *WishlistCommandHandler) HandleUpdateWishlist(ctx context.Context, cmd UpdateWishlistCommand) (*domain.Wishlist, error) {
	wishlist, err := h.wishlistRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find wishlist: %w", err)
	}
	if wishlist == nil {
		return nil, domain.ErrWishlistNotFound
	}

	if err := wishlist.Rename(cmd.Name); err != nil {
		return nil, err
	}
	wishlist.SetPublic(cmd.IsPublic)

	if err := h.wishlistRepo.Update(ctx, wishlist); err != nil {
		return nil, fmt.Errorf("failed to update wishlist: %w", err)
	}

	return wishlist, nil
}

func (h *WishlistCommandHandler) HandleDeleteWishlist(ctx context.Context, cmd DeleteWishlistCommand) error {
	wishlist, err := h.wishlistRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find wishlist: %w", err)
	}
	if wishlist == nil {
		return domain.ErrWishlistNotFound
	}

	if wishlist.IsDefault {
		return domain.ErrCannotDeleteDefault
	}

	if err := h.wishlistRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	return nil
}

func (h *WishlistCommandHandler) HandleSetDefaultWishlist(ctx context.Context, cmd SetDefaultWishlistCommand) (*domain.Wishlist, error) {
	wishlist, err := h.wishlistRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find wishlist: %w", err)
	}
	if wishlist == nil {
		return nil, domain.ErrWishlistNotFound
	}

	if wishlist.CustomerID != cmd.CustomerID {
		return nil, fmt.Errorf("wishlist does not belong to customer")
	}

	existingDefault, err := h.wishlistRepo.FindDefaultByCustomerID(ctx, cmd.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing default wishlist: %w", err)
	}
	if existingDefault != nil && existingDefault.ID != cmd.ID {
		existingDefault.IsDefault = false
		if err := h.wishlistRepo.Update(ctx, existingDefault); err != nil {
			return nil, fmt.Errorf("failed to update existing default wishlist: %w", err)
		}
	}

	wishlist.IsDefault = true
	if err := h.wishlistRepo.Update(ctx, wishlist); err != nil {
		return nil, fmt.Errorf("failed to update wishlist: %w", err)
	}

	return wishlist, nil
}

func (h *WishlistCommandHandler) HandleAddItem(ctx context.Context, cmd AddItemCommand) (*domain.WishlistItem, error) {
	wishlist, err := h.wishlistRepo.FindByID(ctx, cmd.WishlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to find wishlist: %w", err)
	}
	if wishlist == nil {
		return nil, domain.ErrWishlistNotFound
	}

	exists, err := h.wishlistItemRepo.ExistsByWishlistAndProduct(ctx, cmd.WishlistID, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if item exists: %w", err)
	}
	if exists {
		return nil, domain.ErrItemAlreadyInWishlist
	}

	item, err := domain.NewWishlistItem(cmd.WishlistID, cmd.ProductID, cmd.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to create wishlist item: %w", err)
	}

	item.SKUID = cmd.SKUID
	item.Notes = cmd.Notes

	if cmd.Priority > 0 {
		if err := item.SetPriority(cmd.Priority); err != nil {
			return nil, err
		}
	}

	if err := h.wishlistItemRepo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save wishlist item: %w", err)
	}

	return item, nil
}

func (h *WishlistCommandHandler) HandleUpdateItem(ctx context.Context, cmd UpdateItemCommand) (*domain.WishlistItem, error) {
	item, err := h.wishlistItemRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find wishlist item: %w", err)
	}
	if item == nil {
		return nil, domain.ErrWishlistItemNotFound
	}

	if cmd.Quantity > 0 {
		if err := item.UpdateQuantity(cmd.Quantity); err != nil {
			return nil, err
		}
	}

	if cmd.Priority > 0 {
		if err := item.SetPriority(cmd.Priority); err != nil {
			return nil, err
		}
	}

	item.Notes = cmd.Notes

	if err := h.wishlistItemRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to update wishlist item: %w", err)
	}

	return item, nil
}

func (h *WishlistCommandHandler) HandleRemoveItem(ctx context.Context, cmd RemoveItemCommand) error {
	item, err := h.wishlistItemRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find wishlist item: %w", err)
	}
	if item == nil {
		return domain.ErrWishlistItemNotFound
	}

	if err := h.wishlistItemRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete wishlist item: %w", err)
	}

	return nil
}

func (h *WishlistCommandHandler) HandleMoveItem(ctx context.Context, cmd MoveItemCommand) (*domain.WishlistItem, error) {
	item, err := h.wishlistItemRepo.FindByID(ctx, cmd.ItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find wishlist item: %w", err)
	}
	if item == nil {
		return nil, domain.ErrWishlistItemNotFound
	}

	targetWishlist, err := h.wishlistRepo.FindByID(ctx, cmd.TargetWishlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to find target wishlist: %w", err)
	}
	if targetWishlist == nil {
		return nil, domain.ErrWishlistNotFound
	}

	exists, err := h.wishlistItemRepo.ExistsByWishlistAndProduct(ctx, cmd.TargetWishlistID, item.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if item exists in target: %w", err)
	}
	if exists {
		return nil, domain.ErrItemAlreadyInWishlist
	}

	item.WishlistID = cmd.TargetWishlistID
	if err := h.wishlistItemRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to move wishlist item: %w", err)
	}

	return item, nil
}
