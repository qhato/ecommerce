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

// PostgresSKURepository implements the SKURepository interface
type PostgresSKURepository struct {
	db *database.DB
}

// NewPostgresSKURepository creates a new PostgreSQL SKU repository
func NewPostgresSKURepository(db *database.DB) *PostgresSKURepository {
	return &PostgresSKURepository{db: db}
}

// Create creates a new SKU
func (r *PostgresSKURepository) Create(ctx context.Context, sku *domain.SKU) error {
	query := `
		INSERT INTO blc_sku (
			sku_id, active_end_date, active_start_date, available_flag,
			cost, description, container_shape, depth, dimension_unit_of_measure,
			girth, height, container_size, width, discountable_flag, display_template,
			external_id, fulfillment_type, inventory_type, is_machine_sortable,
			long_description, name, price, retail_price,
			sale_price, taxable_flag, tax_code, upc, url_key, weight,
			weight_unit_of_measure, currency_code, default_product_id, addl_product_id
		) VALUES (
			nextval('blc_sku_seq'), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27,
			$28, $29, $30, $31, $32
		) RETURNING sku_id`

	availableFlag := "N"
	if sku.Available {
		availableFlag = "Y"
	}
	discountableFlag := "N"
	if sku.Discountable {
		discountableFlag = "Y"
	}
	taxableFlag := "N"
	if sku.Taxable {
		taxableFlag = "Y"
	}

	err := r.db.QueryRow(ctx, query,
		sku.ActiveEndDate,
		sku.ActiveStartDate,
		availableFlag,
		sku.Cost,
		sku.Description,
		sku.ContainerShape,
		sku.Depth,
		sku.DimensionUnitOfMeasure,
		sku.Girth,
		sku.Height,
		sku.ContainerSize,
		sku.Width,
		discountableFlag,
		sku.DisplayTemplate,
		sku.ExternalID,
		sku.FulfillmentType,
		sku.InventoryType,
		sku.IsMachineSortable,
		sku.LongDescription,
		sku.Name,
		sku.RetailPrice, // Using RetailPrice for 'price' column
		sku.RetailPrice,
		sku.SalePrice,
		taxableFlag,
		sku.TaxCode,
		sku.UPC,
		sku.URLKey,
		sku.Weight,
		sku.WeightUnitOfMeasure,
		sku.CurrencyCode,
		sku.DefaultProductID,
		sku.AdditionalProductID,
	).Scan(&sku.ID)

	if err != nil {
		return errors.InternalWrap(err, "failed to create SKU")
	}

	return nil
}

// Update updates an existing SKU
func (r *PostgresSKURepository) Update(ctx context.Context, sku *domain.SKU) error {
	query := `
		UPDATE blc_sku SET
			active_end_date = $1,
			active_start_date = $2,
			available_flag = $3,
			cost = $4,
			description = $5,
			container_shape = $6,
			depth = $7,
			dimension_unit_of_measure = $8,
			girth = $9,
			height = $10,
			container_size = $11,
			width = $12,
			discountable_flag = $13,
			display_template = $14,
			external_id = $15,
			fulfillment_type = $16,
			inventory_type = $17,
			is_machine_sortable = $18,
			long_description = $19,
			name = $20,
			price = $21,
			retail_price = $22,
			sale_price = $23,
			taxable_flag = $24,
			tax_code = $25,
			upc = $26,
			url_key = $27,
			weight = $28,
			weight_unit_of_measure = $29,
			currency_code = $30,
			default_product_id = $31,
			addl_product_id = $32
		WHERE sku_id = $33`

	availableFlag := "N"
	if sku.Available {
		availableFlag = "Y"
	}
	discountableFlag := "N"
	if sku.Discountable {
		discountableFlag = "Y"
	}
	taxableFlag := "N"
	if sku.Taxable {
		taxableFlag = "Y"
	}

	// Using Pool().Exec to get RowsAffected
	tag, err := r.db.Pool().Exec(ctx, query,
		sku.ActiveEndDate,
		sku.ActiveStartDate,
		availableFlag,
		sku.Cost,
		sku.Description,
		sku.ContainerShape,
		sku.Depth,
		sku.DimensionUnitOfMeasure,
		sku.Girth,
		sku.Height,
		sku.ContainerSize,
		sku.Width,
		discountableFlag,
		sku.DisplayTemplate,
		sku.ExternalID,
		sku.FulfillmentType,
		sku.InventoryType,
		sku.IsMachineSortable,
		sku.LongDescription,
		sku.Name,
		sku.RetailPrice,
		sku.RetailPrice,
		sku.SalePrice,
		taxableFlag,
		sku.TaxCode,
		sku.UPC,
		sku.URLKey,
		sku.Weight,
		sku.WeightUnitOfMeasure,
		sku.CurrencyCode,
		sku.DefaultProductID,
		sku.AdditionalProductID,
		sku.ID,
	)

	if err != nil {
		return errors.InternalWrap(err, "failed to update SKU")
	}

	if tag.RowsAffected() == 0 {
		return errors.NotFound("SKU not found")
	}

	return nil
}

