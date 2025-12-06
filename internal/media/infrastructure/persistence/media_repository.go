package persistence

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
	"github.com/qhato/ecommerce/internal/media/domain"
)

type PostgresMediaRepository struct {
	db *sql.DB
}

func NewPostgresMediaRepository(db *sql.DB) *PostgresMediaRepository {
	return &PostgresMediaRepository{db: db}
}

func (r *PostgresMediaRepository) Create(ctx context.Context, media *domain.Media) error {
	tagsJSON, _ := json.Marshal(media.Tags)
	metadataJSON, _ := json.Marshal(media.Metadata)

	query := `
		INSERT INTO blc_media (id, name, title, description, media_type, status, mime_type,
			file_size, file_path, url, thumbnail_url, width, height, duration, tags, metadata,
			uploaded_by, entity_type, entity_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`

	_, err := r.db.ExecContext(ctx, query,
		media.ID, media.Name, media.Title, media.Description, media.MediaType, media.Status,
		media.MimeType, media.FileSize, media.FilePath, media.URL, media.ThumbnailURL,
		media.Width, media.Height, media.Duration, tagsJSON, metadataJSON,
		media.UploadedBy, media.EntityType, media.EntityID, media.CreatedAt, media.UpdatedAt,
	)
	return err
}

func (r *PostgresMediaRepository) Update(ctx context.Context, media *domain.Media) error {
	tagsJSON, _ := json.Marshal(media.Tags)
	metadataJSON, _ := json.Marshal(media.Metadata)

	query := `
		UPDATE blc_media
		SET title = $1, description = $2, status = $3, url = $4, thumbnail_url = $5,
		    width = $6, height = $7, duration = $8, tags = $9, metadata = $10,
		    entity_type = $11, entity_id = $12, updated_at = $13
		WHERE id = $14`

	_, err := r.db.ExecContext(ctx, query,
		media.Title, media.Description, media.Status, media.URL, media.ThumbnailURL,
		media.Width, media.Height, media.Duration, tagsJSON, metadataJSON,
		media.EntityType, media.EntityID, media.UpdatedAt, media.ID,
	)
	return err
}

func (r *PostgresMediaRepository) FindByID(ctx context.Context, id string) (*domain.Media, error) {
	query := `
		SELECT id, name, title, description, media_type, status, mime_type, file_size,
		       file_path, url, thumbnail_url, width, height, duration, tags, metadata,
		       uploaded_by, entity_type, entity_id, created_at, updated_at
		FROM blc_media WHERE id = $1`

	return r.scanMedia(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresMediaRepository) FindByEntityID(ctx context.Context, entityType, entityID string) ([]*domain.Media, error) {
	query := `
		SELECT id, name, title, description, media_type, status, mime_type, file_size,
		       file_path, url, thumbnail_url, width, height, duration, tags, metadata,
		       uploaded_by, entity_type, entity_id, created_at, updated_at
		FROM blc_media
		WHERE entity_type = $1 AND entity_id = $2 AND status != 'DELETED'
		ORDER BY created_at DESC`

	return r.queryMedias(ctx, query, entityType, entityID)
}

func (r *PostgresMediaRepository) FindByType(ctx context.Context, mediaType domain.MediaType) ([]*domain.Media, error) {
	query := `
		SELECT id, name, title, description, media_type, status, mime_type, file_size,
		       file_path, url, thumbnail_url, width, height, duration, tags, metadata,
		       uploaded_by, entity_type, entity_id, created_at, updated_at
		FROM blc_media
		WHERE media_type = $1 AND status != 'DELETED'
		ORDER BY created_at DESC`

	return r.queryMedias(ctx, query, mediaType)
}

func (r *PostgresMediaRepository) FindByStatus(ctx context.Context, status domain.MediaStatus) ([]*domain.Media, error) {
	query := `
		SELECT id, name, title, description, media_type, status, mime_type, file_size,
		       file_path, url, thumbnail_url, width, height, duration, tags, metadata,
		       uploaded_by, entity_type, entity_id, created_at, updated_at
		FROM blc_media
		WHERE status = $1
		ORDER BY created_at DESC`

	return r.queryMedias(ctx, query, status)
}

func (r *PostgresMediaRepository) FindByTags(ctx context.Context, tags []string) ([]*domain.Media, error) {
	query := `
		SELECT id, name, title, description, media_type, status, mime_type, file_size,
		       file_path, url, thumbnail_url, width, height, duration, tags, metadata,
		       uploaded_by, entity_type, entity_id, created_at, updated_at
		FROM blc_media
		WHERE tags ?| $1 AND status != 'DELETED'
		ORDER BY created_at DESC`

	return r.queryMedias(ctx, query, pq.Array(tags))
}

func (r *PostgresMediaRepository) FindAll(ctx context.Context, limit, offset int) ([]*domain.Media, error) {
	query := `
		SELECT id, name, title, description, media_type, status, mime_type, file_size,
		       file_path, url, thumbnail_url, width, height, duration, tags, metadata,
		       uploaded_by, entity_type, entity_id, created_at, updated_at
		FROM blc_media
		WHERE status != 'DELETED'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	return r.queryMedias(ctx, query, limit, offset)
}

func (r *PostgresMediaRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM blc_media WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresMediaRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM blc_media WHERE status != 'DELETED'`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *PostgresMediaRepository) scanMedia(row interface {
	Scan(dest ...interface{}) error
}) (*domain.Media, error) {
	media := &domain.Media{}
	var tagsJSON, metadataJSON []byte

	err := row.Scan(
		&media.ID, &media.Name, &media.Title, &media.Description,
		&media.MediaType, &media.Status, &media.MimeType, &media.FileSize,
		&media.FilePath, &media.URL, &media.ThumbnailURL, &media.Width, &media.Height,
		&media.Duration, &tagsJSON, &metadataJSON, &media.UploadedBy,
		&media.EntityType, &media.EntityID, &media.CreatedAt, &media.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(tagsJSON, &media.Tags); err != nil {
		media.Tags = make([]string, 0)
	}
	if err := json.Unmarshal(metadataJSON, &media.Metadata); err != nil {
		media.Metadata = make(map[string]interface{})
	}

	return media, nil
}

func (r *PostgresMediaRepository) queryMedias(ctx context.Context, query string, args ...interface{}) ([]*domain.Media, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medias []*domain.Media
	for rows.Next() {
		media := &domain.Media{}
		var tagsJSON, metadataJSON []byte

		if err := rows.Scan(
			&media.ID, &media.Name, &media.Title, &media.Description,
			&media.MediaType, &media.Status, &media.MimeType, &media.FileSize,
			&media.FilePath, &media.URL, &media.ThumbnailURL, &media.Width, &media.Height,
			&media.Duration, &tagsJSON, &metadataJSON, &media.UploadedBy,
			&media.EntityType, &media.EntityID, &media.CreatedAt, &media.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(tagsJSON, &media.Tags); err != nil {
			media.Tags = make([]string, 0)
		}
		if err := json.Unmarshal(metadataJSON, &media.Metadata); err != nil {
			media.Metadata = make(map[string]interface{})
		}

		medias = append(medias, media)
	}

	return medias, rows.Err()
}
