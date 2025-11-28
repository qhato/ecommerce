package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/qhato/ecommerce/pkg/database"
	"github.com/qhato/ecommerce/pkg/errors"
)

// PostgresCategoryRepository implements the CategoryRepository interface
type PostgresCategoryRepository struct {
	db *database.DB
}

// NewPostgresCategoryRepository creates a new PostgreSQL category repository
func NewPostgresCategoryRepository(db *database.DB) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{db: db}
}

// Create creates a new category
func (r *PostgresCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	query := `
		INSERT INTO blc_category (
			category_id, active_end_date, active_start_date, archived,
			description, display_template, external_id, fulfillment_type,
			inventory_type, long_description, meta_desc, meta_title,
			name, override_generated_url, product_desc_pattern_override,
			product_title_pattern_override, root_display_order, tax_code,
			url, url_key, default_parent_category_id
		) VALUES (
			nextval('blc_category_seq'), $1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		) RETURNING category_id`

	archivedFlag := "N"
	if category.Archived {
		archivedFlag = "Y"
	}

	err := r.db.QueryRowContext(ctx, query,
		category.ActiveEndDate,
		category.ActiveStartDate,
		archivedFlag,
		category.Description,
		category.DisplayTemplate,
		category.ExternalID,
		category.FulfillmentType,
		category.InventoryType,
		category.LongDescription,
		category.MetaDescription,
		category.MetaTitle,
		category.Name,
		category.OverrideGeneratedURL,
		category.ProductDescPattern,
		category.ProductTitlePattern,
		category.RootDisplayOrder,
		category.TaxCode,
		category.URL,
		category.URLKey,
		category.DefaultParentCategoryID,
	).Scan(&category.ID)

	if err != nil {
		return errors.Wrap(err, "failed to create category")
	}

	// Insert attributes
	if len(category.Attributes) > 0 {
		if err := r.insertAttributes(ctx, category.ID, category.Attributes); err != nil {
			return err
		}
	}

	return nil
}

// Update updates an existing category
func (r *PostgresCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `
		UPDATE blc_category SET
			active_end_date = $1,
			active_start_date = $2,
			archived = $3,
			description = $4,
			display_template = $5,
			external_id = $6,
			fulfillment_type = $7,
			inventory_type = $8,
			long_description = $9,
			meta_desc = $10,
			meta_title = $11,
			name = $12,
			override_generated_url = $13,
			product_desc_pattern_override = $14,
			product_title_pattern_override = $15,
			root_display_order = $16,
			tax_code = $17,
			url = $18,
			url_key = $19,
			default_parent_category_id = $20
		WHERE category_id = $21`

	archivedFlag := "N"
	if category.Archived {
		archivedFlag = "Y"
	}

	result, err := r.db.ExecContext(ctx, query,
		category.ActiveEndDate,
		category.ActiveStartDate,
		archivedFlag,
		category.Description,
		category.DisplayTemplate,
		category.ExternalID,
		category.FulfillmentType,
		category.InventoryType,
		category.LongDescription,
		category.MetaDescription,
		category.MetaTitle,
		category.Name,
		category.OverrideGeneratedURL,
		category.ProductDescPattern,
		category.ProductTitlePattern,
		category.RootDisplayOrder,
		category.TaxCode,
		category.URL,
		category.URLKey,
		category.DefaultParentCategoryID,
		category.ID,
	)

	if err != nil {
		return errors.Wrap(err, "failed to update category")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("category not found")
	}

	// Update attributes
	if err := r.deleteAttributes(ctx, category.ID); err != nil {
		return err
	}

	if len(category.Attributes) > 0 {
		if err := r.insertAttributes(ctx, category.ID, category.Attributes); err != nil {
			return err
		}
	}

	return nil
}

// Delete soft deletes a category by marking it as archived
func (r *PostgresCategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE blc_category SET archived = 'Y' WHERE category_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete category")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("category not found")
	}

	return nil
}

