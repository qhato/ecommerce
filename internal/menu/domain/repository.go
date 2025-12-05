package domain

import "context"

// MenuRepository defines the interface for menu persistence
type MenuRepository interface {
	Create(ctx context.Context, menu *Menu) error
	Update(ctx context.Context, menu *Menu) error
	FindByID(ctx context.Context, id int64) (*Menu, error)
	FindByType(ctx context.Context, menuType MenuType) ([]*Menu, error)
	FindAll(ctx context.Context, activeOnly bool) ([]*Menu, error)
	Delete(ctx context.Context, id int64) error
}

// MenuItemRepository defines the interface for menu item persistence
type MenuItemRepository interface {
	Create(ctx context.Context, item *MenuItem) error
	Update(ctx context.Context, item *MenuItem) error
	FindByID(ctx context.Context, id int64) (*MenuItem, error)
	FindByMenuID(ctx context.Context, menuID int64) ([]*MenuItem, error)
	FindByParentID(ctx context.Context, parentID int64) ([]*MenuItem, error)
	Delete(ctx context.Context, id int64) error
	BuildTree(ctx context.Context, menuID int64) ([]MenuItem, error)
}
