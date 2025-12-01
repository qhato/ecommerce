package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/qhato/ecommerce/internal/fulfillment/domain"
	"github.com/qhato/ecommerce/pkg/database"
	"github.com/qhato/ecommerce/pkg/errors"
)

// PostgresShipmentRepository implements the ShipmentRepository interface using PostgreSQL
type PostgresShipmentRepository struct {
	db *database.DB
}

// NewPostgresShipmentRepository creates a new PostgresShipmentRepository
func NewPostgresShipmentRepository(db *database.DB) *PostgresShipmentRepository {
	return &PostgresShipmentRepository{db: db}
}

// Create creates a new shipment
func (r *PostgresShipmentRepository) Create(ctx context.Context, shipment *domain.Shipment) error {
	query := `
		INSERT INTO blc_fulfillment_group (
			order_id, status, tracking_number, carrier, shipping_method,
			shipping_cost, estimated_delivery_date, shipped_date, delivered_date,
			address_name, address_line1, address_line2, city, state,
			postal_code, country, phone, notes, date_created, date_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING fulfillment_group_id
	`

	err := r.db.QueryRow(ctx, query,
		shipment.OrderID,
		shipment.Status,
		shipment.TrackingNumber,
		shipment.Carrier,
		shipment.ShippingMethod,
		shipment.ShippingCost,
		shipment.EstimatedDate,
		shipment.ShippedDate,
		shipment.DeliveredDate,
		shipment.ShippingAddress.Name,
		shipment.ShippingAddress.Line1,
		shipment.ShippingAddress.Line2,
		shipment.ShippingAddress.City,
		shipment.ShippingAddress.State,
		shipment.ShippingAddress.PostalCode,
		shipment.ShippingAddress.Country,
		shipment.ShippingAddress.Phone,
		shipment.Notes,
		shipment.CreatedAt,
		shipment.UpdatedAt,
	).Scan(&shipment.ID)

	if err != nil {
		return errors.InternalWrap(err, "failed to create shipment")
	}

	return nil
}

// Update updates an existing shipment
func (r *PostgresShipmentRepository) Update(ctx context.Context, shipment *domain.Shipment) error {
	query := `
		UPDATE blc_fulfillment_group
		SET order_id = $1, status = $2, tracking_number = $3, carrier = $4,
			shipping_method = $5, shipping_cost = $6, estimated_delivery_date = $7,
			shipped_date = $8, delivered_date = $9, address_name = $10,
			address_line1 = $11, address_line2 = $12, city = $13, state = $14,
			postal_code = $15, country = $16, phone = $17, notes = $18,
			date_updated = $19
		WHERE fulfillment_group_id = $20
	`

	// Using Pool().Exec to get RowsAffected
	tag, err := r.db.Pool().Exec(ctx, query,
		shipment.OrderID,
		shipment.Status,
		shipment.TrackingNumber,
		shipment.Carrier,
		shipment.ShippingMethod,
		shipment.ShippingCost,
		shipment.EstimatedDate,
		shipment.ShippedDate,
		shipment.DeliveredDate,
		shipment.ShippingAddress.Name,
		shipment.ShippingAddress.Line1,
		shipment.ShippingAddress.Line2,
		shipment.ShippingAddress.City,
		shipment.ShippingAddress.State,
		shipment.ShippingAddress.PostalCode,
		shipment.ShippingAddress.Country,
		shipment.ShippingAddress.Phone,
		shipment.Notes,
		shipment.UpdatedAt,
		shipment.ID,
	)

	if err != nil {
		return errors.InternalWrap(err, "failed to update shipment")
	}

	if tag.RowsAffected() == 0 {
		return errors.NotFound(fmt.Sprintf("shipment %d", shipment.ID))
	}

	return nil
}

