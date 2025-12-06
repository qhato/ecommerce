package persistence

import (
	"context"
	"database/sql"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

type PostgresProductRelationshipRepository struct {
	db *sql.DB
}

func NewPostgresProductRelationshipRepository(db *sql.DB) *PostgresProductRelationshipRepository {
	return &PostgresProductRelationshipRepository{db: db}
}

func (r *PostgresProductRelationshipRepository) Create(ctx context.Context, relationship *domain.ProductRelationship) error {
	query := `
		INSERT INTO blc_product_relationship (product_id, related_product_id, relationship_type, sequence, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		relationship.ProductID, relationship.RelatedProductID, relationship.RelationshipType,
		relationship.Sequence, relationship.IsActive, relationship.CreatedAt, relationship.UpdatedAt,
	).Scan(&relationship.ID)
}

func (r *PostgresProductRelationshipRepository) Update(ctx context.Context, relationship *domain.ProductRelationship) error {
	query := `
		UPDATE blc_product_relationship
		SET sequence = $1, is_active = $2, updated_at = $3
		WHERE id = $4`

	_, err := r.db.ExecContext(ctx, query,
		relationship.Sequence, relationship.IsActive, relationship.UpdatedAt, relationship.ID,
	)
	return err
}

func (r *PostgresProductRelationshipRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_product_relationship WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresProductRelationshipRepository) FindByID(ctx context.Context, id int64) (*domain.ProductRelationship, error) {
	query := `
		SELECT id, product_id, related_product_id, relationship_type, sequence, is_active, created_at, updated_at
		FROM blc_product_relationship
		WHERE id = $1`

	return r.scanRelationship(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresProductRelationshipRepository) FindByProductID(ctx context.Context, productID int64, relationshipType *domain.ProductRelationshipType) ([]*domain.ProductRelationship, error) {
	query := `
		SELECT id, product_id, related_product_id, relationship_type, sequence, is_active, created_at, updated_at
		FROM blc_product_relationship
		WHERE product_id = $1 AND is_active = true`

	var args []interface{}
	args = append(args, productID)

	if relationshipType != nil {
		query += " AND relationship_type = $2"
		args = append(args, *relationshipType)
	}

	query += " ORDER BY sequence ASC"

	return r.queryRelationships(ctx, query, args...)
}

func (r *PostgresProductRelationshipRepository) FindCrossSell(ctx context.Context, productID int64) ([]*domain.ProductRelationship, error) {
	relType := domain.RelationshipTypeCrossSell
	return r.FindByProductID(ctx, productID, &relType)
}

func (r *PostgresProductRelationshipRepository) FindUpSell(ctx context.Context, productID int64) ([]*domain.ProductRelationship, error) {
	relType := domain.RelationshipTypeUpSell
	return r.FindByProductID(ctx, productID, &relType)
}

func (r *PostgresProductRelationshipRepository) FindRelated(ctx context.Context, productID int64) ([]*domain.ProductRelationship, error) {
	relType := domain.RelationshipTypeRelated
	return r.FindByProductID(ctx, productID, &relType)
}

func (r *PostgresProductRelationshipRepository) ExistsByProducts(ctx context.Context, productID, relatedProductID int64, relationshipType domain.ProductRelationshipType) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM blc_product_relationship
		WHERE product_id = $1 AND related_product_id = $2 AND relationship_type = $3)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, productID, relatedProductID, relationshipType).Scan(&exists)
	return exists, err
}

func (r *PostgresProductRelationshipRepository) scanRelationship(row interface {
	Scan(dest ...interface{}) error
}) (*domain.ProductRelationship, error) {
	relationship := &domain.ProductRelationship{}
	err := row.Scan(
		&relationship.ID, &relationship.ProductID, &relationship.RelatedProductID,
		&relationship.RelationshipType, &relationship.Sequence, &relationship.IsActive,
		&relationship.CreatedAt, &relationship.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return relationship, nil
}

func (r *PostgresProductRelationshipRepository) queryRelationships(ctx context.Context, query string, args ...interface{}) ([]*domain.ProductRelationship, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relationships []*domain.ProductRelationship
	for rows.Next() {
		relationship := &domain.ProductRelationship{}
		if err := rows.Scan(
			&relationship.ID, &relationship.ProductID, &relationship.RelatedProductID,
			&relationship.RelationshipType, &relationship.Sequence, &relationship.IsActive,
			&relationship.CreatedAt, &relationship.UpdatedAt,
		); err != nil {
			return nil, err
		}
		relationships = append(relationships, relationship)
	}

	return relationships, rows.Err()
}
