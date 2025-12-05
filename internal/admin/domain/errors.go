package domain

import "errors"

// User Errors
var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrUsernameRequired       = errors.New("username is required")
	ErrEmailRequired          = errors.New("email is required")
	ErrPasswordRequired       = errors.New("password is required")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrUserNotActive          = errors.New("user is not active")
	ErrUsernameTaken          = errors.New("username is already taken")
	ErrEmailTaken             = errors.New("email is already taken")
)

// Role Errors
var (
	ErrRoleNotFound           = errors.New("role not found")
	ErrRoleAlreadyExists      = errors.New("role already exists")
	ErrRoleNameRequired       = errors.New("role name is required")
	ErrRoleNotActive          = errors.New("role is not active")
	ErrRoleInUse              = errors.New("role is in use and cannot be deleted")
	ErrRoleNameTaken          = errors.New("role name is already taken")
)

// Permission Errors
var (
	ErrPermissionNotFound         = errors.New("permission not found")
	ErrPermissionAlreadyExists    = errors.New("permission already exists")
	ErrPermissionNameRequired     = errors.New("permission name is required")
	ErrPermissionResourceRequired = errors.New("permission resource is required")
	ErrPermissionActionRequired   = errors.New("permission action is required")
	ErrPermissionNotActive        = errors.New("permission is not active")
	ErrPermissionDenied           = errors.New("permission denied")
	ErrInsufficientPermissions    = errors.New("insufficient permissions")
)

// Authentication Errors
var (
	ErrAuthenticationFailed   = errors.New("authentication failed")
	ErrInvalidToken           = errors.New("invalid token")
	ErrTokenExpired           = errors.New("token expired")
	ErrTokenNotFound          = errors.New("token not found")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
)

// Authorization Errors
var (
	ErrUnauthorized           = errors.New("unauthorized")
	ErrForbidden              = errors.New("forbidden")
	ErrAccessDenied           = errors.New("access denied")
)

// Audit Log Errors
var (
	ErrAuditLogNotFound       = errors.New("audit log not found")
	ErrAuditLogCreationFailed = errors.New("failed to create audit log")
)

// Repository Errors
var (
	ErrRepositoryFailure      = errors.New("repository operation failed")
)
