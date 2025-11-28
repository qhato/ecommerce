package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/qhato/ecommerce/pkg/auth"
	"github.com/qhato/ecommerce/pkg/errors"
)

// contextKey is a type for context keys
type contextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey contextKey = "user_id"
	// UserEmailKey is the context key for user email
	UserEmailKey contextKey = "user_email"
	// UserRolesKey is the context key for user roles
	UserRolesKey contextKey = "user_roles"
)

// JWTAuth creates a middleware that validates JWT tokens
func JWTAuth(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				errors.HandleHTTPError(w, errors.Unauthorized("Missing authorization header"))
				return
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				errors.HandleHTTPError(w, errors.Unauthorized("Invalid authorization header format"))
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				errors.HandleHTTPError(w, errors.Unauthorized("Invalid or expired token"))
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UserRolesKey, claims.Roles)

			// Continue with enriched context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole creates a middleware that checks if user has required role
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get roles from context
			roles, ok := r.Context().Value(UserRolesKey).([]string)
			if !ok {
				errors.HandleHTTPError(w, errors.Forbidden("User roles not found in context"))
				return
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if role == requiredRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				errors.HandleHTTPError(w, errors.Forbidden("Insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// OptionalJWTAuth is like JWTAuth but doesn't fail if no token is provided
func OptionalJWTAuth(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := parts[1]
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UserRolesKey, claims.Roles)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	userID, _ := ctx.Value(UserIDKey).(string)
	return userID
}

// GetUserEmail extracts user email from context
func GetUserEmail(ctx context.Context) string {
	email, _ := ctx.Value(UserEmailKey).(string)
	return email
}

// GetUserRoles extracts user roles from context
func GetUserRoles(ctx context.Context) []string {
	roles, _ := ctx.Value(UserRolesKey).([]string)
	return roles
}
