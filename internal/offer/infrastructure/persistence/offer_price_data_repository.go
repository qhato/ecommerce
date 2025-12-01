package persistence

import (
	"context"

	"github.com/qhato/ecommerce/internal/offer/domain"
	"github.com/qhato/ecommerce/pkg/database"
	"github.com/qhato/ecommerce/pkg/errors"
)

type PostgresOfferPriceDataRepository struct {
	db *database.DB
}

func NewPostgresOfferPriceDataRepository(db *database.DB) *PostgresOfferPriceDataRepository {
	return &PostgresOfferPriceDataRepository{db: db}
}

func (r *PostgresOfferPriceDataRepository) Save(ctx context.Context, priceData *domain.OfferPriceData) error {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresOfferPriceDataRepository) FindByID(ctx context.Context, id int64) (*domain.OfferPriceData, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresOfferPriceDataRepository) FindByOfferID(ctx context.Context, offerID int64) ([]*domain.OfferPriceData, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresOfferPriceDataRepository) FindActiveByOfferID(ctx context.Context, offerID int64) ([]*domain.OfferPriceData, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresOfferPriceDataRepository) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresOfferPriceDataRepository) DeleteByOfferID(ctx context.Context, offerID int64) error {
	query := "DELETE FROM blc_offer_price_data WHERE offer_id = $1"
	err := r.db.Exec(ctx, query, offerID)
	if err != nil {
		return errors.InternalWrap(err, "failed to delete offer price data by offer ID")
	}
	return nil
}
