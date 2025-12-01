package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/qhato/ecommerce/pkg/database"
	"github.com/qhato/ecommerce/pkg/errors"
)

// DBTX define una interfaz común para ejecutar consultas,
// permitiendo que los métodos acepten tanto una conexión de pool como una transacción.
type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

// PostgresProductRepository implements the ProductRepository interface
type PostgresProductRepository struct {
	db *database.DB
}

// NewPostgresProductRepository creates a new PostgreSQL product repository
func NewPostgresProductRepository(db *database.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

// Create creates a new product safely within a transaction
func (r *PostgresProductRepository) Create(ctx context.Context, product *domain.Product) error {
	// 1. Iniciar Transacción
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.InternalWrap(err, "failed to begin transaction")
	}
	// Asegurar rollback en caso de error o pánico
	defer func() { _ = tx.Rollback(ctx) }()

	// 2. Insertar Producto
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

	err = tx.QueryRow(ctx, query,
		archivedFlag,
		product.CanSellWithoutOptions,
		product.CanonicalURL,
		product.DisplayTemplate,
		product.EnableDefaultSKUInInventory,
		product.Manufacture,
		product.MetaDescription,
		product.MetaTitle,
		product.Model,
		product.OverrideGeneratedURL,
		product.URL,
		product.URLKey,
		product.DefaultCategoryID,
		product.DefaultSkuID,
	).Scan(&product.ID)

	if err != nil {
		return errors.InternalWrap(err, "failed to create product")
	}

	// 4. Commit Transacción
	if err := tx.Commit(ctx); err != nil {
		return errors.InternalWrap(err, "failed to commit transaction")
	}

	return nil
}

