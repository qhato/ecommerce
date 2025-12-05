package domain

import (
	"time"
)

// AdminUser represents an administrative user
type AdminUser struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	IsActive     bool
	IsSuper      bool // Super admin has all permissions
	LastLoginAt  *time.Time
	Roles        []Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewAdminUser creates a new admin user
func NewAdminUser(username, email, passwordHash, firstName, lastName string) (*AdminUser, error) {
	if username == "" {
		return nil, ErrUsernameRequired
	}
	if email == "" {
		return nil, ErrEmailRequired
	}
	if passwordHash == "" {
		return nil, ErrPasswordRequired
	}

	now := time.Now()
	return &AdminUser{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		IsActive:     true,
		IsSuper:      false,
		Roles:        make([]Role, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Activate activates the user
func (u *AdminUser) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Deactivate deactivates the user
func (u *AdminUser) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// MakeSuper grants super admin privileges
func (u *AdminUser) MakeSuper() {
	u.IsSuper = true
	u.UpdatedAt = time.Now()
}

// RemoveSuper removes super admin privileges
func (u *AdminUser) RemoveSuper() {
	u.IsSuper = false
	u.UpdatedAt = time.Now()
}

// UpdatePassword updates the user's password hash
func (u *AdminUser) UpdatePassword(newPasswordHash string) error {
	if newPasswordHash == "" {
		return ErrPasswordRequired
	}
	u.PasswordHash = newPasswordHash
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateProfile updates the user's profile information
func (u *AdminUser) UpdateProfile(firstName, lastName, email string) error {
	if email == "" {
		return ErrEmailRequired
	}
	u.FirstName = firstName
	u.LastName = lastName
	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}

// RecordLogin records a successful login
func (u *AdminUser) RecordLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// HasRole checks if the user has a specific role
func (u *AdminUser) HasRole(roleName string) bool {
	if u.IsSuper {
		return true
	}
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// HasPermission checks if the user has a specific permission
func (u *AdminUser) HasPermission(permissionName string) bool {
	if u.IsSuper {
		return true
	}
	for _, role := range u.Roles {
		if role.HasPermission(permissionName) {
			return true
		}
	}
	return false
}

// GetAllPermissions returns all permissions the user has
func (u *AdminUser) GetAllPermissions() []Permission {
	if u.IsSuper {
		// Super admin has all permissions - return empty to indicate this
		return nil
	}

	permissionMap := make(map[string]Permission)
	for _, role := range u.Roles {
		for _, permission := range role.Permissions {
			permissionMap[permission.Name] = permission
		}
	}

	permissions := make([]Permission, 0, len(permissionMap))
	for _, permission := range permissionMap {
		permissions = append(permissions, permission)
	}
	return permissions
}

// FullName returns the user's full name
func (u *AdminUser) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return u.FirstName + " " + u.LastName
}
