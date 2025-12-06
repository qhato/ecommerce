package persistence

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/qhato/ecommerce/internal/promotionmsg/domain"
)

type PostgresPromotionMessageRepository struct {
	db *sql.DB
}

func NewPostgresPromotionMessageRepository(db *sql.DB) *PostgresPromotionMessageRepository {
	return &PostgresPromotionMessageRepository{db: db}
}

func (r *PostgresPromotionMessageRepository) Create(ctx context.Context, message *domain.PromotionMessage) error {
	rulesJSON, _ := json.Marshal(message.Rules)
	triggersJSON, _ := json.Marshal(message.Triggers)
	placementsJSON, _ := json.Marshal(message.Placements)
	metadataJSON, _ := json.Marshal(message.Metadata)

	query := `INSERT INTO blc_promotion_message (
		name, type, priority, status, message, description, rules, triggers, placements,
		start_date, end_date, max_views, view_count, click_count, metadata,
		created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		message.Name, message.Type, message.Priority, message.Status, message.Message,
		message.Description, rulesJSON, triggersJSON, placementsJSON,
		message.StartDate, message.EndDate, message.MaxViews, message.ViewCount,
		message.ClickCount, metadataJSON, message.CreatedAt, message.UpdatedAt,
	).Scan(&message.ID)
}

func (r *PostgresPromotionMessageRepository) Update(ctx context.Context, message *domain.PromotionMessage) error {
	rulesJSON, _ := json.Marshal(message.Rules)
	triggersJSON, _ := json.Marshal(message.Triggers)
	placementsJSON, _ := json.Marshal(message.Placements)
	metadataJSON, _ := json.Marshal(message.Metadata)

	query := `UPDATE blc_promotion_message SET
		name = $1, priority = $2, status = $3, message = $4, description = $5,
		rules = $6, triggers = $7, placements = $8, start_date = $9, end_date = $10,
		max_views = $11, view_count = $12, click_count = $13, metadata = $14, updated_at = $15
	WHERE id = $16`

	_, err := r.db.ExecContext(ctx, query,
		message.Name, message.Priority, message.Status, message.Message, message.Description,
		rulesJSON, triggersJSON, placementsJSON, message.StartDate, message.EndDate,
		message.MaxViews, message.ViewCount, message.ClickCount, metadataJSON,
		message.UpdatedAt, message.ID,
	)
	return err
}

func (r *PostgresPromotionMessageRepository) FindByID(ctx context.Context, id int64) (*domain.PromotionMessage, error) {
	query := `SELECT id, name, type, priority, status, message, description, rules, triggers,
		placements, start_date, end_date, max_views, view_count, click_count, metadata,
		created_at, updated_at
	FROM blc_promotion_message WHERE id = $1`

	return r.scanMessage(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresPromotionMessageRepository) FindByType(ctx context.Context, messageType domain.MessageType) ([]*domain.PromotionMessage, error) {
	query := `SELECT id, name, type, priority, status, message, description, rules, triggers,
		placements, start_date, end_date, max_views, view_count, click_count, metadata,
		created_at, updated_at
	FROM blc_promotion_message WHERE type = $1 ORDER BY priority DESC, created_at DESC`

	return r.queryMessages(ctx, query, messageType)
}

func (r *PostgresPromotionMessageRepository) FindByStatus(ctx context.Context, status domain.MessageStatus) ([]*domain.PromotionMessage, error) {
	query := `SELECT id, name, type, priority, status, message, description, rules, triggers,
		placements, start_date, end_date, max_views, view_count, click_count, metadata,
		created_at, updated_at
	FROM blc_promotion_message WHERE status = $1 ORDER BY priority DESC, created_at DESC`

	return r.queryMessages(ctx, query, status)
}

func (r *PostgresPromotionMessageRepository) FindActive(ctx context.Context, limit int) ([]*domain.PromotionMessage, error) {
	query := `SELECT id, name, type, priority, status, message, description, rules, triggers,
		placements, start_date, end_date, max_views, view_count, click_count, metadata,
		created_at, updated_at
	FROM blc_promotion_message
	WHERE status = 'ACTIVE'
	  AND (start_date IS NULL OR start_date <= NOW())
	  AND (end_date IS NULL OR end_date >= NOW())
	  AND (max_views IS NULL OR view_count < max_views)
	ORDER BY priority DESC, created_at DESC
	LIMIT $1`

	return r.queryMessages(ctx, query, limit)
}

func (r *PostgresPromotionMessageRepository) FindByPlacement(ctx context.Context, placement string) ([]*domain.PromotionMessage, error) {
	query := `SELECT id, name, type, priority, status, message, description, rules, triggers,
		placements, start_date, end_date, max_views, view_count, click_count, metadata,
		created_at, updated_at
	FROM blc_promotion_message
	WHERE status = 'ACTIVE'
	  AND placements @> $1::jsonb
	  AND (start_date IS NULL OR start_date <= NOW())
	  AND (end_date IS NULL OR end_date >= NOW())
	  AND (max_views IS NULL OR view_count < max_views)
	ORDER BY priority DESC, created_at DESC`

	placementJSON, _ := json.Marshal([]string{placement})
	return r.queryMessages(ctx, query, string(placementJSON))
}

func (r *PostgresPromotionMessageRepository) FindByEvent(ctx context.Context, event string) ([]*domain.PromotionMessage, error) {
	query := `SELECT id, name, type, priority, status, message, description, rules, triggers,
		placements, start_date, end_date, max_views, view_count, click_count, metadata,
		created_at, updated_at
	FROM blc_promotion_message
	WHERE status = 'ACTIVE'
	  AND triggers::text LIKE $1
	  AND (start_date IS NULL OR start_date <= NOW())
	  AND (end_date IS NULL OR end_date >= NOW())
	  AND (max_views IS NULL OR view_count < max_views)
	ORDER BY priority DESC, created_at DESC`

	return r.queryMessages(ctx, query, "%"+event+"%")
}

func (r *PostgresPromotionMessageRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_promotion_message WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresPromotionMessageRepository) scanMessage(row interface {
	Scan(dest ...interface{}) error
}) (*domain.PromotionMessage, error) {
	message := &domain.PromotionMessage{}
	var rulesJSON, triggersJSON, placementsJSON, metadataJSON []byte

	err := row.Scan(
		&message.ID, &message.Name, &message.Type, &message.Priority, &message.Status,
		&message.Message, &message.Description, &rulesJSON, &triggersJSON,
		&placementsJSON, &message.StartDate, &message.EndDate, &message.MaxViews,
		&message.ViewCount, &message.ClickCount, &metadataJSON,
		&message.CreatedAt, &message.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(rulesJSON, &message.Rules); err != nil {
		message.Rules = make([]domain.MessageRule, 0)
	}
	if err := json.Unmarshal(triggersJSON, &message.Triggers); err != nil {
		message.Triggers = make([]domain.MessageTrigger, 0)
	}
	if err := json.Unmarshal(placementsJSON, &message.Placements); err != nil {
		message.Placements = make([]string, 0)
	}
	if err := json.Unmarshal(metadataJSON, &message.Metadata); err != nil {
		message.Metadata = make(map[string]interface{})
	}

	return message, nil
}

func (r *PostgresPromotionMessageRepository) queryMessages(ctx context.Context, query string, args ...interface{}) ([]*domain.PromotionMessage, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*domain.PromotionMessage, 0)
	for rows.Next() {
		message := &domain.PromotionMessage{}
		var rulesJSON, triggersJSON, placementsJSON, metadataJSON []byte

		if err := rows.Scan(
			&message.ID, &message.Name, &message.Type, &message.Priority, &message.Status,
			&message.Message, &message.Description, &rulesJSON, &triggersJSON,
			&placementsJSON, &message.StartDate, &message.EndDate, &message.MaxViews,
			&message.ViewCount, &message.ClickCount, &metadataJSON,
			&message.CreatedAt, &message.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(rulesJSON, &message.Rules); err != nil {
			message.Rules = make([]domain.MessageRule, 0)
		}
		if err := json.Unmarshal(triggersJSON, &message.Triggers); err != nil {
			message.Triggers = make([]domain.MessageTrigger, 0)
		}
		if err := json.Unmarshal(placementsJSON, &message.Placements); err != nil {
			message.Placements = make([]string, 0)
		}
		if err := json.Unmarshal(metadataJSON, &message.Metadata); err != nil {
			message.Metadata = make(map[string]interface{})
		}

		messages = append(messages, message)
	}

	return messages, nil
}
