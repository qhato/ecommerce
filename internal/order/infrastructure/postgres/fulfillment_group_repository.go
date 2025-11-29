package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/order/domain"
)

// FulfillmentGroupRepository implements domain.FulfillmentGroupRepository for PostgreSQL persistence.
type FulfillmentGroupRepository struct {
	db *sql.DB
}

// NewFulfillmentGroupRepository creates a new PostgreSQL fulfillment group repository.
func NewFulfillmentGroupRepository(db *sql.DB) *FulfillmentGroupRepository {
	return &FulfillmentGroupRepository{db: db}
}

// Save stores a new fulfillment group or updates an existing one.
func (r *FulfillmentGroupRepository) Save(ctx context.Context, fg *domain.FulfillmentGroup) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle nullable fields
	shippingPriceTaxable := sql.NullBool{Bool: fg.ShippingPriceTaxable, Valid: true}
	isPrimary := sql.NullBool{Bool: fg.IsPrimary, Valid: true}
	shippingOverride := sql.NullBool{Bool: fg.ShippingOverride, Valid: true}

	addressID := sql.NullInt64{Int64: 0, Valid: false}
	if fg.AddressID != nil {
		addressID = sql.NullInt64{Int64: *fg.AddressID, Valid: true}
	}
	fulfillmentOptionID := sql.NullInt64{Int64: 0, Valid: false}
	if fg.FulfillmentOptionID != nil {
		fulfillmentOptionID = sql.NullInt64{Int64: *fg.FulfillmentOptionID, Valid: true}
	}
	personalMessageID := sql.NullInt64{Int64: 0, Valid: false}
	if fg.PersonalMessageID != nil {
		personalMessageID = sql.NullInt64{Int64: *fg.PersonalMessageID, Valid: true}
	}
	phoneID := sql.NullInt64{Int64: 0, Valid: false}
	if fg.PhoneID != nil {
		phoneID = sql.NullInt64{Int64: *fg.PhoneID, Valid: true}
	}

	deliveryInstruction := sql.NullString{String: fg.DeliveryInstruction, Valid: fg.DeliveryInstruction != ""}
	method := sql.NullString{String: fg.Method, Valid: fg.Method != ""}
	referenceNumber := sql.NullString{String: fg.ReferenceNumber, Valid: fg.ReferenceNumber != ""}
	service := sql.NullString{String: fg.Service, Valid: fg.Service != ""}
	status := sql.NullString{String: fg.Status, Valid: fg.Status != ""}
	fgType := sql.NullString{String: fg.Type, Valid: fg.Type != ""}

	if fg.ID == 0 {
		// Insert new fulfillment group
		query := `
			INSERT INTO blc_fulfillment_group (
				order_id, type, shipping_price, shipping_price_taxable, merchandise_total, 
				method, is_primary, reference_number, retail_price, sale_price, 
				fulfillment_group_sequnce, service, shipping_override, status, total, 
				total_fee_tax, total_fg_tax, total_item_tax, total_tax, 
				address_id, fulfillment_option_id, personal_message_id, phone_id, 
				delivery_instruction, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
			) RETURNING fulfillment_group_id`
		err = tx.QueryRowContext(ctx, query,
			fg.OrderID, fgType, fg.ShippingPrice, shippingPriceTaxable, fg.MerchandiseTotal,
			method, isPrimary, referenceNumber, fg.RetailPrice, fg.SalePrice,
			fg.Sequence, service, shippingOverride, status, fg.Total,
			fg.TotalFeeTax, fg.TotalFgTax, fg.TotalItemTax, fg.TotalTax,
			addressID, fulfillmentOptionID, personalMessageID, phoneID,
			deliveryInstruction, fg.CreatedAt, fg.UpdatedAt,
		).Scan(&fg.ID)
		if err != nil {
			return fmt.Errorf("failed to insert fulfillment group: %w", err)
		}
	} else {
		// Update existing fulfillment group
		query := `
			UPDATE blc_fulfillment_group SET
				order_id = $1, type = $2, shipping_price = $3, shipping_price_taxable = $4, 
				merchandise_total = $5, method = $6, is_primary = $7, 
				reference_number = $8, retail_price = $9, sale_price = $10, 
				fulfillment_group_sequnce = $11, service = $12, shipping_override = $13, 
				status = $14, total = $15, total_fee_tax = $16, 
				total_fg_tax = $17, total_item_tax = $18, total_tax = $19, 
				address_id = $20, fulfillment_option_id = $21, personal_message_id = $22, 
				phone_id = $23, delivery_instruction = $24, updated_at = $25
			WHERE fulfillment_group_id = $26`
		_, err = tx.ExecContext(ctx, query,
			fg.OrderID, fgType, fg.ShippingPrice, shippingPriceTaxable, fg.MerchandiseTotal,
			method, isPrimary, referenceNumber, fg.RetailPrice, fg.SalePrice,
			fg.Sequence, service, shippingOverride, status, fg.Total,
			fg.TotalFeeTax, fg.TotalFgTax, fg.TotalItemTax, fg.TotalTax,
			addressID, fulfillmentOptionID, personalMessageID, phoneID,
			deliveryInstruction, fg.UpdatedAt, fg.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update fulfillment group: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves a fulfillment group by its unique identifier.
func (r *FulfillmentGroupRepository) FindByID(ctx context.Context, id int64) (*domain.FulfillmentGroup, error) {
	query := `
		SELECT
			fulfillment_group_id, order_id, type, shipping_price, shipping_price_taxable, 
			merchandise_total, method, is_primary, reference_number, retail_price, 
			sale_price, fulfillment_group_sequnce, service, shipping_override, status, 
			total, total_fee_tax, total_fg_tax, total_item_tax, total_tax, 
			address_id, fulfillment_option_id, personal_message_id, phone_id, 
			delivery_instruction, created_at, updated_at
		FROM blc_fulfillment_group WHERE fulfillment_group_id = $1`

	var fg domain.FulfillmentGroup
	var shippingPriceTaxable sql.NullBool
	var isPrimary sql.NullBool
	var shippingOverride sql.NullBool
	var addressID sql.NullInt64
	var fulfillmentOptionID sql.NullInt64
	var personalMessageID sql.NullInt64
	var phoneID sql.NullInt64
	var deliveryInstruction sql.NullString
	var method sql.NullString
	var referenceNumber sql.NullString
	var service sql.NullString
	var status sql.NullString
	var fgType sql.NullString


	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&fg.ID, &fg.OrderID, &fgType, &fg.ShippingPrice, &shippingPriceTaxable,
		&fg.MerchandiseTotal, &method, &isPrimary, &referenceNumber, &fg.RetailPrice,
		&fg.SalePrice, &fg.Sequence, &service, &shippingOverride, &status,
		&fg.Total, &fg.TotalFeeTax, &fg.TotalFgTax, &fg.TotalItemTax, &fg.TotalTax,
		&addressID, &fulfillmentOptionID, &personalMessageID, &phoneID,
		&deliveryInstruction, &fg.CreatedAt, &fg.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query fulfillment group by ID: %w", err)
	}

	if shippingPriceTaxable.Valid {
		fg.ShippingPriceTaxable = shippingPriceTaxable.Bool
	}
	if isPrimary.Valid {
		fg.IsPrimary = isPrimary.Bool
	}
	if shippingOverride.Valid {
		fg.ShippingOverride = shippingOverride.Bool
	}
	if addressID.Valid {
		fg.AddressID = &addressID.Int64
	}
	if fulfillmentOptionID.Valid {
		fg.FulfillmentOptionID = &fulfillmentOptionID.Int64
	}
	if personalMessageID.Valid {
		fg.PersonalMessageID = &personalMessageID.Int64
	}
	if phoneID.Valid {
		fg.PhoneID = &phoneID.Int64
	}
	if deliveryInstruction.Valid {
		fg.DeliveryInstruction = deliveryInstruction.String
	}
	if method.Valid {
		fg.Method = method.String
	}
	if referenceNumber.Valid {
		fg.ReferenceNumber = referenceNumber.String
	}
	if service.Valid {
		fg.Service = service.String
	}
	if status.Valid {
		fg.Status = status.String
	}
	if fgType.Valid {
		fg.Type = fgType.String
	}

	return &fg, nil
}

// FindByOrderID retrieves all fulfillment groups for a given order ID.
func (r *FulfillmentGroupRepository) FindByOrderID(ctx context.Context, orderID int64) ([]*domain.FulfillmentGroup, error) {
	query := `
		SELECT
			fulfillment_group_id, order_id, type, shipping_price, shipping_price_taxable, 
			merchandise_total, method, is_primary, reference_number, retail_price, 
			sale_price, fulfillment_group_sequnce, service, shipping_override, status, 
			total, total_fee_tax, total_fg_tax, total_item_tax, total_tax, 
			address_id, fulfillment_option_id, personal_message_id, phone_id, 
			delivery_instruction, created_at, updated_at
		FROM blc_fulfillment_group WHERE order_id = $1`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query fulfillment groups by order ID: %w", err)
	}
	defer rows.Close()

	var fulfillmentGroups []*domain.FulfillmentGroup
	for rows.Next() {
		var fg domain.FulfillmentGroup
		var shippingPriceTaxable sql.NullBool
		var isPrimary sql.NullBool
		var shippingOverride sql.NullBool
		var addressID sql.NullInt64
		var fulfillmentOptionID sql.NullInt64
		var personalMessageID sql.NullInt64
		var phoneID sql.NullInt64
		var deliveryInstruction sql.NullString
		var method sql.NullString
		var referenceNumber sql.NullString
		var service sql.NullString
		var status sql.NullString
		var fgType sql.NullString

		err := rows.Scan(
			&fg.ID, &fg.OrderID, &fgType, &fg.ShippingPrice, &shippingPriceTaxable,
			&fg.MerchandiseTotal, &method, &isPrimary, &referenceNumber, &fg.RetailPrice,
			&fg.SalePrice, &fg.Sequence, &service, &shippingOverride, &status,
			&fg.Total, &fg.TotalFeeTax, &fg.TotalFgTax, &fg.TotalItemTax, &fg.TotalTax,
			&addressID, &fulfillmentOptionID, &personalMessageID, &phoneID,
			&deliveryInstruction, &fg.CreatedAt, &fg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fulfillment group row: %w", err)
		}

		if shippingPriceTaxable.Valid {
			fg.ShippingPriceTaxable = shippingPriceTaxable.Bool
		}
		if isPrimary.Valid {
			fg.IsPrimary = isPrimary.Bool
		}
		if shippingOverride.Valid {
			fg.ShippingOverride = shippingOverride.Bool
		}
		if addressID.Valid {
			fg.AddressID = &addressID.Int64
		}
		if fulfillmentOptionID.Valid {
			fg.FulfillmentOptionID = &fulfillmentOptionID.Int64
		}
		if personalMessageID.Valid {
			fg.PersonalMessageID = &personalMessageID.Int64
		}
		if phoneID.Valid {
			fg.PhoneID = &phoneID.Int64
		}
		if deliveryInstruction.Valid {
			fg.DeliveryInstruction = deliveryInstruction.String
		}
		if method.Valid {
			fg.Method = method.String
		}
		if referenceNumber.Valid {
			fg.ReferenceNumber = referenceNumber.String
		}
		if service.Valid {
			fg.Service = service.String
		}
		if status.Valid {
			fg.Status = status.String
		}
		if fgType.Valid {
			fg.Type = fgType.String
		}
		fulfillmentGroups = append(fulfillmentGroups, &fg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for fulfillment groups: %w", err)
	}

	return fulfillmentGroups, nil
}

// Delete removes a fulfillment group by its unique identifier.
func (r *FulfillmentGroupRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_fulfillment_group WHERE fulfillment_group_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete fulfillment group: %w", err)
	}
	return nil
}

// DeleteByOrderID removes all fulfillment groups for a given order ID.
func (r *FulfillmentGroupRepository) DeleteByOrderID(ctx context.Context, orderID int64) error {
	query := `DELETE FROM blc_fulfillment_group WHERE order_id = $1`
	_, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete fulfillment groups by order ID: %w", err)
	}
	return nil
}
