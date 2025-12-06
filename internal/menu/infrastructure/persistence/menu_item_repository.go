package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/menu/domain"
)

type PostgresMenuItemRepository struct {
	db *sql.DB
}

func NewPostgresMenuItemRepository(db *sql.DB) *PostgresMenuItemRepository {
	return &PostgresMenuItemRepository{db: db}
}

func (r *PostgresMenuItemRepository) Create(ctx context.Context, item *domain.MenuItem) error {
	query := `
		INSERT INTO blc_menu_item (menu_id, parent_id, title, url, target, icon, css_class, sort_order, is_active, permissions, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		item.MenuID,
		item.ParentID,
		item.Title,
		item.URL,
		item.Target,
		item.Icon,
		item.CSSClass,
		item.SortOrder,
		item.IsActive,
		item.Permissions,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(&item.ID)
}

func (r *PostgresMenuItemRepository) Update(ctx context.Context, item *domain.MenuItem) error {
	query := `
		UPDATE blc_menu_item
		SET parent_id = $1, title = $2, url = $3, target = $4, icon = $5, css_class = $6,
		    sort_order = $7, is_active = $8, permissions = $9, updated_at = $10
		WHERE id = $11`

	_, err := r.db.ExecContext(ctx, query,
		item.ParentID,
		item.Title,
		item.URL,
		item.Target,
		item.Icon,
		item.CSSClass,
		item.SortOrder,
		item.IsActive,
		item.Permissions,
		item.UpdatedAt,
		item.ID,
	)
	return err
}

func (r *PostgresMenuItemRepository) FindByID(ctx context.Context, id int64) (*domain.MenuItem, error) {
	query := `
		SELECT id, menu_id, parent_id, title, url, target, icon, css_class, sort_order, is_active, permissions, created_at, updated_at
		FROM blc_menu_item
		WHERE id = $1`

	item := &domain.MenuItem{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.MenuID,
		&item.ParentID,
		&item.Title,
		&item.URL,
		&item.Target,
		&item.Icon,
		&item.CSSClass,
		&item.SortOrder,
		&item.IsActive,
		&item.Permissions,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *PostgresMenuItemRepository) FindByMenuID(ctx context.Context, menuID int64) ([]*domain.MenuItem, error) {
	query := `
		SELECT id, menu_id, parent_id, title, url, target, icon, css_class, sort_order, is_active, permissions, created_at, updated_at
		FROM blc_menu_item
		WHERE menu_id = $1
		ORDER BY sort_order, title`

	return r.queryItems(ctx, query, menuID)
}

func (r *PostgresMenuItemRepository) FindByParentID(ctx context.Context, parentID int64) ([]*domain.MenuItem, error) {
	query := `
		SELECT id, menu_id, parent_id, title, url, target, icon, css_class, sort_order, is_active, permissions, created_at, updated_at
		FROM blc_menu_item
		WHERE parent_id = $1
		ORDER BY sort_order, title`

	return r.queryItems(ctx, query, parentID)
}

func (r *PostgresMenuItemRepository) FindHierarchy(ctx context.Context, menuID int64) ([]*domain.MenuItem, error) {
	query := `
		SELECT id, menu_id, parent_id, title, url, target, icon, css_class, sort_order, is_active, permissions, created_at, updated_at
		FROM blc_menu_item
		WHERE menu_id = $1
		ORDER BY parent_id NULLS FIRST, sort_order, title`

	return r.queryItems(ctx, query, menuID)
}

func (r *PostgresMenuItemRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_menu_item WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresMenuItemRepository) HasItems(ctx context.Context, menuID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blc_menu_item WHERE menu_id = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, menuID).Scan(&exists)
	return exists, err
}

func (r *PostgresMenuItemRepository) HasChildren(ctx context.Context, parentID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blc_menu_item WHERE parent_id = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, parentID).Scan(&exists)
	return exists, err
}

func (r *PostgresMenuItemRepository) BuildTree(ctx context.Context, menuID int64) ([]domain.MenuItem, error) {
	// Get all items for the menu
	items, err := r.FindHierarchy(ctx, menuID)
	if err != nil {
		return nil, err
	}

	// Build a map for quick lookup
	itemMap := make(map[int64]*domain.MenuItem)
	for _, item := range items {
		itemCopy := *item
		itemCopy.Children = []domain.MenuItem{}
		itemMap[item.ID] = &itemCopy
	}

	// Build the tree structure
	var rootItems []domain.MenuItem
	for _, item := range items {
		if item.ParentID == nil {
			rootItems = append(rootItems, *itemMap[item.ID])
		} else {
			parent, exists := itemMap[*item.ParentID]
			if exists {
				parent.Children = append(parent.Children, *itemMap[item.ID])
			}
		}
	}

	return rootItems, nil
}

func (r *PostgresMenuItemRepository) queryItems(ctx context.Context, query string, args ...interface{}) ([]*domain.MenuItem, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.MenuItem
	for rows.Next() {
		item := &domain.MenuItem{}
		if err := rows.Scan(
			&item.ID,
			&item.MenuID,
			&item.ParentID,
			&item.Title,
			&item.URL,
			&item.Target,
			&item.Icon,
			&item.CSSClass,
			&item.SortOrder,
			&item.IsActive,
			&item.Permissions,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
