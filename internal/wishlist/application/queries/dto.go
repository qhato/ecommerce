package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/wishlist/domain"
)

type WishlistDTO struct {
	ID         string            `json:"id"`
	CustomerID string            `json:"customer_id"`
	Name       string            `json:"name"`
	IsDefault  bool              `json:"is_default"`
	IsPublic   bool              `json:"is_public"`
	Items      []WishlistItemDTO `json:"items"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type WishlistItemDTO struct {
	ID         string     `json:"id"`
	WishlistID string     `json:"wishlist_id"`
	ProductID  string     `json:"product_id"`
	SKUID      *string    `json:"sku_id,omitempty"`
	Quantity   int        `json:"quantity"`
	Priority   int        `json:"priority"`
	Notes      string     `json:"notes"`
	AddedAt    time.Time  `json:"added_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func ToWishlistDTO(wishlist *domain.Wishlist) *WishlistDTO {
	items := make([]WishlistItemDTO, len(wishlist.Items))
	for i, item := range wishlist.Items {
		items[i] = ToWishlistItemDTO(&item)
	}

	return &WishlistDTO{
		ID:         wishlist.ID,
		CustomerID: wishlist.CustomerID,
		Name:       wishlist.Name,
		IsDefault:  wishlist.IsDefault,
		IsPublic:   wishlist.IsPublic,
		Items:      items,
		CreatedAt:  wishlist.CreatedAt,
		UpdatedAt:  wishlist.UpdatedAt,
	}
}

func ToWishlistItemDTO(item *domain.WishlistItem) WishlistItemDTO {
	return WishlistItemDTO{
		ID:         item.ID,
		WishlistID: item.WishlistID,
		ProductID:  item.ProductID,
		SKUID:      item.SKUID,
		Quantity:   item.Quantity,
		Priority:   item.Priority,
		Notes:      item.Notes,
		AddedAt:    item.AddedAt,
		UpdatedAt:  item.UpdatedAt,
	}
}
