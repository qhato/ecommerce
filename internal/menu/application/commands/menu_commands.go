package commands

// CreateMenuCommand creates a new menu
type CreateMenuCommand struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

// UpdateMenuCommand updates an existing menu
type UpdateMenuCommand struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

// DeleteMenuCommand deletes a menu
type DeleteMenuCommand struct {
	ID int64 `json:"id"`
}

// CreateMenuItemCommand creates a new menu item
type CreateMenuItemCommand struct {
	MenuID      int64   `json:"menu_id"`
	ParentID    *int64  `json:"parent_id"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Target      string  `json:"target"`
	Icon        string  `json:"icon"`
	CSSClass    string  `json:"css_class"`
	SortOrder   int     `json:"sort_order"`
	Permissions *string `json:"permissions"`
}

// UpdateMenuItemCommand updates an existing menu item
type UpdateMenuItemCommand struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Target      string  `json:"target"`
	Icon        string  `json:"icon"`
	CSSClass    string  `json:"css_class"`
	SortOrder   int     `json:"sort_order"`
	Permissions *string `json:"permissions"`
}

// DeleteMenuItemCommand deletes a menu item
type DeleteMenuItemCommand struct {
	ID int64 `json:"id"`
}

// MoveMenuItemCommand moves a menu item to a new parent or position
type MoveMenuItemCommand struct {
	ID        int64  `json:"id"`
	ParentID  *int64 `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
}
