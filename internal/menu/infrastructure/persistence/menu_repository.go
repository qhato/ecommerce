package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/menu/domain"
)

type PostgresMenuRepository struct {
	db *sql.DB
}

func NewPostgresMenuRepository(db *sql.DB) *PostgresMenuRepository {
	return &PostgresMenuRepository{db: db}
}

func (r *PostgresMenuRepository) Create(ctx context.Context, menu *domain.Menu) error {
	query := `
		INSERT INTO blc_menu (name, slug, type, description, location, is_active, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		menu.Name,
		menu.Slug,
		menu.Type,
		menu.Description,
		menu.Location,
		menu.IsActive,
		menu.SortOrder,
		menu.CreatedAt,
		menu.UpdatedAt,
	).Scan(&menu.ID)
}

func (r *PostgresMenuRepository) Update(ctx context.Context, menu *domain.Menu) error {
	query := `
		UPDATE blc_menu
		SET name = $1, description = $2, location = $3, type = $4, is_active = $5, sort_order = $6, updated_at = $7
		WHERE id = $8`

	_, err := r.db.ExecContext(ctx, query,
		menu.Name,
		menu.Description,
		menu.Location,
		menu.Type,
		menu.IsActive,
		menu.SortOrder,
		menu.UpdatedAt,
		menu.ID,
	)
	return err
}

func (r *PostgresMenuRepository) FindByID(ctx context.Context, id int64) (*domain.Menu, error) {
	query := `
		SELECT id, name, slug, type, description, location, is_active, sort_order, created_at, updated_at
		FROM blc_menu
		WHERE id = $1`

	menu := &domain.Menu{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&menu.ID,
		&menu.Name,
		&menu.Slug,
		&menu.Type,
		&menu.Description,
		&menu.Location,
		&menu.IsActive,
		&menu.SortOrder,
		&menu.CreatedAt,
		&menu.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (r *PostgresMenuRepository) FindBySlug(ctx context.Context, slug string) (*domain.Menu, error) {
	query := `
		SELECT id, name, slug, type, description, location, is_active, sort_order, created_at, updated_at
		FROM blc_menu
		WHERE slug = $1`

	menu := &domain.Menu{}
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&menu.ID,
		&menu.Name,
		&menu.Slug,
		&menu.Type,
		&menu.Description,
		&menu.Location,
		&menu.IsActive,
		&menu.SortOrder,
		&menu.CreatedAt,
		&menu.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (r *PostgresMenuRepository) FindByLocation(ctx context.Context, location string) (*domain.Menu, error) {
	query := `
		SELECT id, name, slug, type, description, location, is_active, sort_order, created_at, updated_at
		FROM blc_menu
		WHERE location = $1 AND is_active = true`

	menu := &domain.Menu{}
	err := r.db.QueryRowContext(ctx, query, location).Scan(
		&menu.ID,
		&menu.Name,
		&menu.Slug,
		&menu.Type,
		&menu.Description,
		&menu.Location,
		&menu.IsActive,
		&menu.SortOrder,
		&menu.CreatedAt,
		&menu.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (r *PostgresMenuRepository) FindByType(ctx context.Context, menuType domain.MenuType) ([]*domain.Menu, error) {
	query := `
		SELECT id, name, slug, type, description, location, is_active, sort_order, created_at, updated_at
		FROM blc_menu
		WHERE type = $1
		ORDER BY sort_order, name`

	rows, err := r.db.QueryContext(ctx, query, menuType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []*domain.Menu
	for rows.Next() {
		menu := &domain.Menu{}
		if err := rows.Scan(
			&menu.ID,
			&menu.Name,
			&menu.Slug,
			&menu.Type,
			&menu.Description,
			&menu.Location,
			&menu.IsActive,
			&menu.SortOrder,
			&menu.CreatedAt,
			&menu.UpdatedAt,
		); err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}

	return menus, rows.Err()
}

func (r *PostgresMenuRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.Menu, error) {
	query := `
		SELECT id, name, slug, type, description, location, is_active, sort_order, created_at, updated_at
		FROM blc_menu`

	if activeOnly {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY sort_order, name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []*domain.Menu
	for rows.Next() {
		menu := &domain.Menu{}
		if err := rows.Scan(
			&menu.ID,
			&menu.Name,
			&menu.Slug,
			&menu.Type,
			&menu.Description,
			&menu.Location,
			&menu.IsActive,
			&menu.SortOrder,
			&menu.CreatedAt,
			&menu.UpdatedAt,
		); err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}

	return menus, rows.Err()
}

func (r *PostgresMenuRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_menu WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresMenuRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blc_menu WHERE slug = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, slug).Scan(&exists)
	return exists, err
}
