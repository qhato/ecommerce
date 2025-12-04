package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrTokenNotFound    = errors.New("token not found")
	ErrInvalidClaims    = errors.New("invalid token claims")
	ErrBlacklistedToken = errors.New("token has been revoked")
)

// Claims represents JWT claims
type Claims struct {
	UserID      int64    `json:"user_id"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	TokenType   string   `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// Manager handles JWT operations
type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
	issuer        string
	blacklist     Blacklist
}

// Config contains JWT configuration
type Config struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	Issuer        string
}

// Blacklist is an interface for token blacklisting
type Blacklist interface {
	Add(ctx context.Context, tokenID string, expiry time.Duration) error
	IsBlacklisted(ctx context.Context, tokenID string) (bool, error)
}

// NewManager creates a new JWT manager
func NewManager(cfg Config, blacklist Blacklist) *Manager {
	return &Manager{
		accessSecret:  []byte(cfg.AccessSecret),
		refreshSecret: []byte(cfg.RefreshSecret),
		accessTTL:     cfg.AccessTTL,
		refreshTTL:    cfg.RefreshTTL,
		issuer:        cfg.Issuer,
		blacklist:     blacklist,
	}
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// GenerateTokenPair generates both access and refresh tokens
func (m *Manager) GenerateTokenPair(userID int64, email string, roles []string) (*TokenPair, error) {
	now := time.Now()
	tokenID := fmt.Sprintf("%d-%d", userID, now.Unix())

	// Generate access token
	accessClaims := &Claims{
		UserID:    userID,
		Email:     email,
		Roles:     roles,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID + "-access",
			Issuer:    m.issuer,
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(m.accessSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshClaims := &Claims{
		UserID:    userID,
		Email:     email,
		Roles:     roles,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID + "-refresh",
			Issuer:    m.issuer,
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTTL)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(m.refreshSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    now.Add(m.accessTTL),
		TokenType:    "Bearer",
	}, nil
}

// ValidateAccessToken validates an access token and returns claims
func (m *Manager) ValidateAccessToken(ctx context.Context, tokenString string) (*Claims, error) {
	return m.validateToken(ctx, tokenString, m.accessSecret, "access")
}

// ValidateRefreshToken validates a refresh token and returns claims
func (m *Manager) ValidateRefreshToken(ctx context.Context, tokenString string) (*Claims, error) {
	return m.validateToken(ctx, tokenString, m.refreshSecret, "refresh")
}

// validateToken validates a token with the given secret
func (m *Manager) validateToken(ctx context.Context, tokenString string, secret []byte, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	// Verify token type
	if claims.TokenType != expectedType {
		return nil, ErrInvalidClaims
	}

	// Check if token is blacklisted
	if m.blacklist != nil {
		blacklisted, err := m.blacklist.IsBlacklisted(ctx, claims.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check blacklist: %w", err)
		}
		if blacklisted {
			return nil, ErrBlacklistedToken
		}
	}

	return claims, nil
}

// RefreshTokens generates new tokens using a refresh token
func (m *Manager) RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := m.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// Generate new token pair
	return m.GenerateTokenPair(claims.UserID, claims.Email, claims.Roles)
}

// RevokeToken adds a token to the blacklist
func (m *Manager) RevokeToken(ctx context.Context, tokenString string) error {
	if m.blacklist == nil {
		return errors.New("blacklist not configured")
	}

	// Parse token to get claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.accessSecret, nil
	})

	if err != nil {
		// Try refresh secret
		token, err = jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return m.refreshSecret, nil
		})
		if err != nil {
			return ErrInvalidToken
		}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return ErrInvalidClaims
	}

	// Calculate time until expiry
	expiry := time.Until(claims.ExpiresAt.Time)
	if expiry <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}

	// Add to blacklist
	return m.blacklist.Add(ctx, claims.ID, expiry)
}

// ExtractUserID extracts user ID from claims
func (c *Claims) ExtractUserID() int64 {
	return c.UserID
}

// HasRole checks if claims contain a specific role
func (c *Claims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole checks if claims contain any of the specified roles
func (c *Claims) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if c.HasRole(role) {
			return true
		}
	}
	return false
}