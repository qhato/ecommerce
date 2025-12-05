package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/qhato/ecommerce/internal/email/domain"
)

// PostgresEmailRepository implements EmailRepository using PostgreSQL
type PostgresEmailRepository struct {
	db *sql.DB
}

// NewPostgresEmailRepository creates a new PostgreSQL email repository
func NewPostgresEmailRepository(db *sql.DB) *PostgresEmailRepository {
	return &PostgresEmailRepository{db: db}
}

// Create creates a new email
func (r *PostgresEmailRepository) Create(ctx context.Context, email *domain.Email) error {
	query := `
		INSERT INTO emails (
			type, status, priority, from_address, to_addresses, cc_addresses,
			bcc_addresses, reply_to, subject, body, html_body, template_name,
			template_data, headers, max_retries, retry_count, scheduled_at,
			order_id, customer_id, created_by, updated_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23
		) RETURNING id
	`

	// Convert template_data to JSON
	templateDataJSON, err := json.Marshal(email.TemplateData)
	if err != nil {
		return fmt.Errorf("failed to marshal template data: %w", err)
	}

	// Convert headers to JSON
	headersJSON, err := json.Marshal(email.Headers)
	if err != nil {
		return fmt.Errorf("failed to marshal headers: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query,
		email.Type,
		email.Status,
		email.Priority,
		email.From,
		pq.Array(email.To),
		pq.Array(email.CC),
		pq.Array(email.BCC),
		email.ReplyTo,
		email.Subject,
		email.Body,
		email.HTMLBody,
		email.TemplateName,
		templateDataJSON,
		headersJSON,
		email.MaxRetries,
		email.RetryCount,
		email.ScheduledAt,
		email.OrderID,
		email.CustomerID,
		email.CreatedBy,
		email.UpdatedBy,
		email.CreatedAt,
		email.UpdatedAt,
	).Scan(&email.ID)

	if err != nil {
		return fmt.Errorf("failed to create email: %w", err)
	}

	// Create attachments if any
	if len(email.Attachments) > 0 {
		if err := r.createAttachments(ctx, email.ID, email.Attachments); err != nil {
			return fmt.Errorf("failed to create attachments: %w", err)
		}
	}

	return nil
}

// Update updates an existing email
func (r *PostgresEmailRepository) Update(ctx context.Context, email *domain.Email) error {
	query := `
		UPDATE emails SET
			status = $1,
			priority = $2,
			retry_count = $3,
			sent_at = $4,
			failed_at = $5,
			error_message = $6,
			updated_at = $7,
			updated_by = $8
		WHERE id = $9
	`

	result, err := r.db.ExecContext(ctx, query,
		email.Status,
		email.Priority,
		email.RetryCount,
		email.SentAt,
		email.FailedAt,
		email.ErrorMessage,
		time.Now(),
		email.UpdatedBy,
		email.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrEmailNotFound
	}

	return nil
}

// FindByID finds an email by ID
func (r *PostgresEmailRepository) FindByID(ctx context.Context, id int64) (*domain.Email, error) {
	query := `
		SELECT
			id, type, status, priority, from_address, to_addresses, cc_addresses,
			bcc_addresses, reply_to, subject, body, html_body, template_name,
			template_data, headers, max_retries, retry_count, scheduled_at,
			sent_at, failed_at, error_message, order_id, customer_id,
			created_by, updated_by, created_at, updated_at
		FROM emails
		WHERE id = $1
	`

	email := &domain.Email{}
	var templateDataJSON, headersJSON []byte
	var toAddresses, ccAddresses, bccAddresses pq.StringArray

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&email.ID,
		&email.Type,
		&email.Status,
		&email.Priority,
		&email.From,
		&toAddresses,
		&ccAddresses,
		&bccAddresses,
		&email.ReplyTo,
		&email.Subject,
		&email.Body,
		&email.HTMLBody,
		&email.TemplateName,
		&templateDataJSON,
		&headersJSON,
		&email.MaxRetries,
		&email.RetryCount,
		&email.ScheduledAt,
		&email.SentAt,
		&email.FailedAt,
		&email.ErrorMessage,
		&email.OrderID,
		&email.CustomerID,
		&email.CreatedBy,
		&email.UpdatedBy,
		&email.CreatedAt,
		&email.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrEmailNotFound
		}
		return nil, fmt.Errorf("failed to find email: %w", err)
	}

	email.To = toAddresses
	email.CC = ccAddresses
	email.BCC = bccAddresses

	// Unmarshal template_data
	if len(templateDataJSON) > 0 {
		if err := json.Unmarshal(templateDataJSON, &email.TemplateData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal template data: %w", err)
		}
	} else {
		email.TemplateData = make(map[string]interface{})
	}

	// Unmarshal headers
	if len(headersJSON) > 0 {
		if err := json.Unmarshal(headersJSON, &email.Headers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal headers: %w", err)
		}
	} else {
		email.Headers = make(map[string]string)
	}

	// Load attachments
	attachments, err := r.findAttachmentsByEmailID(ctx, email.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load attachments: %w", err)
	}
	email.Attachments = attachments

	return email, nil
}

