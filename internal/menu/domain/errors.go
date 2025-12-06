package domain

import "errors"

var (
	ErrMenuNotFound          = errors.New("menu not found")
	ErrMenuNameRequired      = errors.New("menu name is required")
	ErrMenuSlugRequired      = errors.New("menu slug is required")
	ErrMenuSlugTaken         = errors.New("menu slug is already taken")
	ErrMenuHasItems          = errors.New("menu has items")
	ErrMenuItemNotFound      = errors.New("menu item not found")
	ErrMenuItemTitleRequired = errors.New("menu item title is required")
	ErrMenuItemHasChildren   = errors.New("menu item has children")
	ErrCircularReference     = errors.New("circular reference detected in menu hierarchy")
)
