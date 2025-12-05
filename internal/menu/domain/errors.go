package domain

import "errors"

var (
	ErrMenuNotFound           = errors.New("menu not found")
	ErrMenuNameRequired       = errors.New("menu name is required")
	ErrMenuItemNotFound       = errors.New("menu item not found")
	ErrMenuItemTitleRequired  = errors.New("menu item title is required")
	ErrCircularReference      = errors.New("circular reference detected in menu hierarchy")
)