// FindByID finds a shipment by ID
func (r *PostgresShipmentRepository) FindByID(ctx context.Context, id int64) (*domain.Shipment, error) {
	query := `
		SELECT fulfillment_group_id, order_id, status, tracking_number, carrier,
			   shipping_method, shipping_cost, estimated_delivery_date, shipped_date,
			   delivered_date, address_name, address_line1, address_line2, city,
			   state, postal_code, country, phone, notes, date_created, date_updated
		FROM blc_fulfillment_group
		WHERE fulfillment_group_id = $1
	`

	shipment := &domain.Shipment{}
	var (
		trackingNumber sql.NullString
		estimatedDate  sql.NullTime
		shippedDate    sql.NullTime
		deliveredDate  sql.NullTime
		addressLine2   sql.NullString
		phone          sql.NullString
		notes          sql.NullString
	)

	err := r.db.QueryRow(ctx, query, id).Scan(
		&shipment.ID,
		&shipment.OrderID,
		&shipment.Status,
		&trackingNumber,
		&shipment.Carrier,
		&shipment.ShippingMethod,
		&shipment.ShippingCost,
		&estimatedDate,
		&shippedDate,
		&deliveredDate,
		&shipment.ShippingAddress.Name,
		&shipment.ShippingAddress.Line1,
		&addressLine2,
		&shipment.ShippingAddress.City,
		&shipment.ShippingAddress.State,
		&shipment.ShippingAddress.PostalCode,
		&shipment.ShippingAddress.Country,
		&phone,
		&notes,
		&shipment.CreatedAt,
		&shipment.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find shipment by ID")
	}

	// Handle nullable fields
	if trackingNumber.Valid {
		shipment.TrackingNumber = trackingNumber.String
	}
	if estimatedDate.Valid {
		shipment.EstimatedDate = &estimatedDate.Time
	}
	if shippedDate.Valid {
		shipment.ShippedDate = &shippedDate.Time
	}
	if deliveredDate.Valid {
		shipment.DeliveredDate = &deliveredDate.Time
	}
	if addressLine2.Valid {
		shipment.ShippingAddress.Line2 = addressLine2.String
	}
	if phone.Valid {
		shipment.ShippingAddress.Phone = phone.String
	}
	if notes.Valid {
		shipment.Notes = notes.String
	}

	return shipment, nil
}