// FindByID retrieves a category by ID
func (r *PostgresCategoryRepository) FindByID(ctx context.Context, id int64) (*domain.Category, error) {
	query := `
		SELECT
			category_id, active_end_date, active_start_date, archived,
			description, display_template, external_id, fulfillment_type,
			inventory_type, long_description, meta_desc, meta_title,
			name, override_generated_url, product_desc_pattern_override,
			product_title_pattern_override, root_display_order, tax_code,
			url, url_key, default_parent_category_id
		FROM blc_category
		WHERE category_id = $1`

	category := &domain.Category{}
	var archivedFlag string
	var activeEndDate, activeStartDate sql.NullTime
	var parentID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&activeEndDate,
		&activeStartDate,
		&archivedFlag,
		&category.Description,
		&category.DisplayTemplate,
		&category.ExternalID,
		&category.FulfillmentType,
		&category.InventoryType,
		&category.LongDescription,
		&category.MetaDescription,
		&category.MetaTitle,
		&category.Name,
		&category.OverrideGeneratedURL,
		&category.ProductDescPattern,
		&category.ProductTitlePattern,
		&category.RootDisplayOrder,
		&category.TaxCode,
		&category.URL,
		&category.URLKey,
		&parentID,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("category not found")
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find category")
	}

	category.Archived = archivedFlag == "Y"
	if activeEndDate.Valid {
		category.ActiveEndDate = &activeEndDate.Time
	}
	if activeStartDate.Valid {
		category.ActiveStartDate = &activeStartDate.Time
	}
	if parentID.Valid {
		category.DefaultParentCategoryID = &parentID.Int64
	}

	// Load attributes
	attributes, err := r.findAttributes(ctx, id)
	if err != nil {
		return nil, err
	}
	category.Attributes = attributes

	return category, nil
}

// FindByURL retrieves a category by URL
func (r *PostgresCategoryRepository) FindByURL(ctx context.Context, url string) (*domain.Category, error) {
	query := `
		SELECT category_id
		FROM blc_category
		WHERE url = $1 AND archived = 'N'
		LIMIT 1`

	var id int64
	err := r.db.QueryRowContext(ctx, query, url).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("category not found")
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find category by URL")
	}

	return r.FindByID(ctx, id)
}

// FindByURLKey retrieves a category by URL key
func (r *PostgresCategoryRepository) FindByURLKey(ctx context.Context, urlKey string) (*domain.Category, error) {
	query := `
		SELECT category_id
		FROM blc_category
		WHERE url_key = $1 AND archived = 'N'
		LIMIT 1`

	var id int64
	err := r.db.QueryRowContext(ctx, query, urlKey).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("category not found")
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find category by URL key")
	}

	return r.FindByID(ctx, id)
}

// FindAll retrieves all categories with pagination
func (r *PostgresCategoryRepository) FindAll(ctx context.Context, filter *domain.CategoryFilter) ([]*domain.Category, int64, error) {
	// Build where clause
	whereClause := r.buildWhereClause(filter)

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM blc_category %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "failed to count categories")
	}

	// Build main query
	orderByClause := r.buildOrderByClause(filter.SortBy, filter.SortOrder)
	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(`
		SELECT category_id
		FROM blc_category
		%s
		%s
		LIMIT $1 OFFSET $2`,
		whereClause,
		orderByClause,
	)

	rows, err := r.db.QueryContext(ctx, query, filter.PageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to list categories")
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, errors.Wrap(err, "failed to scan category ID")
		}

		category, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	return categories, total, nil
}

// FindByParentID retrieves child categories by parent ID
func (r *PostgresCategoryRepository) FindByParentID(ctx context.Context, parentID int64, filter *domain.CategoryFilter) ([]*domain.Category, int64, error) {
	// Build where clause
	whereClause := fmt.Sprintf("WHERE default_parent_category_id = %d", parentID)
	if !filter.IncludeArchived {
		whereClause += " AND archived = 'N'"
	}
	if filter.ActiveOnly {
		whereClause += " AND (active_start_date IS NULL OR active_start_date <= NOW())"
		whereClause += " AND (active_end_date IS NULL OR active_end_date >= NOW())"
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM blc_category %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "failed to count child categories")
	}

	// Build main query
	orderByClause := r.buildOrderByClause(filter.SortBy, filter.SortOrder)
	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(`
		SELECT category_id
		FROM blc_category
		%s
		%s
		LIMIT $1 OFFSET $2`,
		whereClause,
		orderByClause,
	)

	rows, err := r.db.QueryContext(ctx, query, filter.PageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to list child categories")
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, errors.Wrap(err, "failed to scan category ID")
		}

		category, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	return categories, total, nil
}

