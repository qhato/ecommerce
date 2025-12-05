package services

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/admin/domain"
)

// AuthorizationService handles authorization operations
type AuthorizationService struct {
	userRepo domain.AdminUserRepository
	roleRepo domain.RoleRepository
}

// NewAuthorizationService creates a new authorization service
func NewAuthorizationService(
	userRepo domain.AdminUserRepository,
	roleRepo domain.RoleRepository,
) *AuthorizationService {
	return &AuthorizationService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// CheckPermission checks if a user has a specific permission
func (s *AuthorizationService) CheckPermission(ctx context.Context, userID int64, permissionName string) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return false, domain.ErrUserNotFound
	}

	// Super admin has all permissions
	if user.IsSuper {
		return true, nil
	}

	// Get user roles with permissions
	roles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Load permissions for each role
	for _, role := range roles {
		permissions, err := s.roleRepo.GetRolePermissions(ctx, role.ID)
		if err != nil {
			continue
		}
		for _, permission := range permissions {
			if permission.Name == permissionName && permission.IsActive {
				return true, nil
			}
		}
	}

	return false, nil
}

// CheckResourcePermission checks if a user has permission for a specific resource and action
func (s *AuthorizationService) CheckResourcePermission(ctx context.Context, userID int64, resource domain.PermissionResource, action domain.PermissionAction) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return false, domain.ErrUserNotFound
	}

	// Super admin has all permissions
	if user.IsSuper {
		return true, nil
	}

	// Get user roles with permissions
	roles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Check if any role has the required permission
	for _, role := range roles {
		permissions, err := s.roleRepo.GetRolePermissions(ctx, role.ID)
		if err != nil {
			continue
		}
		for _, permission := range permissions {
			if permission.Matches(resource, action) {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetUserPermissions returns all permissions for a user
func (s *AuthorizationService) GetUserPermissions(ctx context.Context, userID int64) ([]domain.Permission, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	// Super admin has all permissions - return empty array to indicate this
	if user.IsSuper {
		return nil, nil
	}

	roles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Collect all unique permissions
	permissionMap := make(map[int64]domain.Permission)
	for _, role := range roles {
		permissions, err := s.roleRepo.GetRolePermissions(ctx, role.ID)
		if err != nil {
			continue
		}
		for _, permission := range permissions {
			if permission.IsActive {
				permissionMap[permission.ID] = permission
			}
		}
	}

	// Convert map to slice
	permissions := make([]domain.Permission, 0, len(permissionMap))
	for _, permission := range permissionMap {
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// RequirePermission returns an error if user doesn't have the permission
func (s *AuthorizationService) RequirePermission(ctx context.Context, userID int64, permissionName string) error {
	hasPermission, err := s.CheckPermission(ctx, userID, permissionName)
	if err != nil {
		return err
	}
	if !hasPermission {
		return domain.ErrPermissionDenied
	}
	return nil
}

// RequireResourcePermission returns an error if user doesn't have the resource permission
func (s *AuthorizationService) RequireResourcePermission(ctx context.Context, userID int64, resource domain.PermissionResource, action domain.PermissionAction) error {
	hasPermission, err := s.CheckResourcePermission(ctx, userID, resource, action)
	if err != nil {
		return err
	}
	if !hasPermission {
		return domain.ErrPermissionDenied
	}
	return nil
}
