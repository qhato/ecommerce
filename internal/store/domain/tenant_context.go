package domain

import (
	"context"
	"errors"
)

type contextKey string

const (
	TenantContextKey contextKey = "tenant_id"
	StoreContextKey  contextKey = "store_id"
)

// WithTenantID adds tenant ID to context
func WithTenantID(ctx context.Context, tenantID int64) context.Context {
	return context.WithValue(ctx, TenantContextKey, tenantID)
}

// GetTenantID retrieves tenant ID from context
func GetTenantID(ctx context.Context) (int64, error) {
	tenantID, ok := ctx.Value(TenantContextKey).(int64)
	if !ok {
		return 0, errors.New("tenant ID not found in context")
	}
	return tenantID, nil
}

// WithStoreID adds store ID to context
func WithStoreID(ctx context.Context, storeID int64) context.Context {
	return context.WithValue(ctx, StoreContextKey, storeID)
}

// GetStoreID retrieves store ID from context
func GetStoreID(ctx context.Context) (int64, error) {
	storeID, ok := ctx.Value(StoreContextKey).(int64)
	if !ok {
		return 0, errors.New("store ID not found in context")
	}
	return storeID, nil
}

// TenantContext represents the full tenant context
type TenantContext struct {
	TenantID int64
	StoreID  *int64
}

// NewTenantContext creates a new tenant context
func NewTenantContext(tenantID int64, storeID *int64) *TenantContext {
	return &TenantContext{
		TenantID: tenantID,
		StoreID:  storeID,
	}
}

// WithTenantContext adds full tenant context to context
func WithTenantContext(ctx context.Context, tc *TenantContext) context.Context {
	ctx = WithTenantID(ctx, tc.TenantID)
	if tc.StoreID != nil {
		ctx = WithStoreID(ctx, *tc.StoreID)
	}
	return ctx
}

// GetTenantContext retrieves full tenant context from context
func GetTenantContext(ctx context.Context) (*TenantContext, error) {
	tenantID, err := GetTenantID(ctx)
	if err != nil {
		return nil, err
	}

	tc := &TenantContext{
		TenantID: tenantID,
	}

	if storeID, err := GetStoreID(ctx); err == nil {
		tc.StoreID = &storeID
	}

	return tc, nil
}
