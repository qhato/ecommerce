package domain

import "context"

// MenuRepository defines the interface for menu persistence
type MenuRepository interface {
	Create(ctx context.Context, menu *Menu) error
	Update(ctx context.Context, menu *Menu) error
	FindByID(ctx context.Context, id int64) (*Menu, error)
	FindBySlug(ctx context.Context, slug string) (*Menu, error)
	FindByLocation(ctx context.Context, location string) (*Menu, error)
	FindByType(ctx context.Context, menuType MenuType) ([]*Menu, error)
	FindAll(ctx context.Context, activeOnly bool) ([]*Menu, error)
	Delete(ctx context.Context, id int64) error
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
}

// MenuItemRepository defines the interface for menu item persistence
type MenuItemRepository interface {
	Create(ctx context.Context, item *MenuItem) error
	Update(ctx context.Context, item *MenuItem) error
	FindByID(ctx context.Context, id int64) (*MenuItem, error)
	FindByMenuID(ctx context.Context, menuID int64) ([]*MenuItem, error)
	FindByParentID(ctx context.Context, parentID int64) ([]*MenuItem, error)
	FindHierarchy(ctx context.Context, menuID int64) ([]*MenuItem, error)
	Delete(ctx context.Context, id int64) error
	HasItems(ctx context.Context, menuID int64) (bool, error)
	HasChildren(ctx context.Context, parentID int64) (bool, error)
	BuildTree(ctx context.Context, menuID int64) ([]MenuItem, error)
}
