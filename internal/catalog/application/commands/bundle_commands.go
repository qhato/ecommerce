package commands

import "github.com/shopspring/decimal"

type CreateProductBundleCommand struct {
	Name        string
	Description string
	BundlePrice decimal.Decimal
	Items       []BundleItemInput
}

type UpdateProductBundleCommand struct {
	ID          int64
	Name        string
	Description string
	BundlePrice decimal.Decimal
}

type AddBundleItemCommand struct {
	BundleID  int64
	ProductID *int64
	SKUID     *int64
	Quantity  int
	SortOrder int
}

type ActivateBundleCommand struct {
	ID int64
}

type DeactivateBundleCommand struct {
	ID int64
}

type DeleteBundleCommand struct {
	ID int64
}

type BundleItemInput struct {
	ProductID *int64
	SKUID     *int64
	Quantity  int
	SortOrder int
}
