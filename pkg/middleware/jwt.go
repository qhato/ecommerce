package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/qhato/ecommerce/pkg/jwt"
)

// contextKey is the key type for context values
type contextKey string

const (
	// ClaimsContextKey is the context key for JWT claims
	ClaimsContextKey contextKey = "jwt_claims"
	// UserIDContextKey is the context key for user ID
	UserIDContextKey contextKey = "user_id"
)

// JWTAuth is a middleware that validates JWT tokens
func JWTAuth(jwtManager *jwt.Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check Bearer prefix
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := jwtManager.ValidateAccessToken(r.Context(), tokenString)
			if err != nil {
				switch err {
				case jwt.ErrExpiredToken:
					http.Error(w, "Token has expired", http.StatusUnauthorized)
				case jwt.ErrBlacklistedToken:
					http.Error(w, "Token has been revoked", http.StatusUnauthorized)
				default:
					http.Error(w, "Invalid token", http.StatusUnauthorized)
				}
				return
			}

			// Add claims and user ID to context
			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)

			// Call next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalJWTAuth is a middleware that validates JWT tokens but doesn't require them
func OptionalJWTAuth(jwtManager *jwt.Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No token, continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			// Check Bearer prefix
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				// Invalid format, continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := jwtManager.ValidateAccessToken(r.Context(), tokenString)
			if err == nil {
				// Valid token, add to context
				ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
				ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)
				r = r.WithContext(ctx)
			}

			// Continue regardless of validation result
			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole is a middleware that requires specific roles
func RequireRole(roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get claims from context
			claims := GetClaims(r.Context())
			if claims == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has any of the required roles
			if !claims.HasAnyRole(roles...) {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetClaims retrieves JWT claims from context
func GetClaims(ctx context.Context) *jwt.Claims {
	claims, ok := ctx.Value(ClaimsContextKey).(*jwt.Claims)
	if !ok {
		return nil
	}
	return claims
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int64)
	return userID, ok
}

// MustGetUserID retrieves user ID from context or panics
func MustGetUserID(ctx context.Context) int64 {
	userID, ok := GetUserID(ctx)
	if !ok {
		panic("user ID not found in context")
	}
	return userID
}
