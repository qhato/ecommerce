package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/qhato/ecommerce/internal/offer/domain"
	"github.com/qhato/ecommerce/pkg/database"
	"github.com/qhato/ecommerce/pkg/errors"
)

// PostgresOfferRepository implements the OfferRepository interface
type PostgresOfferRepository struct {
	db *database.DB
}

// NewPostgresOfferRepository creates a new PostgresOfferRepository
func NewPostgresOfferRepository(db *database.DB) *PostgresOfferRepository {
	return &PostgresOfferRepository{db: db}
}

// Save stores a new offer or updates an existing one.
func (r *PostgresOfferRepository) Save(ctx context.Context, offer *domain.Offer) error {
	if offer.ID == 0 {
		return r.create(ctx, offer)
	}
	return r.update(ctx, offer)
}

func (r *PostgresOfferRepository) create(ctx context.Context, offer *domain.Offer) error {
	query := `
		INSERT INTO blc_offer (
			offer_id, offer_name, offer_type, offer_value, adjustment_type,
			apply_to_child_items, apply_to_sale_price, archived, automatically_added,
			combinable_with_other_offers, offer_description, offer_discount_type,
			end_date, marketing_message, max_uses_per_customer, max_uses,
			max_uses_strategy, minimum_days_per_usage, offer_item_qualifier_rule,
			offer_item_target_rule, order_min_total, offer_priority,
			qualifying_item_min_total, requires_related_tar_qual, start_date,
			target_min_total, target_system, totalitarian_offer, use_list_for_discounts,
			date_created, date_updated
		) VALUES (
			nextval('blc_offer_seq'), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25,
			$26, $27, $28, $29, $30
		) RETURNING offer_id`

	archivedFlag := "N"
	if offer.Archived {
		archivedFlag = "Y"
	}

	err := r.db.QueryRow(ctx, query,
		offer.Name, offer.OfferType, offer.OfferValue, offer.AdjustmentType,
		offer.ApplyToChildItems, offer.ApplyToSalePrice, archivedFlag, offer.AutomaticallyAdded,
		offer.CombinableWithOtherOffers, offer.OfferDescription, offer.OfferDiscountType,
		offer.EndDate, offer.MarketingMessage, offer.MaxUsesPerCustomer, offer.MaxUses,
		offer.MaxUsesStrategy, offer.MinimumDaysPerUsage, offer.OfferItemQualifierRule,
		offer.OfferItemTargetRule, offer.OrderMinTotal, offer.OfferPriority,
		offer.QualifyingItemMinTotal, offer.RequiresRelatedTarQual, offer.StartDate,
		offer.TargetMinTotal, offer.TargetSystem, offer.TotalitarianOffer, offer.UseListForDiscounts,
		offer.CreatedAt, offer.UpdatedAt,
	).Scan(&offer.ID)

	if err != nil {
		return errors.InternalWrap(err, "failed to create offer")
	}
	return nil
}

func (r *PostgresOfferRepository) update(ctx context.Context, offer *domain.Offer) error {
	query := `
		UPDATE blc_offer SET
			offer_name = $1, offer_type = $2, offer_value = $3, adjustment_type = $4,
			apply_to_child_items = $5, apply_to_sale_price = $6, archived = $7, automatically_added = $8,
			combinable_with_other_offers = $9, offer_description = $10, offer_discount_type = $11,
			end_date = $12, marketing_message = $13, max_uses_per_customer = $14, max_uses = $15,
			max_uses_strategy = $16, minimum_days_per_usage = $17, offer_item_qualifier_rule = $18,
			offer_item_target_rule = $19, order_min_total = $20, offer_priority = $21,
			qualifying_item_min_total = $22, requires_related_tar_qual = $23, start_date = $24,
			target_min_total = $25, target_system = $26, totalitarian_offer = $27, use_list_for_discounts = $28,
			date_updated = $29
		WHERE offer_id = $30`

	archivedFlag := "N"
	if offer.Archived {
		archivedFlag = "Y"
	}

	tag, err := r.db.Pool().Exec(ctx, query,
		offer.Name, offer.OfferType, offer.OfferValue, offer.AdjustmentType,
		offer.ApplyToChildItems, offer.ApplyToSalePrice, archivedFlag, offer.AutomaticallyAdded,
		offer.CombinableWithOtherOffers, offer.OfferDescription, offer.OfferDiscountType,
		offer.EndDate, offer.MarketingMessage, offer.MaxUsesPerCustomer, offer.MaxUses,
		offer.MaxUsesStrategy, offer.MinimumDaysPerUsage, offer.OfferItemQualifierRule,
		offer.OfferItemTargetRule, offer.OrderMinTotal, offer.OfferPriority,
		offer.QualifyingItemMinTotal, offer.RequiresRelatedTarQual, offer.StartDate,
		offer.TargetMinTotal, offer.TargetSystem, offer.TotalitarianOffer, offer.UseListForDiscounts,
		offer.UpdatedAt, offer.ID,
	)

	if err != nil {
		return errors.InternalWrap(err, "failed to update offer")
	}

	if tag.RowsAffected() == 0 {
		return errors.NotFound("offer not found")
	}
	return nil
}

