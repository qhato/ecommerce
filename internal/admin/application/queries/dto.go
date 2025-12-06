package queries

import (
	"time"

	"github.com/qhato/ecommerce/internal/admin/domain"
)

// AdminUserDTO represents an admin user for API responses
type AdminUserDTO struct {
	ID          int64       `json:"id"`
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	IsActive    bool        `json:"is_active"`
	IsSuper     bool        `json:"is_super"`
	LastLoginAt *time.Time  `json:"last_login_at,omitempty"`
	Roles       []RoleDTO   `json:"roles"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// RoleDTO represents a role for API responses
type RoleDTO struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	IsActive    bool            `json:"is_active"`
	Permissions []PermissionDTO `json:"permissions,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// PermissionDTO represents a permission for API responses
type PermissionDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AuditLogDTO represents an audit log for API responses
type AuditLogDTO struct {
	ID          int64                  `json:"id"`
	UserID      int64                  `json:"user_id"`
	Username    string                 `json:"username"`
	Action      string                 `json:"action"`
	Resource    string                 `json:"resource"`
	ResourceID  string                 `json:"resource_id"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"`
	IPAddress   string                 `json:"ip_address,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Success     bool                   `json:"success"`
	ErrorMsg    *string                `json:"error_msg,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// ToAdminUserDTO converts domain AdminUser to AdminUserDTO
func ToAdminUserDTO(user *domain.AdminUser) *AdminUserDTO {
	roles := make([]RoleDTO, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = RoleDTO{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			IsActive:    role.IsActive,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		}
	}

	return &AdminUserDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		IsActive:    user.IsActive,
		IsSuper:     user.IsSuper,
		LastLoginAt: user.LastLoginAt,
		Roles:       roles,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// ToRoleDTO converts domain Role to RoleDTO
func ToRoleDTO(role *domain.Role) *RoleDTO {
	permissions := make([]PermissionDTO, len(role.Permissions))
	for i, perm := range role.Permissions {
		permissions[i] = PermissionDTO{
			ID:          perm.ID,
			Name:        perm.Name,
			Description: perm.Description,
			Resource:    string(perm.Resource),
			Action:      string(perm.Action),
			IsActive:    perm.IsActive,
			CreatedAt:   perm.CreatedAt,
			UpdatedAt:   perm.UpdatedAt,
		}
	}

	return &RoleDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsActive:    role.IsActive,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// ToPermissionDTO converts domain Permission to PermissionDTO
func ToPermissionDTO(permission *domain.Permission) *PermissionDTO {
	return &PermissionDTO{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		Resource:    string(permission.Resource),
		Action:      string(permission.Action),
		IsActive:    permission.IsActive,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}
}

// ToAuditLogDTO converts domain AuditLog to AuditLogDTO
func ToAuditLogDTO(log *domain.AuditLog) *AuditLogDTO {
	return &AuditLogDTO{
		ID:          log.ID,
		UserID:      log.UserID,
		Username:    log.Username,
		Action:      string(log.Action),
		Resource:    log.Resource,
		ResourceID:  log.ResourceID,
		Description: log.Description,
		Severity:    string(log.Severity),
		IPAddress:   log.IPAddress,
		UserAgent:   log.UserAgent,
		Details:     log.Details,
		Success:     log.Success,
		ErrorMsg:    log.ErrorMsg,
		CreatedAt:   log.CreatedAt,
	}
}