// FindPendingEmails finds emails that are pending to be sent
func (r *PostgresEmailRepository) FindPendingEmails(ctx context.Context, limit int) ([]*domain.Email, error) {
	query := `
		SELECT id FROM emails
		WHERE status IN ('PENDING', 'QUEUED')
		AND (scheduled_at IS NULL OR scheduled_at <= $1)
		ORDER BY priority DESC, created_at ASC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, time.Now(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending emails: %w", err)
	}
	defer rows.Close()

	var emails []*domain.Email
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan email id: %w", err)
		}

		email, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

// FindByStatus finds emails by status
func (r *PostgresEmailRepository) FindByStatus(ctx context.Context, status domain.EmailStatus, offset, limit int) ([]*domain.Email, error) {
	query := `
		SELECT id FROM emails
		WHERE status = $1
		ORDER BY created_at DESC
		OFFSET $2 LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, status, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query emails by status: %w", err)
	}
	defer rows.Close()

	var emails []*domain.Email
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan email id: %w", err)
		}

		email, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

// FindByType finds emails by type
func (r *PostgresEmailRepository) FindByType(ctx context.Context, emailType domain.EmailType, offset, limit int) ([]*domain.Email, error) {
	query := `
		SELECT id FROM emails
		WHERE type = $1
		ORDER BY created_at DESC
		OFFSET $2 LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, emailType, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query emails by type: %w", err)
	}
	defer rows.Close()

	var emails []*domain.Email
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan email id: %w", err)
		}

		email, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

// FindByOrderID finds emails associated with an order
func (r *PostgresEmailRepository) FindByOrderID(ctx context.Context, orderID int64) ([]*domain.Email, error) {
	query := `
		SELECT id FROM emails
		WHERE order_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query emails by order: %w", err)
	}
	defer rows.Close()

	var emails []*domain.Email
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan email id: %w", err)
		}

		email, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

// FindByCustomerID finds emails associated with a customer
func (r *PostgresEmailRepository) FindByCustomerID(ctx context.Context, customerID int64, offset, limit int) ([]*domain.Email, error) {
	query := `
		SELECT id FROM emails
		WHERE customer_id = $1
		ORDER BY created_at DESC
		OFFSET $2 LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, customerID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query emails by customer: %w", err)
	}
	defer rows.Close()

	var emails []*domain.Email
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan email id: %w", err)
		}

		email, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

// FindScheduledEmails finds emails scheduled to be sent
func (r *PostgresEmailRepository) FindScheduledEmails(ctx context.Context) ([]*domain.Email, error) {
	query := `
		SELECT id FROM emails
		WHERE status = 'PENDING'
		AND scheduled_at IS NOT NULL
		AND scheduled_at <= $1
		ORDER BY scheduled_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to query scheduled emails: %w", err)
	}
	defer rows.Close()

	var emails []*domain.Email
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan email id: %w", err)
		}

		email, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

// FindFailedEmailsForRetry finds emails that failed and can be retried
func (r *PostgresEmailRepository) FindFailedEmailsForRetry(ctx context.Context, limit int) ([]*domain.Email, error) {
	query := `
		SELECT id FROM emails
		WHERE status = 'FAILED'
		AND retry_count < max_retries
		ORDER BY priority DESC, failed_at ASC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query failed emails: %w", err)
	}
	defer rows.Close()

	var emails []*domain.Email
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan email id: %w", err)
		}

		email, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

// Delete deletes an email by ID
func (r *PostgresEmailRepository) Delete(ctx context.Context, id int64) error {
	// Delete attachments first
	if _, err := r.db.ExecContext(ctx, "DELETE FROM email_attachments WHERE email_id = $1", id); err != nil {
		return fmt.Errorf("failed to delete attachments: %w", err)
	}

	// Delete email
	result, err := r.db.ExecContext(ctx, "DELETE FROM emails WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete email: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrEmailNotFound
	}

	return nil
}

// Count counts total emails
func (r *PostgresEmailRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM emails").Scan(&count)
	return count, err
}

// CountByStatus counts emails by status
func (r *PostgresEmailRepository) CountByStatus(ctx context.Context, status domain.EmailStatus) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM emails WHERE status = $1", status).Scan(&count)
	return count, err
}

// createAttachments creates email attachments
func (r *PostgresEmailRepository) createAttachments(ctx context.Context, emailID int64, attachments []domain.EmailAttachment) error {
	query := `
		INSERT INTO email_attachments (
			email_id, filename, content_type, content, size, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, att := range attachments {
		_, err := r.db.ExecContext(ctx, query,
			emailID,
			att.Filename,
			att.ContentType,
			att.Content,
			att.Size,
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to create attachment: %w", err)
		}
	}

	return nil
}

// findAttachmentsByEmailID finds attachments for an email
func (r *PostgresEmailRepository) findAttachmentsByEmailID(ctx context.Context, emailID int64) ([]domain.EmailAttachment, error) {
	query := `
		SELECT id, email_id, filename, content_type, content, size, created_at
		FROM email_attachments
		WHERE email_id = $1
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query, emailID)
	if err != nil {
		return nil, fmt.Errorf("failed to query attachments: %w", err)
	}
	defer rows.Close()

	var attachments []domain.EmailAttachment
	for rows.Next() {
		var att domain.EmailAttachment
		if err := rows.Scan(
			&att.ID,
			&att.EmailID,
			&att.Filename,
			&att.ContentType,
			&att.Content,
			&att.Size,
			&att.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan attachment: %w", err)
		}
		attachments = append(attachments, att)
	}

	return attachments, nil
}
