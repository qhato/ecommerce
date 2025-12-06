package domain

import "context"

type WishlistRepository interface {
	Create(ctx context.Context, wishlist *Wishlist) error
	Update(ctx context.Context, wishlist *Wishlist) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*Wishlist, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*Wishlist, error)
	FindDefaultByCustomerID(ctx context.Context, customerID string) (*Wishlist, error)
	FindPublicByID(ctx context.Context, id string) (*Wishlist, error)
}

type WishlistItemRepository interface {
	Create(ctx context.Context, item *WishlistItem) error
	Update(ctx context.Context, item *WishlistItem) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*WishlistItem, error)
	FindByWishlistID(ctx context.Context, wishlistID string) ([]*WishlistItem, error)
	ExistsByWishlistAndProduct(ctx context.Context, wishlistID, productID string) (bool, error)
}