// FindByOrderID finds shipments by order ID
func (r *PostgresShipmentRepository) FindByOrderID(ctx context.Context, orderID int64) ([]*domain.Shipment, error) {
	query := `
		SELECT fulfillment_group_id, order_id, status, tracking_number, carrier,
			   shipping_method, shipping_cost, estimated_delivery_date, shipped_date,
			   delivered_date, address_name, address_line1, address_line2, city,
			   state, postal_code, country, phone, notes, date_created, date_updated
		FROM blc_fulfillment_group
		WHERE order_id = $1
		ORDER BY date_created DESC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find shipments by order")
	}
	defer rows.Close()

	return r.scanShipments(rows)
}

// FindByTrackingNumber finds a shipment by tracking number
func (r *PostgresShipmentRepository) FindByTrackingNumber(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
	query := `
		SELECT fulfillment_group_id, order_id, status, tracking_number, carrier,
			   shipping_method, shipping_cost, estimated_delivery_date, shipped_date,
			   delivered_date, address_name, address_line1, address_line2, city,
			   state, postal_code, country, phone, notes, date_created, date_updated
		FROM blc_fulfillment_group
		WHERE tracking_number = $1
	`

	shipment := &domain.Shipment{}
	var (
		trackNum      sql.NullString
		estimatedDate sql.NullTime
		shippedDate   sql.NullTime
		deliveredDate sql.NullTime
		addressLine2  sql.NullString
		phone         sql.NullString
		notes         sql.NullString
	)

	err := r.db.QueryRow(ctx, query, trackingNumber).Scan(
		&shipment.ID,
		&shipment.OrderID,
		&shipment.Status,
		&trackNum,
		&shipment.Carrier,
		&shipment.ShippingMethod,
		&shipment.ShippingCost,
		&estimatedDate,
		&shippedDate,
		&deliveredDate,
		&shipment.ShippingAddress.Name,
		&shipment.ShippingAddress.Line1,
		&addressLine2,
		&shipment.ShippingAddress.City,
		&shipment.ShippingAddress.State,
		&shipment.ShippingAddress.PostalCode,
		&shipment.ShippingAddress.Country,
		&phone,
		&notes,
		&shipment.CreatedAt,
		&shipment.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find shipment by tracking number")
	}

	// Handle nullable fields
	if trackNum.Valid {
		shipment.TrackingNumber = trackNum.String
	}
	if estimatedDate.Valid {
		shipment.EstimatedDate = &estimatedDate.Time
	}
	if shippedDate.Valid {
		shipment.ShippedDate = &shippedDate.Time
	}
	if deliveredDate.Valid {
		shipment.DeliveredDate = &deliveredDate.Time
	}
	if addressLine2.Valid {
		shipment.ShippingAddress.Line2 = addressLine2.String
	}
	if phone.Valid {
		shipment.ShippingAddress.Phone = phone.String
	}
	if notes.Valid {
		shipment.Notes = notes.String
	}

	return shipment, nil
}

// FindAll finds all shipments
func (r *PostgresShipmentRepository) FindAll(ctx context.Context, filter *domain.ShipmentFilter) ([]*domain.Shipment, int64, error) {
	query := `
		SELECT fulfillment_group_id, order_id, status, tracking_number, carrier,
			   shipping_method, shipping_cost, estimated_delivery_date, shipped_date,
			   delivered_date, address_name, address_line1, address_line2, city,
			   state, postal_code, country, phone, notes, date_created, date_updated
		FROM blc_fulfillment_group
		WHERE 1=1
	`

	args := make([]interface{}, 0)
	argIndex := 1

	// Add filters
	if filter != nil {
		if filter.Status != "" {
			query += fmt.Sprintf(" AND status = $%d", argIndex)
			args = append(args, filter.Status)
			argIndex++
		}
		if filter.Carrier != "" {
			query += fmt.Sprintf(" AND carrier = $%d", argIndex)
			args = append(args, filter.Carrier)
			argIndex++
		}
		if filter.OrderID > 0 {
			query += fmt.Sprintf(" AND order_id = $%d", argIndex)
			args = append(args, filter.OrderID)
			argIndex++
		}
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM blc_fulfillment_group WHERE 1=1"
	countArgs := make([]interface{}, 0)
	countArgIndex := 1
	if filter != nil {
		if filter.Status != "" {
			countQuery += fmt.Sprintf(" AND status = $%d", countArgIndex)
			countArgs = append(countArgs, filter.Status)
			countArgIndex++
		}
		if filter.Carrier != "" {
			countQuery += fmt.Sprintf(" AND carrier = $%d", countArgIndex)
			countArgs = append(countArgs, filter.Carrier)
			countArgIndex++
		}
		if filter.OrderID > 0 {
			countQuery += fmt.Sprintf(" AND order_id = $%d", countArgIndex)
			countArgs = append(countArgs, filter.OrderID)
		}
	}

	var total int64
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to count shipments")
	}

	// Add sorting
	if filter != nil && filter.SortBy != "" {
		sortOrder := "ASC"
		if filter.SortOrder == "DESC" {
			sortOrder = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, sortOrder)
	} else {
		query += " ORDER BY date_created DESC"
	}

	// Add pagination
	if filter != nil && filter.PageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.InternalWrap(err, "failed to find all shipments")
	}
	defer rows.Close()

	shipments, err := r.scanShipments(rows)
	return shipments, total, err
}

// scanShipments scans shipment rows
func (r *PostgresShipmentRepository) scanShipments(rows pgx.Rows) ([]*domain.Shipment, error) {
	shipments := make([]*domain.Shipment, 0)

	for rows.Next() {
		shipment := &domain.Shipment{}
		var (
			trackingNumber sql.NullString
			estimatedDate  sql.NullTime
			shippedDate    sql.NullTime
			deliveredDate  sql.NullTime
			addressLine2   sql.NullString
			phone          sql.NullString
			notes          sql.NullString
		)

		err := rows.Scan(
			&shipment.ID,
			&shipment.OrderID,
			&shipment.Status,
			&trackingNumber,
			&shipment.Carrier,
			&shipment.ShippingMethod,
			&shipment.ShippingCost,
			&estimatedDate,
			&shippedDate,
			&deliveredDate,
			&shipment.ShippingAddress.Name,
			&shipment.ShippingAddress.Line1,
			&addressLine2,
			&shipment.ShippingAddress.City,
			&shipment.ShippingAddress.State,
			&shipment.ShippingAddress.PostalCode,
			&shipment.ShippingAddress.Country,
			&phone,
			&notes,
			&shipment.CreatedAt,
			&shipment.UpdatedAt,
		)
		if err != nil {
			return nil, errors.InternalWrap(err, "failed to scan shipment")
		}

		// Handle nullable fields
		if trackingNumber.Valid {
			shipment.TrackingNumber = trackingNumber.String
		}
		if estimatedDate.Valid {
			shipment.EstimatedDate = &estimatedDate.Time
		}
		if shippedDate.Valid {
			shipment.ShippedDate = &shippedDate.Time
		}
		if deliveredDate.Valid {
			shipment.DeliveredDate = &deliveredDate.Time
		}
		if addressLine2.Valid {
			shipment.ShippingAddress.Line2 = addressLine2.String
		}
		if phone.Valid {
			shipment.ShippingAddress.Phone = phone.String
		}
		if notes.Valid {
			shipment.Notes = notes.String
		}

		shipments = append(shipments, shipment)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.InternalWrap(err, "failed to iterate shipments")
	}

	return shipments, nil
}