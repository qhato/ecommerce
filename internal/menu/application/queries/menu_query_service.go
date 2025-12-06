package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/menu/domain"
)

// MenuQueryService handles menu-related queries
type MenuQueryService struct {
	menuRepo     domain.MenuRepository
	menuItemRepo domain.MenuItemRepository
}

// NewMenuQueryService creates a new menu query service
func NewMenuQueryService(
	menuRepo domain.MenuRepository,
	menuItemRepo domain.MenuItemRepository,
) *MenuQueryService {
	return &MenuQueryService{
		menuRepo:     menuRepo,
		menuItemRepo: menuItemRepo,
	}
}

// GetMenu retrieves a menu by ID
func (s *MenuQueryService) GetMenu(ctx context.Context, id int64) (*MenuDTO, error) {
	menu, err := s.menuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu: %w", err)
	}
	if menu == nil {
		return nil, domain.ErrMenuNotFound
	}

	// Load menu items
	items, err := s.menuItemRepo.FindByMenuID(ctx, menu.ID)
	if err == nil {
		menu.Items = make([]domain.MenuItem, len(items))
		for i, item := range items {
			menu.Items[i] = *item
		}
	}

	return ToMenuDTO(menu), nil
}

// GetMenuBySlug retrieves a menu by slug
func (s *MenuQueryService) GetMenuBySlug(ctx context.Context, slug string) (*MenuDTO, error) {
	menu, err := s.menuRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu: %w", err)
	}
	if menu == nil {
		return nil, domain.ErrMenuNotFound
	}

	// Load menu items with hierarchy
	items, err := s.menuItemRepo.FindHierarchy(ctx, menu.ID)
	if err == nil {
		menu.Items = make([]domain.MenuItem, len(items))
		for i, item := range items {
			menu.Items[i] = *item
		}
	}

	return ToMenuDTO(menu), nil
}

// GetAllMenus retrieves all menus
func (s *MenuQueryService) GetAllMenus(ctx context.Context, activeOnly bool) ([]*MenuDTO, error) {
	menus, err := s.menuRepo.FindAll(ctx, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find menus: %w", err)
	}

	dtos := make([]*MenuDTO, len(menus))
	for i, menu := range menus {
		dtos[i] = ToMenuDTO(menu)
	}

	return dtos, nil
}

// GetMenuByLocation retrieves a menu by location
func (s *MenuQueryService) GetMenuByLocation(ctx context.Context, location string) (*MenuDTO, error) {
	menu, err := s.menuRepo.FindByLocation(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu: %w", err)
	}
	if menu == nil {
		return nil, domain.ErrMenuNotFound
	}

	// Load menu items with hierarchy
	items, err := s.menuItemRepo.FindHierarchy(ctx, menu.ID)
	if err == nil {
		menu.Items = make([]domain.MenuItem, len(items))
		for i, item := range items {
			menu.Items[i] = *item
		}
	}

	return ToMenuDTO(menu), nil
}

// GetMenuItem retrieves a menu item by ID
func (s *MenuQueryService) GetMenuItem(ctx context.Context, id int64) (*MenuItemDTO, error) {
	item, err := s.menuItemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu item: %w", err)
	}
	if item == nil {
		return nil, domain.ErrMenuItemNotFound
	}

	return ToMenuItemDTO(item), nil
}

// GetMenusByType retrieves menus by type
func (s *MenuQueryService) GetMenusByType(ctx context.Context, menuType string) ([]*MenuDTO, error) {
	menus, err := s.menuRepo.FindByType(ctx, domain.MenuType(menuType))
	if err != nil {
		return nil, fmt.Errorf("failed to find menus by type: %w", err)
	}

	dtos := make([]*MenuDTO, len(menus))
	for i, menu := range menus {
		dtos[i] = ToMenuDTO(menu)
	}

	return dtos, nil
}

// GetMenuItems retrieves all items for a menu
func (s *MenuQueryService) GetMenuItems(ctx context.Context, menuID int64) ([]*MenuItemDTO, error) {
	items, err := s.menuItemRepo.FindByMenuID(ctx, menuID)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu items: %w", err)
	}

	dtos := make([]*MenuItemDTO, len(items))
	for i, item := range items {
		dtos[i] = ToMenuItemDTO(item)
	}

	return dtos, nil
}

// GetMenuTree retrieves menu items as a hierarchical tree
func (s *MenuQueryService) GetMenuTree(ctx context.Context, menuID int64) ([]MenuItemTreeDTO, error) {
	items, err := s.menuItemRepo.BuildTree(ctx, menuID)
	if err != nil {
		return nil, fmt.Errorf("failed to build menu tree: %w", err)
	}

	dtos := make([]MenuItemTreeDTO, len(items))
	for i, item := range items {
		dtos[i] = ToMenuItemTreeDTO(&item)
	}

	return dtos, nil
}

// ToMenuItemTreeDTO converts a MenuItem with children to DTO
func ToMenuItemTreeDTO(item *domain.MenuItem) MenuItemTreeDTO {
	children := make([]MenuItemTreeDTO, len(item.Children))
	for i, child := range item.Children {
		childCopy := child
		children[i] = ToMenuItemTreeDTO(&childCopy)
	}

	return MenuItemTreeDTO{
		ID:          item.ID,
		MenuID:      item.MenuID,
		ParentID:    item.ParentID,
		Title:       item.Title,
		URL:         item.URL,
		Target:      item.Target,
		Icon:        item.Icon,
		CSSClass:    item.CSSClass,
		SortOrder:   item.SortOrder,
		IsActive:    item.IsActive,
		Permissions: item.Permissions,
		Children:    children,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
