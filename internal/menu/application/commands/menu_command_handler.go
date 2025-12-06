package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/menu/domain"
)

// MenuCommandHandler handles menu-related commands
type MenuCommandHandler struct {
	menuRepo     domain.MenuRepository
	menuItemRepo domain.MenuItemRepository
}

// NewMenuCommandHandler creates a new menu command handler
func NewMenuCommandHandler(
	menuRepo domain.MenuRepository,
	menuItemRepo domain.MenuItemRepository,
) *MenuCommandHandler {
	return &MenuCommandHandler{
		menuRepo:     menuRepo,
		menuItemRepo: menuItemRepo,
	}
}

// HandleCreateMenu handles creating a new menu
func (h *MenuCommandHandler) HandleCreateMenu(ctx context.Context, cmd CreateMenuCommand) (*domain.Menu, error) {
	// Check if slug already exists
	exists, err := h.menuRepo.ExistsBySlug(ctx, cmd.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug: %w", err)
	}
	if exists {
		return nil, domain.ErrMenuSlugTaken
	}

	menu, err := domain.NewMenu(cmd.Name, cmd.Slug, cmd.Description, cmd.Location)
	if err != nil {
		return nil, err
	}

	if err := h.menuRepo.Create(ctx, menu); err != nil {
		return nil, fmt.Errorf("failed to create menu: %w", err)
	}

	return menu, nil
}

// HandleUpdateMenu handles updating a menu
func (h *MenuCommandHandler) HandleUpdateMenu(ctx context.Context, cmd UpdateMenuCommand) (*domain.Menu, error) {
	menu, err := h.menuRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu: %w", err)
	}
	if menu == nil {
		return nil, domain.ErrMenuNotFound
	}

	if err := menu.Update(cmd.Name, cmd.Description, cmd.Location); err != nil {
		return nil, err
	}

	if err := h.menuRepo.Update(ctx, menu); err != nil {
		return nil, fmt.Errorf("failed to update menu: %w", err)
	}

	return menu, nil
}

// HandleDeleteMenu handles deleting a menu
func (h *MenuCommandHandler) HandleDeleteMenu(ctx context.Context, cmd DeleteMenuCommand) error {
	menu, err := h.menuRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find menu: %w", err)
	}
	if menu == nil {
		return domain.ErrMenuNotFound
	}

	// Check if menu has items
	hasItems, err := h.menuItemRepo.HasItems(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to check menu items: %w", err)
	}
	if hasItems {
		return domain.ErrMenuHasItems
	}

	if err := h.menuRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete menu: %w", err)
	}

	return nil
}

// HandleCreateMenuItem handles creating a new menu item
func (h *MenuCommandHandler) HandleCreateMenuItem(ctx context.Context, cmd CreateMenuItemCommand) (*domain.MenuItem, error) {
	// Verify menu exists
	menu, err := h.menuRepo.FindByID(ctx, cmd.MenuID)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu: %w", err)
	}
	if menu == nil {
		return nil, domain.ErrMenuNotFound
	}

	// Verify parent exists if provided
	if cmd.ParentID != nil {
		parent, err := h.menuItemRepo.FindByID(ctx, *cmd.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to find parent item: %w", err)
		}
		if parent == nil {
			return nil, domain.ErrMenuItemNotFound
		}
	}

	menuItem, err := domain.NewMenuItem(cmd.MenuID, cmd.ParentID, cmd.Title, cmd.URL)
	if err != nil {
		return nil, err
	}

	menuItem.Target = cmd.Target
	menuItem.Icon = cmd.Icon
	menuItem.CSSClass = cmd.CSSClass
	menuItem.SortOrder = cmd.SortOrder
	menuItem.Permissions = cmd.Permissions

	if err := h.menuItemRepo.Create(ctx, menuItem); err != nil {
		return nil, fmt.Errorf("failed to create menu item: %w", err)
	}

	return menuItem, nil
}

// HandleUpdateMenuItem handles updating a menu item
func (h *MenuCommandHandler) HandleUpdateMenuItem(ctx context.Context, cmd UpdateMenuItemCommand) (*domain.MenuItem, error) {
	menuItem, err := h.menuItemRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu item: %w", err)
	}
	if menuItem == nil {
		return nil, domain.ErrMenuItemNotFound
	}

	if err := menuItem.Update(cmd.Title, cmd.URL, cmd.Target, cmd.Icon); err != nil {
		return nil, err
	}

	menuItem.CSSClass = cmd.CSSClass
	menuItem.SortOrder = cmd.SortOrder
	menuItem.Permissions = cmd.Permissions

	if err := h.menuItemRepo.Update(ctx, menuItem); err != nil {
		return nil, fmt.Errorf("failed to update menu item: %w", err)
	}

	return menuItem, nil
}

// HandleDeleteMenuItem handles deleting a menu item
func (h *MenuCommandHandler) HandleDeleteMenuItem(ctx context.Context, cmd DeleteMenuItemCommand) error {
	menuItem, err := h.menuItemRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to find menu item: %w", err)
	}
	if menuItem == nil {
		return domain.ErrMenuItemNotFound
	}

	// Check if item has children
	hasChildren, err := h.menuItemRepo.HasChildren(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to check children: %w", err)
	}
	if hasChildren {
		return domain.ErrMenuItemHasChildren
	}

	if err := h.menuItemRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete menu item: %w", err)
	}

	return nil
}

// HandleMoveMenuItem handles moving a menu item
func (h *MenuCommandHandler) HandleMoveMenuItem(ctx context.Context, cmd MoveMenuItemCommand) (*domain.MenuItem, error) {
	menuItem, err := h.menuItemRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu item: %w", err)
	}
	if menuItem == nil {
		return nil, domain.ErrMenuItemNotFound
	}

	// Verify parent exists if provided
	if cmd.ParentID != nil {
		parent, err := h.menuItemRepo.FindByID(ctx, *cmd.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to find parent item: %w", err)
		}
		if parent == nil {
			return nil, domain.ErrMenuItemNotFound
		}
	}

	menuItem.ParentID = cmd.ParentID
	menuItem.SortOrder = cmd.SortOrder

	if err := h.menuItemRepo.Update(ctx, menuItem); err != nil {
		return nil, fmt.Errorf("failed to move menu item: %w", err)
	}

	return menuItem, nil
}
