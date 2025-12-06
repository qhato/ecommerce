package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/menu/domain"
)

// MenuDTO represents a menu for API responses
type MenuDTO struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Description string        `json:"description"`
	Location    string        `json:"location"`
	IsActive    bool          `json:"is_active"`
	Items       []MenuItemDTO `json:"items,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// MenuItemDTO represents a menu item for API responses
type MenuItemDTO struct {
	ID          int64         `json:"id"`
	MenuID      int64         `json:"menu_id"`
	ParentID    *int64        `json:"parent_id,omitempty"`
	Title       string        `json:"title"`
	URL         string        `json:"url"`
	Target      string        `json:"target,omitempty"`
	Icon        string        `json:"icon,omitempty"`
	CSSClass    string        `json:"css_class,omitempty"`
	SortOrder   int           `json:"sort_order"`
	IsActive    bool          `json:"is_active"`
	Permissions *string       `json:"permissions,omitempty"`
	Children    []MenuItemDTO `json:"children,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// ToMenuDTO converts domain Menu to MenuDTO
func ToMenuDTO(menu *domain.Menu) *MenuDTO {
	dto := &MenuDTO{
		ID:          menu.ID,
		Name:        menu.Name,
		Slug:        menu.Slug,
		Description: menu.Description,
		Location:    menu.Location,
		IsActive:    menu.IsActive,
		CreatedAt:   menu.CreatedAt,
		UpdatedAt:   menu.UpdatedAt,
	}

	if len(menu.Items) > 0 {
		dto.Items = make([]MenuItemDTO, len(menu.Items))
		for i, item := range menu.Items {
			dto.Items[i] = *ToMenuItemDTO(&item)
		}
	}

	return dto
}

// ToMenuItemDTO converts domain MenuItem to MenuItemDTO
func ToMenuItemDTO(item *domain.MenuItem) *MenuItemDTO {
	dto := &MenuItemDTO{
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
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	if len(item.Children) > 0 {
		dto.Children = make([]MenuItemDTO, len(item.Children))
		for i, child := range item.Children {
			dto.Children[i] = *ToMenuItemDTO(&child)
		}
	}

	return dto
}

// MenuItemTreeDTO represents a hierarchical menu item tree
type MenuItemTreeDTO struct {
	ID          int64             `json:"id"`
	MenuID      int64             `json:"menu_id"`
	ParentID    *int64            `json:"parent_id,omitempty"`
	Title       string            `json:"title"`
	URL         string            `json:"url"`
	Target      string            `json:"target,omitempty"`
	Icon        string            `json:"icon,omitempty"`
	CSSClass    string            `json:"css_class,omitempty"`
	SortOrder   int               `json:"sort_order"`
	IsActive    bool              `json:"is_active"`
	Permissions *string           `json:"permissions,omitempty"`
	Children    []MenuItemTreeDTO `json:"children,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
