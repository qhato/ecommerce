package domain

import "time"

// Role represents an administrative role
type Role struct {
	ID          int64
	Name        string
	Description string
	IsActive    bool
	Permissions []Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewRole creates a new role
func NewRole(name, description string) (*Role, error) {
	if name == "" {
		return nil, ErrRoleNameRequired
	}

	now := time.Now()
	return &Role{
		Name:        name,
		Description: description,
		IsActive:    true,
		Permissions: make([]Permission, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Activate activates the role
func (r *Role) Activate() {
	r.IsActive = true
	r.UpdatedAt = time.Now()
}

// Deactivate deactivates the role
func (r *Role) Deactivate() {
	r.IsActive = false
	r.UpdatedAt = time.Now()
}

// UpdateInfo updates the role information
func (r *Role) UpdateInfo(name, description string) error {
	if name == "" {
		return ErrRoleNameRequired
	}
	r.Name = name
	r.Description = description
	r.UpdatedAt = time.Now()
	return nil
}

// HasPermission checks if the role has a specific permission
func (r *Role) HasPermission(permissionName string) bool {
	for _, permission := range r.Permissions {
		if permission.Name == permissionName {
			return true
		}
	}
	return false
}

// AddPermission adds a permission to the role (in-memory only)
func (r *Role) AddPermission(permission Permission) {
	if !r.HasPermission(permission.Name) {
		r.Permissions = append(r.Permissions, permission)
		r.UpdatedAt = time.Now()
	}
}

// RemovePermission removes a permission from the role (in-memory only)
func (r *Role) RemovePermission(permissionName string) {
	for i, permission := range r.Permissions {
		if permission.Name == permissionName {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			r.UpdatedAt = time.Now()
			break
		}
	}
}
