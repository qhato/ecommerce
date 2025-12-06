package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/admin/domain"
)

type PostgresRoleRepository struct {
	db *sql.DB
}

func NewPostgresRoleRepository(db *sql.DB) *PostgresRoleRepository {
	return &PostgresRoleRepository{db: db}
}

func (r *PostgresRoleRepository) Create(ctx context.Context, role *domain.Role) error {
	query := `INSERT INTO blc_admin_role (name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, role.Name, role.Description, role.IsActive, role.CreatedAt, role.UpdatedAt).Scan(&role.ID)
}

func (r *PostgresRoleRepository) Update(ctx context.Context, role *domain.Role) error {
	query := `UPDATE blc_admin_role SET name = $1, description = $2, is_active = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, role.Name, role.Description, role.IsActive, role.UpdatedAt, role.ID)
	return err
}

func (r *PostgresRoleRepository) FindByID(ctx context.Context, id int64) (*domain.Role, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM blc_admin_role WHERE id = $1`
	role := &domain.Role{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&role.ID, &role.Name, &role.Description, &role.IsActive, &role.CreatedAt, &role.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return role, err
}

func (r *PostgresRoleRepository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM blc_admin_role WHERE name = $1`
	role := &domain.Role{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(&role.ID, &role.Name, &role.Description, &role.IsActive, &role.CreatedAt, &role.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return role, err
}

func (r *PostgresRoleRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.Role, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM blc_admin_role`
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY name ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]*domain.Role, 0)
	for rows.Next() {
		role := &domain.Role{}
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.IsActive, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *PostgresRoleRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM blc_admin_role WHERE id = $1`, id)
	return err
}

func (r *PostgresRoleRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM blc_admin_role WHERE name = $1)`, name).Scan(&exists)
	return exists, err
}

func (r *PostgresRoleRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]domain.Permission, error) {
	query := `SELECT p.id, p.name, p.description, p.resource, p.action, p.is_active, p.created_at, p.updated_at
		FROM blc_admin_permission p
		INNER JOIN blc_admin_role_permission rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1 ORDER BY p.name ASC`

	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := make([]domain.Permission, 0)
	for rows.Next() {
		perm := domain.Permission{}
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Description, &perm.Resource, &perm.Action, &perm.IsActive, &perm.CreatedAt, &perm.UpdatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	return perms, nil
}

func (r *PostgresRoleRepository) GrantPermission(ctx context.Context, roleID, permissionID int64) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO blc_admin_role_permission (role_id, permission_id, granted_at) VALUES ($1, $2, NOW()) ON CONFLICT DO NOTHING`, roleID, permissionID)
	return err
}

func (r *PostgresRoleRepository) RevokePermission(ctx context.Context, roleID, permissionID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM blc_admin_role_permission WHERE role_id = $1 AND permission_id = $2`, roleID, permissionID)
	return err
}

func (r *PostgresRoleRepository) IsRoleInUse(ctx context.Context, roleID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM blc_admin_user_role WHERE role_id = $1)`, roleID).Scan(&exists)
	return exists, err
}
