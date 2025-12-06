package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/qhato/ecommerce/internal/cms/domain"
)

type PostgresContentRepository struct {
	db *sql.DB
}

func NewPostgresContentRepository(db *sql.DB) *PostgresContentRepository {
	return &PostgresContentRepository{db: db}
}

func (r *PostgresContentRepository) Create(ctx context.Context, content *domain.Content) error {
	customFieldsJSON, err := json.Marshal(content.CustomFields)
	if err != nil {
		return fmt.Errorf("failed to marshal custom fields: %w", err)
	}

	query := `INSERT INTO blc_cms_content (
		title, slug, type, status, body, excerpt, featured_image,
		meta_title, meta_description, meta_keywords, template,
		author_id, parent_id, sort_order, locale, published_at,
		scheduled_for, expires_at, custom_fields, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		content.Title, content.Slug, content.Type, content.Status, content.Body,
		content.Excerpt, content.FeaturedImage, content.MetaTitle, content.MetaDescription,
		content.MetaKeywords, content.Template, content.AuthorID, content.ParentID,
		content.SortOrder, content.Locale, content.PublishedAt, content.ScheduledFor,
		content.ExpiresAt, customFieldsJSON, content.CreatedAt, content.UpdatedAt,
	).Scan(&content.ID)
}

func (r *PostgresContentRepository) Update(ctx context.Context, content *domain.Content) error {
	customFieldsJSON, err := json.Marshal(content.CustomFields)
	if err != nil {
		return fmt.Errorf("failed to marshal custom fields: %w", err)
	}

	query := `UPDATE blc_cms_content SET
		title = $1, slug = $2, status = $3, body = $4, excerpt = $5,
		featured_image = $6, meta_title = $7, meta_description = $8,
		meta_keywords = $9, template = $10, parent_id = $11, sort_order = $12,
		published_at = $13, scheduled_for = $14, expires_at = $15,
		custom_fields = $16, updated_at = $17
	WHERE id = $18`

	_, err = r.db.ExecContext(ctx, query,
		content.Title, content.Slug, content.Status, content.Body, content.Excerpt,
		content.FeaturedImage, content.MetaTitle, content.MetaDescription,
		content.MetaKeywords, content.Template, content.ParentID, content.SortOrder,
		content.PublishedAt, content.ScheduledFor, content.ExpiresAt,
		customFieldsJSON, content.UpdatedAt, content.ID,
	)
	return err
}

func (r *PostgresContentRepository) FindByID(ctx context.Context, id int64) (*domain.Content, error) {
	query := `SELECT id, title, slug, type, status, body, excerpt, featured_image,
		meta_title, meta_description, meta_keywords, template, author_id, parent_id,
		sort_order, view_count, locale, published_at, scheduled_for, expires_at,
		custom_fields, created_at, updated_at
	FROM blc_cms_content WHERE id = $1`

	return r.scanContent(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresContentRepository) FindBySlug(ctx context.Context, slug, locale string) (*domain.Content, error) {
	query := `SELECT id, title, slug, type, status, body, excerpt, featured_image,
		meta_title, meta_description, meta_keywords, template, author_id, parent_id,
		sort_order, view_count, locale, published_at, scheduled_for, expires_at,
		custom_fields, created_at, updated_at
	FROM blc_cms_content WHERE slug = $1 AND locale = $2`

	return r.scanContent(r.db.QueryRowContext(ctx, query, slug, locale))
}

func (r *PostgresContentRepository) FindAll(ctx context.Context, locale string, publishedOnly bool) ([]*domain.Content, error) {
	query := `SELECT id, title, slug, type, status, body, excerpt, featured_image,
		meta_title, meta_description, meta_keywords, template, author_id, parent_id,
		sort_order, view_count, locale, published_at, scheduled_for, expires_at,
		custom_fields, created_at, updated_at
	FROM blc_cms_content WHERE locale = $1`

	if publishedOnly {
		query += " AND status = 'PUBLISHED'"
	}
	query += " ORDER BY sort_order ASC, created_at DESC"

	return r.queryContents(ctx, query, locale)
}

func (r *PostgresContentRepository) FindByType(ctx context.Context, contentType domain.ContentType, locale string, publishedOnly bool) ([]*domain.Content, error) {
	query := `SELECT id, title, slug, type, status, body, excerpt, featured_image,
		meta_title, meta_description, meta_keywords, template, author_id, parent_id,
		sort_order, view_count, locale, published_at, scheduled_for, expires_at,
		custom_fields, created_at, updated_at
	FROM blc_cms_content WHERE type = $1 AND locale = $2`

	if publishedOnly {
		query += " AND status = 'PUBLISHED'"
	}
	query += " ORDER BY sort_order ASC, created_at DESC"

	return r.queryContents(ctx, query, contentType, locale)
}

func (r *PostgresContentRepository) FindHierarchy(ctx context.Context, contentType domain.ContentType, locale string) ([]*domain.Content, error) {
	// Get all root level content
	query := `SELECT id, title, slug, type, status, body, excerpt, featured_image,
		meta_title, meta_description, meta_keywords, template, author_id, parent_id,
		sort_order, view_count, locale, published_at, scheduled_for, expires_at,
		custom_fields, created_at, updated_at
	FROM blc_cms_content
	WHERE type = $1 AND locale = $2 AND parent_id IS NULL
	ORDER BY sort_order ASC`

	roots, err := r.queryContents(ctx, query, contentType, locale)
	if err != nil {
		return nil, err
	}

	// Load children for each root
	for _, root := range roots {
		children, err := r.FindChildren(ctx, root.ID)
		if err == nil {
			root.Children = make([]domain.Content, len(children))
			for i, child := range children {
				root.Children[i] = *child
			}
		}
	}

	return roots, nil
}

func (r *PostgresContentRepository) FindChildren(ctx context.Context, parentID int64) ([]*domain.Content, error) {
	query := `SELECT id, title, slug, type, status, body, excerpt, featured_image,
		meta_title, meta_description, meta_keywords, template, author_id, parent_id,
		sort_order, view_count, locale, published_at, scheduled_for, expires_at,
		custom_fields, created_at, updated_at
	FROM blc_cms_content WHERE parent_id = $1 ORDER BY sort_order ASC`

	return r.queryContents(ctx, query, parentID)
}

func (r *PostgresContentRepository) Search(ctx context.Context, query, locale string, publishedOnly bool) ([]*domain.Content, error) {
	searchQuery := `SELECT id, title, slug, type, status, body, excerpt, featured_image,
		meta_title, meta_description, meta_keywords, template, author_id, parent_id,
		sort_order, view_count, locale, published_at, scheduled_for, expires_at,
		custom_fields, created_at, updated_at
	FROM blc_cms_content
	WHERE locale = $1 AND (
		title ILIKE $2 OR
		body ILIKE $2 OR
		excerpt ILIKE $2 OR
		meta_description ILIKE $2
	)`

	if publishedOnly {
		searchQuery += " AND status = 'PUBLISHED'"
	}
	searchQuery += " ORDER BY created_at DESC LIMIT 100"

	searchTerm := "%" + query + "%"
	return r.queryContents(ctx, searchQuery, locale, searchTerm)
}

func (r *PostgresContentRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM blc_cms_content WHERE id = $1`, id)
	return err
}

