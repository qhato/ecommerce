package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/qhato/ecommerce/internal/admin/domain"
)

type PostgresAuditLogRepository struct {
	db *sql.DB
}

func NewPostgresAuditLogRepository(db *sql.DB) *PostgresAuditLogRepository {
	return &PostgresAuditLogRepository{db: db}
}

func (r *PostgresAuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	// Convert Details map to JSON
	var detailsJSON []byte
	var err error
	if log.Details != nil {
		detailsJSON, err = json.Marshal(log.Details)
		if err != nil {
			return fmt.Errorf("failed to marshal details: %w", err)
		}
	}

	query := `INSERT INTO blc_admin_audit_log (
		user_id, username, action, resource, resource_id, description,
		severity, ip_address, user_agent, details, success, error_msg, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		log.UserID, log.Username, log.Action, log.Resource, log.ResourceID,
		log.Description, log.Severity, log.IPAddress, log.UserAgent,
		detailsJSON, log.Success, log.ErrorMsg, log.CreatedAt,
	).Scan(&log.ID)
}

func (r *PostgresAuditLogRepository) FindByID(ctx context.Context, id int64) (*domain.AuditLog, error) {
	query := `SELECT id, user_id, username, action, resource, resource_id, description,
		severity, ip_address, user_agent, details, success, error_msg, created_at
		FROM blc_admin_audit_log WHERE id = $1`

	log := &domain.AuditLog{}
	var detailsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID, &log.UserID, &log.Username, &log.Action, &log.Resource, &log.ResourceID,
		&log.Description, &log.Severity, &log.IPAddress, &log.UserAgent,
		&detailsJSON, &log.Success, &log.ErrorMsg, &log.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal details JSON
	if detailsJSON != nil {
		if err := json.Unmarshal(detailsJSON, &log.Details); err != nil {
			return nil, fmt.Errorf("failed to unmarshal details: %w", err)
		}
	}

	return log, nil
}

func (r *PostgresAuditLogRepository) FindByUserID(ctx context.Context, userID int64, limit int) ([]*domain.AuditLog, error) {
	query := `SELECT id, user_id, username, action, resource, resource_id, description,
		severity, ip_address, user_agent, details, success, error_msg, created_at
		FROM blc_admin_audit_log WHERE user_id = $1
		ORDER BY created_at DESC LIMIT $2`

	return r.queryAuditLogs(ctx, query, userID, limit)
}

func (r *PostgresAuditLogRepository) FindRecent(ctx context.Context, limit int) ([]*domain.AuditLog, error) {
	query := `SELECT id, user_id, username, action, resource, resource_id, description,
		severity, ip_address, user_agent, details, success, error_msg, created_at
		FROM blc_admin_audit_log
		ORDER BY created_at DESC LIMIT $1`

	return r.queryAuditLogs(ctx, query, limit)
}

func (r *PostgresAuditLogRepository) FindSecurityEvents(ctx context.Context, limit int) ([]*domain.AuditLog, error) {
	query := `SELECT id, user_id, username, action, resource, resource_id, description,
		severity, ip_address, user_agent, details, success, error_msg, created_at
		FROM blc_admin_audit_log
		WHERE severity IN ('HIGH', 'CRITICAL')
		OR action IN ('LOGIN_FAILED', 'PASSWORD_CHANGED', 'PERMISSION_GRANTED', 'PERMISSION_REVOKED')
		ORDER BY created_at DESC LIMIT $1`

	return r.queryAuditLogs(ctx, query, limit)
}

func (r *PostgresAuditLogRepository) FindFailedLogins(ctx context.Context, username string, limit int) ([]*domain.AuditLog, error) {
	query := `SELECT id, user_id, username, action, resource, resource_id, description,
		severity, ip_address, user_agent, details, success, error_msg, created_at
		FROM blc_admin_audit_log
		WHERE username = $1 AND action = 'LOGIN_FAILED'
		ORDER BY created_at DESC LIMIT $2`

	return r.queryAuditLogs(ctx, query, username, limit)
}

func (r *PostgresAuditLogRepository) FindByAction(ctx context.Context, action domain.AuditAction, limit int) ([]*domain.AuditLog, error) {
	query := `SELECT id, user_id, username, action, resource, resource_id, description,
		severity, ip_address, user_agent, details, success, error_msg, created_at
		FROM blc_admin_audit_log WHERE action = $1
		ORDER BY created_at DESC LIMIT $2`

	return r.queryAuditLogs(ctx, query, action, limit)
}

func (r *PostgresAuditLogRepository) FindByResource(ctx context.Context, resource, resourceID string, limit int) ([]*domain.AuditLog, error) {
	query := `SELECT id, user_id, username, action, resource, resource_id, description,
		severity, ip_address, user_agent, details, success, error_msg, created_at
		FROM blc_admin_audit_log WHERE resource = $1 AND resource_id = $2
		ORDER BY created_at DESC LIMIT $3`

	return r.queryAuditLogs(ctx, query, resource, resourceID, limit)
}

func (r *PostgresAuditLogRepository) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	query := `DELETE FROM blc_admin_audit_log WHERE created_at < NOW() - INTERVAL '1 day' * $1`
	result, err := r.db.ExecContext(ctx, query, days)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Helper method to query and scan audit logs
func (r *PostgresAuditLogRepository) queryAuditLogs(ctx context.Context, query string, args ...interface{}) ([]*domain.AuditLog, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]*domain.AuditLog, 0)
	for rows.Next() {
		log := &domain.AuditLog{}
		var detailsJSON []byte

		if err := rows.Scan(
			&log.ID, &log.UserID, &log.Username, &log.Action, &log.Resource, &log.ResourceID,
			&log.Description, &log.Severity, &log.IPAddress, &log.UserAgent,
			&detailsJSON, &log.Success, &log.ErrorMsg, &log.CreatedAt,
		); err != nil {
			return nil, err
		}

		// Unmarshal details JSON
		if detailsJSON != nil {
			if err := json.Unmarshal(detailsJSON, &log.Details); err != nil {
				return nil, fmt.Errorf("failed to unmarshal details: %w", err)
			}
		}

		logs = append(logs, log)
	}

	return logs, nil
}
