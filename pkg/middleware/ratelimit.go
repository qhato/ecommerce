package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/qhato/ecommerce/pkg/ratelimit"
)

// RateLimit is a middleware that enforces rate limiting
func RateLimit(limiter ratelimit.Limiter, keyFunc KeyFunc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get rate limit key
			key := keyFunc(r)

			// Check rate limit
			allowed, err := limiter.Allow(r.Context(), key)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// KeyFunc is a function that extracts a rate limit key from a request
type KeyFunc func(r *http.Request) string

// IPKeyFunc returns a key based on client IP address
func IPKeyFunc(r *http.Request) string {
	ip := GetClientIP(r)
	return fmt.Sprintf("ip:%s", ip)
}

// UserKeyFunc returns a key based on authenticated user ID
func UserKeyFunc(r *http.Request) string {
	userID, ok := GetUserID(r.Context())
	if !ok {
		// Fallback to IP if not authenticated
		return IPKeyFunc(r)
	}
	return fmt.Sprintf("user:%d", userID)
}

// EndpointKeyFunc returns a key based on endpoint and IP
func EndpointKeyFunc(r *http.Request) string {
	ip := GetClientIP(r)
	return fmt.Sprintf("endpoint:%s:%s", r.URL.Path, ip)
}

// UserEndpointKeyFunc returns a key based on user and endpoint
func UserEndpointKeyFunc(r *http.Request) string {
	userID, ok := GetUserID(r.Context())
	if !ok {
		return EndpointKeyFunc(r)
	}
	return fmt.Sprintf("user-endpoint:%d:%s", userID, r.URL.Path)
}

// GetClientIP extracts the client IP from the request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, use the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
