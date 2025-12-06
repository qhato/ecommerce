package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/qhato/ecommerce/internal/inventory/domain"
	"github.com/qhato/ecommerce/pkg/database"
)

// PostgresInventoryReservationRepository implements the InventoryReservationRepository interface
type PostgresInventoryReservationRepository struct {
	db *database.DB
}

// NewPostgresInventoryReservationRepository creates a new PostgresInventoryReservationRepository
func NewPostgresInventoryReservationRepository(db *database.DB) *PostgresInventoryReservationRepository {
	return &PostgresInventoryReservationRepository{db: db}
}

// Save stores a new reservation or updates an existing one
func (r *PostgresInventoryReservationRepository) Save(ctx context.Context, reservation *domain.InventoryReservation) error {
	if reservation.CreatedAt.IsZero() {
		return r.create(ctx, reservation)
	}
	return r.update(ctx, reservation)
}

func (r *PostgresInventoryReservationRepository) create(ctx context.Context, reservation *domain.InventoryReservation) error {
	query := `
		INSERT INTO blc_inventory_reservation (
			id, sku_id, quantity, order_id, order_item_id, status,
			reserved_at, expires_at, released_at, fulfilled_at,
			reservation_ref, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	err := r.db.Exec(ctx, query,
		reservation.ID,
		reservation.SKUID,
		reservation.Quantity,
		reservation.OrderID,
		reservation.OrderItemID,
		reservation.Status,
		reservation.ReservedAt,
		reservation.ExpiresAt,
		reservation.ReleasedAt,
		reservation.FulfilledAt,
		reservation.ReservationRef,
		reservation.CreatedAt,
		reservation.UpdatedAt,
	)

	return err
}

func (r *PostgresInventoryReservationRepository) update(ctx context.Context, reservation *domain.InventoryReservation) error {
	query := `
		UPDATE blc_inventory_reservation SET
			quantity = $1,
			status = $2,
			expires_at = $3,
			released_at = $4,
			fulfilled_at = $5,
			updated_at = $6
		WHERE id = $7`

	err := r.db.Exec(ctx, query,
		reservation.Quantity,
		reservation.Status,
		reservation.ExpiresAt,
		reservation.ReleasedAt,
		reservation.FulfilledAt,
		reservation.UpdatedAt,
		reservation.ID,
	)

	return err
}

// FindByID retrieves a reservation by its unique identifier
func (r *PostgresInventoryReservationRepository) FindByID(ctx context.Context, id string) (*domain.InventoryReservation, error) {
	query := `
		SELECT id, sku_id, quantity, order_id, order_item_id, status,
		       reserved_at, expires_at, released_at, fulfilled_at,
		       reservation_ref, created_at, updated_at
		FROM blc_inventory_reservation
		WHERE id = $1`

	reservation := &domain.InventoryReservation{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.SKUID,
		&reservation.Quantity,
		&reservation.OrderID,
		&reservation.OrderItemID,
		&reservation.Status,
		&reservation.ReservedAt,
		&reservation.ExpiresAt,
		&reservation.ReleasedAt,
		&reservation.FulfilledAt,
		&reservation.ReservationRef,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return reservation, err
}

// FindByOrderID retrieves all reservations for an order
func (r *PostgresInventoryReservationRepository) FindByOrderID(ctx context.Context, orderID string) ([]*domain.InventoryReservation, error) {
	query := `
		SELECT id, sku_id, quantity, order_id, order_item_id, status,
		       reserved_at, expires_at, released_at, fulfilled_at,
		       reservation_ref, created_at, updated_at
		FROM blc_inventory_reservation
		WHERE order_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReservations(rows)
}

// FindExpired retrieves all expired reservations
func (r *PostgresInventoryReservationRepository) FindExpired(ctx context.Context) ([]*domain.InventoryReservation, error) {
	query := `
		SELECT id, sku_id, quantity, order_id, order_item_id, status,
		       reserved_at, expires_at, released_at, fulfilled_at,
		       reservation_ref, created_at, updated_at
		FROM blc_inventory_reservation
		WHERE status IN ($1, $2)
		  AND expires_at IS NOT NULL
		  AND expires_at < $3
		ORDER BY expires_at ASC`

	now := time.Now()
	rows, err := r.db.Query(ctx, query,
		domain.ReservationStatusPending,
		domain.ReservationStatusConfirmed,
		now,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReservations(rows)
}

// Delete removes a reservation by its unique identifier
func (r *PostgresInventoryReservationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM blc_inventory_reservation WHERE id = $1`
	err := r.db.Exec(ctx, query, id)
	return err
}

// Helper method to scan multiple reservations
func (r *PostgresInventoryReservationRepository) scanReservations(rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}) ([]*domain.InventoryReservation, error) {
	var reservations []*domain.InventoryReservation

	for rows.Next() {
		reservation := &domain.InventoryReservation{}
		err := rows.Scan(
			&reservation.ID,
			&reservation.SKUID,
			&reservation.Quantity,
			&reservation.OrderID,
			&reservation.OrderItemID,
			&reservation.Status,
			&reservation.ReservedAt,
			&reservation.ExpiresAt,
			&reservation.ReleasedAt,
			&reservation.FulfilledAt,
			&reservation.ReservationRef,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}