// FindByID retrieves an offer by its unique identifier.
func (r *PostgresOfferRepository) FindByID(ctx context.Context, id int64) (*domain.Offer, error) {
	query := `
		SELECT
			offer_id, offer_name, offer_type, offer_value, adjustment_type,
			apply_to_child_items, apply_to_sale_price, archived, automatically_added,
			combinable_with_other_offers, offer_description, offer_discount_type,
			end_date, marketing_message, max_uses_per_customer, max_uses,
			max_uses_strategy, minimum_days_per_usage, offer_item_qualifier_rule,
			offer_item_target_rule, order_min_total, offer_priority,
			qualifying_item_min_total, requires_related_tar_qual, start_date,
			target_min_total, target_system, totalitarian_offer, use_list_for_discounts,
			date_created, date_updated
		FROM blc_offer
		WHERE offer_id = $1`

	offer := &domain.Offer{}
	var (
		archivedFlag                    string
		endDate, startDate              sql.NullTime
		maxUsesPerCustomer              sql.NullInt64
		maxUses                         sql.NullInt32
		applyToChildItems               sql.NullBool
		applyToSalePrice                sql.NullBool
		automaticallyAdded              sql.NullBool
		combinableWithOtherOffers       sql.NullBool
		requiresRelatedTarQual          sql.NullBool
		totalitarianOffer               sql.NullBool
		useListForDiscounts             sql.NullBool
		minimumDaysPerUsage             sql.NullInt64
		offerItemQualifierRule          sql.NullString
		offerItemTargetRule             sql.NullString
		marketingMessage                sql.NullString
		maxUsesStrategy                 sql.NullString
		offerDescription                sql.NullString
		targetSystem                    sql.NullString
		offerDiscountType               sql.NullString
		offerType                       sql.NullString
		adjustmentType                  sql.NullString
	)

	err := r.db.QueryRow(ctx, query, id).Scan(
		&offer.ID,
		&offer.Name,
		&offerType,
		&offer.OfferValue,
		&adjustmentType,
		&applyToChildItems,
		&applyToSalePrice,
		&archivedFlag,
		&automaticallyAdded,
		&combinableWithOtherOffers,
		&offerDescription,
		&offerDiscountType,
		&endDate,
		&marketingMessage,
		&maxUsesPerCustomer,
		&maxUses,
		&maxUsesStrategy,
		&minimumDaysPerUsage,
		&offerItemQualifierRule,
		&offerItemTargetRule,
		&offer.OrderMinTotal,
		&offer.OfferPriority,
		&offer.QualifyingItemMinTotal,
		&requiresRelatedTarQual,
		&startDate,
		&offer.TargetMinTotal,
		&targetSystem,
		&totalitarianOffer,
		&useListForDiscounts,
		&offer.CreatedAt,
		&offer.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find offer")
	}

	offer.Archived = (archivedFlag == "Y")
	offer.OfferType = domain.OfferType(offerType.String)
	offer.AdjustmentType = domain.OfferAdjustmentType(adjustmentType.String)
	offer.OfferDiscountType = domain.OfferDiscountType(offerDiscountType.String)

	if applyToChildItems.Valid {
		offer.ApplyToChildItems = applyToChildItems.Bool
	}
	if applyToSalePrice.Valid {
		offer.ApplyToSalePrice = applyToSalePrice.Bool
	}
	if automaticallyAdded.Valid {
		offer.AutomaticallyAdded = automaticallyAdded.Bool
	}
	if combinableWithOtherOffers.Valid {
		offer.CombinableWithOtherOffers = combinableWithOtherOffers.Bool
	}
	if requiresRelatedTarQual.Valid {
		offer.RequiresRelatedTarQual = requiresRelatedTarQual.Bool
	}
	if totalitarianOffer.Valid {
		offer.TotalitarianOffer = totalitarianOffer.Bool
	}
	if useListForDiscounts.Valid {
		offer.UseListForDiscounts = useListForDiscounts.Bool
	}
	if endDate.Valid {
		offer.EndDate = &endDate.Time
	}
	if startDate.Valid {
		offer.StartDate = startDate.Time
	}
	if maxUsesPerCustomer.Valid {
		offer.MaxUsesPerCustomer = &maxUsesPerCustomer.Int64
	}
	if maxUses.Valid {
		maxUsesInt := int(maxUses.Int32)
		offer.MaxUses = &maxUsesInt
	}
	if minimumDaysPerUsage.Valid {
		offer.MinimumDaysPerUsage = &minimumDaysPerUsage.Int64
	}
	if offerItemQualifierRule.Valid {
		offer.OfferItemQualifierRule = offerItemQualifierRule.String
	}
	if offerItemTargetRule.Valid {
		offer.OfferItemTargetRule = offerItemTargetRule.String
	}
	if marketingMessage.Valid {
		offer.MarketingMessage = marketingMessage.String
	}
	if maxUsesStrategy.Valid {
		offer.MaxUsesStrategy = maxUsesStrategy.String
	}
	if offerDescription.Valid {
		offer.OfferDescription = offerDescription.String
	}
	if targetSystem.Valid {
		offer.TargetSystem = targetSystem.String
	}

	return offer, nil
}

