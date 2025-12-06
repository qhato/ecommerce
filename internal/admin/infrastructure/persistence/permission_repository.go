package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/admin/domain"
)

type PostgresPermissionRepository struct {
	db *sql.DB
}

func NewPostgresPermissionRepository(db *sql.DB) *PostgresPermissionRepository {
	return &PostgresPermissionRepository{db: db}
}

func (r *PostgresPermissionRepository) Create(ctx context.Context, permission *domain.Permission) error {
	query := `INSERT INTO blc_admin_permission (name, description, resource, action, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRowContext(ctx, query,
		permission.Name, permission.Description, permission.Resource, permission.Action,
		permission.IsActive, permission.CreatedAt, permission.UpdatedAt).Scan(&permission.ID)
}

func (r *PostgresPermissionRepository) Update(ctx context.Context, permission *domain.Permission) error {
	query := `UPDATE blc_admin_permission SET name = $1, description = $2, resource = $3,
		action = $4, is_active = $5, updated_at = $6 WHERE id = $7`
	_, err := r.db.ExecContext(ctx, query,
		permission.Name, permission.Description, permission.Resource, permission.Action,
		permission.IsActive, permission.UpdatedAt, permission.ID)
	return err
}

func (r *PostgresPermissionRepository) FindByID(ctx context.Context, id int64) (*domain.Permission, error) {
	query := `SELECT id, name, description, resource, action, is_active, created_at, updated_at
		FROM blc_admin_permission WHERE id = $1`
	permission := &domain.Permission{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&permission.ID, &permission.Name, &permission.Description, &permission.Resource,
		&permission.Action, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return permission, err
}

func (r *PostgresPermissionRepository) FindByName(ctx context.Context, name string) (*domain.Permission, error) {
	query := `SELECT id, name, description, resource, action, is_active, created_at, updated_at
		FROM blc_admin_permission WHERE name = $1`
	permission := &domain.Permission{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&permission.ID, &permission.Name, &permission.Description, &permission.Resource,
		&permission.Action, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return permission, err
}

func (r *PostgresPermissionRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.Permission, error) {
	query := `SELECT id, name, description, resource, action, is_active, created_at, updated_at
		FROM blc_admin_permission`
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY resource ASC, action ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]*domain.Permission, 0)
	for rows.Next() {
		permission := &domain.Permission{}
		if err := rows.Scan(
			&permission.ID, &permission.Name, &permission.Description, &permission.Resource,
			&permission.Action, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt,
		); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (r *PostgresPermissionRepository) FindByResource(ctx context.Context, resource domain.PermissionResource, activeOnly bool) ([]*domain.Permission, error) {
	query := `SELECT id, name, description, resource, action, is_active, created_at, updated_at
		FROM blc_admin_permission WHERE resource = $1`
	if activeOnly {
		query += " AND is_active = true"
	}
	query += " ORDER BY action ASC"

	rows, err := r.db.QueryContext(ctx, query, resource)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]*domain.Permission, 0)
	for rows.Next() {
		permission := &domain.Permission{}
		if err := rows.Scan(
			&permission.ID, &permission.Name, &permission.Description, &permission.Resource,
			&permission.Action, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt,
		); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (r *PostgresPermissionRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM blc_admin_permission WHERE id = $1`, id)
	return err
}

func (r *PostgresPermissionRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM blc_admin_permission WHERE name = $1)`, name).Scan(&exists)
	return exists, err
}

func (r *PostgresPermissionRepository) FindByResourceAndAction(ctx context.Context, resource domain.PermissionResource, action domain.PermissionAction) (*domain.Permission, error) {
	query := `SELECT id, name, description, resource, action, is_active, created_at, updated_at
		FROM blc_admin_permission WHERE resource = $1 AND action = $2`
	permission := &domain.Permission{}
	err := r.db.QueryRowContext(ctx, query, resource, action).Scan(
		&permission.ID, &permission.Name, &permission.Description, &permission.Resource,
		&permission.Action, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return permission, err
}

func (r *PostgresPermissionRepository) IsPermissionInUse(ctx context.Context, permissionID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM blc_admin_role_permission WHERE permission_id = $1)`, permissionID).Scan(&exists)
	return exists, err
}
