package domain

import "time"

// AuditAction represents the type of action performed
type AuditAction string

const (
	// User actions
	AuditActionUserCreated      AuditAction = "USER_CREATED"
	AuditActionUserUpdated      AuditAction = "USER_UPDATED"
	AuditActionUserDeleted      AuditAction = "USER_DELETED"
	AuditActionUserActivated    AuditAction = "USER_ACTIVATED"
	AuditActionUserDeactivated  AuditAction = "USER_DEACTIVATED"
	AuditActionUserLogin        AuditAction = "USER_LOGIN"
	AuditActionUserLogout       AuditAction = "USER_LOGOUT"
	AuditActionUserLoginFailed  AuditAction = "USER_LOGIN_FAILED"
	AuditActionPasswordChanged  AuditAction = "PASSWORD_CHANGED"

	// Role actions
	AuditActionRoleCreated      AuditAction = "ROLE_CREATED"
	AuditActionRoleUpdated      AuditAction = "ROLE_UPDATED"
	AuditActionRoleDeleted      AuditAction = "ROLE_DELETED"
	AuditActionRoleAssigned     AuditAction = "ROLE_ASSIGNED"
	AuditActionRoleUnassigned   AuditAction = "ROLE_UNASSIGNED"

	// Permission actions
	AuditActionPermissionCreated  AuditAction = "PERMISSION_CREATED"
	AuditActionPermissionUpdated  AuditAction = "PERMISSION_UPDATED"
	AuditActionPermissionDeleted  AuditAction = "PERMISSION_DELETED"
	AuditActionPermissionGranted  AuditAction = "PERMISSION_GRANTED"
	AuditActionPermissionRevoked  AuditAction = "PERMISSION_REVOKED"

	// Resource actions
	AuditActionResourceCreated   AuditAction = "RESOURCE_CREATED"
	AuditActionResourceUpdated   AuditAction = "RESOURCE_UPDATED"
	AuditActionResourceDeleted   AuditAction = "RESOURCE_DELETED"
	AuditActionResourceViewed    AuditAction = "RESOURCE_VIEWED"
	AuditActionResourceExported  AuditAction = "RESOURCE_EXPORTED"
	AuditActionResourceImported  AuditAction = "RESOURCE_IMPORTED"

	// Settings actions
	AuditActionSettingUpdated    AuditAction = "SETTING_UPDATED"
)

// AuditSeverity represents the severity level of the audit log
type AuditSeverity string

const (
	AuditSeverityInfo     AuditSeverity = "INFO"
	AuditSeverityWarning  AuditSeverity = "WARNING"
	AuditSeverityError    AuditSeverity = "ERROR"
	AuditSeverityCritical AuditSeverity = "CRITICAL"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID          int64
	UserID      int64
	Username    string
	Action      AuditAction
	Resource    string // e.g., "USER", "PRODUCT", "ORDER"
	ResourceID  string // ID of the affected resource
	Description string
	Severity    AuditSeverity
	IPAddress   string
	UserAgent   string
	Details     map[string]interface{} // Additional context
	Success     bool
	ErrorMsg    *string
	CreatedAt   time.Time
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(
	userID int64,
	username string,
	action AuditAction,
	resource, resourceID string,
	description string,
	severity AuditSeverity,
	ipAddress, userAgent string,
) *AuditLog {
	return &AuditLog{
		UserID:      userID,
		Username:    username,
		Action:      action,
		Resource:    resource,
		ResourceID:  resourceID,
		Description: description,
		Severity:    severity,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Details:     make(map[string]interface{}),
		Success:     true,
		CreatedAt:   time.Now(),
	}
}

// MarkAsFailure marks the audit log as a failed operation
func (a *AuditLog) MarkAsFailure(errorMsg string) {
	a.Success = false
	a.ErrorMsg = &errorMsg
}

// AddDetail adds additional context to the audit log
func (a *AuditLog) AddDetail(key string, value interface{}) {
	if a.Details == nil {
		a.Details = make(map[string]interface{})
	}
	a.Details[key] = value
}

// IsSecurityEvent checks if this is a security-related event
func (a *AuditLog) IsSecurityEvent() bool {
	securityActions := []AuditAction{
		AuditActionUserLogin,
		AuditActionUserLogout,
		AuditActionUserLoginFailed,
		AuditActionPasswordChanged,
		AuditActionUserActivated,
		AuditActionUserDeactivated,
		AuditActionRoleAssigned,
		AuditActionRoleUnassigned,
		AuditActionPermissionGranted,
		AuditActionPermissionRevoked,
	}

	for _, action := range securityActions {
		if a.Action == action {
			return true
		}
	}
	return false
}
