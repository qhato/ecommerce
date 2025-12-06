package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/admin/domain"
)

// AdminQueryService handles admin-related queries
type AdminQueryService struct {
	userRepo       domain.AdminUserRepository
	roleRepo       domain.RoleRepository
	permissionRepo domain.PermissionRepository
	auditLogRepo   domain.AuditLogRepository
}

// NewAdminQueryService creates a new admin query service
func NewAdminQueryService(
	userRepo domain.AdminUserRepository,
	roleRepo domain.RoleRepository,
	permissionRepo domain.PermissionRepository,
	auditLogRepo domain.AuditLogRepository,
) *AdminQueryService {
	return &AdminQueryService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		auditLogRepo:   auditLogRepo,
	}
}

// GetUser retrieves a user by ID
func (s *AdminQueryService) GetUser(ctx context.Context, id int64) (*AdminUserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	// Load user roles
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	user.Roles = roles

	return ToAdminUserDTO(user), nil
}

// GetUserByUsername retrieves a user by username
func (s *AdminQueryService) GetUserByUsername(ctx context.Context, username string) (*AdminUserDTO, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	// Load user roles
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	user.Roles = roles

	return ToAdminUserDTO(user), nil
}

// GetAllUsers retrieves all users
func (s *AdminQueryService) GetAllUsers(ctx context.Context, activeOnly bool) ([]*AdminUserDTO, error) {
	users, err := s.userRepo.FindAll(ctx, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	dtos := make([]*AdminUserDTO, len(users))
	for i, user := range users {
		// Load user roles
		roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
		if err == nil {
			user.Roles = roles
		}
		dtos[i] = ToAdminUserDTO(user)
	}

	return dtos, nil
}

// GetRole retrieves a role by ID
func (s *AdminQueryService) GetRole(ctx context.Context, id int64) (*RoleDTO, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	if role == nil {
		return nil, domain.ErrRoleNotFound
	}

	// Load role permissions
	permissions, err := s.roleRepo.GetRolePermissions(ctx, role.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	role.Permissions = permissions

	return ToRoleDTO(role), nil
}

// GetAllRoles retrieves all roles
func (s *AdminQueryService) GetAllRoles(ctx context.Context, activeOnly bool) ([]*RoleDTO, error) {
	roles, err := s.roleRepo.FindAll(ctx, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find roles: %w", err)
	}

	dtos := make([]*RoleDTO, len(roles))
	for i, role := range roles {
		// Load role permissions
		permissions, err := s.roleRepo.GetRolePermissions(ctx, role.ID)
		if err == nil {
			role.Permissions = permissions
		}
		dtos[i] = ToRoleDTO(role)
	}

	return dtos, nil
}

// GetPermission retrieves a permission by ID
func (s *AdminQueryService) GetPermission(ctx context.Context, id int64) (*PermissionDTO, error) {
	permission, err := s.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find permission: %w", err)
	}
	if permission == nil {
		return nil, domain.ErrPermissionNotFound
	}

	return ToPermissionDTO(permission), nil
}

// GetAllPermissions retrieves all permissions
func (s *AdminQueryService) GetAllPermissions(ctx context.Context, activeOnly bool) ([]*PermissionDTO, error) {
	permissions, err := s.permissionRepo.FindAll(ctx, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find permissions: %w", err)
	}

	dtos := make([]*PermissionDTO, len(permissions))
	for i, permission := range permissions {
		dtos[i] = ToPermissionDTO(permission)
	}

	return dtos, nil
}

// GetPermissionsByResource retrieves permissions by resource
func (s *AdminQueryService) GetPermissionsByResource(ctx context.Context, resource string, activeOnly bool) ([]*PermissionDTO, error) {
	permissions, err := s.permissionRepo.FindByResource(ctx, domain.PermissionResource(resource), activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find permissions: %w", err)
	}

	dtos := make([]*PermissionDTO, len(permissions))
	for i, permission := range permissions {
		dtos[i] = ToPermissionDTO(permission)
	}

	return dtos, nil
}

// GetAuditLog retrieves an audit log by ID
func (s *AdminQueryService) GetAuditLog(ctx context.Context, id int64) (*AuditLogDTO, error) {
	log, err := s.auditLogRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find audit log: %w", err)
	}
	if log == nil {
		return nil, domain.ErrAuditLogNotFound
	}

	return ToAuditLogDTO(log), nil
}

// GetAuditLogsByUser retrieves audit logs by user ID
func (s *AdminQueryService) GetAuditLogsByUser(ctx context.Context, userID int64, limit int) ([]*AuditLogDTO, error) {
	logs, err := s.auditLogRepo.FindByUserID(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find audit logs: %w", err)
	}

	dtos := make([]*AuditLogDTO, len(logs))
	for i, log := range logs {
		dtos[i] = ToAuditLogDTO(log)
	}

	return dtos, nil
}

// GetRecentAuditLogs retrieves recent audit logs
func (s *AdminQueryService) GetRecentAuditLogs(ctx context.Context, limit int) ([]*AuditLogDTO, error) {
	logs, err := s.auditLogRepo.FindRecent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find audit logs: %w", err)
	}

	dtos := make([]*AuditLogDTO, len(logs))
	for i, log := range logs {
		dtos[i] = ToAuditLogDTO(log)
	}

	return dtos, nil
}

// GetSecurityEvents retrieves security-related audit logs
func (s *AdminQueryService) GetSecurityEvents(ctx context.Context, limit int) ([]*AuditLogDTO, error) {
	logs, err := s.auditLogRepo.FindSecurityEvents(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find security events: %w", err)
	}

	dtos := make([]*AuditLogDTO, len(logs))
	for i, log := range logs {
		dtos[i] = ToAuditLogDTO(log)
	}

	return dtos, nil
}

// GetFailedLogins retrieves failed login attempts for a username
func (s *AdminQueryService) GetFailedLogins(ctx context.Context, username string, limit int) ([]*AuditLogDTO, error) {
	logs, err := s.auditLogRepo.FindFailedLogins(ctx, username, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find failed logins: %w", err)
	}

	dtos := make([]*AuditLogDTO, len(logs))
	for i, log := range logs {
		dtos[i] = ToAuditLogDTO(log)
	}

	return dtos, nil
}
