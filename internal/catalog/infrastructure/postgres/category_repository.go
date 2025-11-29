package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// CategoryRepository implements domain.CategoryRepository for PostgreSQL persistence.
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new PostgreSQL category repository.
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Save stores a new category or updates an existing one.
func (r *CategoryRepository) Save(ctx context.Context, category *domain.Category) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	archivedChar := "N"
	if category.Archived {
		archivedChar = "Y"
	}

	// Handle nullable fields for insert/update
	activeEndDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if category.ActiveEndDate != nil {
		activeEndDate = sql.NullTime{Time: *category.ActiveEndDate, Valid: true}
	}
	activeStartDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if category.ActiveStartDate != nil {
		activeStartDate = sql.NullTime{Time: *category.ActiveStartDate, Valid: true}
	}
	longDescription := sql.NullString{String: category.LongDescription, Valid: category.LongDescription != ""}
	metaDescription := sql.NullString{String: category.MetaDescription, Valid: category.MetaDescription != ""}
	metaTitle := sql.NullString{String: category.MetaTitle, Valid: category.MetaTitle != ""}
	productDescPattern := sql.NullString{String: category.ProductDescPattern, Valid: category.ProductDescPattern != ""}
	productTitlePattern := sql.NullString{String: category.ProductTitlePattern, Valid: category.ProductTitlePattern != ""}
	rootDisplayOrder := sql.NullFloat64{Float64: category.RootDisplayOrder, Valid: true}
	// Broadleaf's default for root_display_order is NULL, so only set Valid if it's a non-default value you want to persist.
	if category.RootDisplayOrder == 0.0 { // Assuming 0.0 is the default/unset value in your domain
		rootDisplayOrder.Valid = false
	}
	taxCode := sql.NullString{String: category.TaxCode, Valid: category.TaxCode != ""}
	url := sql.NullString{String: category.URL, Valid: category.URL != ""}
	urlKey := sql.NullString{String: category.URLKey, Valid: category.URLKey != ""}
	defaultParentCategoryID := sql.NullInt64{Int64: 0, Valid: false}
	if category.DefaultParentCategoryID != nil {
		defaultParentCategoryID = sql.NullInt64{Int64: *category.DefaultParentCategoryID, Valid: true}
	}

	if category.ID == 0 {
		// Insert new category
		query := `
			INSERT INTO blc_category (
				active_end_date, active_start_date, archived, description, display_template, 
				external_id, fulfillment_type, inventory_type, long_description, meta_desc, 
				meta_title, name, override_generated_url, product_desc_pattern_override, 
				product_title_pattern_override, root_display_order, tax_code, url, url_key, 
				default_parent_category_id, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
			) RETURNING category_id`
		err = tx.QueryRowContext(ctx, query,
			activeEndDate, activeStartDate, archivedChar, category.Description, category.DisplayTemplate,
			category.ExternalID, category.FulfillmentType, category.InventoryType, longDescription, metaDescription,
			metaTitle, category.Name, category.OverrideGeneratedURL, productDescPattern,
			productTitlePattern, rootDisplayOrder, taxCode, url, urlKey,
			defaultParentCategoryID, category.CreatedAt, category.UpdatedAt,
		).Scan(&category.ID)
		if err != nil {
			return fmt.Errorf("failed to insert category: %w", err)
		}
	} else {
		// Update existing category
		query := `
			UPDATE blc_category SET
				active_end_date = $1, active_start_date = $2, archived = $3, description = $4, 
				display_template = $5, external_id = $6, fulfillment_type = $7, 
				inventory_type = $8, long_description = $9, meta_desc = $10, 
				meta_title = $11, name = $12, override_generated_url = $13, 
				product_desc_pattern_override = $14, product_title_pattern_override = $15, 
				root_display_order = $16, tax_code = $17, url = $18, url_key = $19, 
				default_parent_category_id = $20, updated_at = $21
			WHERE category_id = $22`
		_, err = tx.ExecContext(ctx, query,
			activeEndDate, activeStartDate, archivedChar, category.Description,
			category.DisplayTemplate, category.ExternalID, category.FulfillmentType,
			category.InventoryType, longDescription, metaDescription,
			metaTitle, category.Name, category.OverrideGeneratedURL,
			productDescPattern, productTitlePattern,
			rootDisplayOrder, taxCode, url, urlKey,
			defaultParentCategoryID, category.UpdatedAt, category.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update category: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a category by its unique identifier.
func (r *CategoryRepository) FindByID(ctx context.Context, id int64) (*domain.Category, error) {
	query := `
		SELECT
			category_id, active_end_date, active_start_date, archived, description, 
			display_template, external_id, fulfillment_type, inventory_type, 
			long_description, meta_desc, meta_title, name, override_generated_url, 
			product_desc_pattern_override, product_title_pattern_override, 
			root_display_order, tax_code, url, url_key, default_parent_category_id, 
			created_at, updated_at
		FROM blc_category WHERE category_id = $1`

	var category domain.Category
	var activeEndDate sql.NullTime
	var activeStartDate sql.NullTime
	var archivedChar string
	var longDescription sql.NullString
	var metaDescription sql.NullString
	var metaTitle sql.NullString
	var productDescPattern sql.NullString
	var productTitlePattern sql.NullString
	var rootDisplayOrder sql.NullFloat64
	var taxCode sql.NullString
	var url sql.NullString
	var urlKey sql.NullString
	var defaultParentCategoryID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&category.ID, &activeEndDate, &activeStartDate, &archivedChar, &category.Description,
		&category.DisplayTemplate, &category.ExternalID, &category.FulfillmentType,
		&category.InventoryType, &longDescription, &metaDescription, &metaTitle,
		&category.Name, &category.OverrideGeneratedURL, &productDescPattern,
		&productTitlePattern, &rootDisplayOrder, &taxCode, &url, &urlKey,
		&defaultParentCategoryID, &category.CreatedAt, &category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query category by ID: %w", err)
	}

	if activeEndDate.Valid {
		category.ActiveEndDate = &activeEndDate.Time
	}
	if activeStartDate.Valid {
		category.ActiveStartDate = &activeStartDate.Time
	}
	category.Archived = (archivedChar == "Y")
	if longDescription.Valid {
		category.LongDescription = longDescription.String
	}
	if metaDescription.Valid {
		category.MetaDescription = metaDescription.String
	}
	if metaTitle.Valid {
		category.MetaTitle = metaTitle.String
	}
	if productDescPattern.Valid {
		category.ProductDescPattern = productDescPattern.String
	}
	if productTitlePattern.Valid {
		category.ProductTitlePattern = productTitlePattern.String
	}
	if rootDisplayOrder.Valid {
		category.RootDisplayOrder = rootDisplayOrder.Float64
	}
	if taxCode.Valid {
		category.TaxCode = taxCode.String
	}
	if url.Valid {
		category.URL = url.String
	}
	if urlKey.Valid {
		category.URLKey = urlKey.String
	}
	if defaultParentCategoryID.Valid {
		category.DefaultParentCategoryID = &defaultParentCategoryID.Int64
	}

	return &category, nil
}

// FindByURL retrieves a category by URL.
func (r *CategoryRepository) FindByURL(ctx context.Context, url string) (*domain.Category, error) {
	query := `
		SELECT
			category_id, active_end_date, active_start_date, archived, description, 
			display_template, external_id, fulfillment_type, inventory_type, 
			long_description, meta_desc, meta_title, name, override_generated_url, 
			product_desc_pattern_override, product_title_pattern_override, 
			root_display_order, tax_code, url, url_key, default_parent_category_id, 
			created_at, updated_at
		FROM blc_category WHERE url = $1`

	var category domain.Category
	var activeEndDate sql.NullTime
	var activeStartDate sql.NullTime
	var archivedChar string
	var longDescription sql.NullString
	var metaDescription sql.NullString
	var metaTitle sql.NullString
	var productDescPattern sql.NullString
	var productTitlePattern sql.NullString
	var rootDisplayOrder sql.NullFloat64
	var taxCode sql.NullString
	var urlScan sql.NullString // Use different variable name to avoid conflict
	var urlKey sql.NullString
	var defaultParentCategoryID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, url)
	err := row.Scan(
		&category.ID, &activeEndDate, &activeStartDate, &archivedChar, &category.Description,
		&category.DisplayTemplate, &category.ExternalID, &category.FulfillmentType,
		&category.InventoryType, &longDescription, &metaDescription, &metaTitle,
		&category.Name, &category.OverrideGeneratedURL, &productDescPattern,
		&productTitlePattern, &rootDisplayOrder, &taxCode, &urlScan, &urlKey,
		&defaultParentCategoryID, &category.CreatedAt, &category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query category by URL: %w", err)
	}

	if activeEndDate.Valid {
		category.ActiveEndDate = &activeEndDate.Time
	}
	if activeStartDate.Valid {
		category.ActiveStartDate = &activeStartDate.Time
	}
	category.Archived = (archivedChar == "Y")
	if longDescription.Valid {
		category.LongDescription = longDescription.String
	}
	if metaDescription.Valid {
		category.MetaDescription = metaDescription.String
	}
	if metaTitle.Valid {
		category.MetaTitle = metaTitle.String
	}
	if productDescPattern.Valid {
		category.ProductDescPattern = productDescPattern.String
	}
	if productTitlePattern.Valid {
		category.ProductTitlePattern = productTitlePattern.String
	}
	if rootDisplayOrder.Valid {
		category.RootDisplayOrder = rootDisplayOrder.Float64
	}
	if taxCode.Valid {
		category.TaxCode = taxCode.String
	}
	if urlScan.Valid {
		category.URL = urlScan.String
	}
	if urlKey.Valid {
		category.URLKey = urlKey.String
	}
	if defaultParentCategoryID.Valid {
		category.DefaultParentCategoryID = &defaultParentCategoryID.Int64
	}

	return &category, nil
}