func (r *PostgresContentRepository) ExistsBySlug(ctx context.Context, slug, locale string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM blc_cms_content WHERE slug = $1 AND locale = $2)`,
		slug, locale).Scan(&exists)
	return exists, err
}

func (r *PostgresContentRepository) HasChildren(ctx context.Context, parentID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM blc_cms_content WHERE parent_id = $1)`,
		parentID).Scan(&exists)
	return exists, err
}

func (r *PostgresContentRepository) IncrementViewCount(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE blc_cms_content SET view_count = view_count + 1 WHERE id = $1`, id)
	return err
}

func (r *PostgresContentRepository) UpdateSortOrder(ctx context.Context, id int64, sortOrder int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE blc_cms_content SET sort_order = $1 WHERE id = $2`, sortOrder, id)
	return err
}

// Helper methods

func (r *PostgresContentRepository) scanContent(row interface {
	Scan(dest ...interface{}) error
}) (*domain.Content, error) {
	content := &domain.Content{}
	var customFieldsJSON []byte

	err := row.Scan(
		&content.ID, &content.Title, &content.Slug, &content.Type, &content.Status,
		&content.Body, &content.Excerpt, &content.FeaturedImage, &content.MetaTitle,
		&content.MetaDescription, &content.MetaKeywords, &content.Template,
		&content.AuthorID, &content.ParentID, &content.SortOrder, &content.ViewCount,
		&content.Locale, &content.PublishedAt, &content.ScheduledFor, &content.ExpiresAt,
		&customFieldsJSON, &content.CreatedAt, &content.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal custom fields
	if customFieldsJSON != nil {
		if err := json.Unmarshal(customFieldsJSON, &content.CustomFields); err != nil {
			return nil, fmt.Errorf("failed to unmarshal custom fields: %w", err)
		}
	}

	return content, nil
}

func (r *PostgresContentRepository) queryContents(ctx context.Context, query string, args ...interface{}) ([]*domain.Content, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contents := make([]*domain.Content, 0)
	for rows.Next() {
		content := &domain.Content{}
		var customFieldsJSON []byte

		if err := rows.Scan(
			&content.ID, &content.Title, &content.Slug, &content.Type, &content.Status,
			&content.Body, &content.Excerpt, &content.FeaturedImage, &content.MetaTitle,
			&content.MetaDescription, &content.MetaKeywords, &content.Template,
			&content.AuthorID, &content.ParentID, &content.SortOrder, &content.ViewCount,
			&content.Locale, &content.PublishedAt, &content.ScheduledFor, &content.ExpiresAt,
			&customFieldsJSON, &content.CreatedAt, &content.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Unmarshal custom fields
		if customFieldsJSON != nil {
			if err := json.Unmarshal(customFieldsJSON, &content.CustomFields); err != nil {
				return nil, fmt.Errorf("failed to unmarshal custom fields: %w", err)
			}
		}

		contents = append(contents, content)
	}

	return contents, nil
}
