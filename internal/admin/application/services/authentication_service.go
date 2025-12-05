package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/qhato/ecommerce/internal/admin/domain"
)

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds
}

// AuthenticationService handles authentication operations
type AuthenticationService struct {
	userRepo      domain.AdminUserRepository
	auditLogRepo  domain.AuditLogRepository
	tokenService  TokenService
}

// TokenService defines token generation and validation interface
type TokenService interface {
	GenerateAccessToken(userID int64, username string, roles []string) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
	ValidateAccessToken(token string) (*TokenClaims, error)
	ValidateRefreshToken(token string) (int64, error)
}

// TokenClaims represents decoded token claims
type TokenClaims struct {
	UserID   int64
	Username string
	Roles    []string
	IssuedAt int64
	ExpiresAt int64
}

// NewAuthenticationService creates a new authentication service
func NewAuthenticationService(
	userRepo domain.AdminUserRepository,
	auditLogRepo domain.AuditLogRepository,
	tokenService TokenService,
) *AuthenticationService {
	return &AuthenticationService{
		userRepo:     userRepo,
		auditLogRepo: auditLogRepo,
		tokenService: tokenService,
	}
}

// Login authenticates a user and returns token pair
func (s *AuthenticationService) Login(ctx context.Context, username, password, ipAddress, userAgent string) (*TokenPair, *domain.AdminUser, error) {
	// Find user
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		s.logFailedLogin(ctx, username, "user not found", ipAddress, userAgent)
		return nil, nil, domain.ErrInvalidCredentials
	}

	if user == nil {
		s.logFailedLogin(ctx, username, "user not found", ipAddress, userAgent)
		return nil, nil, domain.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		s.logFailedLogin(ctx, username, "user not active", ipAddress, userAgent)
		return nil, nil, domain.ErrUserNotActive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.logFailedLogin(ctx, username, "invalid password", ipAddress, userAgent)
		return nil, nil, domain.ErrInvalidCredentials
	}

	// Get user roles
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	user.Roles = roles

	// Generate tokens
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Username, roleNames)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Record login
	user.RecordLogin()
	if err := s.userRepo.Update(ctx, user); err != nil {
		// Log but don't fail
		fmt.Printf("Failed to update last login time: %v\n", err)
	}

	// Create audit log
	auditLog := domain.NewAuditLog(
		user.ID,
		user.Username,
		domain.AuditActionUserLogin,
		"USER",
		fmt.Sprintf("%d", user.ID),
		fmt.Sprintf("User %s logged in successfully", user.Username),
		domain.AuditSeverityInfo,
		ipAddress,
		userAgent,
	)
	s.auditLogRepo.Create(ctx, auditLog)

	tokenPair := &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
	}

	return tokenPair, user, nil
}

// RefreshToken generates a new access token using refresh token
func (s *AuthenticationService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Validate refresh token
	userID, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	if !user.IsActive {
		return nil, domain.ErrUserNotActive
	}

	// Get user roles
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// Generate new tokens
	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Username, roleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokenPair := &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600, // 1 hour
	}

	return tokenPair, nil
}

// ValidateToken validates an access token and returns claims
func (s *AuthenticationService) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	claims, err := s.tokenService.ValidateAccessToken(token)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Check if token is expired
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, domain.ErrTokenExpired
	}

	return claims, nil
}

// HashPassword hashes a password using bcrypt
func (s *AuthenticationService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// GenerateRandomPassword generates a random password
func (s *AuthenticationService) GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		length = 8
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// logFailedLogin logs a failed login attempt
func (s *AuthenticationService) logFailedLogin(ctx context.Context, username, reason, ipAddress, userAgent string) {
	auditLog := domain.NewAuditLog(
		0, // No user ID for failed login
		username,
		domain.AuditActionUserLoginFailed,
		"USER",
		"",
		fmt.Sprintf("Failed login attempt for user %s: %s", username, reason),
		domain.AuditSeverityWarning,
		ipAddress,
		userAgent,
	)
	auditLog.MarkAsFailure(reason)
	s.auditLogRepo.Create(ctx, auditLog)
}
