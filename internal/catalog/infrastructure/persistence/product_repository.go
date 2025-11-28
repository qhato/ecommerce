package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/qhato/ecommerce/pkg/database"
	"github.com/qhato/ecommerce/pkg/errors"
)

// PostgresProductRepository implements the ProductRepository interface
type PostgresProductRepository struct {
	db *database.DB
}

// NewPostgresProductRepository creates a new PostgreSQL product repository
func NewPostgresProductRepository(db *database.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

// Create creates a new product
func (r *PostgresProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO blc_product (
			product_id, archived, can_sell_without_options, canonical_url,
			display_template, enable_default_sku_in_inventory, manufacture,
			meta_desc, meta_title, model, override_generated_url,
			url, url_key, default_category_id, default_sku_id
		) VALUES (
			nextval('blc_product_seq'), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING product_id`

	archivedFlag := "N"
	if product.Archived {
		archivedFlag = "Y"
	}

	err := r.db.QueryRow(ctx, query,
		archivedFlag,
		product.CanSellWithoutOptions,
		product.CanonicalURL,
		product.DisplayTemplate,
		product.EnableDefaultSKU,
		product.Manufacture,
		product.MetaDescription,
		product.MetaTitle,
		product.Model,
		product.OverrideGeneratedURL,
		product.URL,
		product.URLKey,
		product.DefaultCategoryID,
		product.DefaultSKUID,
	).Scan(&product.ID)

	if err != nil {
		return errors.InternalWrap(err, "failed to create product")
	}

	// Insert attributes
	if len(product.Attributes) > 0 {
		if err := r.insertAttributes(ctx, product.ID, product.Attributes); err != nil {
			return err
		}
	}

	return nil
}

// Update updates an existing product
func (r *PostgresProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE blc_product SET
			archived = $1,
			can_sell_without_options = $2,
			canonical_url = $3,
			display_template = $4,
			enable_default_sku_in_inventory = $5,
			manufacture = $6,
			meta_desc = $7,
			meta_title = $8,
			model = $9,
			override_generated_url = $10,
			url = $11,
			url_key = $12,
			default_category_id = $13,
			default_sku_id = $14
		WHERE product_id = $15`

	archivedFlag := "N"
	if product.Archived {
		archivedFlag = "Y"
	}

	err := r.db.Exec(ctx, query,
		archivedFlag,
		product.CanSellWithoutOptions,
		product.CanonicalURL,
		product.DisplayTemplate,
		product.EnableDefaultSKU,
		product.Manufacture,
		product.MetaDescription,
		product.MetaTitle,
		product.Model,
		product.OverrideGeneratedURL,
		product.URL,
		product.URLKey,
		product.DefaultCategoryID,
		product.DefaultSKUID,
		product.ID,
	)

	if err != nil {
	tag, err := r.db.Pool().Exec(ctx, query,
		archivedFlag,
		product.CanSellWithoutOptions,
		product.CanonicalURL,
		product.DisplayTemplate,
		product.EnableDefaultSKU,
		product.Manufacture,
		product.MetaDescription,
		product.MetaTitle,
		product.Model,
		product.OverrideGeneratedURL,
		product.URL,
		product.URLKey,
		product.DefaultCategoryID,
		product.DefaultSKUID,
		product.ID,
	)
	if err != nil {
		return errors.InternalWrap(err, "failed to update product")
	}

	if tag.RowsAffected() == 0 {
		return errors.NotFound("product not found")
	}

	// Update attributes (delete and re-insert)
	if err := r.deleteAttributes(ctx, product.ID); err != nil {
		return err
	}

	if len(product.Attributes) > 0 {
		if err := r.insertAttributes(ctx, product.ID, product.Attributes); err != nil {
			return err
		}
	}

	return nil
}

// Delete soft deletes a product by marking it as archived
func (r *PostgresProductRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE blc_product SET archived = 'Y' WHERE product_id = $1`

	tag, err := r.db.Pool().Exec(ctx, query, id)
	if err != nil {
		return errors.InternalWrap(err, "failed to delete product")
	}

	if tag.RowsAffected() == 0 {
		return errors.NotFound("product not found")
	}

	return nil
}

// FindByID retrieves a product by ID
func (r *PostgresProductRepository) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT
			product_id, archived, can_sell_without_options, canonical_url,
			display_template, enable_default_sku_in_inventory, manufacture,
			meta_desc, meta_title, model, override_generated_url,
			url, url_key, default_category_id, default_sku_id
		FROM blc_product
		WHERE product_id = $1`

	product := &domain.Product{}
	var archivedFlag string
	var defaultCategoryID, defaultSKUID sql.NullInt64

	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&archivedFlag,
		&product.CanSellWithoutOptions,
		&product.CanonicalURL,
		&product.DisplayTemplate,
		&product.EnableDefaultSKU,
		&product.Manufacture,
		&product.MetaDescription,
		&product.MetaTitle,
		&product.Model,
		&product.OverrideGeneratedURL,
		&product.URL,
		&product.URLKey,
		&defaultCategoryID,
		&defaultSKUID,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("product not found")
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find product")
	}

	product.Archived = archivedFlag == "Y"
	if defaultCategoryID.Valid {
		product.DefaultCategoryID = &defaultCategoryID.Int64
	}
	if defaultSKUID.Valid {
		product.DefaultSKUID = &defaultSKUID.Int64
	}

	// Load attributes
	attributes, err := r.findAttributes(ctx, id)
	if err != nil {
		return nil, err
	}
	product.Attributes = attributes

	return product, nil
}

// FindByURL retrieves a product by URL
func (r *PostgresProductRepository) FindByURL(ctx context.Context, url string) (*domain.Product, error) {
	query := `
		SELECT product_id
		FROM blc_product
		WHERE url = $1 AND archived = 'N'
		LIMIT 1`

	var id int64
	err := r.db.QueryRow(ctx, query, url).Scan(&id)
	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("product not found")
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find product by URL")
	}

	return r.FindByID(ctx, id)
}

// FindByURLKey retrieves a product by URL key
func (r *PostgresProductRepository) FindByURLKey(ctx context.Context, urlKey string) (*domain.Product, error) {
	query := `
		SELECT product_id
		FROM blc_product
		WHERE url_key = $1 AND archived = 'N'
		LIMIT 1`

	var id int64
	err := r.db.QueryRow(ctx, query, urlKey).Scan(&id)
	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("product not found")
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find product by URL key")
	}

	return r.FindByID(ctx, id)
}