// FindRootCategories retrieves root categories
func (r *PostgresCategoryRepository) FindRootCategories(ctx context.Context, filter *domain.CategoryFilter) ([]*domain.Category, int64, error) {
	// Build where clause
	whereClause := "WHERE default_parent_category_id IS NULL"
	if !filter.IncludeArchived {
		whereClause += " AND archived = 'N'"
	}
	if filter.ActiveOnly {
		whereClause += " AND (active_start_date IS NULL OR active_start_date <= NOW())"
		whereClause += " AND (active_end_date IS NULL OR active_end_date >= NOW())"
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM blc_category %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "failed to count root categories")
	}

	// Build main query
	orderByClause := r.buildOrderByClause(filter.SortBy, filter.SortOrder)
	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(`
		SELECT category_id
		FROM blc_category
		%s
		%s
		LIMIT $1 OFFSET $2`,
		whereClause,
		orderByClause,
	)

	rows, err := r.db.QueryContext(ctx, query, filter.PageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to list root categories")
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, errors.Wrap(err, "failed to scan category ID")
		}

		category, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	return categories, total, nil
}

// GetCategoryPath retrieves the full path from root to category
func (r *PostgresCategoryRepository) GetCategoryPath(ctx context.Context, categoryID int64) ([]*domain.Category, error) {
	var path []*domain.Category
	currentID := categoryID

	for currentID != 0 {
		category, err := r.FindByID(ctx, currentID)
		if err != nil {
			return nil, err
		}

		// Prepend to path (to get root -> child order)
		path = append([]*domain.Category{category}, path...)

		// Move to parent
		if category.DefaultParentCategoryID == nil {
			break
		}
		currentID = *category.DefaultParentCategoryID
	}

	return path, nil
}

// Helper methods

func (r *PostgresCategoryRepository) insertAttributes(ctx context.Context, categoryID int64, attributes []domain.CategoryAttribute) error {
	query := `
		INSERT INTO blc_category_attribute (category_attribute_id, name, value, category_id)
		VALUES (nextval('blc_category_attribute_seq'), $1, $2, $3)`

	for _, attr := range attributes {
		_, err := r.db.ExecContext(ctx, query, attr.Name, attr.Value, categoryID)
		if err != nil {
			return errors.Wrap(err, "failed to insert category attribute")
		}
	}

	return nil
}

func (r *PostgresCategoryRepository) deleteAttributes(ctx context.Context, categoryID int64) error {
	query := `DELETE FROM blc_category_attribute WHERE category_id = $1`
	_, err := r.db.ExecContext(ctx, query, categoryID)
	if err != nil {
		return errors.Wrap(err, "failed to delete category attributes")
	}
	return nil
}

func (r *PostgresCategoryRepository) findAttributes(ctx context.Context, categoryID int64) ([]domain.CategoryAttribute, error) {
	query := `
		SELECT category_attribute_id, name, value, category_id
		FROM blc_category_attribute
		WHERE category_id = $1`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find category attributes")
	}
	defer rows.Close()

	var attributes []domain.CategoryAttribute
	for rows.Next() {
		var attr domain.CategoryAttribute
		if err := rows.Scan(&attr.ID, &attr.Name, &attr.Value, &attr.CategoryID); err != nil {
			return nil, errors.Wrap(err, "failed to scan category attribute")
		}
		attributes = append(attributes, attr)
	}

	return attributes, nil
}

func (r *PostgresCategoryRepository) buildWhereClause(filter *domain.CategoryFilter) string {
	conditions := []string{}

	if !filter.IncludeArchived {
		conditions = append(conditions, "archived = 'N'")
	}

	if filter.ActiveOnly {
		conditions = append(conditions, "(active_start_date IS NULL OR active_start_date <= NOW())")
		conditions = append(conditions, "(active_end_date IS NULL OR active_end_date >= NOW())")
	}

	if len(conditions) == 0 {
		return ""
	}

	whereClause := "WHERE " + conditions[0]
	for i := 1; i < len(conditions); i++ {
		whereClause += " AND " + conditions[i]
	}

	return whereClause
}

func (r *PostgresCategoryRepository) buildOrderByClause(sortBy, sortOrder string) string {
	validColumns := map[string]string{
		"name":          "name",
		"display_order": "root_display_order",
		"created_at":    "category_id",
	}

	column, ok := validColumns[sortBy]
	if !ok {
		column = "root_display_order"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	return fmt.Sprintf("ORDER BY %s %s", column, sortOrder)
}
