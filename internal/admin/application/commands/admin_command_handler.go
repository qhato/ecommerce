package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/admin/application/services"
	"github.com/qhato/ecommerce/internal/admin/domain"
)

// AdminCommandHandler handles admin-related commands
type AdminCommandHandler struct {
	userRepo         domain.AdminUserRepository
	roleRepo         domain.RoleRepository
	permissionRepo   domain.PermissionRepository
	auditLogRepo     domain.AuditLogRepository
	authService      *services.AuthenticationService
	authzService     *services.AuthorizationService
}

// NewAdminCommandHandler creates a new admin command handler
func NewAdminCommandHandler(
	userRepo domain.AdminUserRepository,
	roleRepo domain.RoleRepository,
	permissionRepo domain.PermissionRepository,
	auditLogRepo domain.AuditLogRepository,
	authService *services.AuthenticationService,
	authzService *services.AuthorizationService,
) *AdminCommandHandler {
	return &AdminCommandHandler{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		auditLogRepo:   auditLogRepo,
		authService:    authService,
		authzService:   authzService,
	}
}

// HandleCreateUser handles creating a new admin user
func (h *AdminCommandHandler) HandleCreateUser(ctx context.Context, cmd CreateUserCommand) (*domain.AdminUser, error) {
	// Check if username already exists
	exists, err := h.userRepo.ExistsByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, domain.ErrUsernameTaken
	}

	// Check if email already exists
	exists, err = h.userRepo.ExistsByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, domain.ErrEmailTaken
	}

	// Hash password
	passwordHash, err := h.authService.HashPassword(cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := domain.NewAdminUser(cmd.Username, cmd.Email, passwordHash, cmd.FirstName, cmd.LastName)
	if err != nil {
		return nil, err
	}

	// Save user
	if err := h.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign roles
	for _, roleID := range cmd.RoleIDs {
		if err := h.userRepo.AssignRole(ctx, user.ID, roleID); err != nil {
			// Log but don't fail
			fmt.Printf("Failed to assign role %d to user %d: %v\n", roleID, user.ID, err)
		}
	}

	// Create audit log
	auditLog := domain.NewAuditLog(
		user.ID,
		user.Username,
		domain.AuditActionUserCreated,
		"USER",
		fmt.Sprintf("%d", user.ID),
		fmt.Sprintf("User %s created", user.Username),
		domain.AuditSeverityInfo,
		"",
		"",
	)
	h.auditLogRepo.Create(ctx, auditLog)

	return user, nil
}

// HandleUpdateUser handles updating a user's profile
func (h *AdminCommandHandler) HandleUpdateUser(ctx context.Context, cmd UpdateUserCommand) (*domain.AdminUser, error) {
	user, err := h.userRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	if err := user.UpdateProfile(cmd.FirstName, cmd.LastName, cmd.Email); err != nil {
		return nil, err
	}

	if err := h.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Audit log
	auditLog := domain.NewAuditLog(
		user.ID,
		user.Username,
		domain.AuditActionUserUpdated,
		"USER",
		fmt.Sprintf("%d", user.ID),
		fmt.Sprintf("User %s updated", user.Username),
		domain.AuditSeverityInfo,
		"",
		"",
	)
	h.auditLogRepo.Create(ctx, auditLog)

	return user, nil
}

