package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/payment/domain"
)

type PostgresWebhookEventRepository struct {
	db *sql.DB
}

func NewPostgresWebhookEventRepository(db *sql.DB) *PostgresWebhookEventRepository {
	return &PostgresWebhookEventRepository{db: db}
}

func (r *PostgresWebhookEventRepository) Create(ctx context.Context, event *domain.WebhookEvent) error {
	query := `
		INSERT INTO blc_webhook_event (id, gateway_name, event_type, event_id, payload, status,
			processed_at, error_msg, signature, ip_address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.ExecContext(ctx, query,
		event.ID, event.GatewayName, event.EventType, event.EventID, event.Payload, event.Status,
		event.ProcessedAt, event.ErrorMsg, event.Signature, event.IPAddress,
		event.CreatedAt, event.UpdatedAt,
	)
	return err
}

func (r *PostgresWebhookEventRepository) Update(ctx context.Context, event *domain.WebhookEvent) error {
	query := `
		UPDATE blc_webhook_event
		SET status = $1, processed_at = $2, error_msg = $3, updated_at = $4
		WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query,
		event.Status, event.ProcessedAt, event.ErrorMsg, event.UpdatedAt, event.ID,
	)
	return err
}

func (r *PostgresWebhookEventRepository) FindByID(ctx context.Context, id string) (*domain.WebhookEvent, error) {
	query := `
		SELECT id, gateway_name, event_type, event_id, payload, status, processed_at,
		       error_msg, signature, ip_address, created_at, updated_at
		FROM blc_webhook_event
		WHERE id = $1`

	return r.scanWebhookEvent(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresWebhookEventRepository) FindByEventID(ctx context.Context, gatewayName, eventID string) (*domain.WebhookEvent, error) {
	query := `
		SELECT id, gateway_name, event_type, event_id, payload, status, processed_at,
		       error_msg, signature, ip_address, created_at, updated_at
		FROM blc_webhook_event
		WHERE gateway_name = $1 AND event_id = $2
		ORDER BY created_at DESC
		LIMIT 1`

	return r.scanWebhookEvent(r.db.QueryRowContext(ctx, query, gatewayName, eventID))
}

func (r *PostgresWebhookEventRepository) FindPending(ctx context.Context, limit int) ([]*domain.WebhookEvent, error) {
	query := `
		SELECT id, gateway_name, event_type, event_id, payload, status, processed_at,
		       error_msg, signature, ip_address, created_at, updated_at
		FROM blc_webhook_event
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
		LIMIT $1`

	return r.queryWebhookEvents(ctx, query, limit)
}

func (r *PostgresWebhookEventRepository) FindByStatus(ctx context.Context, status domain.WebhookStatus, limit int) ([]*domain.WebhookEvent, error) {
	query := `
		SELECT id, gateway_name, event_type, event_id, payload, status, processed_at,
		       error_msg, signature, ip_address, created_at, updated_at
		FROM blc_webhook_event
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2`

	return r.queryWebhookEvents(ctx, query, status, limit)
}

func (r *PostgresWebhookEventRepository) scanWebhookEvent(row interface {
	Scan(dest ...interface{}) error
}) (*domain.WebhookEvent, error) {
	event := &domain.WebhookEvent{}
	err := row.Scan(
		&event.ID, &event.GatewayName, &event.EventType, &event.EventID, &event.Payload,
		&event.Status, &event.ProcessedAt, &event.ErrorMsg, &event.Signature, &event.IPAddress,
		&event.CreatedAt, &event.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *PostgresWebhookEventRepository) queryWebhookEvents(ctx context.Context, query string, args ...interface{}) ([]*domain.WebhookEvent, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.WebhookEvent
	for rows.Next() {
		event := &domain.WebhookEvent{}
		if err := rows.Scan(
			&event.ID, &event.GatewayName, &event.EventType, &event.EventID, &event.Payload,
			&event.Status, &event.ProcessedAt, &event.ErrorMsg, &event.Signature, &event.IPAddress,
			&event.CreatedAt, &event.UpdatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}
