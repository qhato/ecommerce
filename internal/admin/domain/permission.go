package domain

import "time"

// PermissionResource represents the resource type
type PermissionResource string

const (
	ResourceProduct      PermissionResource = "PRODUCT"
	ResourceCategory     PermissionResource = "CATEGORY"
	ResourceOrder        PermissionResource = "ORDER"
	ResourceCustomer     PermissionResource = "CUSTOMER"
	ResourcePromotion    PermissionResource = "PROMOTION"
	ResourceCoupon       PermissionResource = "COUPON"
	ResourceContent      PermissionResource = "CONTENT"
	ResourceMenu         PermissionResource = "MENU"
	ResourceWorkflow     PermissionResource = "WORKFLOW"
	ResourceAdmin        PermissionResource = "ADMIN"
	ResourceRole         PermissionResource = "ROLE"
	ResourcePermission   PermissionResource = "PERMISSION"
	ResourceReport       PermissionResource = "REPORT"
	ResourceSettings     PermissionResource = "SETTINGS"
	ResourceAudit        PermissionResource = "AUDIT"
)

// PermissionAction represents the action type
type PermissionAction string

const (
	ActionCreate PermissionAction = "CREATE"
	ActionRead   PermissionAction = "READ"
	ActionUpdate PermissionAction = "UPDATE"
	ActionDelete PermissionAction = "DELETE"
	ActionList   PermissionAction = "LIST"
	ActionExport PermissionAction = "EXPORT"
	ActionImport PermissionAction = "IMPORT"
	ActionAll    PermissionAction = "ALL"
)

// Permission represents a permission
type Permission struct {
	ID          int64
	Name        string // e.g., "PRODUCT_CREATE", "ORDER_READ"
	Description string
	Resource    PermissionResource
	Action      PermissionAction
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewPermission creates a new permission
func NewPermission(name, description string, resource PermissionResource, action PermissionAction) (*Permission, error) {
	if name == "" {
		return nil, ErrPermissionNameRequired
	}
	if resource == "" {
		return nil, ErrPermissionResourceRequired
	}
	if action == "" {
		return nil, ErrPermissionActionRequired
	}

	now := time.Now()
	return &Permission{
		Name:        name,
		Description: description,
		Resource:    resource,
		Action:      action,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Activate activates the permission
func (p *Permission) Activate() {
	p.IsActive = true
	p.UpdatedAt = time.Now()
}

// Deactivate deactivates the permission
func (p *Permission) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

// UpdateInfo updates the permission information
func (p *Permission) UpdateInfo(name, description string, resource PermissionResource, action PermissionAction) error {
	if name == "" {
		return ErrPermissionNameRequired
	}
	if resource == "" {
		return ErrPermissionResourceRequired
	}
	if action == "" {
		return ErrPermissionActionRequired
	}

	p.Name = name
	p.Description = description
	p.Resource = resource
	p.Action = action
	p.UpdatedAt = time.Now()
	return nil
}

// Matches checks if the permission matches the resource and action
func (p *Permission) Matches(resource PermissionResource, action PermissionAction) bool {
	if !p.IsActive {
		return false
	}

	// Check if resource matches
	resourceMatches := p.Resource == resource

	// Check if action matches (ActionAll matches any action)
	actionMatches := p.Action == action || p.Action == ActionAll

	return resourceMatches && actionMatches
}
