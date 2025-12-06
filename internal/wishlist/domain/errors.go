package domain

import "errors"

var (
	ErrCustomerIDRequired      = errors.New("customer ID is required")
	ErrWishlistIDRequired      = errors.New("wishlist ID is required")
	ErrProductIDRequired       = errors.New("product ID is required")
	ErrWishlistNameRequired    = errors.New("wishlist name is required")
	ErrInvalidQuantity         = errors.New("quantity must be greater than 0")
	ErrInvalidPriority         = errors.New("priority must be between 1 and 5")
	ErrWishlistNotFound        = errors.New("wishlist not found")
	ErrWishlistItemNotFound    = errors.New("wishlist item not found")
	ErrItemAlreadyInWishlist   = errors.New("item already exists in wishlist")
	ErrCannotDeleteDefault     = errors.New("cannot delete default wishlist")
)
