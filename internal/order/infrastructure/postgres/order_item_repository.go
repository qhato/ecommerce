package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qhato/ecommerce/internal/order/domain"
)

// OrderItemRepository implements domain.OrderItemRepository for PostgreSQL persistence.
type OrderItemRepository struct {
	db *sql.DB
}

// NewOrderItemRepository creates a new PostgreSQL order item repository.
func NewOrderItemRepository(db *sql.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

// Save stores a new order item or updates an existing one.
func (r *OrderItemRepository) Save(ctx context.Context, item *domain.OrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	categoryID := sql.NullInt64{Int64: 0, Valid: false}
	if item.CategoryID != nil {
		categoryID = sql.NullInt64{Int64: *item.CategoryID, Valid: true}
	}
	giftWrapItemID := sql.NullInt64{Int64: 0, Valid: false}
	if item.GiftWrapItemID != nil {
		giftWrapItemID = sql.NullInt64{Int64: *item.GiftWrapItemID, Valid: true}
	}
	parentOrderItemID := sql.NullInt64{Int64: 0, Valid: false}
	if item.ParentOrderItemID != nil {
		parentOrderItemID = sql.NullInt64{Int64: *item.ParentOrderItemID, Valid: true}
	}
	personalMessageID := sql.NullInt64{Int64: 0, Valid: false}
	if item.PersonalMessageID != nil {
		personalMessageID = sql.NullInt64{Int64: *item.PersonalMessageID, Valid: true}
	}

	discountsAllowed := sql.NullBool{Bool: item.DiscountsAllowed, Valid: true}
	hasValidationErrors := sql.NullBool{Bool: item.HasValidationErrors, Valid: true}
	itemTaxableFlag := sql.NullBool{Bool: item.ItemTaxableFlag, Valid: true}
	retailPriceOverride := sql.NullBool{Bool: item.RetailPriceOverride, Valid: true}
	salePriceOverride := sql.NullBool{Bool: item.SalePriceOverride, Valid: true}

	name := sql.NullString{String: item.Name, Valid: item.Name != ""}
	orderItemType := sql.NullString{String: item.OrderItemType, Valid: item.OrderItemType != ""}
	taxCategory := sql.NullString{String: item.TaxCategory, Valid: item.TaxCategory != ""}

	if item.ID == 0 {
		// Insert new order item
		query := `
			INSERT INTO blc_order_item (
				order_id, sku_id, product_id, name, quantity, retail_price, sale_price, price, total_price, 
				tax_amount, tax_category, shipping_amount, discounts_allowed, has_validation_errors, 
				item_taxable_flag, order_item_type, retail_price_override, sale_price_override, 
				category_id, gift_wrap_item_id, parent_order_item_id, personal_message_id, 
				created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
			) RETURNING order_item_id`
		err = tx.QueryRowContext(ctx, query,
			item.OrderID, item.SKUID, item.ProductID, name, item.Quantity, item.RetailPrice, item.SalePrice, item.Price, item.TotalPrice,
			item.TaxAmount, taxCategory, item.ShippingAmount, discountsAllowed, hasValidationErrors,
			itemTaxableFlag, orderItemType, retailPriceOverride, salePriceOverride,
			categoryID, giftWrapItemID, parentOrderItemID, personalMessageID,
			item.CreatedAt, item.UpdatedAt,
		).Scan(&item.ID)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	} else {
		// Update existing order item
		query := `
			UPDATE blc_order_item SET
				order_id = $1, sku_id = $2, product_id = $3, name = $4, quantity = $5, retail_price = $6, sale_price = $7, 
				price = $8, total_price = $9, tax_amount = $10, tax_category = $11, shipping_amount = $12, 
				discounts_allowed = $13, has_validation_errors = $14, item_taxable_flag = $15, 
				order_item_type = $16, retail_price_override = $17, sale_price_override = $18, 
				category_id = $19, gift_wrap_item_id = $20, parent_order_item_id = $21, 
				personal_message_id = $22, updated_at = $23
			WHERE order_item_id = $24`
		_, err = tx.ExecContext(ctx, query,
			item.OrderID, item.SKUID, item.ProductID, name, item.Quantity, item.RetailPrice, item.SalePrice, item.Price, item.TotalPrice,
			item.TaxAmount, taxCategory, item.ShippingAmount, discountsAllowed, hasValidationErrors,
			itemTaxableFlag, orderItemType, retailPriceOverride, salePriceOverride,
			categoryID, giftWrapItemID, parentOrderItemID, personalMessageID,
			item.UpdatedAt, item.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update order item: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves an order item by its unique identifier.
func (r *OrderItemRepository) FindByID(ctx context.Context, id int64) (*domain.OrderItem, error) {
	query := `
		SELECT
			order_item_id, order_id, sku_id, product_id, name, quantity, retail_price, sale_price, price, total_price, 
			tax_amount, tax_category, shipping_amount, discounts_allowed, has_validation_errors, 
			item_taxable_flag, order_item_type, retail_price_override, sale_price_override, 
			category_id, gift_wrap_item_id, parent_order_item_id, personal_message_id, 
			created_at, updated_at
		FROM blc_order_item WHERE order_item_id = $1`

	var item domain.OrderItem
	var categoryID sql.NullInt64
	var giftWrapItemID sql.NullInt64
	var parentOrderItemID sql.NullInt64
	var personalMessageID sql.NullInt64
	var name sql.NullString
	var taxCategory sql.NullString
	var orderItemType sql.NullString
	var discountsAllowed sql.NullBool
	var hasValidationErrors sql.NullBool
	var itemTaxableFlag sql.NullBool
	var retailPriceOverride sql.NullBool
	var salePriceOverride sql.NullBool

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&item.ID, &item.OrderID, &item.SKUID, &item.ProductID, &name, &item.Quantity, &item.RetailPrice, &item.SalePrice, &item.Price, &item.TotalPrice,
		&item.TaxAmount, &taxCategory, &item.ShippingAmount, &discountsAllowed, &hasValidationErrors,
		&itemTaxableFlag, &orderItemType, &retailPriceOverride, &salePriceOverride,
		&categoryID, &giftWrapItemID, &parentOrderItemID, &personalMessageID,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query order item by ID: %w", err)
	}

	if categoryID.Valid {
		item.CategoryID = &categoryID.Int64
	}
	if giftWrapItemID.Valid {
		item.GiftWrapItemID = &giftWrapItemID.Int64
	}
	if parentOrderItemID.Valid {
		item.ParentOrderItemID = &parentOrderItemID.Int64
	}
	if personalMessageID.Valid {
		item.PersonalMessageID = &personalMessageID.Int64
	}
	if name.Valid {
		item.Name = name.String
	}
	if taxCategory.Valid {
		item.TaxCategory = taxCategory.String
	}
	if orderItemType.Valid {
		item.OrderItemType = orderItemType.String
	}
	if discountsAllowed.Valid {
		item.DiscountsAllowed = discountsAllowed.Bool
	}
	if hasValidationErrors.Valid {
		item.HasValidationErrors = hasValidationErrors.Bool
	}
	if itemTaxableFlag.Valid {
		item.ItemTaxableFlag = itemTaxableFlag.Bool
	}
	if retailPriceOverride.Valid {
		item.RetailPriceOverride = retailPriceOverride.Bool
	}
	if salePriceOverride.Valid {
		item.SalePriceOverride = salePriceOverride.Bool
	}

	return &item, nil
}

// FindByOrderID retrieves all order items for a given order ID.
func (r *OrderItemRepository) FindByOrderID(ctx context.Context, orderID int64) ([]*domain.OrderItem, error) {
	query := `
		SELECT
			order_item_id, order_id, sku_id, product_id, name, quantity, retail_price, sale_price, price, total_price, 
			tax_amount, tax_category, shipping_amount, discounts_allowed, has_validation_errors, 
			item_taxable_flag, order_item_type, retail_price_override, sale_price_override, 
			category_id, gift_wrap_item_id, parent_order_item_id, personal_message_id, 
			created_at, updated_at
		FROM blc_order_item WHERE order_id = $1`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items by order ID: %w", err)
	}
	defer rows.Close()

	var items []*domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		var categoryID sql.NullInt64
		var giftWrapItemID sql.NullInt64
		var parentOrderItemID sql.NullInt64
		var personalMessageID sql.NullInt64
		var name sql.NullString
		var taxCategory sql.NullString
		var orderItemType sql.NullString
		var discountsAllowed sql.NullBool
		var hasValidationErrors sql.NullBool
		var itemTaxableFlag sql.NullBool
		var retailPriceOverride sql.NullBool
		var salePriceOverride sql.NullBool

		err := rows.Scan(
			&item.ID, &item.OrderID, &item.SKUID, &item.ProductID, &name, &item.Quantity, &item.RetailPrice, &item.SalePrice, &item.Price, &item.TotalPrice,
			&item.TaxAmount, &taxCategory, &item.ShippingAmount, &discountsAllowed, &hasValidationErrors,
			&itemTaxableFlag, &orderItemType, &retailPriceOverride, &salePriceOverride,
			&categoryID, &giftWrapItemID, &parentOrderItemID, &personalMessageID,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item row: %w", err)
		}

		if categoryID.Valid {
			item.CategoryID = &categoryID.Int64
		}
		if giftWrapItemID.Valid {
			item.GiftWrapItemID = &giftWrapItemID.Int64
		}
		if parentOrderItemID.Valid {
			item.ParentOrderItemID = &parentOrderItemID.Int64
		}
		if personalMessageID.Valid {
			item.PersonalMessageID = &personalMessageID.Int64
		}
		if name.Valid {
			item.Name = name.String
		}
		if taxCategory.Valid {
			item.TaxCategory = taxCategory.String
		}
		if orderItemType.Valid {
			item.OrderItemType = orderItemType.String
		}
		if discountsAllowed.Valid {
			item.DiscountsAllowed = discountsAllowed.Bool
		}
		if hasValidationErrors.Valid {
			item.HasValidationErrors = hasValidationErrors.Bool
		}
		if itemTaxableFlag.Valid {
			item.ItemTaxableFlag = itemTaxableFlag.Bool
		}
		if retailPriceOverride.Valid {
			item.RetailPriceOverride = retailPriceOverride.Bool
		}
		if salePriceOverride.Valid {
			item.SalePriceOverride = salePriceOverride.Bool
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for order items: %w", err)
	}

	return items, nil
}

// Delete removes an order item by its unique identifier.
func (r *OrderItemRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_order_item WHERE order_item_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order item: %w", err)
	}
	return nil
}

// DeleteByOrderID removes all order items for a given order ID.
func (r *OrderItemRepository) DeleteByOrderID(ctx context.Context, orderID int64) error {
	query := `DELETE FROM blc_order_item WHERE order_id = $1`
	_, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order items by order ID: %w", err)
	}
	return nil
}
