package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/wishlist/domain"
)

type WishlistQueryService struct {
	wishlistRepo     domain.WishlistRepository
	wishlistItemRepo domain.WishlistItemRepository
}

func NewWishlistQueryService(
	wishlistRepo domain.WishlistRepository,
	wishlistItemRepo domain.WishlistItemRepository,
) *WishlistQueryService {
	return &WishlistQueryService{
		wishlistRepo:     wishlistRepo,
		wishlistItemRepo: wishlistItemRepo,
	}
}

func (s *WishlistQueryService) GetWishlist(ctx context.Context, id string) (*WishlistDTO, error) {
	wishlist, err := s.wishlistRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find wishlist: %w", err)
	}
	if wishlist == nil {
		return nil, domain.ErrWishlistNotFound
	}

	items, err := s.wishlistItemRepo.FindByWishlistID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find wishlist items: %w", err)
	}

	wishlist.Items = make([]domain.WishlistItem, len(items))
	for i, item := range items {
		wishlist.Items[i] = *item
	}

	return ToWishlistDTO(wishlist), nil
}

func (s *WishlistQueryService) GetCustomerWishlists(ctx context.Context, customerID string) ([]*WishlistDTO, error) {
	wishlists, err := s.wishlistRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find customer wishlists: %w", err)
	}

	dtos := make([]*WishlistDTO, len(wishlists))
	for i, wishlist := range wishlists {
		items, err := s.wishlistItemRepo.FindByWishlistID(ctx, wishlist.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to find wishlist items: %w", err)
		}

		wishlist.Items = make([]domain.WishlistItem, len(items))
		for j, item := range items {
			wishlist.Items[j] = *item
		}

		dtos[i] = ToWishlistDTO(wishlist)
	}

	return dtos, nil
}

func (s *WishlistQueryService) GetDefaultWishlist(ctx context.Context, customerID string) (*WishlistDTO, error) {
	wishlist, err := s.wishlistRepo.FindDefaultByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find default wishlist: %w", err)
	}
	if wishlist == nil {
		return nil, domain.ErrWishlistNotFound
	}

	items, err := s.wishlistItemRepo.FindByWishlistID(ctx, wishlist.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find wishlist items: %w", err)
	}

	wishlist.Items = make([]domain.WishlistItem, len(items))
	for i, item := range items {
		wishlist.Items[i] = *item
	}

	return ToWishlistDTO(wishlist), nil
}

func (s *WishlistQueryService) GetPublicWishlist(ctx context.Context, id string) (*WishlistDTO, error) {
	wishlist, err := s.wishlistRepo.FindPublicByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find public wishlist: %w", err)
	}
	if wishlist == nil {
		return nil, domain.ErrWishlistNotFound
	}

	items, err := s.wishlistItemRepo.FindByWishlistID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find wishlist items: %w", err)
	}

	wishlist.Items = make([]domain.WishlistItem, len(items))
	for i, item := range items {
		wishlist.Items[i] = *item
	}

	return ToWishlistDTO(wishlist), nil
}

func (s *WishlistQueryService) GetWishlistItem(ctx context.Context, id string) (*WishlistItemDTO, error) {
	item, err := s.wishlistItemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find wishlist item: %w", err)
	}
	if item == nil {
		return nil, domain.ErrWishlistItemNotFound
	}

	dto := ToWishlistItemDTO(item)
	return &dto, nil
}

func (s *WishlistQueryService) GetWishlistItems(ctx context.Context, wishlistID string) ([]WishlistItemDTO, error) {
	items, err := s.wishlistItemRepo.FindByWishlistID(ctx, wishlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to find wishlist items: %w", err)
	}

	dtos := make([]WishlistItemDTO, len(items))
	for i, item := range items {
		dtos[i] = ToWishlistItemDTO(item)
	}

	return dtos, nil
}
