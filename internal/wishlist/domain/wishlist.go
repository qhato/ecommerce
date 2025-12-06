package domain

import (
	"time"

	"github.com/google/uuid"
)

// Wishlist represents a customer's wishlist
type Wishlist struct {
	ID         string
	CustomerID string
	Name       string
	IsDefault  bool
	IsPublic   bool
	Items      []WishlistItem
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// WishlistItem represents an item in a wishlist
type WishlistItem struct {
	ID          string
	WishlistID  string
	ProductID   string
	SKUID       *string
	Quantity    int
	Priority    int // 1-5, 5 being highest
	Notes       string
	AddedAt     time.Time
	UpdatedAt   time.Time
}

// NewWishlist creates a new wishlist
func NewWishlist(customerID, name string, isDefault, isPublic bool) (*Wishlist, error) {
	if customerID == "" {
		return nil, ErrCustomerIDRequired
	}
	if name == "" {
		name = "My Wishlist"
	}

	now := time.Now()
	return &Wishlist{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Name:       name,
		IsDefault:  isDefault,
		IsPublic:   isPublic,
		Items:      make([]WishlistItem, 0),
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// NewWishlistItem creates a new wishlist item
func NewWishlistItem(wishlistID, productID string, quantity int) (*WishlistItem, error) {
	if wishlistID == "" {
		return nil, ErrWishlistIDRequired
	}
	if productID == "" {
		return nil, ErrProductIDRequired
	}
	if quantity <= 0 {
		quantity = 1
	}

	now := time.Now()
	return &WishlistItem{
		ID:         uuid.New().String(),
		WishlistID: wishlistID,
		ProductID:  productID,
		Quantity:   quantity,
		Priority:   3,
		AddedAt:    now,
		UpdatedAt:  now,
	}, nil
}

// AddItem adds an item to the wishlist
func (w *Wishlist) AddItem(item WishlistItem) {
	w.Items = append(w.Items, item)
	w.UpdatedAt = time.Now()
}

// UpdateQuantity updates item quantity
func (wi *WishlistItem) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	wi.Quantity = quantity
	wi.UpdatedAt = time.Now()
	return nil
}

// SetPriority sets item priority
func (wi *WishlistItem) SetPriority(priority int) error {
	if priority < 1 || priority > 5 {
		return ErrInvalidPriority
	}
	wi.Priority = priority
	wi.UpdatedAt = time.Now()
	return nil
}

// SetPublic makes wishlist public
func (w *Wishlist) SetPublic(isPublic bool) {
	w.IsPublic = isPublic
	w.UpdatedAt = time.Now()
}

// Rename renames the wishlist
func (w *Wishlist) Rename(name string) error {
	if name == "" {
		return ErrWishlistNameRequired
	}
	w.Name = name
	w.UpdatedAt = time.Now()
	return nil
}
