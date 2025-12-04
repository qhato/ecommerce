package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/domain"
)

// SKURepository implements domain.SKURepository for PostgreSQL persistence.
type SKURepository struct {
	db *sql.DB
}

// NewSKURepository creates a new PostgreSQL SKU repository.
func NewSKURepository(db *sql.DB) *SKURepository {
	return &SKURepository{db: db}
}

// Save stores a new SKU or updates an existing one.
func (r *SKURepository) Save(ctx context.Context, sku *domain.SKU) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Convert bool flags to BPCHAR(1)
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

	// Handle nullable time fields
	activeEndDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if sku.ActiveEndDate != nil {
		activeEndDate = sql.NullTime{Time: *sku.ActiveEndDate, Valid: true}
	}
	activeStartDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if sku.ActiveStartDate != nil {
		activeStartDate = sql.NullTime{Time: *sku.ActiveStartDate, Valid: true}
	}

	// Handle nullable numeric fields (cost, dimensions, prices, weight)
	cost := sql.NullFloat64{Float64: sku.Cost, Valid: sku.Cost != 0.0}
	depth := sql.NullFloat64{Float64: sku.Depth, Valid: sku.Depth != 0.0}
	girth := sql.NullFloat64{Float64: sku.Girth, Valid: sku.Girth != 0.0}
	height := sql.NullFloat64{Float64: sku.Height, Valid: sku.Height != 0.0}
	width := sql.NullFloat64{Float64: sku.Width, Valid: sku.Width != 0.0}
	retailPrice := sql.NullFloat64{Float64: sku.RetailPrice, Valid: sku.RetailPrice != 0.0}
	salePrice := sql.NullFloat64{Float64: sku.SalePrice, Valid: sku.SalePrice != 0.0}
	weight := sql.NullFloat64{Float64: sku.Weight, Valid: sku.Weight != 0.0}

	// Handle nullable string fields
	description := sql.NullString{String: sku.Description, Valid: sku.Description != ""}
	longDescription := sql.NullString{String: sku.LongDescription, Valid: sku.LongDescription != ""}
	containerShape := sql.NullString{String: sku.ContainerShape, Valid: sku.ContainerShape != ""}
	dimensionUnitOfMeasure := sql.NullString{String: sku.DimensionUnitOfMeasure, Valid: sku.DimensionUnitOfMeasure != ""}
	containerSize := sql.NullString{String: sku.ContainerSize, Valid: sku.ContainerSize != ""}
	displayTemplate := sql.NullString{String: sku.DisplayTemplate, Valid: sku.DisplayTemplate != ""}
	externalID := sql.NullString{String: sku.ExternalID, Valid: sku.ExternalID != ""}
	fulfillmentType := sql.NullString{String: sku.FulfillmentType, Valid: sku.FulfillmentType != ""}
	inventoryType := sql.NullString{String: sku.InventoryType, Valid: sku.InventoryType != ""}
	name := sql.NullString{String: sku.Name, Valid: sku.Name != ""}
	taxCode := sql.NullString{String: sku.TaxCode, Valid: sku.TaxCode != ""}
	upc := sql.NullString{String: sku.UPC, Valid: sku.UPC != ""}
	urlKey := sql.NullString{String: sku.URLKey, Valid: sku.URLKey != ""}
	weightUnitOfMeasure := sql.NullString{String: sku.WeightUnitOfMeasure, Valid: sku.WeightUnitOfMeasure != ""}
	currencyCode := sql.NullString{String: sku.CurrencyCode, Valid: sku.CurrencyCode != ""}

	// Handle nullable int64 (foreign keys)
	defaultProductID := sql.NullInt64{Int64: 0, Valid: false}
	if sku.DefaultProductID != nil {
		defaultProductID = sql.NullInt64{Int64: *sku.DefaultProductID, Valid: true}
	}
	additionalProductID := sql.NullInt64{Int64: 0, Valid: false}
	if sku.AdditionalProductID != nil {
		additionalProductID = sql.NullInt64{Int64: *sku.AdditionalProductID, Valid: true}
	}

	if sku.ID == 0 {
		// Insert new SKU
		query := `
			INSERT INTO blc_sku (
				active_end_date, active_start_date, available_flag, cost, description, 
				container_shape, depth, dimension_unit_of_measure, girth, height, 
				container_size, width, discountable_flag, display_template, external_id, 
				fulfillment_type, inventory_type, is_machine_sortable, long_description, 
				name, quantity_available, retail_price, sale_price, tax_code, 
				taxable_flag, upc, url_key, weight, weight_unit_of_measure, 
				currency_code, default_product_id, addl_product_id, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34
			) RETURNING sku_id`
		err = tx.QueryRowContext(ctx, query,
			activeEndDate, activeStartDate, availableFlag, cost, description,
			containerShape, depth, dimensionUnitOfMeasure, girth, height,
			containerSize, width, discountableFlag, displayTemplate, externalID,
			fulfillmentType, inventoryType, sku.IsMachineSortable, longDescription,
			name, nil, retailPrice, salePrice, taxCode,
			taxableFlag, upc, urlKey, weight, weightUnitOfMeasure,
			currencyCode, defaultProductID, additionalProductID, sku.CreatedAt, sku.UpdatedAt,
		).Scan(&sku.ID)
		if err != nil {
			return fmt.Errorf("failed to insert SKU: %w", err)
		}
	} else {
		// Update existing SKU
		query := `
			UPDATE blc_sku SET
				active_end_date = $1, active_start_date = $2, available_flag = $3, cost = $4, 
				description = $5, container_shape = $6, depth = $7, 
				dimension_unit_of_measure = $8, girth = $9, height = $10, 
				container_size = $11, width = $12, discountable_flag = $13, 
				display_template = $14, external_id = $15, fulfillment_type = $16, 
				inventory_type = $17, is_machine_sortable = $18, long_description = $19, 
				name = $20, quantity_available = $21, retail_price = $22, 
				sale_price = $23, tax_code = $24, taxable_flag = $25, 
				upc = $26, url_key = $27, weight = $28, 
				weight_unit_of_measure = $29, currency_code = $30, 
				default_product_id = $31, addl_product_id = $32, updated_at = $33
			WHERE sku_id = $34`
		_, err = tx.ExecContext(ctx, query,
			activeEndDate, activeStartDate, availableFlag, cost, description,
			containerShape, depth, dimensionUnitOfMeasure, girth, height,
			containerSize, width, discountableFlag, displayTemplate, externalID,
			fulfillmentType, inventoryType, sku.IsMachineSortable, longDescription,
			name, nil, retailPrice, salePrice, taxCode,
			taxableFlag, upc, urlKey, weight, weightUnitOfMeasure,
			currencyCode, defaultProductID, additionalProductID, sku.UpdatedAt, sku.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update SKU: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a SKU by its unique identifier.
func (r *SKURepository) FindByID(ctx context.Context, id int64) (*domain.SKU, error) {
	query := `
		SELECT
			sku_id, active_end_date, active_start_date, available_flag, cost, 
			description, container_shape, depth, dimension_unit_of_measure, girth, 
			height, container_size, width, discountable_flag, display_template, 
			external_id, fulfillment_type, inventory_type, is_machine_sortable, long_description, 
			name, quantity_available, retail_price, sale_price, tax_code, 
			taxable_flag, upc, url_key, weight, weight_unit_of_measure, 
			currency_code, default_product_id, addl_product_id, created_at, updated_at
		FROM blc_sku WHERE sku_id = $1`

	var sku domain.SKU
	var activeEndDate sql.NullTime
	var activeStartDate sql.NullTime
	var availableFlagChar string
	var cost sql.NullFloat64
	var description sql.NullString
	var containerShape sql.NullString
	var depth sql.NullFloat64
	var dimensionUnitOfMeasure sql.NullString
	var girth sql.NullFloat64
	var height sql.NullFloat64
	var containerSize sql.NullString
	var width sql.NullFloat64
	var discountableFlagChar string
	var displayTemplate sql.NullString
	var externalID sql.NullString
	var fulfillmentType sql.NullString
	var inventoryType sql.NullString
	var longDescription sql.NullString
	var name sql.NullString
	var quantityAvailable sql.NullInt32
	var retailPrice sql.NullFloat64
	var salePrice sql.NullFloat64
	var taxCode sql.NullString
	var taxableFlagChar string
	var upc sql.NullString
	var urlKey sql.NullString
	var weight sql.NullFloat64
	var weightUnitOfMeasure sql.NullString
	var currencyCode sql.NullString
	var defaultProductID sql.NullInt64
	var additionalProductID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&sku.ID, &activeEndDate, &activeStartDate, &availableFlagChar, &cost,
		&description, &containerShape, &depth, &dimensionUnitOfMeasure, &girth,
		&height, &containerSize, &width, &discountableFlagChar, &displayTemplate,
		&externalID, &fulfillmentType, &inventoryType, &sku.IsMachineSortable, &longDescription,
		&name, &quantityAvailable, &retailPrice, &salePrice, &taxCode,
		&taxableFlagChar, &upc, &urlKey, &weight, &weightUnitOfMeasure,
		&currencyCode, &defaultProductID, &additionalProductID, &sku.CreatedAt, &sku.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query SKU by ID: %w", err)
	}

	if activeEndDate.Valid {
		sku.ActiveEndDate = &activeEndDate.Time
	}
	if activeStartDate.Valid {
		sku.ActiveStartDate = &activeStartDate.Time
	}
	sku.Available = (availableFlagChar == "Y")
	if cost.Valid {
		sku.Cost = cost.Float64
	}
	if description.Valid {
		sku.Description = description.String
	}
	if containerShape.Valid {
		sku.ContainerShape = containerShape.String
	}
	if depth.Valid {
		sku.Depth = depth.Float64
	}
	if dimensionUnitOfMeasure.Valid {
		sku.DimensionUnitOfMeasure = dimensionUnitOfMeasure.String
	}
	if girth.Valid {
		sku.Girth = girth.Float64
	}
	if height.Valid {
		sku.Height = height.Float64
	}
	if containerSize.Valid {
		sku.ContainerSize = containerSize.String
	}
	if width.Valid {
		sku.Width = width.Float64
	}
	sku.Discountable = (discountableFlagChar == "Y")
	if displayTemplate.Valid {
		sku.DisplayTemplate = displayTemplate.String
	}
	if externalID.Valid {
		sku.ExternalID = externalID.String
	}
	if fulfillmentType.Valid {
		sku.FulfillmentType = fulfillmentType.String
	}
	if inventoryType.Valid {
		sku.InventoryType = inventoryType.String
	}
	if longDescription.Valid {
		sku.LongDescription = longDescription.String
	}
	if name.Valid {
		sku.Name = name.String
	}
	// quantityAvailable is intentionally not mapped to sku.QuantityAvailable if it's meant to be managed by Inventory.
	if retailPrice.Valid {
		sku.RetailPrice = retailPrice.Float64
	}
	if salePrice.Valid {
		sku.SalePrice = salePrice.Float64
	}
	if taxCode.Valid {
		sku.TaxCode = taxCode.String
	}
	sku.Taxable = (taxableFlagChar == "Y")
	if upc.Valid {
		sku.UPC = upc.String
	}
	if urlKey.Valid {
		sku.URLKey = urlKey.String
	}
	if weight.Valid {
		sku.Weight = weight.Float64
	}
	if weightUnitOfMeasure.Valid {
		sku.WeightUnitOfMeasure = weightUnitOfMeasure.String
	}
	if currencyCode.Valid {
		sku.CurrencyCode = currencyCode.String
	}
	if defaultProductID.Valid {
		sku.DefaultProductID = &defaultProductID.Int64
	}
	if additionalProductID.Valid {
		sku.AdditionalProductID = &additionalProductID.Int64
	}

	return &sku, nil
}

// FindByUPC retrieves a SKU by its Universal Product Code.
func (r *SKURepository) FindByUPC(ctx context.Context, upc string) (*domain.SKU, error) {
	query := `
		SELECT
			sku_id, active_end_date, active_start_date, available_flag, cost, 
			description, container_shape, depth, dimension_unit_of_measure, girth, 
			height, container_size, width, discountable_flag, display_template, 
			external_id, fulfillment_type, inventory_type, is_machine_sortable, long_description, 
			name, quantity_available, retail_price, sale_price, tax_code, 
			taxable_flag, upc, url_key, weight, weight_unit_of_measure, 
			currency_code, default_product_id, addl_product_id, created_at, updated_at
		FROM blc_sku WHERE upc = $1`

	var sku domain.SKU
	var activeEndDate sql.NullTime
	var activeStartDate sql.NullTime
	var availableFlagChar string
	var cost sql.NullFloat64
	var description sql.NullString
	var containerShape sql.NullString
	var depth sql.NullFloat64
	var dimensionUnitOfMeasure sql.NullString
	var girth sql.NullFloat64
	var height sql.NullFloat64
	var containerSize sql.NullString
	var width sql.NullFloat64
	var discountableFlagChar string
	var displayTemplate sql.NullString
	var externalID sql.NullString
	var fulfillmentType sql.NullString
	var inventoryType sql.NullString
	var longDescription sql.NullString
	var name sql.NullString
	var quantityAvailable sql.NullInt32
	var retailPrice sql.NullFloat64
	var salePrice sql.NullFloat64
	var taxCode sql.NullString
	var taxableFlagChar string
	var upcScan sql.NullString // Use a different name to avoid conflict
	var urlKey sql.NullString
	var weight sql.NullFloat64
	var weightUnitOfMeasure sql.NullString
	var currencyCode sql.NullString
	var defaultProductID sql.NullInt64
	var additionalProductID sql.NullInt64

	row := r.db.QueryRowContext(ctx, query, upc)
	err := row.Scan(
		&sku.ID, &activeEndDate, &activeStartDate, &availableFlagChar, &cost,
		&description, &containerShape, &depth, &dimensionUnitOfMeasure, &girth,
		&height, &containerSize, &width, &discountableFlagChar, &displayTemplate,
		&externalID, &fulfillmentType, &inventoryType, &sku.IsMachineSortable, &longDescription,
		&name, &quantityAvailable, &retailPrice, &salePrice, &taxCode,
		&taxableFlagChar, &upcScan, &urlKey, &weight, &weightUnitOfMeasure,
		&currencyCode, &defaultProductID, &additionalProductID, &sku.CreatedAt, &sku.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query SKU by UPC: %w", err)
	}

	if activeEndDate.Valid {
		sku.ActiveEndDate = &activeEndDate.Time
	}
	if activeStartDate.Valid {
		sku.ActiveStartDate = &activeStartDate.Time
	}
	sku.Available = (availableFlagChar == "Y")
	if cost.Valid {
		sku.Cost = cost.Float64
	}
	if description.Valid {
		sku.Description = description.String
	}
	if containerShape.Valid {
		sku.ContainerShape = containerShape.String
	}
	if depth.Valid {
		sku.Depth = depth.Float64
	}
	if dimensionUnitOfMeasure.Valid {
		sku.DimensionUnitOfMeasure = dimensionUnitOfMeasure.String
	}
	if girth.Valid {
		sku.Girth = girth.Float64
	}
	if height.Valid {
		sku.Height = height.Float64
	}
	if containerSize.Valid {
		sku.ContainerSize = containerSize.String
	}
	if width.Valid {
		sku.Width = width.Float64
	}
	sku.Discountable = (discountableFlagChar == "Y")
	if displayTemplate.Valid {
		sku.DisplayTemplate = displayTemplate.String
	}
	if externalID.Valid {
		sku.ExternalID = externalID.String
	}
	if fulfillmentType.Valid {
		sku.FulfillmentType = fulfillmentType.String
	}
	if inventoryType.Valid {
		sku.InventoryType = inventoryType.String
	}
	if longDescription.Valid {
		sku.LongDescription = longDescription.String
	}
	if name.Valid {
		sku.Name = name.String
	}
	// quantityAvailable is intentionally not mapped to sku.QuantityAvailable if it's meant to be managed by Inventory.
	if retailPrice.Valid {
		sku.RetailPrice = retailPrice.Float64
	}
	if salePrice.Valid {
		sku.SalePrice = salePrice.Float64
	}
	if taxCode.Valid {
		sku.TaxCode = taxCode.String
	}
	sku.Taxable = (taxableFlagChar == "Y")
	if upcScan.Valid {
		sku.UPC = upcScan.String
	}
	if urlKey.Valid {
		sku.URLKey = urlKey.String
	}
	if weight.Valid {
		sku.Weight = weight.Float64
	}
	if weightUnitOfMeasure.Valid {
		sku.WeightUnitOfMeasure = weightUnitOfMeasure.String
	}
	if currencyCode.Valid {
		sku.CurrencyCode = currencyCode.String
	}
	if defaultProductID.Valid {
		sku.DefaultProductID = &defaultProductID.Int64
	}
	if additionalProductID.Valid {
		sku.AdditionalProductID = &additionalProductID.Int64
	}

	return &sku, nil
}

// FindByProductID retrieves SKUs by product ID
func (r *SKURepository) FindByProductID(ctx context.Context, productID int64) ([]*domain.SKU, error) {
	query := `
		SELECT
			sku_id, active_end_date, active_start_date, available_flag, cost, 
			description, container_shape, depth, dimension_unit_of_measure, girth, 
			height, container_size, width, discountable_flag, display_template, 
			external_id, fulfillment_type, inventory_type, is_machine_sortable, long_description, 
			name, quantity_available, retail_price, sale_price, tax_code, 
			taxable_flag, upc, url_key, weight, weight_unit_of_measure, 
			currency_code, default_product_id, addl_product_id, created_at, updated_at
		FROM blc_sku WHERE default_product_id = $1 OR addl_product_id = $1`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query SKUs by product ID: %w", err)
	}
	defer rows.Close()

	var skus []*domain.SKU
	for rows.Next() {
		var sku domain.SKU
		var activeEndDate sql.NullTime
		var activeStartDate sql.NullTime
		var availableFlagChar string
		var cost sql.NullFloat64
		var description sql.NullString
		var containerShape sql.NullString
		var depth sql.NullFloat64
		var dimensionUnitOfMeasure sql.NullString
		var girth sql.NullFloat64
		var height sql.NullFloat64
		var containerSize sql.NullString
		var width sql.NullFloat64
		var discountableFlagChar string
		var displayTemplate sql.NullString
		var externalID sql.NullString
		var fulfillmentType sql.NullString
		var inventoryType sql.NullString
		var longDescription sql.NullString
		var name sql.NullString
		var quantityAvailable sql.NullInt32
		var retailPrice sql.NullFloat64
		var salePrice sql.NullFloat64
		var taxCode sql.NullString
		var taxableFlagChar string
		var upc sql.NullString
		var urlKey sql.NullString
		var weight sql.NullFloat64
		var weightUnitOfMeasure sql.NullString
		var currencyCode sql.NullString
		var defaultProductID sql.NullInt64
		var additionalProductID sql.NullInt64

		err := rows.Scan(
			&sku.ID, &activeEndDate, &activeStartDate, &availableFlagChar, &cost,
			&description, &containerShape, &depth, &dimensionUnitOfMeasure, &girth,
			&height, &containerSize, &width, &discountableFlagChar, &displayTemplate,
			&externalID, &fulfillmentType, &inventoryType, &sku.IsMachineSortable, &longDescription,
			&name, &quantityAvailable, &retailPrice, &salePrice, &taxCode,
			&taxableFlagChar, &upc, &urlKey, &weight, &weightUnitOfMeasure,
			&currencyCode, &defaultProductID, &additionalProductID, &sku.CreatedAt, &sku.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan SKU row: %w", err)
		}

		if activeEndDate.Valid {
			sku.ActiveEndDate = &activeEndDate.Time
		}
		if activeStartDate.Valid {
			sku.ActiveStartDate = &activeStartDate.Time
		}
		sku.Available = (availableFlagChar == "Y")
		if cost.Valid {
			sku.Cost = cost.Float64
		}
		if description.Valid {
			sku.Description = description.String
	}
	if containerShape.Valid {
		sku.ContainerShape = containerShape.String
	}
	if depth.Valid {
		sku.Depth = depth.Float64
	}
	if dimensionUnitOfMeasure.Valid {
		sku.DimensionUnitOfMeasure = dimensionUnitOfMeasure.String
	}
	if girth.Valid {
		sku.Girth = girth.Float64
	}
	if height.Valid {
		sku.Height = height.Float64
	}
	if containerSize.Valid {
		sku.ContainerSize = containerSize.String
	}
	if width.Valid {
		sku.Width = width.Float64
	}
	sku.Discountable = (discountableFlagChar == "Y")
	if displayTemplate.Valid {
		sku.DisplayTemplate = displayTemplate.String
	}
	if externalID.Valid {
		sku.ExternalID = externalID.String
	}
	if fulfillmentType.Valid {
		sku.FulfillmentType = fulfillmentType.String
	}
	if inventoryType.Valid {
		sku.InventoryType = inventoryType.String
	}
	if longDescription.Valid {
		sku.LongDescription = longDescription.String
	}
	if name.Valid {
		sku.Name = name.String
	}
	// quantityAvailable is intentionally not mapped to sku.QuantityAvailable if it's meant to be managed by Inventory.
	if retailPrice.Valid {
		sku.RetailPrice = retailPrice.Float64
	}
	if salePrice.Valid {
		sku.SalePrice = salePrice.Float64
	}
	if taxCode.Valid {
		sku.TaxCode = taxCode.String
	}
	sku.Taxable = (taxableFlagChar == "Y")
	if upc.Valid {
		sku.UPC = upc.String
	}
	if urlKey.Valid {
		sku.URLKey = urlKey.String
	}
	if weight.Valid {
		sku.Weight = weight.Float64
	}
	if weightUnitOfMeasure.Valid {
		sku.WeightUnitOfMeasure = weightUnitOfMeasure.String
	}
	if currencyCode.Valid {
		sku.CurrencyCode = currencyCode.String
	}
	if defaultProductID.Valid {
		sku.DefaultProductID = &defaultProductID.Int64
	}
	if additionalProductID.Valid {
		sku.AdditionalProductID = &additionalProductID.Int64
	}
	skus = append(skus, &sku)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return skus, nil
}

// FindAll retrieves all SKUs with pagination
func (r *SKURepository) FindAll(ctx context.Context, filter *domain.SKUFilter) ([]*domain.SKU, int64, error) {
	// Base query
	countQuery := `SELECT COUNT(*) FROM blc_sku`
	query := `SELECT sku_id, active_end_date, active_start_date, available_flag, cost, 
			description, container_shape, depth, dimension_unit_of_measure, girth, 
			height, container_size, width, discountable_flag, display_template, 
			external_id, fulfillment_type, inventory_type, is_machine_sortable, long_description, 
			name, quantity_available, retail_price, sale_price, tax_code, 
			taxable_flag, upc, url_key, weight, weight_unit_of_measure, 
			currency_code, default_product_id, addl_product_id, created_at, updated_at
		FROM blc_sku`

	var args []interface{}
	argIdx := 1

	// Build WHERE clause
	whereClauses := []string{}

	if filter.AvailableOnly {
		whereClauses = append(whereClauses, fmt.Sprintf("available_flag = $%d", argIdx))
		args = append(args, "Y")
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
		return nil, 0, fmt.Errorf("failed to count SKUs: %w", err)
	}

	// Apply sorting
	if filter.SortBy != "" {
		orderBy := map[string]string{
			"name":       "name",
			"price":      "retail_price", // Or sale_price, depending on desired sorting
			"created_at": "created_at",
			"updated_at": "updated_at",
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
		return nil, 0, fmt.Errorf("failed to query all SKUs: %w", err)
	}
	defer rows.Close()

	var skus []*domain.SKU
	for rows.Next() {
		var sku domain.SKU
		var activeEndDate sql.NullTime
		var activeStartDate sql.NullTime
		var availableFlagChar string
		var cost sql.NullFloat64
		var description sql.NullString
		var containerShape sql.NullString
		var depth sql.NullFloat64
		var dimensionUnitOfMeasure sql.NullString
		var girth sql.NullFloat64
		var height sql.NullFloat64
		var containerSize sql.NullString
		var width sql.NullFloat64
		var discountableFlagChar string
		var displayTemplate sql.NullString
		var externalID sql.NullString
		var fulfillmentType sql.NullString
		var inventoryType sql.NullString
		var longDescription sql.NullString
		var name sql.NullString
		var quantityAvailable sql.NullInt32
		var retailPrice sql.NullFloat64
		var salePrice sql.NullFloat64
		var taxCode sql.NullString
		var taxableFlagChar string
		var upc sql.NullString
		var urlKey sql.NullString
		var weight sql.NullFloat64
		var weightUnitOfMeasure sql.NullString
		var currencyCode sql.NullString
		var defaultProductID sql.NullInt64
		var additionalProductID sql.NullInt64

		err := rows.Scan(
			&sku.ID, &activeEndDate, &activeStartDate, &availableFlagChar, &cost,
			&description, &containerShape, &depth, &dimensionUnitOfMeasure, &girth,
			&height, &containerSize, &width, &discountableFlagChar, &displayTemplate,
			&externalID, &fulfillmentType, &inventoryType, &sku.IsMachineSortable, &longDescription,
			&name, &quantityAvailable, &retailPrice, &salePrice, &taxCode,
			&taxableFlagChar, &upc, &urlKey, &weight, &weightUnitOfMeasure,
			&currencyCode, &defaultProductID, &additionalProductID, &sku.CreatedAt, &sku.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan SKU row: %w", err)
		}

		if activeEndDate.Valid {
			sku.ActiveEndDate = &activeEndDate.Time
		}
		if activeStartDate.Valid {
			sku.ActiveStartDate = &activeStartDate.Time
		}
		sku.Available = (availableFlagChar == "Y")
		if cost.Valid {
			sku.Cost = cost.Float64
		}
		if description.Valid {
			sku.Description = description.String
		}
		if containerShape.Valid {
			sku.ContainerShape = containerShape.String
		}
		if depth.Valid {
			sku.Depth = depth.Float64
		}
		if dimensionUnitOfMeasure.Valid {
			sku.DimensionUnitOfMeasure = dimensionUnitOfMeasure.String
		}
		if girth.Valid {
			sku.Girth = girth.Float64
		}
		if height.Valid {
			sku.Height = height.Float64
		}
		if containerSize.Valid {
			sku.ContainerSize = containerSize.String
		}
		if width.Valid {
			sku.Width = width.Float64
		}
		sku.Discountable = (discountableFlagChar == "Y")
		if displayTemplate.Valid {
		sku.DisplayTemplate = displayTemplate.String
	}
		if externalID.Valid {
			sku.ExternalID = externalID.String
		}
		if fulfillmentType.Valid {
			sku.FulfillmentType = fulfillmentType.String
		}
		if inventoryType.Valid {
			sku.InventoryType = inventoryType.String
		}
		if longDescription.Valid {
			sku.LongDescription = longDescription.String
		}
		if name.Valid {
			sku.Name = name.String
		}
		// quantityAvailable is intentionally not mapped to sku.QuantityAvailable if it's meant to be managed by Inventory.
		if retailPrice.Valid {
			sku.RetailPrice = retailPrice.Float64
		}
		if salePrice.Valid {
			sku.SalePrice = salePrice.Float64
		}
		if taxCode.Valid {
			sku.TaxCode = taxCode.String
		}
		sku.Taxable = (taxableFlagChar == "Y")
		if upc.Valid {
			sku.UPC = upc.String
		}
		if urlKey.Valid {
			sku.URLKey = urlKey.String
		}
		if weight.Valid {
			sku.Weight = weight.Float64
		}
		if weightUnitOfMeasure.Valid {
			sku.WeightUnitOfMeasure = weightUnitOfMeasure.String
		}
		if currencyCode.Valid {
			sku.CurrencyCode = currencyCode.String
		}
		if defaultProductID.Valid {
			sku.DefaultProductID = &defaultProductID.Int64
		}
		if additionalProductID.Valid {
			sku.AdditionalProductID = &additionalProductID.Int64
		}
		skus = append(skus, &sku)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during rows iteration: %w", err)
	}

	return skus, totalCount, nil
}

// UpdateAvailability updates the availability of a SKU.
func (r *SKURepository) UpdateAvailability(ctx context.Context, id int64, available bool) error {
	query := `UPDATE blc_sku SET available_flag = $1, updated_at = $2 WHERE sku_id = $3`
	availableFlag := "N"
	if available {
		availableFlag = "Y"
	}
	_, err := r.db.ExecContext(ctx, query, availableFlag, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update SKU availability: %w", err)
	}
	return nil
}

// Delete deletes a SKU by its unique identifier.
func (r *SKURepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_sku WHERE sku_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete SKU: %w", err)
	}
	return nil
}