// FindAll retrieves all offers based on a filter.
func (r *PostgresOfferRepository) FindAll(ctx context.Context, filter *domain.OfferFilter) ([]*domain.Offer, error) {
	var offers []*domain.Offer
	query := `
		SELECT
			offer_id, offer_name, offer_type, offer_value, adjustment_type,
			apply_to_child_items, apply_to_sale_price, archived, automatically_added,
			combinable_with_other_offers, offer_description, offer_discount_type,
			end_date, marketing_message, max_uses_per_customer, max_uses,
			max_uses_strategy, minimum_days_per_usage, offer_item_qualifier_rule,
			offer_item_target_rule, order_min_total, offer_priority,
			qualifying_item_min_total, requires_related_tar_qual, start_date,
			target_min_total, target_system, totalitarian_offer, use_list_for_discounts,
			date_created, date_updated
		FROM blc_offer
		WHERE 1=1`

	args := []interface{}{}
	argCounter := 1

	if filter != nil {
		if filter.ActiveOnly {
			query += fmt.Sprintf(" AND archived = 'N' AND start_date <= NOW() AND (end_date IS NULL OR end_date >= NOW())")
		}
		if !filter.IncludeArchived {
			query += fmt.Sprintf(" AND archived = 'N'")
		}
		if filter.OfferType != nil {
			query += fmt.Sprintf(" AND offer_type = $%d", argCounter)
			args = append(args, string(*filter.OfferType))
			argCounter++
		}
		// Add other filters as needed
	}

	// Add sorting
	if filter != nil && filter.SortBy != "" {
		sortOrder := "ASC"
		if filter.SortOrder == "DESC" {
			sortOrder = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, sortOrder)
	} else {
		query += " ORDER BY offer_priority ASC, date_created DESC"
	}

	// Add pagination
	if filter != nil && filter.PageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
		args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.InternalWrap(err, "failed to find offers")
	}
	defer rows.Close()

	for rows.Next() {
		offer := &domain.Offer{}
		var (
			archivedFlag                    string
			endDate, startDate              sql.NullTime
			maxUsesPerCustomer              sql.NullInt64
			maxUses                         sql.NullInt32
			applyToChildItems               sql.NullBool
			applyToSalePrice                sql.NullBool
			automaticallyAdded              sql.NullBool
			combinableWithOtherOffers       sql.NullBool
			requiresRelatedTarQual          sql.NullBool
			totalitarianOffer               sql.NullBool
			useListForDiscounts             sql.NullBool
			minimumDaysPerUsage             sql.NullInt64
			offerItemQualifierRule          sql.NullString
			offerItemTargetRule             sql.NullString
			marketingMessage                sql.NullString
			maxUsesStrategy                 sql.NullString
			offerDescription                sql.NullString
			targetSystem                    sql.NullString
			offerDiscountType               sql.NullString
			offerType                       sql.NullString
			adjustmentType                  sql.NullString
		)

		err := rows.Scan(
			&offer.ID,
			&offer.Name,
			&offerType,
			&offer.OfferValue,
			&adjustmentType,
			&applyToChildItems,
			&applyToSalePrice,
			&archivedFlag,
			&automaticallyAdded,
			&combinableWithOtherOffers,
			&offerDescription,
			&offerDiscountType,
			&endDate,
			&marketingMessage,
			&maxUsesPerCustomer,
			&maxUses,
			&maxUsesStrategy,
			&minimumDaysPerUsage,
			&offerItemQualifierRule,
			&offerItemTargetRule,
			&offer.OrderMinTotal,
			&offer.OfferPriority,
			&offer.QualifyingItemMinTotal,
			&requiresRelatedTarQual,
			&startDate,
			&offer.TargetMinTotal,
			&targetSystem,
			&totalitarianOffer,
			&useListForDiscounts,
			&offer.CreatedAt,
			&offer.UpdatedAt,
		)
		if err != nil {
			return nil, errors.InternalWrap(err, "failed to scan offer")
		}

		offer.Archived = (archivedFlag == "Y")
		offer.OfferType = domain.OfferType(offerType.String)
		offer.AdjustmentType = domain.OfferAdjustmentType(adjustmentType.String)
		offer.OfferDiscountType = domain.OfferDiscountType(offerDiscountType.String)

		if applyToChildItems.Valid {
			offer.ApplyToChildItems = applyToChildItems.Bool
		}
		if applyToSalePrice.Valid {
			offer.ApplyToSalePrice = applyToSalePrice.Bool
		}
		if automaticallyAdded.Valid {
			offer.AutomaticallyAdded = automaticallyAdded.Bool
		}
		if combinableWithOtherOffers.Valid {
			offer.CombinableWithOtherOffers = combinableWithOtherOffers.Bool
		}
		if requiresRelatedTarQual.Valid {
			offer.RequiresRelatedTarQual = requiresRelatedTarQual.Bool
		}
		if totalitarianOffer.Valid {
			offer.TotalitarianOffer = totalitarianOffer.Bool
		}
		if useListForDiscounts.Valid {
			offer.UseListForDiscounts = useListForDiscounts.Bool
		}
		if endDate.Valid {
			offer.EndDate = &endDate.Time
		}
		if startDate.Valid {
			offer.StartDate = startDate.Time
		}
		if maxUsesPerCustomer.Valid {
			offer.MaxUsesPerCustomer = &maxUsesPerCustomer.Int64
		}
		if maxUses.Valid {
			maxUsesInt := int(maxUses.Int32)
			offer.MaxUses = &maxUsesInt
		}
		if minimumDaysPerUsage.Valid {
			offer.MinimumDaysPerUsage = &minimumDaysPerUsage.Int64
		}
		if offerItemQualifierRule.Valid {
			offer.OfferItemQualifierRule = offerItemQualifierRule.String
		}
		if offerItemTargetRule.Valid {
			offer.OfferItemTargetRule = offerItemTargetRule.String
		}
		if marketingMessage.Valid {
			offer.MarketingMessage = marketingMessage.String
		}
		if maxUsesStrategy.Valid {
			offer.MaxUsesStrategy = maxUsesStrategy.String
		}
		if offerDescription.Valid {
			offer.OfferDescription = offerDescription.String
		}
		if targetSystem.Valid {
			offer.TargetSystem = targetSystem.String
		}

		offers = append(offers, offer)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.InternalWrap(err, "failed to iterate offers")
	}

	return offers, nil
}

// FindActiveOffers retrieves all currently active offers.
func (r *PostgresOfferRepository) FindActiveOffers(ctx context.Context) ([]*domain.Offer, error) {
	// Reusing FindAll with ActiveOnly filter
	return r.FindAll(ctx, &domain.OfferFilter{
		ActiveOnly: true,
	})
}

// Delete removes an offer by its unique identifier.
func (r *PostgresOfferRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_offer WHERE offer_id = $1`
	tag, err := r.db.Pool().Exec(ctx, query, id)
	if err != nil {
		return errors.InternalWrap(err, "failed to delete offer")
	}
	if tag.RowsAffected() == 0 {
		return errors.NotFound("offer not found")
	}
	return nil
}