// FindByURLKey retrieves a category by URL key.
func (r *CategoryRepository) FindByURLKey(ctx context.Context, urlKey string) (*domain.Category, error) {
	query := `
		SELECT
			category_id, active_end_date, active_start_date, archived, description, 
			display_template, external_id, fulfillment_type, inventory_type, 
			long_description, meta_desc, meta_title, name, override_generated_url, 
			product_desc_pattern_override, product_title_pattern_override, 
			root_display_order, tax_code, url, url_key, default_parent_category_id, 
			created_at, updated_at
		FROM blc_category WHERE url_key = $1`

	var category domain.Category
	var activeEndDate sql.NullTime
	var activeStartDate sql.NullTime
	var archivedChar string
	var longDescription sql.NullString
	var metaDescription sql.NullString
	var metaTitle sql.NullString
	var productDescPattern sql.NullString
	var productTitlePattern sql.NullString
	var rootDisplayOrder sql.NullFloat64
	var taxCode sql.NullString
	var url sql.NullString
	var urlKeyScan sql.NullString // Use different variable name to avoid conflict
	var defaultParentCategoryID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, urlKey)
	err := row.Scan(
		&category.ID, &activeEndDate, &activeStartDate, &archivedChar, &category.Description,
		&category.DisplayTemplate, &category.ExternalID, &category.FulfillmentType,
		&category.InventoryType, &longDescription, &metaDescription, &metaTitle,
		&category.Name, &category.OverrideGeneratedURL, &productDescPattern,
		&productTitlePattern, &rootDisplayOrder, &taxCode, &url, &urlKeyScan,
		&defaultParentCategoryID, &category.CreatedAt, &category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query category by URL key: %w", err)
	}

	if activeEndDate.Valid {
		category.ActiveEndDate = &activeEndDate.Time
	}
	if activeStartDate.Valid {
		category.ActiveStartDate = &activeStartDate.Time
	}
	category.Archived = (archivedChar == "Y")
	if longDescription.Valid {
		category.LongDescription = longDescription.String
	}
	if metaDescription.Valid {
		category.MetaDescription = metaDescription.String
	}
	if metaTitle.Valid {
		category.MetaTitle = metaTitle.String
	}
	if productDescPattern.Valid {
		category.ProductDescPattern = productDescPattern.String
	}
	if productTitlePattern.Valid {
		category.ProductTitlePattern = productTitlePattern.String
	}
	if rootDisplayOrder.Valid {
		category.RootDisplayOrder = rootDisplayOrder.Float64
	}
	if taxCode.Valid {
		category.TaxCode = taxCode.String
	}
	if url.Valid {
		category.URL = url.String
	}
	if urlKeyScan.Valid {
		category.URLKey = urlKeyScan.String
	}
	if defaultParentCategoryID.Valid {
		category.DefaultParentCategoryID = &defaultParentCategoryID.Int64
	}

	return &category, nil
}