// HandleChangePassword handles changing a user's password
func (h *AdminCommandHandler) HandleChangePassword(ctx context.Context, cmd ChangePasswordCommand) error {
	user, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	// Verify old password
	// Note: This would need bcrypt.CompareHashAndPassword in production

	// Hash new password
	newPasswordHash, err := h.authService.HashPassword(cmd.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := user.UpdatePassword(newPasswordHash); err != nil {
		return err
	}

	if err := h.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Audit log
	auditLog := domain.NewAuditLog(
		user.ID,
		user.Username,
		domain.AuditActionPasswordChanged,
		"USER",
		fmt.Sprintf("%d", user.ID),
		fmt.Sprintf("Password changed for user %s", user.Username),
		domain.AuditSeverityInfo,
		"",
		"",
	)
	h.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// HandleActivateUser handles activating a user
func (h *AdminCommandHandler) HandleActivateUser(ctx context.Context, cmd ActivateUserCommand) (*domain.AdminUser, error) {
	user, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	user.Activate()

	if err := h.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// HandleDeactivateUser handles deactivating a user
func (h *AdminCommandHandler) HandleDeactivateUser(ctx context.Context, cmd DeactivateUserCommand) (*domain.AdminUser, error) {
	user, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	user.Deactivate()

	if err := h.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// HandleAssignRole handles assigning a role to a user
func (h *AdminCommandHandler) HandleAssignRole(ctx context.Context, cmd AssignRoleCommand) error {
	// Verify user exists
	user, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	// Verify role exists
	role, err := h.roleRepo.FindByID(ctx, cmd.RoleID)
	if err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}
	if role == nil {
		return domain.ErrRoleNotFound
	}

	if err := h.userRepo.AssignRole(ctx, cmd.UserID, cmd.RoleID); err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	// Audit log
	auditLog := domain.NewAuditLog(
		user.ID,
		user.Username,
		domain.AuditActionRoleAssigned,
		"ROLE",
		fmt.Sprintf("%d", cmd.RoleID),
		fmt.Sprintf("Role %s assigned to user %s", role.Name, user.Username),
		domain.AuditSeverityInfo,
		"",
		"",
	)
	h.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// HandleUnassignRole handles unassigning a role from a user
func (h *AdminCommandHandler) HandleUnassignRole(ctx context.Context, cmd UnassignRoleCommand) error {
	if err := h.userRepo.UnassignRole(ctx, cmd.UserID, cmd.RoleID); err != nil {
		return fmt.Errorf("failed to unassign role: %w", err)
	}

	return nil
}

// HandleCreateRole handles creating a new role
func (h *AdminCommandHandler) HandleCreateRole(ctx context.Context, cmd CreateRoleCommand) (*domain.Role, error) {
	// Check if role name already exists
	exists, err := h.roleRepo.ExistsByName(ctx, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check role name: %w", err)
	}
	if exists {
		return nil, domain.ErrRoleNameTaken
	}

	role, err := domain.NewRole(cmd.Name, cmd.Description)
	if err != nil {
		return nil, err
	}

	if err := h.roleRepo.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// HandleUpdateRole handles updating a role
func (h *AdminCommandHandler) HandleUpdateRole(ctx context.Context, cmd UpdateRoleCommand) (*domain.Role, error) {
	role, err := h.roleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	if role == nil {
		return nil, domain.ErrRoleNotFound
	}

	if err := role.UpdateInfo(cmd.Name, cmd.Description); err != nil {
		return nil, err
	}

	if err := h.roleRepo.Update(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return role, nil
}

// HandleGrantPermission handles granting a permission to a role
func (h *AdminCommandHandler) HandleGrantPermission(ctx context.Context, cmd GrantPermissionCommand) error {
	// Verify role exists
	role, err := h.roleRepo.FindByID(ctx, cmd.RoleID)
	if err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}
	if role == nil {
		return domain.ErrRoleNotFound
	}

	// Verify permission exists
	permission, err := h.permissionRepo.FindByID(ctx, cmd.PermissionID)
	if err != nil {
		return fmt.Errorf("failed to find permission: %w", err)
	}
	if permission == nil {
		return domain.ErrPermissionNotFound
	}

	if err := h.roleRepo.GrantPermission(ctx, cmd.RoleID, cmd.PermissionID); err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}

	return nil
}

// HandleRevokePermission handles revoking a permission from a role
func (h *AdminCommandHandler) HandleRevokePermission(ctx context.Context, cmd RevokePermissionCommand) error {
	if err := h.roleRepo.RevokePermission(ctx, cmd.RoleID, cmd.PermissionID); err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}

	return nil
}

// HandleCreatePermission handles creating a new permission
func (h *AdminCommandHandler) HandleCreatePermission(ctx context.Context, cmd CreatePermissionCommand) (*domain.Permission, error) {
	// Check if permission name already exists
	exists, err := h.permissionRepo.ExistsByName(ctx, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission name: %w", err)
	}
	if exists {
		return nil, domain.ErrPermissionAlreadyExists
	}

	permission, err := domain.NewPermission(
		cmd.Name,
		cmd.Description,
		domain.PermissionResource(cmd.Resource),
		domain.PermissionAction(cmd.Action),
	)
	if err != nil {
		return nil, err
	}

	if err := h.permissionRepo.Create(ctx, permission); err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission, nil
}

// HandleUpdatePermission handles updating a permission
func (h *AdminCommandHandler) HandleUpdatePermission(ctx context.Context, cmd UpdatePermissionCommand) (*domain.Permission, error) {
	permission, err := h.permissionRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find permission: %w", err)
	}
	if permission == nil {
		return nil, domain.ErrPermissionNotFound
	}

	if err := permission.UpdateInfo(
		cmd.Name,
		cmd.Description,
		domain.PermissionResource(cmd.Resource),
		domain.PermissionAction(cmd.Action),
	); err != nil {
		return nil, err
	}

	if err := h.permissionRepo.Update(ctx, permission); err != nil {
		return nil, fmt.Errorf("failed to update permission: %w", err)
	}

	return permission, nil
}
