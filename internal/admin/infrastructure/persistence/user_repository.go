package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/admin/domain"
)

// PostgresAdminUserRepository implements domain.AdminUserRepository
type PostgresAdminUserRepository struct {
	db *sql.DB
}

// NewPostgresAdminUserRepository creates a new repository
func NewPostgresAdminUserRepository(db *sql.DB) *PostgresAdminUserRepository {
	return &PostgresAdminUserRepository{db: db}
}

// Create creates a new admin user
func (r *PostgresAdminUserRepository) Create(ctx context.Context, user *domain.AdminUser) error {
	query := `
		INSERT INTO blc_admin_user (
			username, email, password_hash, first_name, last_name,
			is_active, is_super, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.IsActive, user.IsSuper, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update updates an existing admin user
func (r *PostgresAdminUserRepository) Update(ctx context.Context, user *domain.AdminUser) error {
	query := `
		UPDATE blc_admin_user SET
			email = $1, password_hash = $2, first_name = $3, last_name = $4,
			is_active = $5, is_super = $6, last_login_at = $7, updated_at = $8
		WHERE id = $9`

	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.IsActive, user.IsSuper, user.LastLoginAt, user.UpdatedAt, user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// FindByID finds an admin user by ID
func (r *PostgresAdminUserRepository) FindByID(ctx context.Context, id int64) (*domain.AdminUser, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name,
			   is_active, is_super, last_login_at, created_at, updated_at
		FROM blc_admin_user WHERE id = $1`

	user := &domain.AdminUser{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.IsActive, &user.IsSuper,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByUsername finds an admin user by username
func (r *PostgresAdminUserRepository) FindByUsername(ctx context.Context, username string) (*domain.AdminUser, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name,
			   is_active, is_super, last_login_at, created_at, updated_at
		FROM blc_admin_user WHERE username = $1`

	user := &domain.AdminUser{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.IsActive, &user.IsSuper,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByEmail finds an admin user by email
func (r *PostgresAdminUserRepository) FindByEmail(ctx context.Context, email string) (*domain.AdminUser, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name,
			   is_active, is_super, last_login_at, created_at, updated_at
		FROM blc_admin_user WHERE email = $1`

	user := &domain.AdminUser{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.IsActive, &user.IsSuper,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindAll finds all admin users
func (r *PostgresAdminUserRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.AdminUser, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name,
			   is_active, is_super, last_login_at, created_at, updated_at
		FROM blc_admin_user`

	if activeOnly {
		query += " WHERE is_active = true"
	}

	query += " ORDER BY username ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	users := make([]*domain.AdminUser, 0)
	for rows.Next() {
		user := &domain.AdminUser{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.FirstName, &user.LastName, &user.IsActive, &user.IsSuper,
			&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// Delete deletes an admin user
func (r *PostgresAdminUserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_admin_user WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ExistsByUsername checks if a username already exists
func (r *PostgresAdminUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blc_admin_user WHERE username = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return exists, nil
}

// ExistsByEmail checks if an email already exists
func (r *PostgresAdminUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blc_admin_user WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// GetUserRoles gets all roles for a user
func (r *PostgresAdminUserRepository) GetUserRoles(ctx context.Context, userID int64) ([]domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.is_active, r.created_at, r.updated_at
		FROM blc_admin_role r
		INNER JOIN blc_admin_user_role ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.name ASC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user roles: %w", err)
	}
	defer rows.Close()

	roles := make([]domain.Role, 0)
	for rows.Next() {
		role := domain.Role{}
		err := rows.Scan(
			&role.ID, &role.Name, &role.Description,
			&role.IsActive, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// AssignRole assigns a role to a user
func (r *PostgresAdminUserRepository) AssignRole(ctx context.Context, userID, roleID int64) error {
	query := `
		INSERT INTO blc_admin_user_role (user_id, role_id, assigned_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id, role_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// UnassignRole unassigns a role from a user
func (r *PostgresAdminUserRepository) UnassignRole(ctx context.Context, userID, roleID int64) error {
	query := `DELETE FROM blc_admin_user_role WHERE user_id = $1 AND role_id = $2`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to unassign role: %w", err)
	}

	return nil
}