// FindAll retrieves all products with pagination
func (r *PostgresProductRepository) FindAll(ctx context.Context, filter *domain.ProductFilter) ([]*domain.Product, int64, error) {
	// Build query
	whereClause := ""
	if !filter.IncludeArchived {
		whereClause = "WHERE archived = 'N'"
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM blc_product %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to count products")
	}

	// Build main query with pagination
	orderByClause := r.buildOrderByClause(filter.SortBy, filter.SortOrder)
	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(`
		SELECT product_id
		FROM blc_product
		%s
		%s
		LIMIT $1 OFFSET $2`,
		whereClause,
		orderByClause,
	)

	rows, err := r.db.Query(ctx, query, filter.PageSize, offset)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to list products")
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, errors.Wrap(err, "failed to scan product ID")
		}

		product, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, nil
}

// FindByCategoryID retrieves products by category ID
func (r *PostgresProductRepository) FindByCategoryID(ctx context.Context, categoryID int64, filter *domain.ProductFilter) ([]*domain.Product, int64, error) {
	// Build where clause
	whereClause := "WHERE xref.category_id = $1"
	if !filter.IncludeArchived {
		whereClause += " AND p.archived = 'N'"
	}

	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT p.product_id)
		FROM blc_product p
		INNER JOIN blc_category_product_xref xref ON p.product_id = xref.product_id
		%s`, whereClause)

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, categoryID).Scan(&total); err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to count products by category")
	}

	// Build main query
	orderByClause := r.buildOrderByClause(filter.SortBy, filter.SortOrder)
	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(`
		SELECT DISTINCT p.product_id
		FROM blc_product p
		INNER JOIN blc_category_product_xref xref ON p.product_id = xref.product_id
		%s
		%s
		LIMIT $2 OFFSET $3`,
		whereClause,
		orderByClause,
	)

	rows, err := r.db.Query(ctx, query, categoryID, filter.PageSize, offset)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to list products by category")
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, errors.Wrap(err, "failed to scan product ID")
		}

		product, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, nil
}

// Search searches products by query
func (r *PostgresProductRepository) Search(ctx context.Context, query string, filter *domain.ProductFilter) ([]*domain.Product, int64, error) {
	// Build where clause with full-text search
	whereClause := fmt.Sprintf(`
		WHERE (
			model ILIKE '%%%s%%' OR
			manufacture ILIKE '%%%s%%' OR
			meta_title ILIKE '%%%s%%' OR
			meta_desc ILIKE '%%%s%%'
		)`, query, query, query, query)

	if !filter.IncludeArchived {
		whereClause += " AND archived = 'N'"
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM blc_product %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to count search results")
	}

	// Build main query
	orderByClause := r.buildOrderByClause(filter.SortBy, filter.SortOrder)
	offset := (filter.Page - 1) * filter.PageSize

	searchQuery := fmt.Sprintf(`
		SELECT product_id
		FROM blc_product
		%s
		%s
		LIMIT $1 OFFSET $2`,
		whereClause,
		orderByClause,
	)

	rows, err := r.db.Query(ctx, searchQuery, filter.PageSize, offset)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to search products")
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, errors.Wrap(err, "failed to scan product ID")
		}

		product, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, nil
}

// AddToCategory adds a product to a category
func (r *PostgresProductRepository) AddToCategory(ctx context.Context, productID, categoryID int64) error {
	query := `
		INSERT INTO blc_category_product_xref (category_product_id, product_id, category_id)
		VALUES (nextval('blc_category_product_xref_seq'), $1, $2)
		ON CONFLICT DO NOTHING`

	err := r.db.Exec(ctx, query, productID, categoryID)
	if err != nil {
		return errors.InternalWrap(err, "failed to add product to category")
	}

	return nil
}

// RemoveFromCategory removes a product from a category
func (r *PostgresProductRepository) RemoveFromCategory(ctx context.Context, productID, categoryID int64) error {
	query := `
		DELETE FROM blc_category_product_xref
		WHERE product_id = $1 AND category_id = $2`

	err := r.db.Exec(ctx, query, productID, categoryID)
	if err != nil {
		return errors.InternalWrap(err, "failed to remove product from category")
	}

	return nil
}

// Helper methods

func (r *PostgresProductRepository) insertAttributes(ctx context.Context, productID int64, attributes []domain.ProductAttribute) error {
	query := `
		INSERT INTO blc_product_attribute (product_attribute_id, name, value, product_id)
		VALUES (nextval('blc_product_attribute_seq'), $1, $2, $3)`

	for _, attr := range attributes {
		err := r.db.Exec(ctx, query, attr.Name, attr.Value, productID)
		if err != nil {
			return errors.InternalWrap(err, "failed to insert product attribute")
		}
	}

	return nil
}

func (r *PostgresProductRepository) deleteAttributes(ctx context.Context, productID int64) error {
	query := `DELETE FROM blc_product_attribute WHERE product_id = $1`
	err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return errors.InternalWrap(err, "failed to delete product attributes")
	}
	return nil
}

func (r *PostgresProductRepository) findAttributes(ctx context.Context, productID int64) ([]domain.ProductAttribute, error) {
	query := `
		SELECT product_attribute_id, name, value, product_id
		FROM blc_product_attribute
		WHERE product_id = $1`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find product attributes")
	}
	defer rows.Close()

	var attributes []domain.ProductAttribute
	for rows.Next() {
		var attr domain.ProductAttribute
		if err := rows.Scan(&attr.ID, &attr.Name, &attr.Value, &attr.ProductID); err != nil {
			return nil, errors.InternalWrap(err, "failed to scan product attribute")
		}
		attributes = append(attributes, attr)
	}

	return attributes, nil
}

func (r *PostgresProductRepository) buildOrderByClause(sortBy, sortOrder string) string {
	validColumns := map[string]string{
		"name":       "model",
		"created_at": "product_id",
		"updated_at": "product_id",
		"price":      "product_id",
	}

	column, ok := validColumns[sortBy]
	if !ok {
		column = "product_id"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	return fmt.Sprintf("ORDER BY %s %s", column, sortOrder)
}
