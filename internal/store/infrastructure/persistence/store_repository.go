package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/qhato/ecommerce/internal/store/domain"
)

type PostgresStoreRepository struct {
	db *sql.DB
}

func NewPostgresStoreRepository(db *sql.DB) *PostgresStoreRepository {
	return &PostgresStoreRepository{db: db}
}

func (r *PostgresStoreRepository) Create(ctx context.Context, store *domain.Store) error {
	addressJSON, _ := json.Marshal(store.Address)
	settingsJSON, _ := json.Marshal(store.Settings)
	metadataJSON, _ := json.Marshal(store.Metadata)

	query := `INSERT INTO blc_store (
		code, name, description, type, status, email, phone, website,
		address, timezone, currency, locale, tax_id, settings, metadata,
		parent_store_id, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		store.Code, store.Name, store.Description, store.Type, store.Status,
		store.Email, store.Phone, store.Website, addressJSON, store.Timezone,
		store.Currency, store.Locale, store.TaxID, settingsJSON, metadataJSON,
		store.ParentStoreID, store.CreatedAt, store.UpdatedAt,
	).Scan(&store.ID)
}

func (r *PostgresStoreRepository) Update(ctx context.Context, store *domain.Store) error {
	addressJSON, _ := json.Marshal(store.Address)
	settingsJSON, _ := json.Marshal(store.Settings)
	metadataJSON, _ := json.Marshal(store.Metadata)

	query := `UPDATE blc_store SET
		name = $1, description = $2, status = $3, email = $4, phone = $5,
		website = $6, address = $7, timezone = $8, currency = $9, locale = $10,
		tax_id = $11, settings = $12, metadata = $13, updated_at = $14
	WHERE id = $15`

	_, err := r.db.ExecContext(ctx, query,
		store.Name, store.Description, store.Status, store.Email, store.Phone,
		store.Website, addressJSON, store.Timezone, store.Currency, store.Locale,
		store.TaxID, settingsJSON, metadataJSON, store.UpdatedAt, store.ID,
	)
	return err
}

func (r *PostgresStoreRepository) FindByID(ctx context.Context, id int64) (*domain.Store, error) {
	query := `SELECT id, code, name, description, type, status, email, phone, website,
		address, timezone, currency, locale, tax_id, settings, metadata, parent_store_id,
		created_at, updated_at
	FROM blc_store WHERE id = $1`

	return r.scanStore(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresStoreRepository) FindByCode(ctx context.Context, code string) (*domain.Store, error) {
	query := `SELECT id, code, name, description, type, status, email, phone, website,
		address, timezone, currency, locale, tax_id, settings, metadata, parent_store_id,
		created_at, updated_at
	FROM blc_store WHERE code = $1`

	return r.scanStore(r.db.QueryRowContext(ctx, query, code))
}

func (r *PostgresStoreRepository) FindByStatus(ctx context.Context, status domain.StoreStatus) ([]*domain.Store, error) {
	query := `SELECT id, code, name, description, type, status, email, phone, website,
		address, timezone, currency, locale, tax_id, settings, metadata, parent_store_id,
		created_at, updated_at
	FROM blc_store WHERE status = $1 ORDER BY name`

	return r.queryStores(ctx, query, status)
}

func (r *PostgresStoreRepository) FindByType(ctx context.Context, storeType domain.StoreType) ([]*domain.Store, error) {
	query := `SELECT id, code, name, description, type, status, email, phone, website,
		address, timezone, currency, locale, tax_id, settings, metadata, parent_store_id,
		created_at, updated_at
	FROM blc_store WHERE type = $1 ORDER BY name`

	return r.queryStores(ctx, query, storeType)
}

func (r *PostgresStoreRepository) FindAll(ctx context.Context, limit int) ([]*domain.Store, error) {
	query := `SELECT id, code, name, description, type, status, email, phone, website,
		address, timezone, currency, locale, tax_id, settings, metadata, parent_store_id,
		created_at, updated_at
	FROM blc_store ORDER BY name LIMIT $1`

	return r.queryStores(ctx, query, limit)
}

func (r *PostgresStoreRepository) FindNearby(ctx context.Context, lat, lng float64, radiusKm float64, limit int) ([]*domain.Store, error) {
	query := `SELECT s.id, s.code, s.name, s.description, s.type, s.status, s.email, s.phone, s.website,
		s.address, s.timezone, s.currency, s.locale, s.tax_id, s.settings, s.metadata, s.parent_store_id,
		s.created_at, s.updated_at,
		( 6371 * acos( cos( radians($1) ) * cos( radians( (address->>'latitude')::float ) ) *
		  cos( radians( (address->>'longitude')::float ) - radians($2) ) + sin( radians($1) ) *
		  sin( radians( (address->>'latitude')::float ) ) ) ) AS distance
	FROM blc_store s
	WHERE status = 'ACTIVE'
	  AND address->>'latitude' IS NOT NULL
	  AND address->>'longitude' IS NOT NULL
	  AND ( 6371 * acos( cos( radians($1) ) * cos( radians( (address->>'latitude')::float ) ) *
		  cos( radians( (address->>'longitude')::float ) - radians($2) ) + sin( radians($1) ) *
		  sin( radians( (address->>'latitude')::float ) ) ) ) <= $3
	ORDER BY distance
	LIMIT $4`

	return r.queryStoresWithDistance(ctx, query, lat, lng, radiusKm, limit)
}

func (r *PostgresStoreRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_store WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresStoreRepository) scanStore(row interface {
	Scan(dest ...interface{}) error
}) (*domain.Store, error) {
	store := &domain.Store{}
	var addressJSON, settingsJSON, metadataJSON []byte

	err := row.Scan(
		&store.ID, &store.Code, &store.Name, &store.Description, &store.Type,
		&store.Status, &store.Email, &store.Phone, &store.Website, &addressJSON,
		&store.Timezone, &store.Currency, &store.Locale, &store.TaxID,
		&settingsJSON, &metadataJSON, &store.ParentStoreID,
		&store.CreatedAt, &store.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(addressJSON, &store.Address); err != nil {
		return nil, fmt.Errorf("failed to unmarshal address: %w", err)
	}
	if err := json.Unmarshal(settingsJSON, &store.Settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}
	if err := json.Unmarshal(metadataJSON, &store.Metadata); err != nil {
		store.Metadata = make(map[string]interface{})
	}

	return store, nil
}

func (r *PostgresStoreRepository) queryStores(ctx context.Context, query string, args ...interface{}) ([]*domain.Store, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := make([]*domain.Store, 0)
	for rows.Next() {
		store := &domain.Store{}
		var addressJSON, settingsJSON, metadataJSON []byte

		if err := rows.Scan(
			&store.ID, &store.Code, &store.Name, &store.Description, &store.Type,
			&store.Status, &store.Email, &store.Phone, &store.Website, &addressJSON,
			&store.Timezone, &store.Currency, &store.Locale, &store.TaxID,
			&settingsJSON, &metadataJSON, &store.ParentStoreID,
			&store.CreatedAt, &store.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(addressJSON, &store.Address); err != nil {
			return nil, fmt.Errorf("failed to unmarshal address: %w", err)
		}
		if err := json.Unmarshal(settingsJSON, &store.Settings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
		}
		if err := json.Unmarshal(metadataJSON, &store.Metadata); err != nil {
			store.Metadata = make(map[string]interface{})
		}

		stores = append(stores, store)
	}

	return stores, nil
}

func (r *PostgresStoreRepository) queryStoresWithDistance(ctx context.Context, query string, args ...interface{}) ([]*domain.Store, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := make([]*domain.Store, 0)
	for rows.Next() {
		store := &domain.Store{}
		var addressJSON, settingsJSON, metadataJSON []byte
		var distance float64

		if err := rows.Scan(
			&store.ID, &store.Code, &store.Name, &store.Description, &store.Type,
			&store.Status, &store.Email, &store.Phone, &store.Website, &addressJSON,
			&store.Timezone, &store.Currency, &store.Locale, &store.TaxID,
			&settingsJSON, &metadataJSON, &store.ParentStoreID,
			&store.CreatedAt, &store.UpdatedAt, &distance,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(addressJSON, &store.Address); err != nil {
			return nil, fmt.Errorf("failed to unmarshal address: %w", err)
		}
		if err := json.Unmarshal(settingsJSON, &store.Settings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
		}
		if err := json.Unmarshal(metadataJSON, &store.Metadata); err != nil {
			store.Metadata = make(map[string]interface{})
		}

		stores = append(stores, store)
	}

	return stores, nil
}
