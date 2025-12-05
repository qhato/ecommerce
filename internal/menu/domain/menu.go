package domain

import "time"

// MenuType represents the type of menu
type MenuType string

const (
	MenuTypeHeader MenuType = "HEADER"
	MenuTypeFooter MenuType = "FOOTER"
	MenuTypeSidebar MenuType = "SIDEBAR"
	MenuTypeMobile MenuType = "MOBILE"
)

// Menu represents a navigation menu
type Menu struct {
	ID          int64
	Name        string
	Type        MenuType
	Description string
	IsActive    bool
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// MenuItem represents an item in a menu
type MenuItem struct {
	ID          int64
	MenuID      int64
	ParentID    *int64
	Title       string
	URL         string
	Target      string // _self, _blank
	Icon        string
	CSSClass    string
	SortOrder   int
	IsActive    bool
	Children    []MenuItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewMenu creates a new menu
func NewMenu(name string, menuType MenuType) (*Menu, error) {
	if name == "" {
		return nil, ErrMenuNameRequired
	}

	now := time.Now()
	return &Menu{
		Name:      name,
		Type:      menuType,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// NewMenuItem creates a new menu item
func NewMenuItem(menuID int64, title, url string) (*MenuItem, error) {
	if title == "" {
		return nil, ErrMenuItemTitleRequired
	}

	now := time.Now()
	return &MenuItem{
		MenuID:    menuID,
		Title:     title,
		URL:       url,
		Target:    "_self",
		IsActive:  true,
		Children:  make([]MenuItem, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Activate activates the menu
func (m *Menu) Activate() {
	m.IsActive = true
	m.UpdatedAt = time.Now()
}

// Deactivate deactivates the menu
func (m *Menu) Deactivate() {
	m.IsActive = false
	m.UpdatedAt = time.Now()
}

// AddChild adds a child menu item
func (mi *MenuItem) AddChild(child MenuItem) {
	mi.Children = append(mi.Children, child)
}
