package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/cms/domain"
)

type PostgresContentVersionRepository struct {
	db *sql.DB
}

func NewPostgresContentVersionRepository(db *sql.DB) *PostgresContentVersionRepository {
	return &PostgresContentVersionRepository{db: db}
}

func (r *PostgresContentVersionRepository) Create(ctx context.Context, version *domain.ContentVersion) error {
	query := `INSERT INTO blc_cms_content_version (
		content_id, version_number, body, created_by, comment, created_at
	) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		version.ContentID, version.VersionNumber, version.Body,
		version.CreatedBy, version.Comment, version.CreatedAt,
	).Scan(&version.ID)
}

func (r *PostgresContentVersionRepository) FindByID(ctx context.Context, id int64) (*domain.ContentVersion, error) {
	query := `SELECT id, content_id, version_number, body, created_by, comment, created_at
	FROM blc_cms_content_version WHERE id = $1`

	version := &domain.ContentVersion{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&version.ID, &version.ContentID, &version.VersionNumber,
		&version.Body, &version.CreatedBy, &version.Comment, &version.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return version, err
}

func (r *PostgresContentVersionRepository) FindByContentID(ctx context.Context, contentID int64) ([]*domain.ContentVersion, error) {
	query := `SELECT id, content_id, version_number, body, created_by, comment, created_at
	FROM blc_cms_content_version
	WHERE content_id = $1
	ORDER BY version_number DESC`

	rows, err := r.db.QueryContext(ctx, query, contentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make([]*domain.ContentVersion, 0)
	for rows.Next() {
		version := &domain.ContentVersion{}
		if err := rows.Scan(
			&version.ID, &version.ContentID, &version.VersionNumber,
			&version.Body, &version.CreatedBy, &version.Comment, &version.CreatedAt,
		); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, nil
}

func (r *PostgresContentVersionRepository) GetNextVersionNumber(ctx context.Context, contentID int64) (int, error) {
	var maxVersion sql.NullInt64
	err := r.db.QueryRowContext(ctx,
		`SELECT MAX(version_number) FROM blc_cms_content_version WHERE content_id = $1`,
		contentID).Scan(&maxVersion)

	if err != nil {
		return 0, err
	}

	if maxVersion.Valid {
		return int(maxVersion.Int64) + 1, nil
	}
	return 1, nil
}

func (r *PostgresContentVersionRepository) DeleteOldVersions(ctx context.Context, contentID int64, keepCount int) error {
	query := `DELETE FROM blc_cms_content_version
	WHERE content_id = $1 AND id NOT IN (
		SELECT id FROM blc_cms_content_version
		WHERE content_id = $1
		ORDER BY version_number DESC
		LIMIT $2
	)`

	_, err := r.db.ExecContext(ctx, query, contentID, keepCount)
	return err
}