// Delete deletes a SKU by ID
func (r *PostgresSKURepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_sku WHERE sku_id = $1`

	tag, err := r.db.Pool().Exec(ctx, query, id)
	if err != nil {
		return errors.InternalWrap(err, "failed to delete SKU")
	}

	if tag.RowsAffected() == 0 {
		return errors.NotFound("SKU not found")
	}

	return nil
}

// FindByID retrieves a SKU by ID
func (r *PostgresSKURepository) FindByID(ctx context.Context, id int64) (*domain.SKU, error) {
	query := `
		SELECT
			sku_id, active_end_date, active_start_date, available_flag,
			cost, description, container_shape, depth, dimension_unit_of_measure,
			girth, height, container_size, width, discountable_flag, display_template,
			external_id, fulfillment_type, inventory_type, is_machine_sortable,
			long_description, name, override_generated_url, price, retail_price,
			sale_price, taxable_flag, tax_code, upc, url_key, weight,
			weight_unit_of_measure, currency_code, default_product_id, addl_product_id
		FROM blc_sku
		WHERE sku_id = $1`

	sku := &domain.SKU{}
	var availableFlag, discountableFlag, taxableFlag string
	var activeEndDate, activeStartDate sql.NullTime
	var defaultProductID, additionalProductID sql.NullInt64
	var overrideGeneratedURL bool // Ignored but scanned
	var price float64             // Ignored but scanned

	err := r.db.QueryRow(ctx, query, id).Scan(
		&sku.ID,
		&activeEndDate,
		&activeStartDate,
		&availableFlag,
		&sku.Cost,
		&sku.Description,
		&sku.ContainerShape,
		&sku.Depth,
		&sku.DimensionUnitOfMeasure,
		&sku.Girth,
		&sku.Height,
		&sku.ContainerSize,
		&sku.Width,
		&discountableFlag,
		&sku.DisplayTemplate,
		&sku.ExternalID,
		&sku.FulfillmentType,
		&sku.InventoryType,
		&sku.IsMachineSortable,
		&sku.LongDescription,
		&sku.Name,
		&overrideGeneratedURL, // Scanned but unused
		&price,                // Scanned but unused
		&sku.RetailPrice,
		&sku.SalePrice,
		&taxableFlag,
		&sku.TaxCode,
		&sku.UPC,
		&sku.URLKey,
		&sku.Weight,
		&sku.WeightUnitOfMeasure,
		&sku.CurrencyCode,
		&defaultProductID,
		&additionalProductID,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("SKU not found")
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find SKU")
	}

	sku.Available = availableFlag == "Y"
	sku.Discountable = discountableFlag == "Y"
	sku.Taxable = taxableFlag == "Y"

	if activeEndDate.Valid {
		sku.ActiveEndDate = &activeEndDate.Time
	}
	if activeStartDate.Valid {
		sku.ActiveStartDate = &activeStartDate.Time
	}
	if defaultProductID.Valid {
		sku.DefaultProductID = &defaultProductID.Int64
	}
	if additionalProductID.Valid {
		sku.AdditionalProductID = &additionalProductID.Int64
	}

	return sku, nil
}

// FindByUPC retrieves a SKU by UPC
func (r *PostgresSKURepository) FindByUPC(ctx context.Context, upc string) (*domain.SKU, error) {
	query := `SELECT sku_id FROM blc_sku WHERE upc = $1 LIMIT 1`

	var id int64
	err := r.db.QueryRow(ctx, query, upc).Scan(&id)
	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("SKU not found")
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find SKU by UPC")
	}

	return r.FindByID(ctx, id)
}

// FindByProductID retrieves SKUs by product ID
func (r *PostgresSKURepository) FindByProductID(ctx context.Context, productID int64) ([]*domain.SKU, error) {
	query := `
		SELECT sku_id
		FROM blc_sku
		WHERE default_product_id = $1 OR addl_product_id = $1
		ORDER BY sku_id`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find SKUs by product")
	}
	defer rows.Close()

	var skus []*domain.SKU
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, errors.InternalWrap(err, "failed to scan SKU ID")
		}

		sku, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		skus = append(skus, sku)
	}

	return skus, nil
}

// FindAll retrieves all SKUs with pagination
func (r *PostgresSKURepository) FindAll(ctx context.Context, filter *domain.SKUFilter) ([]*domain.SKU, int64, error) {
	// Build where clause
	whereClause := r.buildWhereClause(filter)

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM blc_sku %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to count SKUs")
	}

	// Build main query
	orderByClause := r.buildOrderByClause(filter.SortBy, filter.SortOrder)
	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(`
		SELECT sku_id
		FROM blc_sku
		%s
		%s
		LIMIT $1 OFFSET $2`,
		whereClause,
		orderByClause,
	)

	rows, err := r.db.Query(ctx, query, filter.PageSize, offset)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to list SKUs")
	}
	defer rows.Close()

	var skus []*domain.SKU
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, errors.InternalWrap(err, "failed to scan SKU ID")
		}

		sku, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		skus = append(skus, sku)
	}

	return skus, total, nil
}

// UpdateAvailability updates the availability of a SKU
func (r *PostgresSKURepository) UpdateAvailability(ctx context.Context, id int64, available bool) error {
	availableFlag := "N"
	if available {
		availableFlag = "Y"
	}

	query := `UPDATE blc_sku SET available_flag = $1 WHERE sku_id = $2`

	tag, err := r.db.Pool().Exec(ctx, query, availableFlag, id)
	if err != nil {
		return errors.InternalWrap(err, "failed to update SKU availability")
	}

	if tag.RowsAffected() == 0 {
		return errors.NotFound("SKU not found")
	}

	return nil
}

func (r *PostgresSKURepository) buildWhereClause(filter *domain.SKUFilter) string {
	conditions := []string{}

	if filter.AvailableOnly {
		conditions = append(conditions, "available_flag = 'Y'")
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

func (r *PostgresSKURepository) buildOrderByClause(sortBy, sortOrder string) string {
	validColumns := map[string]string{
		"name":       "name",
		"price":      "price",
		"created_at": "sku_id",
	}

	column, ok := validColumns[sortBy]
	if !ok {
		column = "name"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	return fmt.Sprintf("ORDER BY %s %s", column, sortOrder)
}