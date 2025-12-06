package domain

import "errors"

var (
	ErrStoreNotFound         = errors.New("store not found")
	ErrStoreCodeTaken        = errors.New("store code already taken")
	ErrInvalidStoreCode      = errors.New("invalid store code")
	ErrInsufficientInventory = errors.New("insufficient inventory")
	ErrInventoryNotFound     = errors.New("inventory not found")
	ErrStoreInactive         = errors.New("store is inactive")
	ErrStoreClosed           = errors.New("store is closed")
)