// Update updates an existing product safely within a transaction
func (r *PostgresProductRepository) Update(ctx context.Context, product *domain.Product) error {
	// 1. Iniciar Transacción
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.InternalWrap(err, "failed to begin transaction")
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 2. Actualizar Producto base
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

	tag, err := tx.Exec(ctx, query,
		archivedFlag,
		product.CanSellWithoutOptions,
		product.CanonicalURL,
		product.DisplayTemplate,
		product.EnableDefaultSKUInInventory,
		product.Manufacture,
		product.MetaDescription,
		product.MetaTitle,
		product.Model,
		product.OverrideGeneratedURL,
		product.URL,
		product.URLKey,
		product.DefaultCategoryID,
		product.DefaultSkuID,
		product.ID,
	)

	if err != nil {
		return errors.InternalWrap(err, "failed to update product")
	}

	if tag.RowsAffected() == 0 {
		return errors.NotFound("product not found")
	}

	// 4. Commit Transacción
	if err := tx.Commit(ctx); err != nil {
		return errors.InternalWrap(err, "failed to commit transaction")
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

	// Usamos r.db.Pool() directamente ya que es una lectura simple
	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&archivedFlag,
		&product.CanSellWithoutOptions,
		&product.CanonicalURL,
		&product.DisplayTemplate,
		&product.EnableDefaultSKUInInventory,
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
		product.DefaultSkuID = &defaultSKUID.Int64
	}

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

// FindAll retrieves all products with pagination (Optimized for N+1)
func (r *PostgresProductRepository) FindAll(ctx context.Context, filter *domain.ProductFilter) ([]*domain.Product, int64, error) {
	whereClause := ""
	if !filter.IncludeArchived {
		whereClause = "WHERE archived = 'N'"
	}

	// 1. Contar total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM blc_product %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to count products")
	}

	// 2. Obtener productos (solo datos base)
	orderByClause := r.buildOrderByClause(filter.SortBy, filter.SortOrder)
	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(`
		SELECT
			product_id, archived, can_sell_without_options, canonical_url,
			display_template, enable_default_sku_in_inventory, manufacture,
			meta_desc, meta_title, model, override_generated_url,
			url, url_key, default_category_id, default_sku_id
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

	products, _, err := r.scanProducts(rows)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// FindByCategoryID retrieves products by category ID (Optimized for N+1)
func (r *PostgresProductRepository) FindByCategoryID(ctx context.Context, categoryID int64, filter *domain.ProductFilter) ([]*domain.Product, int64, error) {
	whereClause := "WHERE xref.category_id = $1"
	if !filter.IncludeArchived {
		whereClause += " AND p.archived = 'N'"
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT p.product_id)
		FROM blc_product p
		INNER JOIN blc_category_product_xref xref ON p.product_id = xref.product_id
		%s`, whereClause)

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, categoryID).Scan(&total); err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to count products by category")
	}

	orderByClause := r.buildOrderByClause(filter.SortBy, filter.SortOrder)
	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(`
		SELECT DISTINCT
			p.product_id, p.archived, p.can_sell_without_options, p.canonical_url,
			p.display_template, p.enable_default_sku_in_inventory, p.manufacture,
			p.meta_desc, p.meta_title, p.model, p.override_generated_url,
			p.url, p.url_key, p.default_category_id, p.default_sku_id
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

	products, _, err := r.scanProducts(rows)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Search searches products by query (Optimized and Secure)
func (r *PostgresProductRepository) Search(ctx context.Context, queryTerm string, filter *domain.ProductFilter) ([]*domain.Product, int64, error) {
	whereClause := `
		WHERE (
			model ILIKE $1 OR
			manufacture ILIKE $1 OR
			meta_title ILIKE $1 OR
			meta_desc ILIKE $1
		)`

	if !filter.IncludeArchived {
		whereClause += " AND archived = 'N'"
	}

	searchTerm := "%" + queryTerm + "%"

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM blc_product %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, searchTerm).Scan(&total); err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to count search results")
	}

	orderByClause := r.buildOrderByClause(filter.SortBy, filter.SortOrder)
	offset := (filter.Page - 1) * filter.PageSize

	searchQuery := fmt.Sprintf(`
		SELECT
			product_id, archived, can_sell_without_options, canonical_url,
			display_template, enable_default_sku_in_inventory, manufacture,
			meta_desc, meta_title, model, override_generated_url,
			url, url_key, default_category_id, default_sku_id
		FROM blc_product
		%s
		%s
		LIMIT $2 OFFSET $3`,
		whereClause,
		orderByClause,
	)

	rows, err := r.db.Query(ctx, searchQuery, searchTerm, filter.PageSize, offset)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to search products")
	}
	defer rows.Close()

	products, _, err := r.scanProducts(rows)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

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

// scanProducts escanea las filas en objetos Product y retorna también la lista de IDs
func (r *PostgresProductRepository) scanProducts(rows pgx.Rows) ([]*domain.Product, []int64, error) {
	var products []*domain.Product
	var ids []int64

	for rows.Next() {
		product := &domain.Product{}
		var archivedFlag string
		var defaultCategoryID, defaultSKUID sql.NullInt64

		err := rows.Scan(
			&product.ID,
			&archivedFlag,
			&product.CanSellWithoutOptions,
			&product.CanonicalURL,
			&product.DisplayTemplate,
			&product.EnableDefaultSKUInInventory,
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
		if err != nil {
			return nil, nil, errors.InternalWrap(err, "failed to scan product")
		}

		product.Archived = archivedFlag == "Y"
		if defaultCategoryID.Valid {
			product.DefaultCategoryID = &defaultCategoryID.Int64
		}
		if defaultSKUID.Valid {
			product.DefaultSkuID = &defaultSKUID.Int64
		}

		products = append(products, product)
		ids = append(ids, product.ID)
	}
	return products, ids, nil
}

func (r *PostgresProductRepository) buildOrderByClause(sortBy, sortOrder string) string {
	validColumns := map[string]string{
		"name":       "model",
		"created_at": "product_id", // Fallback seguro, idealmente tener fecha de creación real
		"updated_at": "product_id",
		"price":      "product_id", // Fallback hasta implementar joins con SKU
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
