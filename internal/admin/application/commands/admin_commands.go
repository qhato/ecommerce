package commands

// User Commands

type CreateUserCommand struct {
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
	RoleIDs   []int64
}

type UpdateUserCommand struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
}

type ChangePasswordCommand struct {
	UserID      int64
	OldPassword string
	NewPassword string
}

type ResetPasswordCommand struct {
	UserID      int64
	NewPassword string
}

type ActivateUserCommand struct {
	UserID int64
}

type DeactivateUserCommand struct {
	UserID int64
}

type MakeUserSuperCommand struct {
	UserID int64
}

type RemoveUserSuperCommand struct {
	UserID int64
}

type AssignRoleCommand struct {
	UserID int64
	RoleID int64
}

type UnassignRoleCommand struct {
	UserID int64
	RoleID int64
}

type DeleteUserCommand struct {
	UserID int64
}

// Role Commands

type CreateRoleCommand struct {
	Name        string
	Description string
}

type UpdateRoleCommand struct {
	ID          int64
	Name        string
	Description string
}

type ActivateRoleCommand struct {
	ID int64
}

type DeactivateRoleCommand struct {
	ID int64
}

type GrantPermissionCommand struct {
	RoleID       int64
	PermissionID int64
}

type RevokePermissionCommand struct {
	RoleID       int64
	PermissionID int64
}

type DeleteRoleCommand struct {
	ID int64
}

// Permission Commands

type CreatePermissionCommand struct {
	Name        string
	Description string
	Resource    string
	Action      string
}

type UpdatePermissionCommand struct {
	ID          int64
	Name        string
	Description string
	Resource    string
	Action      string
}

type ActivatePermissionCommand struct {
	ID int64
}

type DeactivatePermissionCommand struct {
	ID int64
}

type DeletePermissionCommand struct {
	ID int64
}

// Login Command

type LoginCommand struct {
	Username  string
	Password  string
	IPAddress string
	UserAgent string
}

type LogoutCommand struct {
	UserID    int64
	IPAddress string
	UserAgent string
}
