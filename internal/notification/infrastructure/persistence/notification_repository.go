package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/qhato/ecommerce/internal/notification/domain"
)

type PostgresNotificationRepository struct {
	db *sql.DB
}

func NewPostgresNotificationRepository(db *sql.DB) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{db: db}
}

func (r *PostgresNotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	templateDataJSON, err := json.Marshal(notification.TemplateData)
	if err != nil {
		return fmt.Errorf("failed to marshal template data: %w", err)
	}

	query := `INSERT INTO blc_notification (
		type, status, priority, recipient_id, recipient_email, recipient_phone,
		subject, body, template_id, template_data, scheduled_for, sent_at, failed_at,
		error_msg, retry_count, max_retries, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		notification.Type, notification.Status, notification.Priority,
		notification.RecipientID, notification.RecipientEmail, notification.RecipientPhone,
		notification.Subject, notification.Body, notification.TemplateID, templateDataJSON,
		notification.ScheduledFor, notification.SentAt, notification.FailedAt,
		notification.ErrorMsg, notification.RetryCount, notification.MaxRetries,
		notification.CreatedAt, notification.UpdatedAt,
	).Scan(&notification.ID)
}

func (r *PostgresNotificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	templateDataJSON, err := json.Marshal(notification.TemplateData)
	if err != nil {
		return fmt.Errorf("failed to marshal template data: %w", err)
	}

	query := `UPDATE blc_notification SET
		status = $1, priority = $2, recipient_email = $3, recipient_phone = $4,
		subject = $5, body = $6, template_id = $7, template_data = $8,
		scheduled_for = $9, sent_at = $10, failed_at = $11, error_msg = $12,
		retry_count = $13, max_retries = $14, updated_at = $15
	WHERE id = $16`

	_, err = r.db.ExecContext(ctx, query,
		notification.Status, notification.Priority,
		notification.RecipientEmail, notification.RecipientPhone,
		notification.Subject, notification.Body, notification.TemplateID, templateDataJSON,
		notification.ScheduledFor, notification.SentAt, notification.FailedAt,
		notification.ErrorMsg, notification.RetryCount, notification.MaxRetries,
		notification.UpdatedAt, notification.ID,
	)
	return err
}

func (r *PostgresNotificationRepository) FindByID(ctx context.Context, id int64) (*domain.Notification, error) {
	query := `SELECT id, type, status, priority, recipient_id, recipient_email, recipient_phone,
		subject, body, template_id, template_data, scheduled_for, sent_at, failed_at,
		error_msg, retry_count, max_retries, created_at, updated_at
	FROM blc_notification WHERE id = $1`

	return r.scanNotification(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresNotificationRepository) FindByRecipientID(ctx context.Context, recipientID string, limit int) ([]*domain.Notification, error) {
	query := `SELECT id, type, status, priority, recipient_id, recipient_email, recipient_phone,
		subject, body, template_id, template_data, scheduled_for, sent_at, failed_at,
		error_msg, retry_count, max_retries, created_at, updated_at
	FROM blc_notification WHERE recipient_id = $1 ORDER BY created_at DESC LIMIT $2`

	return r.queryNotifications(ctx, query, recipientID, limit)
}

func (r *PostgresNotificationRepository) FindByStatus(ctx context.Context, status domain.NotificationStatus, limit int) ([]*domain.Notification, error) {
	query := `SELECT id, type, status, priority, recipient_id, recipient_email, recipient_phone,
		subject, body, template_id, template_data, scheduled_for, sent_at, failed_at,
		error_msg, retry_count, max_retries, created_at, updated_at
	FROM blc_notification WHERE status = $1 ORDER BY created_at DESC LIMIT $2`

	return r.queryNotifications(ctx, query, status, limit)
}

func (r *PostgresNotificationRepository) FindPending(ctx context.Context, limit int) ([]*domain.Notification, error) {
	query := `SELECT id, type, status, priority, recipient_id, recipient_email, recipient_phone,
		subject, body, template_id, template_data, scheduled_for, sent_at, failed_at,
		error_msg, retry_count, max_retries, created_at, updated_at
	FROM blc_notification WHERE status = $1 AND (scheduled_for IS NULL OR scheduled_for <= NOW())
	ORDER BY priority DESC, created_at ASC LIMIT $2`

	return r.queryNotifications(ctx, query, domain.NotificationStatusPending, limit)
}

func (r *PostgresNotificationRepository) FindScheduled(ctx context.Context, limit int) ([]*domain.Notification, error) {
	query := `SELECT id, type, status, priority, recipient_id, recipient_email, recipient_phone,
		subject, body, template_id, template_data, scheduled_for, sent_at, failed_at,
		error_msg, retry_count, max_retries, created_at, updated_at
	FROM blc_notification WHERE status = $1 AND scheduled_for > NOW()
	ORDER BY scheduled_for ASC LIMIT $2`

	return r.queryNotifications(ctx, query, domain.NotificationStatusPending, limit)
}

func (r *PostgresNotificationRepository) FindFailed(ctx context.Context, limit int) ([]*domain.Notification, error) {
	query := `SELECT id, type, status, priority, recipient_id, recipient_email, recipient_phone,
		subject, body, template_id, template_data, scheduled_for, sent_at, failed_at,
		error_msg, retry_count, max_retries, created_at, updated_at
	FROM blc_notification WHERE status = $1 ORDER BY failed_at DESC LIMIT $2`

	return r.queryNotifications(ctx, query, domain.NotificationStatusFailed, limit)
}

func (r *PostgresNotificationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_notification WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresNotificationRepository) scanNotification(row interface {
	Scan(dest ...interface{}) error
}) (*domain.Notification, error) {
	notification := &domain.Notification{}
	var templateDataJSON []byte

	err := row.Scan(
		&notification.ID, &notification.Type, &notification.Status, &notification.Priority,
		&notification.RecipientID, &notification.RecipientEmail, &notification.RecipientPhone,
		&notification.Subject, &notification.Body, &notification.TemplateID, &templateDataJSON,
		&notification.ScheduledFor, &notification.SentAt, &notification.FailedAt,
		&notification.ErrorMsg, &notification.RetryCount, &notification.MaxRetries,
		&notification.CreatedAt, &notification.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(templateDataJSON, &notification.TemplateData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal template data: %w", err)
	}

	return notification, nil
}

func (r *PostgresNotificationRepository) queryNotifications(ctx context.Context, query string, args ...interface{}) ([]*domain.Notification, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]*domain.Notification, 0)
	for rows.Next() {
		notification := &domain.Notification{}
		var templateDataJSON []byte

		if err := rows.Scan(
			&notification.ID, &notification.Type, &notification.Status, &notification.Priority,
			&notification.RecipientID, &notification.RecipientEmail, &notification.RecipientPhone,
			&notification.Subject, &notification.Body, &notification.TemplateID, &templateDataJSON,
			&notification.ScheduledFor, &notification.SentAt, &notification.FailedAt,
			&notification.ErrorMsg, &notification.RetryCount, &notification.MaxRetries,
			&notification.CreatedAt, &notification.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(templateDataJSON, &notification.TemplateData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal template data: %w", err)
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}