// FindAll retrieves all categories with pagination
func (r *CategoryRepository) FindAll(ctx context.Context, filter *domain.CategoryFilter) ([]*domain.Category, int64, error) {
	// Base query
	countQuery := `SELECT COUNT(*) FROM blc_category`
	query := `SELECT category_id, active_end_date, active_start_date, archived, description, 
			display_template, external_id, fulfillment_type, inventory_type, 
			long_description, meta_desc, meta_title, name, override_generated_url, 
			product_desc_pattern_override, product_title_pattern_override, 
			root_display_order, tax_code, url, url_key, default_parent_category_id, 
			created_at, updated_at
		FROM blc_category`

	var args []interface{}
	argIdx := 1

	// Build WHERE clause
	whereClauses := []string{}

	if !filter.IncludeArchived {
		whereClauses = append(whereClauses, fmt.Sprintf("archived = $%d", argIdx))
		args = append(args, "N")
		argIdx++
	}
	if filter.ActiveOnly {
		whereClauses = append(whereClauses, fmt.Sprintf("(active_start_date IS NULL OR active_start_date <= NOW()) AND (active_end_date IS NULL OR active_end_date >= NOW() ) "))
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
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	// Apply sorting
	if filter.SortBy != "" {
		orderBy := map[string]string{
			"name":          "name",
			"display_order": "root_display_order",
			"created_at":    "created_at",
		}
		sortColumn, ok := orderBy[filter.SortBy]
		if !ok {
			sortColumn = "name"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, strings.ToUpper(filter.SortOrder))
	}

	// Apply pagination
	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, (filter.Page-1)*filter.PageSize, filter.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query all categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		var activeEndDate sql.NullTime
		var activeStartDate sql.NullTime
		var archivedChar string
		var longDescription sql.NullString
		var metaDescription sql.NullString
		var metaTitle sql.NullString
		var productDescPattern sql.NullString
		var productTitlePattern sql.NullString
		var rootDisplayOrder sql.NullFloat64
		var taxCode sql.NullString
		var url sql.NullString
		var urlKey sql.NullString
		var defaultParentCategoryID sql.NullInt64

		err := rows.Scan(
			&category.ID, &activeEndDate, &activeStartDate, &archivedChar, &category.Description,
			&category.DisplayTemplate, &category.ExternalID, &category.FulfillmentType,
			&category.InventoryType, &longDescription, &metaDescription, &metaTitle,
			&category.Name, &category.OverrideGeneratedURL, &productDescPattern,
			&productTitlePattern, &rootDisplayOrder, &taxCode, &url, &urlKey,
			&defaultParentCategoryID, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan category row: %w", err)
		}

		if activeEndDate.Valid {
			category.ActiveEndDate = &activeEndDate.Time
		}
		if activeStartDate.Valid {
			category.ActiveStartDate = &activeStartDate.Time
		}
		category.Archived = (archivedChar == "Y")
		if longDescription.Valid {
			category.LongDescription = longDescription.String
		}
		if metaDescription.Valid {
			category.MetaDescription = metaDescription.String
		}
		if metaTitle.Valid {
			category.MetaTitle = metaTitle.String
		}
		if productDescPattern.Valid {
			category.ProductDescPattern = productDescPattern.String
		}
		if productTitlePattern.Valid {
			category.ProductTitlePattern = productTitlePattern.String
		}
		if rootDisplayOrder.Valid {
			category.RootDisplayOrder = rootDisplayOrder.Float64
		}
		if taxCode.Valid {
			category.TaxCode = taxCode.String
		}
		if url.Valid {
			category.URL = url.String
		}
		if urlKey.Valid {
			category.URLKey = urlKey.String
		}
		if defaultParentCategoryID.Valid {
			category.DefaultParentCategoryID = &defaultParentCategoryID.Int64
		}
		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during rows iteration: %w", err)
	}

	return categories, totalCount, nil
}

