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
	Slug        string
	Type        MenuType
	Description string
	Location    string
	IsActive    bool
	SortOrder   int
	Items       []MenuItem
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
	Permissions *string
	Children    []MenuItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewMenu creates a new menu
func NewMenu(name, slug, description, location string) (*Menu, error) {
	if name == "" {
		return nil, ErrMenuNameRequired
	}
	if slug == "" {
		return nil, ErrMenuSlugRequired
	}

	now := time.Now()
	return &Menu{
		Name:        name,
		Slug:        slug,
		Description: description,
		Location:    location,
		IsActive:    true,
		Items:       make([]MenuItem, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// NewMenuItem creates a new menu item
func NewMenuItem(menuID int64, parentID *int64, title, url string) (*MenuItem, error) {
	if title == "" {
		return nil, ErrMenuItemTitleRequired
	}

	now := time.Now()
	return &MenuItem{
		MenuID:    menuID,
		ParentID:  parentID,
		Title:     title,
		URL:       url,
		Target:    "_self",
		IsActive:  true,
		Children:  make([]MenuItem, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update updates menu information
func (m *Menu) Update(name, description, location string) error {
	if name == "" {
		return ErrMenuNameRequired
	}
	m.Name = name
	m.Description = description
	m.Location = location
	m.UpdatedAt = time.Now()
	return nil
}

// Update updates menu item information
func (mi *MenuItem) Update(title, url, target, icon string) error {
	if title == "" {
		return ErrMenuItemTitleRequired
	}
	mi.Title = title
	mi.URL = url
	mi.Target = target
	mi.Icon = icon
	mi.UpdatedAt = time.Now()
	return nil
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
