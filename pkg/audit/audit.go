package audit

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	AuditActionCreate AuditAction = "CREATE"
	AuditActionUpdate AuditAction = "UPDATE"
	AuditActionDelete AuditAction = "DELETE"
	AuditActionRead   AuditAction = "READ"
	AuditActionLogin  AuditAction = "LOGIN"
	AuditActionLogout AuditAction = "LOGOUT"
)

// AuditEntry represents an audit log entry
type AuditEntry struct {
	ID           string
	EntityType   string // e.g., "Product", "Order", "Customer"
	EntityID     string
	Action       AuditAction
	UserID       *string
	Username     *string
	IPAddress    *string
	UserAgent    *string
	Changes      map[string]interface{} // Old/new values
	Metadata     map[string]interface{} // Additional context
	Timestamp    time.Time
}

// AuditLogger defines the interface for audit logging
type AuditLogger interface {
	Log(ctx context.Context, entry *AuditEntry) error
	Query(ctx context.Context, filter *AuditFilter) ([]*AuditEntry, error)
}

// AuditFilter represents filter criteria for querying audit logs
type AuditFilter struct {
	EntityType *string
	EntityID   *string
	UserID     *string
	Action     *AuditAction
	StartTime  *time.Time
	EndTime    *time.Time
	Limit      int
	Offset     int
}

// DefaultAuditLogger is a simple in-memory audit logger
type DefaultAuditLogger struct {
	entries []*AuditEntry
}

// NewDefaultAuditLogger creates a new default audit logger
func NewDefaultAuditLogger() *DefaultAuditLogger {
	return &DefaultAuditLogger{
		entries: make([]*AuditEntry, 0),
	}
}

func (l *DefaultAuditLogger) Log(ctx context.Context, entry *AuditEntry) error {
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	l.entries = append(l.entries, entry)
	return nil
}

func (l *DefaultAuditLogger) Query(ctx context.Context, filter *AuditFilter) ([]*AuditEntry, error) {
	result := make([]*AuditEntry, 0)

	for _, entry := range l.entries {
		if filter.EntityType != nil && *filter.EntityType != entry.EntityType {
			continue
		}
		if filter.EntityID != nil && *filter.EntityID != entry.EntityID {
			continue
		}
		if filter.UserID != nil && (entry.UserID == nil || *filter.UserID != *entry.UserID) {
			continue
		}
		if filter.Action != nil && *filter.Action != entry.Action {
			continue
		}
		if filter.StartTime != nil && entry.Timestamp.Before(*filter.StartTime) {
			continue
		}
		if filter.EndTime != nil && entry.Timestamp.After(*filter.EndTime) {
			continue
		}

		result = append(result, entry)
	}

	// Apply limit and offset
	start := filter.Offset
	if start > len(result) {
		return []*AuditEntry{}, nil
	}

	end := start + filter.Limit
	if filter.Limit == 0 || end > len(result) {
		end = len(result)
	}

	return result[start:end], nil
}

// AuditService provides high-level audit logging functionality
type AuditService struct {
	logger AuditLogger
}

// NewAuditService creates a new audit service
func NewAuditService(logger AuditLogger) *AuditService {
	return &AuditService{logger: logger}
}

// LogCreate logs a create action
func (s *AuditService) LogCreate(
	ctx context.Context,
	entityType, entityID string,
	userID *string,
	changes map[string]interface{},
) error {
	entry := &AuditEntry{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     AuditActionCreate,
		UserID:     userID,
		Changes:    changes,
		Timestamp:  time.Now(),
	}

	return s.logger.Log(ctx, entry)
}

// LogUpdate logs an update action
func (s *AuditService) LogUpdate(
	ctx context.Context,
	entityType, entityID string,
	userID *string,
	changes map[string]interface{},
) error {
	entry := &AuditEntry{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     AuditActionUpdate,
		UserID:     userID,
		Changes:    changes,
		Timestamp:  time.Now(),
	}

	return s.logger.Log(ctx, entry)
}

// LogDelete logs a delete action
func (s *AuditService) LogDelete(
	ctx context.Context,
	entityType, entityID string,
	userID *string,
) error {
	entry := &AuditEntry{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     AuditActionDelete,
		UserID:     userID,
		Timestamp:  time.Now(),
	}

	return s.logger.Log(ctx, entry)
}

// LogCustomAction logs a custom action
func (s *AuditService) LogCustomAction(
	ctx context.Context,
	action AuditAction,
	entityType, entityID string,
	userID *string,
	metadata map[string]interface{},
) error {
	entry := &AuditEntry{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		UserID:     userID,
		Metadata:   metadata,
		Timestamp:  time.Now(),
	}

	return s.logger.Log(ctx, entry)
}

// GetAuditTrail retrieves audit trail for an entity
func (s *AuditService) GetAuditTrail(
	ctx context.Context,
	entityType, entityID string,
) ([]*AuditEntry, error) {
	filter := &AuditFilter{
		EntityType: &entityType,
		EntityID:   &entityID,
		Limit:      100,
	}

	return s.logger.Query(ctx, filter)
}