// FindByParentID retrieves child categories by parent ID
func (r *CategoryRepository) FindByParentID(ctx context.Context, parentID int64, filter *domain.CategoryFilter) ([]*domain.Category, int64, error) {
	// Base query
	countQuery := `SELECT COUNT(*) FROM blc_category WHERE default_parent_category_id = $1`
	query := `SELECT category_id, active_end_date, active_start_date, archived, description, 
			display_template, external_id, fulfillment_type, inventory_type, 
			long_description, meta_desc, meta_title, name, override_generated_url, 
			product_desc_pattern_override, product_title_pattern_override, 
			root_display_order, tax_code, url, url_key, default_parent_category_id, 
			created_at, updated_at
		FROM blc_category WHERE default_parent_category_id = $1`

	var args []interface{}
	args = append(args, parentID)
	argIdx := 2

	// Build WHERE clause
	whereClauses := []string{}

	if !filter.IncludeArchived {
		whereClauses = append(whereClauses, fmt.Sprintf("archived = $%d", argIdx))
		args = append(args, "N")
		argIdx++
	}
	if filter.ActiveOnly {
		whereClauses = append(whereClauses, fmt.Sprintf("(active_start_date IS NULL OR active_start_date <= NOW()) AND (active_end_date IS NULL OR active_end_date >= NOW() ) "))
	}

	// Apply WHERE clauses
	if len(whereClauses) > 0 {
		countQuery += " AND " + strings.Join(whereClauses, " AND ")
		query += " AND " + strings.Join(whereClauses, " AND ")
	}

	// Count total results
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count child categories: %w", err)
	}

	// Apply sorting
	if filter.SortBy != "" {
		orderBy := map[string]string{
			"name":          "name",
			"display_order": "root_display_order",
			"created_at":    "created_at",
		}
		sortColumn, ok := orderBy[filter.SortBy]
		if !ok {
			sortColumn = "name"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, strings.ToUpper(filter.SortOrder))
	}

	// Apply pagination
	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, (filter.Page-1)*filter.PageSize, filter.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query child categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		var activeEndDate sql.NullTime
		var activeStartDate sql.NullTime
		var archivedChar string
		var longDescription sql.NullString
		var metaDescription sql.NullString
		var metaTitle sql.NullString
		var productDescPattern sql.NullString
		var productTitlePattern sql.NullString
		var rootDisplayOrder sql.NullFloat64
		var taxCode sql.NullString
		var url sql.NullString
		var urlKey sql.NullString
		var defaultParentCategoryID sql.NullInt64

		err := rows.Scan(
			&category.ID, &activeEndDate, &activeStartDate, &archivedChar, &category.Description,
			&category.DisplayTemplate, &category.ExternalID, &category.FulfillmentType,
			&category.InventoryType, &longDescription, &metaDescription, &metaTitle,
			&category.Name, &category.OverrideGeneratedURL, &productDescPattern,
			&productTitlePattern, &rootDisplayOrder, &taxCode, &url, &urlKey,
			&defaultParentCategoryID, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan category row: %w", err)
		}

		if activeEndDate.Valid {
			category.ActiveEndDate = &activeEndDate.Time
		}
		if activeStartDate.Valid {
			category.ActiveStartDate = &activeStartDate.Time
		}
		category.Archived = (archivedChar == "Y")
		if longDescription.Valid {
			category.LongDescription = longDescription.String
		}
		if metaDescription.Valid {
			category.MetaDescription = metaDescription.String
		}
		if metaTitle.Valid {
			category.MetaTitle = metaTitle.String
		}
		if productDescPattern.Valid {
			category.ProductDescPattern = productDescPattern.String
		}
		if productTitlePattern.Valid {
			category.ProductTitlePattern = productTitlePattern.String
		}
		if rootDisplayOrder.Valid {
			category.RootDisplayOrder = rootDisplayOrder.Float64
		}
		if taxCode.Valid {
			category.TaxCode = taxCode.String
		}
		if url.Valid {
			category.URL = url.String
		}
		if urlKey.Valid {
			category.URLKey = urlKey.String
		}
		if defaultParentCategoryID.Valid {
			category.DefaultParentCategoryID = &defaultParentCategoryID.Int64
		}
		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during rows iteration: %w", err)
	}

	return categories, totalCount, nil
}

