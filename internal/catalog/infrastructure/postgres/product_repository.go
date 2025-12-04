package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// ProductRepository implements domain.ProductRepository for PostgreSQL persistence.
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new PostgreSQL product repository.
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Save stores a new product or updates an existing one.
func (r *ProductRepository) Save(ctx context.Context, product *domain.Product) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	archivedChar := "N"
	if product.Archived {
		archivedChar = "Y"
	}

	canSellWithoutOptions := sql.NullBool{Bool: product.CanSellWithoutOptions, Valid: true}
	enableDefaultSKUInInventory := sql.NullBool{Bool: product.EnableDefaultSKUInInventory, Valid: true}

	if product.ID == 0 {
		// Insert new product
		query := `
			INSERT INTO blc_product (
				archived, can_sell_without_options, canonical_url, display_template, 
				enable_default_sku_in_inventory, manufacture, meta_desc, meta_title, 
				model, override_generated_url, url, url_key, 
				default_category_id, default_sku_id, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
			) RETURNING product_id`
		err = tx.QueryRowContext(ctx, query,
			archivedChar, canSellWithoutOptions, product.CanonicalURL, product.DisplayTemplate,
			enableDefaultSKUInInventory, product.Manufacture, product.MetaDescription, product.MetaTitle,
			product.Model, product.OverrideGeneratedURL, product.URL, product.URLKey,
			product.DefaultCategoryID, product.DefaultSkuID, product.CreatedAt, product.UpdatedAt,
		).Scan(&product.ID)
		if err != nil {
			return fmt.Errorf("failed to insert product: %w", err)
		}
	} else {
		// Update existing product
		query := `
			UPDATE blc_product SET
				archived = $1, can_sell_without_options = $2, canonical_url = $3, 
				display_template = $4, enable_default_sku_in_inventory = $5, 
				manufacture = $6, meta_desc = $7, meta_title = $8, model = $9, 
				override_generated_url = $10, url = $11, url_key = $12, 
				default_category_id = $13, default_sku_id = $14, updated_at = $15
			WHERE product_id = $16`
		_, err = tx.ExecContext(ctx, query,
			archivedChar, canSellWithoutOptions, product.CanonicalURL, product.DisplayTemplate,
			enableDefaultSKUInInventory, product.Manufacture, product.MetaDescription, product.MetaTitle,
			product.Model, product.OverrideGeneratedURL, product.URL, product.URLKey,
			product.DefaultCategoryID, product.DefaultSkuID, product.UpdatedAt, product.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a product by its unique identifier.
func (r *ProductRepository) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT
			product_id, archived, can_sell_without_options, canonical_url, display_template, 
			enable_default_sku_in_inventory, manufacture, meta_desc, meta_title, 
			model, override_generated_url, url, url_key, 
			default_category_id, default_sku_id, created_at, updated_at
		FROM blc_product WHERE product_id = $1`

	var product domain.Product
	var archivedChar string
	var canSellWithoutOptions sql.NullBool
	var enableDefaultSKUInInventory sql.NullBool
	var defaultCategoryID sql.NullInt64
	var defaultSkuID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&product.ID, &archivedChar, &canSellWithoutOptions, &product.CanonicalURL, &product.DisplayTemplate,
		&enableDefaultSKUInInventory, &product.Manufacture, &product.MetaDescription, &product.MetaTitle,
		&product.Model, &product.OverrideGeneratedURL, &product.URL, &product.URLKey,
		&defaultCategoryID, &defaultSkuID, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query product by ID: %w", err)
	}

	product.Archived = (archivedChar == "Y")
	if canSellWithoutOptions.Valid {
		product.CanSellWithoutOptions = canSellWithoutOptions.Bool
	}
	if enableDefaultSKUInInventory.Valid {
		product.EnableDefaultSKUInInventory = enableDefaultSKUInInventory.Bool
	}
	if defaultCategoryID.Valid {
		product.DefaultCategoryID = &defaultCategoryID.Int64
	}
	if defaultSkuID.Valid {
		product.DefaultSkuID = &defaultSkuID.Int64
	}

	return &product, nil
}

// FindAll retrieves all products.
func (r *ProductRepository) FindAll(ctx context.Context, filter *domain.ProductFilter) ([]*domain.Product, int64, error) {
	// Base query
	countQuery := `SELECT COUNT(*) FROM blc_product`
	query := `SELECT product_id, archived, can_sell_without_options, canonical_url, display_template,
			enable_default_sku_in_inventory, manufacture, meta_desc, meta_title, model, override_generated_url,
			url, url_key, default_category_id, default_sku_id, created_at, updated_at
			FROM blc_product`

	var args []interface{}
	argIdx := 1

	// Build WHERE clause
	whereClauses := []string{}

	// Example: filter by archived status
	if !filter.IncludeArchived {
		whereClauses = append(whereClauses, fmt.Sprintf("archived = $%d", argIdx))
		args = append(args, "N")
		argIdx++
	}

	// Apply WHERE clauses
	if len(whereClauses) > 0 {
		countQuery += " WHERE " + strings.Join(whereClauses, " AND ")
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total results
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Apply sorting
	if filter.SortBy != "" {
		orderBy := map[string]string{
			"name":       "manufacture", // Mapping to a relevant column in blc_product
			"created_at": "created_at",
			"updated_at": "updated_at",
			"model":      "model",
		}
		sortColumn, ok := orderBy[filter.SortBy]
		if !ok {
			sortColumn = "created_at"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, strings.ToUpper(filter.SortOrder))
	}

	// Apply pagination
	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, (filter.Page-1)*filter.PageSize, filter.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query all products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		var archivedChar string
		var canSellWithoutOptions sql.NullBool
		var enableDefaultSKUInInventory sql.NullBool
		var defaultCategoryID sql.NullInt64
		var defaultSkuID sql.NullInt64

		err := rows.Scan(
			&product.ID, &archivedChar, &canSellWithoutOptions, &product.CanonicalURL, &product.DisplayTemplate,
			&enableDefaultSKUInInventory, &product.Manufacture, &product.MetaDescription, &product.MetaTitle,
			&product.Model, &product.OverrideGeneratedURL, &product.URL, &product.URLKey,
			&defaultCategoryID, &defaultSkuID, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product row: %w", err)
		}

		product.Archived = (archivedChar == "Y")
		if canSellWithoutOptions.Valid {
			product.CanSellWithoutOptions = canSellWithoutOptions.Bool
		}
		if enableDefaultSKUInInventory.Valid {
			product.EnableDefaultSKUInInventory = enableDefaultSKUInInventory.Bool
		}
		if defaultCategoryID.Valid {
			product.DefaultCategoryID = &defaultCategoryID.Int64
		}
		if defaultSkuID.Valid {
			product.DefaultSkuID = &defaultSkuID.Int64
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during rows iteration: %w", err)
	}

	return products, totalCount, nil
}

// FindByURL retrieves a product by URL
func (r *ProductRepository) FindByURL(ctx context.Context, url string) (*domain.Product, error) {
	query := `
		SELECT
			product_id, archived, can_sell_without_options, canonical_url, display_template, 
			enable_default_sku_in_inventory, manufacture, meta_desc, meta_title, 
			model, override_generated_url, url, url_key, 
			default_category_id, default_sku_id, created_at, updated_at
		FROM blc_product WHERE url = $1`

	var product domain.Product
	var archivedChar string
	var canSellWithoutOptions sql.NullBool
	var enableDefaultSKUInInventory sql.NullBool
	var defaultCategoryID sql.NullInt64
	var defaultSkuID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, url)
	err := row.Scan(
		&product.ID, &archivedChar, &canSellWithoutOptions, &product.CanonicalURL, &product.DisplayTemplate,
		&enableDefaultSKUInInventory, &product.Manufacture, &product.MetaDescription, &product.MetaTitle,
		&product.Model, &product.OverrideGeneratedURL, &product.URL, &product.URLKey,
		&defaultCategoryID, &defaultSkuID, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query product by URL: %w", err)
	}

	product.Archived = (archivedChar == "Y")
	if canSellWithoutOptions.Valid {
		product.CanSellWithoutOptions = canSellWithoutOptions.Bool
	}
	if enableDefaultSKUInInventory.Valid {
		product.EnableDefaultSKUInInventory = enableDefaultSKUInInventory.Bool
	}
	if defaultCategoryID.Valid {
		product.DefaultCategoryID = &defaultCategoryID.Int64
	}
	if defaultSkuID.Valid {
		product.DefaultSkuID = &defaultSkuID.Int64
	}

	return &product, nil
}

// FindByURLKey retrieves a product by URL key
func (r *ProductRepository) FindByURLKey(ctx context.Context, urlKey string) (*domain.Product, error) {
	query := `
		SELECT
			product_id, archived, can_sell_without_options, canonical_url, display_template, 
			enable_default_sku_in_inventory, manufacture, meta_desc, meta_title, 
			model, override_generated_url, url, url_key, 
			default_category_id, default_sku_id, created_at, updated_at
		FROM blc_product WHERE url_key = $1`

	var product domain.Product
	var archivedChar string
	var canSellWithoutOptions sql.NullBool
	var enableDefaultSKUInInventory sql.NullBool
	var defaultCategoryID sql.NullInt64
	var defaultSkuID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, urlKey)
	err := row.Scan(
		&product.ID, &archivedChar, &canSellWithoutOptions, &product.CanonicalURL, &product.DisplayTemplate,
		&enableDefaultSKUInInventory, &product.Manufacture, &product.MetaDescription, &product.MetaTitle,
		&product.Model, &product.OverrideGeneratedURL, &product.URL, &product.URLKey,
		&defaultCategoryID, &defaultSkuID, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query product by URL key: %w", err)
	}

	product.Archived = (archivedChar == "Y")
	if canSellWithoutOptions.Valid {
		product.CanSellWithoutOptions = canSellWithoutOptions.Bool
	}
	if enableDefaultSKUInInventory.Valid {
		product.EnableDefaultSKUInInventory = enableDefaultSKUInInventory.Bool
	}
	if defaultCategoryID.Valid {
		product.DefaultCategoryID = &defaultCategoryID.Int64
	}
	if defaultSkuID.Valid {
		product.DefaultSkuID = &defaultSkuID.Int64
	}

	return &product, nil
}

// FindAll retrieves all products with pagination
func (r *ProductRepository) FindAll(ctx context.Context, filter *domain.ProductFilter) ([]*domain.Product, int64, error) {
	// Base query
	countQuery := `SELECT COUNT(*) FROM blc_product`
	query := `SELECT product_id, archived, can_sell_without_options, canonical_url, display_template,
			enable_default_sku_in_inventory, manufacture, meta_desc, meta_title, model, override_generated_url,
			url, url_key, default_category_id, default_sku_id, created_at, updated_at
			FROM blc_product`

	var args []interface{}
	argIdx := 1

	// Build WHERE clause
	whereClauses := []string{}

	// Example: filter by archived status
	if !filter.IncludeArchived {
		whereClauses = append(whereClauses, fmt.Sprintf("archived = $%d", argIdx))
		args = append(args, "N")
		argIdx++
	}

	// Apply WHERE clauses
	if len(whereClauses) > 0 {
		countQuery += " WHERE " + strings.Join(whereClauses, " AND ")
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total results
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Apply sorting
	if filter.SortBy != "" {
		orderBy := map[string]string{
			"name":       "manufacture", // Mapping to a relevant column in blc_product
			"created_at": "created_at",
			"updated_at": "updated_at",
			"model":      "model",
		}
		sortColumn, ok := orderBy[filter.SortBy]
		if !ok {
			sortColumn = "created_at"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, strings.ToUpper(filter.SortOrder))
	}

	// Apply pagination
	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, (filter.Page-1)*filter.PageSize, filter.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query all products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		var archivedChar string
		var canSellWithoutOptions sql.NullBool
		var enableDefaultSKUInInventory sql.NullBool
		var defaultCategoryID sql.NullInt64
		var defaultSkuID sql.NullInt64

		err := rows.Scan(
			&product.ID, &archivedChar, &canSellWithoutOptions, &product.CanonicalURL, &product.DisplayTemplate,
			&enableDefaultSKUInInventory, &product.Manufacture, &product.MetaDescription, &product.MetaTitle,
			&product.Model, &product.OverrideGeneratedURL, &product.URL, &product.URLKey,
			&defaultCategoryID, &defaultSkuID, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product row: %w", err)
		}

		product.Archived = (archivedChar == "Y")
		if canSellWithoutOptions.Valid {
			product.CanSellWithoutOptions = canSellWithoutOptions.Bool
		}
		if enableDefaultSKUInInventory.Valid {
			product.EnableDefaultSKUInInventory = enableDefaultSKUInInventory.Bool
		}
		if defaultCategoryID.Valid {
			product.DefaultCategoryID = &defaultCategoryID.Int64
		}
		if defaultSkuID.Valid {
			product.DefaultSkuID = &defaultSkuID.Int64
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during rows iteration: %w", err)
	}

	return products, totalCount, nil
}

// FindByCategoryID retrieves products by category ID (now handled by CategoryProductXrefRepository)
func (r *ProductRepository) FindByCategoryID(ctx context.Context, categoryID int64, filter *domain.ProductFilter) ([]*domain.Product, int64, error) {
	return nil, 0, domain.NewDomainError("FindByCategoryID is no longer supported directly by ProductRepository. Use CategoryProductXrefRepository to find product IDs, then retrieve products.")
}

// Search searches products by query. (Simplified implementation, needs to be refined for full-text search).
func (r *ProductRepository) Search(ctx context.Context, query string, filter *domain.ProductFilter) ([]*domain.Product, int64, error) {
	// This is a basic example. A real search would use PostgreSQL's full-text search features.
	searchTerm := "%" + strings.ToLower(query) + "%"

	countQuery := `SELECT COUNT(*) FROM blc_product WHERE LOWER(manufacture) LIKE $1 OR LOWER(model) LIKE $1 OR LOWER(meta_desc) LIKE $1 OR LOWER(meta_title) LIKE $1`
	querySQL := `SELECT product_id, archived, can_sell_without_options, canonical_url, display_template,
			enable_default_sku_in_inventory, manufacture, meta_desc, meta_title, model, override_generated_url,
			url, url_key, default_category_id, default_sku_id, created_at, updated_at
			FROM blc_product WHERE LOWER(manufacture) LIKE $1 OR LOWER(model) LIKE $1 OR LOWER(meta_desc) LIKE $1 OR LOWER(meta_title) LIKE $1`

	var args []interface{}
	args = append(args, searchTerm)
	argIdx := 2

	// Apply archived filter
	if !filter.IncludeArchived {
		querySQL += fmt.Sprintf(" AND archived = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND archived = $%d", argIdx)
		args = append(args, "N")
		argIdx++
	}

	// Apply sorting
	if filter.SortBy != "" {
		orderBy := map[string]string{
			"name":       "manufacture", // Mapping to a relevant column
			"created_at": "created_at",
			"updated_at": "updated_at",
			"model":      "model",
		}
		sortColumn, ok := orderBy[filter.SortBy]
		if !ok {
			sortColumn = "created_at"
		}
		querySQL += fmt.Sprintf(" ORDER BY %s %s", sortColumn, strings.ToUpper(filter.SortOrder))
	}

	// Count total results
	var totalCount int64
	// Use args without pagination for count
	countArgs := args[:len(args)]
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Apply pagination
	querySQL += fmt.Sprintf(" OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, (filter.Page-1)*filter.PageSize, filter.PageSize)

	rows, err := r.db.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		var archivedChar string
		var canSellWithoutOptions sql.NullBool
		var enableDefaultSKUInInventory sql.NullBool
		var defaultCategoryID sql.NullInt64
		var defaultSkuID sql.NullInt64

		err := rows.Scan(
			&product.ID, &archivedChar, &canSellWithoutOptions, &product.CanonicalURL, &product.DisplayTemplate,
			&enableDefaultSKUInInventory, &product.Manufacture, &product.MetaDescription, &product.MetaTitle,
			&product.Model, &product.OverrideGeneratedURL, &product.URL, &product.URLKey,
			&defaultCategoryID, &defaultSkuID, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product row during search: %w", err)
		}

		product.Archived = (archivedChar == "Y")
		if canSellWithoutOptions.Valid {
			product.CanSellWithoutOptions = canSellWithoutOptions.Bool
		}
		if enableDefaultSKUInInventory.Valid {
			product.EnableDefaultSKUInInventory = enableDefaultSKUInInventory.Bool
		}
		if defaultCategoryID.Valid {
			product.DefaultCategoryID = &defaultCategoryID.Int64
		}
		if defaultSkuID.Valid {
			product.DefaultSkuID = &defaultSkuID.Int64
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during search rows iteration: %w", err)
	}

	return products, totalCount, nil
}

// AddToCategory adds a product to a category (now handled by CategoryProductXrefRepository)
func (r *ProductRepository) AddToCategory(ctx context.Context, productID, categoryID int64) error {
	return domain.NewDomainError("AddToCategory is no longer supported directly by ProductRepository. Use CategoryProductXrefRepository instead.")
}

// RemoveFromCategory removes a product from a category (now handled by CategoryProductXrefRepository)
func (r *ProductRepository) RemoveFromCategory(ctx context.Context, productID, categoryID int64) error {
	return domain.NewDomainError("RemoveFromCategory is no longer supported directly by ProductRepository. Use CategoryProductXrefRepository instead.")
}
