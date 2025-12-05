package domain

import "context"

// AdminUserRepository defines the interface for admin user persistence
type AdminUserRepository interface {
	// Create creates a new admin user
	Create(ctx context.Context, user *AdminUser) error

	// Update updates an existing admin user
	Update(ctx context.Context, user *AdminUser) error

	// FindByID finds an admin user by ID
	FindByID(ctx context.Context, id int64) (*AdminUser, error)

	// FindByUsername finds an admin user by username
	FindByUsername(ctx context.Context, username string) (*AdminUser, error)

	// FindByEmail finds an admin user by email
	FindByEmail(ctx context.Context, email string) (*AdminUser, error)

	// FindAll finds all admin users
	FindAll(ctx context.Context, activeOnly bool) ([]*AdminUser, error)

	// Delete deletes an admin user
	Delete(ctx context.Context, id int64) error

	// ExistsByUsername checks if a username already exists
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByEmail checks if an email already exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// GetUserRoles gets all roles for a user
	GetUserRoles(ctx context.Context, userID int64) ([]Role, error)

	// AssignRole assigns a role to a user
	AssignRole(ctx context.Context, userID, roleID int64) error

	// UnassignRole unassigns a role from a user
	UnassignRole(ctx context.Context, userID, roleID int64) error
}

// RoleRepository defines the interface for role persistence
type RoleRepository interface {
	// Create creates a new role
	Create(ctx context.Context, role *Role) error

	// Update updates an existing role
	Update(ctx context.Context, role *Role) error

	// FindByID finds a role by ID
	FindByID(ctx context.Context, id int64) (*Role, error)

	// FindByName finds a role by name
	FindByName(ctx context.Context, name string) (*Role, error)

	// FindAll finds all roles
	FindAll(ctx context.Context, activeOnly bool) ([]*Role, error)

	// Delete deletes a role
	Delete(ctx context.Context, id int64) error

	// ExistsByName checks if a role name already exists
	ExistsByName(ctx context.Context, name string) (bool, error)

	// GetRolePermissions gets all permissions for a role
	GetRolePermissions(ctx context.Context, roleID int64) ([]Permission, error)

	// GrantPermission grants a permission to a role
	GrantPermission(ctx context.Context, roleID, permissionID int64) error

	// RevokePermission revokes a permission from a role
	RevokePermission(ctx context.Context, roleID, permissionID int64) error

	// IsRoleInUse checks if a role is assigned to any users
	IsRoleInUse(ctx context.Context, roleID int64) (bool, error)
}

// PermissionRepository defines the interface for permission persistence
type PermissionRepository interface {
	// Create creates a new permission
	Create(ctx context.Context, permission *Permission) error

	// Update updates an existing permission
	Update(ctx context.Context, permission *Permission) error

	// FindByID finds a permission by ID
	FindByID(ctx context.Context, id int64) (*Permission, error)

	// FindByName finds a permission by name
	FindByName(ctx context.Context, name string) (*Permission, error)

	// FindAll finds all permissions
	FindAll(ctx context.Context, activeOnly bool) ([]*Permission, error)

	// FindByResource finds permissions by resource
	FindByResource(ctx context.Context, resource PermissionResource, activeOnly bool) ([]*Permission, error)

	// Delete deletes a permission
	Delete(ctx context.Context, id int64) error

	// ExistsByName checks if a permission name already exists
	ExistsByName(ctx context.Context, name string) (bool, error)
}

// AuditLogRepository defines the interface for audit log persistence
type AuditLogRepository interface {
	// Create creates a new audit log entry
	Create(ctx context.Context, log *AuditLog) error

	// FindByID finds an audit log by ID
	FindByID(ctx context.Context, id int64) (*AuditLog, error)

	// FindByUserID finds audit logs by user ID
	FindByUserID(ctx context.Context, userID int64, limit int) ([]*AuditLog, error)

	// FindByAction finds audit logs by action
	FindByAction(ctx context.Context, action AuditAction, limit int) ([]*AuditLog, error)

	// FindByResource finds audit logs by resource and resource ID
	FindByResource(ctx context.Context, resource, resourceID string, limit int) ([]*AuditLog, error)

	// FindSecurityEvents finds security-related events
	FindSecurityEvents(ctx context.Context, limit int) ([]*AuditLog, error)

	// FindFailedLogins finds failed login attempts
	FindFailedLogins(ctx context.Context, username string, limit int) ([]*AuditLog, error)

	// FindRecent finds recent audit logs
	FindRecent(ctx context.Context, limit int) ([]*AuditLog, error)

	// FindBetween finds audit logs between two dates
	FindBetween(ctx context.Context, startTime, endTime int64, limit int) ([]*AuditLog, error)
}