// FindRootCategories retrieves root categories (categories with no parent)
func (r *CategoryRepository) FindRootCategories(ctx context.Context, filter *domain.CategoryFilter) ([]*domain.Category, int64, error) {
	// Base query
	countQuery := `SELECT COUNT(*) FROM blc_category WHERE default_parent_category_id IS NULL`
	query := `SELECT category_id, active_end_date, active_start_date, archived, description, 
			display_template, external_id, fulfillment_type, inventory_type, 
			long_description, meta_desc, meta_title, name, override_generated_url, 
			product_desc_pattern_override, product_title_pattern_override, 
			root_display_order, tax_code, url, url_key, default_parent_category_id, 
			created_at, updated_at
		FROM blc_category WHERE default_parent_category_id IS NULL`

	var args []interface{}
	argIdx := 1

	// Build WHERE clause
	whereClauses := []string{}

	if !filter.IncludeArchived {
		whereClauses = append(whereClauses, fmt.Sprintf("archived = $%d", argIdx))
		args = append(args, "N")
		argIdx++
	}
	if filter.ActiveOnly {
		whereClauses = append(whereClauses, fmt.Sprintf("(active_start_date IS NULL OR active_start_date <= NOW()) AND (active_end_date IS NULL OR active_end_date >= NOW() ) "))
	}

	// Apply WHERE clauses
	if len(whereClauses) > 0 {
		countQuery += " AND " + strings.Join(whereClauses, " AND ")
		query += " AND " + strings.Join(whereClauses, " AND ")
	}

	// Count total results
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count root categories: %w", err)
	}

	// Apply sorting
	if filter.SortBy != "" {
		orderBy := map[string]string{
			"name":          "name",
			"display_order": "root_display_order",
			"created_at":    "created_at",
		}
		sortColumn, ok := orderBy[filter.SortBy]
		if !ok {
			sortColumn = "name"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, strings.ToUpper(filter.SortOrder))
	}

	// Apply pagination
	query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", argIdx, argIdx+1)
	args = append(args, (filter.Page-1)*filter.PageSize, filter.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query root categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		var activeEndDate sql.NullTime
		var activeStartDate sql.NullTime
		var archivedChar string
		var longDescription sql.NullString
		var metaDescription sql.NullString
		var metaTitle sql.NullString
		var productDescPattern sql.NullString
		var productTitlePattern sql.NullString
		var rootDisplayOrder sql.NullFloat64
		var taxCode sql.NullString
		var url sql.NullString
		var urlKey sql.NullString
		var defaultParentCategoryID sql.NullInt64

		err := rows.Scan(
			&category.ID, &activeEndDate, &activeStartDate, &archivedChar, &category.Description,
			&category.DisplayTemplate, &category.ExternalID, &category.FulfillmentType,
			&category.InventoryType, &longDescription, &metaDescription, &metaTitle,
			&category.Name, &category.OverrideGeneratedURL, &productDescPattern,
			&productTitlePattern, &rootDisplayOrder, &taxCode, &url, &urlKey,
			&defaultParentCategoryID, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan category row: %w", err)
		}

		if activeEndDate.Valid {
			category.ActiveEndDate = &activeEndDate.Time
		}
		if activeStartDate.Valid {
			category.ActiveStartDate = &activeStartDate.Time
		}
		category.Archived = (archivedChar == "Y")
		if longDescription.Valid {
			category.LongDescription = longDescription.String
		}
		if metaDescription.Valid {
			category.MetaDescription = metaDescription.String
		}
		if metaTitle.Valid {
			category.MetaTitle = metaTitle.String
		}
		if productDescPattern.Valid {
			category.ProductDescPattern = productDescPattern.String
		}
		if productTitlePattern.Valid {
			category.ProductTitlePattern = productTitlePattern.String
		}
		if rootDisplayOrder.Valid {
			category.RootDisplayOrder = rootDisplayOrder.Float64
		}
		if taxCode.Valid {
			category.TaxCode = taxCode.String
		}
		if url.Valid {
			category.URL = url.String
		}
		if urlKey.Valid {
			category.URLKey = urlKey.String
		}
		if defaultParentCategoryID.Valid {
			category.DefaultParentCategoryID = &defaultParentCategoryID.Int64
		}
		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during rows iteration: %w", err)
	}

	return categories, totalCount, nil
}

// GetCategoryPath retrieves the full path from root to category (simplified)
func (r *CategoryRepository) GetCategoryPath(ctx context.Context, categoryID int64) ([]*domain.Category, error) {
	// This is a simplified implementation that only retrieves the direct category itself.
	// A full path would require a recursive query (e.g., using CTEs in PostgreSQL)
	// or multiple queries to build the path. For now, we'll just get the category.
	category, err := r.FindByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category for path: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category with ID %d not found", categoryID)
	}

	path := []*domain.Category{category}
	// To implement full path: fetch parent, then parent's parent, etc. until root.
	// Example (conceptual, not implemented here):
	// currentCategoryID := categoryID
	// for currentCategoryID != 0 && currentCategoryID != nil {
	// 	parentCategory, err := r.FindByID(ctx, *category.DefaultParentCategoryID)
	// 	if err != nil { return nil, err }
	// 	if parentCategory != nil {
	// 		path = append([]*domain.Category{parentCategory}, path...)
	// 		currentCategoryID = *parentCategory.ID
	// 	} else { currentCategoryID = nil }
	// }

	return path, nil
}
