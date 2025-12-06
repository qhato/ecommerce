package domain

import (
	"errors"
	"time"
)

// Tenant represents a multi-tenant organization that can own multiple stores
type Tenant struct {
	ID          int64
	Code        string
	Name        string
	Description string
	Status      TenantStatus
	Plan        TenantPlan
	Email       string
	Phone       string
	Website     string
	Settings    TenantSettings
	Limits      TenantLimits
	Metadata    map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TenantStatus represents the status of a tenant
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "ACTIVE"
	TenantStatusSuspended TenantStatus = "SUSPENDED"
	TenantStatusCanceled  TenantStatus = "CANCELED"
	TenantStatusTrial     TenantStatus = "TRIAL"
)

// TenantPlan represents the subscription plan
type TenantPlan string

const (
	TenantPlanFree       TenantPlan = "FREE"
	TenantPlanStarter    TenantPlan = "STARTER"
	TenantPlanPro        TenantPlan = "PRO"
	TenantPlanEnterprise TenantPlan = "ENTERPRISE"
)

// TenantSettings represents tenant-specific settings
type TenantSettings struct {
	DefaultCurrency     string
	DefaultLocale       string
	DefaultTimezone     string
	AllowMultiStore     bool
	AllowCustomDomain   bool
	EnableSSO           bool
	EnableAPI           bool
	CustomBranding      BrandingSettings
	DataRetentionDays   int
}

// BrandingSettings represents branding configuration
type BrandingSettings struct {
	LogoURL        string
	FaviconURL     string
	PrimaryColor   string
	SecondaryColor string
	CompanyName    string
}

// TenantLimits represents usage limits for the tenant
type TenantLimits struct {
	MaxStores      int
	MaxProducts    int
	MaxOrders      int
	MaxUsers       int
	MaxAPIRequests int
	MaxStorage     int64 // in bytes
}

// NewTenant creates a new tenant
func NewTenant(code, name string, plan TenantPlan) (*Tenant, error) {
	if code == "" {
		return nil, errors.New("tenant code is required")
	}
	if name == "" {
		return nil, errors.New("tenant name is required")
	}

	now := time.Now()
	return &Tenant{
		Code:     code,
		Name:     name,
		Status:   TenantStatusActive,
		Plan:     plan,
		Metadata: make(map[string]interface{}),
		Settings: TenantSettings{
			DefaultCurrency:   "USD",
			DefaultLocale:     "en",
			DefaultTimezone:   "UTC",
			AllowMultiStore:   true,
			EnableAPI:         true,
			DataRetentionDays: 365,
		},
		Limits:    getDefaultLimits(plan),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// getDefaultLimits returns default limits based on plan
func getDefaultLimits(plan TenantPlan) TenantLimits {
	switch plan {
	case TenantPlanFree:
		return TenantLimits{
			MaxStores:      1,
			MaxProducts:    100,
			MaxOrders:      1000,
			MaxUsers:       2,
			MaxAPIRequests: 10000,
			MaxStorage:     1073741824, // 1GB
		}
	case TenantPlanStarter:
		return TenantLimits{
			MaxStores:      3,
			MaxProducts:    1000,
			MaxOrders:      10000,
			MaxUsers:       5,
			MaxAPIRequests: 100000,
			MaxStorage:     5368709120, // 5GB
		}
	case TenantPlanPro:
		return TenantLimits{
			MaxStores:      10,
			MaxProducts:    10000,
			MaxOrders:      100000,
			MaxUsers:       20,
			MaxAPIRequests: 1000000,
			MaxStorage:     21474836480, // 20GB
		}
	case TenantPlanEnterprise:
		return TenantLimits{
			MaxStores:      -1, // unlimited
			MaxProducts:    -1,
			MaxOrders:      -1,
			MaxUsers:       -1,
			MaxAPIRequests: -1,
			MaxStorage:     -1,
		}
	default:
		return getDefaultLimits(TenantPlanFree)
	}
}

// Activate activates the tenant
func (t *Tenant) Activate() {
	t.Status = TenantStatusActive
	t.UpdatedAt = time.Now()
}

// Suspend suspends the tenant
func (t *Tenant) Suspend() {
	t.Status = TenantStatusSuspended
	t.UpdatedAt = time.Now()
}

// Cancel cancels the tenant
func (t *Tenant) Cancel() {
	t.Status = TenantStatusCanceled
	t.UpdatedAt = time.Now()
}

// UpgradePlan upgrades the tenant's plan
func (t *Tenant) UpgradePlan(newPlan TenantPlan) {
	t.Plan = newPlan
	t.Limits = getDefaultLimits(newPlan)
	t.UpdatedAt = time.Now()
}

// CanCreateStore checks if tenant can create a new store
func (t *Tenant) CanCreateStore(currentStoreCount int) bool {
	if t.Status != TenantStatusActive {
		return false
	}
	if t.Limits.MaxStores == -1 {
		return true
	}
	return currentStoreCount < t.Limits.MaxStores
}

// CanCreateProduct checks if tenant can create a new product
func (t *Tenant) CanCreateProduct(currentProductCount int) bool {
	if t.Status != TenantStatusActive {
		return false
	}
	if t.Limits.MaxProducts == -1 {
		return true
	}
	return currentProductCount < t.Limits.MaxProducts
}
