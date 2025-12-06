package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/cms/domain"
)

type PostgresMediaRepository struct {
	db *sql.DB
}

func NewPostgresMediaRepository(db *sql.DB) *PostgresMediaRepository {
	return &PostgresMediaRepository{db: db}
}

func (r *PostgresMediaRepository) Create(ctx context.Context, media *domain.Media) error {
	query := `INSERT INTO blc_cms_media (
		file_name, file_path, mime_type, file_size, title, alt_text,
		caption, uploaded_by, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		media.FileName, media.FilePath, media.MimeType, media.FileSize,
		media.Title, media.AltText, media.Caption, media.UploadedBy,
		media.CreatedAt, media.UpdatedAt,
	).Scan(&media.ID)
}

func (r *PostgresMediaRepository) Update(ctx context.Context, media *domain.Media) error {
	query := `UPDATE blc_cms_media SET
		title = $1, alt_text = $2, caption = $3, updated_at = $4
	WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query,
		media.Title, media.AltText, media.Caption, media.UpdatedAt, media.ID)
	return err
}

func (r *PostgresMediaRepository) FindByID(ctx context.Context, id int64) (*domain.Media, error) {
	query := `SELECT id, file_name, file_path, mime_type, file_size, title,
		alt_text, caption, uploaded_by, created_at, updated_at
	FROM blc_cms_media WHERE id = $1`

	media := &domain.Media{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&media.ID, &media.FileName, &media.FilePath, &media.MimeType,
		&media.FileSize, &media.Title, &media.AltText, &media.Caption,
		&media.UploadedBy, &media.CreatedAt, &media.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return media, err
}

func (r *PostgresMediaRepository) FindAll(ctx context.Context, mimeType string, limit int) ([]*domain.Media, error) {
	query := `SELECT id, file_name, file_path, mime_type, file_size, title,
		alt_text, caption, uploaded_by, created_at, updated_at
	FROM blc_cms_media`

	args := []interface{}{}
	if mimeType != "" {
		query += " WHERE mime_type LIKE $1"
		args = append(args, mimeType+"%")
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		if len(args) > 0 {
			query += " LIMIT $2"
		} else {
			query += " LIMIT $1"
		}
		args = append(args, limit)
	}

	return r.queryMedias(ctx, query, args...)
}

func (r *PostgresMediaRepository) FindByUploader(ctx context.Context, uploaderID int64, limit int) ([]*domain.Media, error) {
	query := `SELECT id, file_name, file_path, mime_type, file_size, title,
		alt_text, caption, uploaded_by, created_at, updated_at
	FROM blc_cms_media WHERE uploaded_by = $1
	ORDER BY created_at DESC LIMIT $2`

	return r.queryMedias(ctx, query, uploaderID, limit)
}

func (r *PostgresMediaRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM blc_cms_media WHERE id = $1`, id)
	return err
}

func (r *PostgresMediaRepository) IsMediaInUse(ctx context.Context, mediaID int64) (bool, error) {
	// Check if media is referenced in any content's featured_image or body
	query := `SELECT EXISTS(
		SELECT 1 FROM blc_cms_content
		WHERE featured_image LIKE $1 OR body LIKE $1
	)`

	mediaPath := "%/media/" + string(mediaID) + "%"
	var exists bool
	err := r.db.QueryRowContext(ctx, query, mediaPath).Scan(&exists)
	return exists, err
}

// Helper method
func (r *PostgresMediaRepository) queryMedias(ctx context.Context, query string, args ...interface{}) ([]*domain.Media, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	medias := make([]*domain.Media, 0)
	for rows.Next() {
		media := &domain.Media{}
		if err := rows.Scan(
			&media.ID, &media.FileName, &media.FilePath, &media.MimeType,
			&media.FileSize, &media.Title, &media.AltText, &media.Caption,
			&media.UploadedBy, &media.CreatedAt, &media.UpdatedAt,
		); err != nil {
			return nil, err
		}
		medias = append(medias, media)
	}

	return medias, nil
}
