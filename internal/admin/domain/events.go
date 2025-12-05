package domain

import "time"

// AdminEventType represents the type of admin event
type AdminEventType string

const (
	// User events
	EventUserCreated          AdminEventType = "admin.user.created"
	EventUserUpdated          AdminEventType = "admin.user.updated"
	EventUserDeleted          AdminEventType = "admin.user.deleted"
	EventUserActivated        AdminEventType = "admin.user.activated"
	EventUserDeactivated      AdminEventType = "admin.user.deactivated"
	EventUserLoggedIn         AdminEventType = "admin.user.logged_in"
	EventUserLoggedOut        AdminEventType = "admin.user.logged_out"
	EventUserLoginFailed      AdminEventType = "admin.user.login_failed"
	EventPasswordChanged      AdminEventType = "admin.user.password_changed"

	// Role events
	EventRoleCreated          AdminEventType = "admin.role.created"
	EventRoleUpdated          AdminEventType = "admin.role.updated"
	EventRoleDeleted          AdminEventType = "admin.role.deleted"
	EventRoleActivated        AdminEventType = "admin.role.activated"
	EventRoleDeactivated      AdminEventType = "admin.role.deactivated"
	EventRoleAssignedToUser   AdminEventType = "admin.role.assigned_to_user"
	EventRoleUnassignedFromUser AdminEventType = "admin.role.unassigned_from_user"

	// Permission events
	EventPermissionCreated       AdminEventType = "admin.permission.created"
	EventPermissionUpdated       AdminEventType = "admin.permission.updated"
	EventPermissionDeleted       AdminEventType = "admin.permission.deleted"
	EventPermissionGrantedToRole AdminEventType = "admin.permission.granted_to_role"
	EventPermissionRevokedFromRole AdminEventType = "admin.permission.revoked_from_role"
)

// AdminEvent is the base event for all admin-related events
type AdminEvent struct {
	EventType  AdminEventType
	EventID    string
	UserID     int64
	Username   string
	OccurredAt time.Time
	IPAddress  string
	Data       interface{}
}

// UserCreatedEvent is published when a user is created
type UserCreatedEvent struct {
	AdminEvent
	CreatedUserID int64
	CreatedUsername string
}

// UserUpdatedEvent is published when a user is updated
type UserUpdatedEvent struct {
	AdminEvent
	UpdatedUserID int64
	Changes       map[string]interface{}
}

// UserDeletedEvent is published when a user is deleted
type UserDeletedEvent struct {
	AdminEvent
	DeletedUserID int64
	DeletedUsername string
}

// UserLoggedInEvent is published when a user logs in
type UserLoggedInEvent struct {
	AdminEvent
	LoginMethod string
	UserAgent   string
}

// UserLoginFailedEvent is published when a login attempt fails
type UserLoginFailedEvent struct {
	AdminEvent
	Username    string
	Reason      string
	UserAgent   string
}

// PasswordChangedEvent is published when a password is changed
type PasswordChangedEvent struct {
	AdminEvent
	ChangedForUserID int64
}

// RoleCreatedEvent is published when a role is created
type RoleCreatedEvent struct {
	AdminEvent
	RoleID   int64
	RoleName string
}

// RoleUpdatedEvent is published when a role is updated
type RoleUpdatedEvent struct {
	AdminEvent
	RoleID  int64
	Changes map[string]interface{}
}

// RoleDeletedEvent is published when a role is deleted
type RoleDeletedEvent struct {
	AdminEvent
	RoleID   int64
	RoleName string
}

// RoleAssignedToUserEvent is published when a role is assigned to a user
type RoleAssignedToUserEvent struct {
	AdminEvent
	TargetUserID int64
	RoleID       int64
	RoleName     string
}

// RoleUnassignedFromUserEvent is published when a role is unassigned from a user
type RoleUnassignedFromUserEvent struct {
	AdminEvent
	TargetUserID int64
	RoleID       int64
	RoleName     string
}

// PermissionGrantedToRoleEvent is published when a permission is granted to a role
type PermissionGrantedToRoleEvent struct {
	AdminEvent
	RoleID         int64
	RoleName       string
	PermissionID   int64
	PermissionName string
}

// PermissionRevokedFromRoleEvent is published when a permission is revoked from a role
type PermissionRevokedFromRoleEvent struct {
	AdminEvent
	RoleID         int64
	RoleName       string
	PermissionID   int64
	PermissionName string
}

// NewAdminEvent creates a new admin event
func NewAdminEvent(eventType AdminEventType, userID int64, username, ipAddress string) AdminEvent {
	return AdminEvent{
		EventType:  eventType,
		EventID:    generateEventID(),
		UserID:     userID,
		Username:   username,
		OccurredAt: time.Now(),
		IPAddress:  ipAddress,
	}
}

func generateEventID() string {
	return "EVT-" + time.Now().Format("20060102150405")
}
