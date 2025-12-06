package commands

type CreateWishlistCommand struct {
	CustomerID string
	Name       string
	IsDefault  bool
	IsPublic   bool
}

type UpdateWishlistCommand struct {
	ID       string
	Name     string
	IsPublic bool
}

type DeleteWishlistCommand struct {
	ID string
}

type SetDefaultWishlistCommand struct {
	ID         string
	CustomerID string
}

type AddItemCommand struct {
	WishlistID string
	ProductID  string
	SKUID      *string
	Quantity   int
	Priority   int
	Notes      string
}

type UpdateItemCommand struct {
	ID       string
	Quantity int
	Priority int
	Notes    string
}

type RemoveItemCommand struct {
	ID string
}

type MoveItemCommand struct {
	ItemID          string
	TargetWishlistID string
}
