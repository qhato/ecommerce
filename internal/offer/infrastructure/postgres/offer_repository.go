package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/offer/domain"
)

// OfferRepository implements domain.OfferRepository for PostgreSQL persistence.
type OfferRepository struct {
	db *sql.DB
}

// NewOfferRepository creates a new PostgreSQL offer repository.
func NewOfferRepository(db *sql.DB) *OfferRepository {
	return &OfferRepository{db: db}
}

// Save stores a new offer or updates an existing one.
func (r *OfferRepository) Save(ctx context.Context, offer *domain.Offer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Handle BPCHAR(1) conversion
	archivedChar := "N"
	if offer.Archived {
		archivedChar = "Y"
	}

	// Handle nullable fields
	applyToChildItems := sql.NullBool{Bool: offer.ApplyToChildItems, Valid: true}
	applyToSalePrice := sql.NullBool{Bool: offer.ApplyToSalePrice, Valid: true}
	automaticallyAdded := sql.NullBool{Bool: offer.AutomaticallyAdded, Valid: true}
	combinableWithOtherOffers := sql.NullBool{Bool: offer.CombinableWithOtherOffers, Valid: true}
	requiresRelatedTarQual := sql.NullBool{Bool: offer.RequiresRelatedTarQual, Valid: true}
	totalitarianOffer := sql.NullBool{Bool: offer.TotalitarianOffer, Valid: true}
	useListForDiscounts := sql.NullBool{Bool: offer.UseListForDiscounts, Valid: true}

	offerDescription := sql.NullString{String: offer.OfferDescription, Valid: offer.OfferDescription != ""}
	offerDiscountType := sql.NullString{String: string(offer.OfferDiscountType), Valid: offer.OfferDiscountType != ""}
	endDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if offer.EndDate != nil {
		endDate = sql.NullTime{Time: *offer.EndDate, Valid: true}
	}
	marketingMessage := sql.NullString{String: offer.MarketingMessage, Valid: offer.MarketingMessage != ""}
	maxUsesPerCustomer := sql.NullInt64{Int64: 0, Valid: false}
	if offer.MaxUsesPerCustomer != nil {
		maxUsesPerCustomer = sql.NullInt64{Int64: *offer.MaxUsesPerCustomer, Valid: true}
	}
	maxUses := sql.NullInt32{Int32: 0, Valid: false}
	if offer.MaxUses != nil {
		maxUses = sql.NullInt32{Int32: int32(*offer.MaxUses), Valid: true}
	}
	maxUsesStrategy := sql.NullString{String: offer.MaxUsesStrategy, Valid: offer.MaxUsesStrategy != ""}
	minimumDaysPerUsage := sql.NullInt64{Int64: 0, Valid: false}
	if offer.MinimumDaysPerUsage != nil {
		minimumDaysPerUsage = sql.NullInt64{Int64: *offer.MinimumDaysPerUsage, Valid: true}
	}
	offerItemQualifierRule := sql.NullString{String: offer.OfferItemQualifierRule, Valid: offer.OfferItemQualifierRule != ""}
	offerItemTargetRule := sql.NullString{String: offer.OfferItemTargetRule, Valid: offer.OfferItemTargetRule != ""}
	orderMinTotal := sql.NullFloat64{Float64: offer.OrderMinTotal, Valid: offer.OrderMinTotal != 0.0}
	qualifyingItemMinTotal := sql.NullFloat64{Float64: offer.QualifyingItemMinTotal, Valid: offer.QualifyingItemMinTotal != 0.0}
	targetMinTotal := sql.NullFloat64{Float64: offer.TargetMinTotal, Valid: offer.TargetMinTotal != 0.0}
	targetSystem := sql.NullString{String: offer.TargetSystem, Valid: offer.TargetSystem != ""}

	if offer.ID == 0 {
		// Insert new offer
		query := `
			INSERT INTO blc_offer (
				offer_adjustment_type, apply_to_child_items, apply_to_sale_price, archived, 
				automatically_added, combinable_with_other_offers, offer_description, 
				offer_discount_type, end_date, marketing_message, max_uses_per_customer, 
				max_uses, max_uses_strategy, minimum_days_per_usage, offer_name, 
				offer_item_qualifier_rule, offer_item_target_rule, order_min_total, 
				offer_priority, qualifying_item_min_total, requires_related_tar_qual, 
				start_date, target_min_total, target_system, totalitarian_offer, 
				offer_type, use_list_for_discounts, offer_value, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, 
				$19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
			) RETURNING offer_id`
		err = tx.QueryRowContext(ctx, query,
			offer.AdjustmentType, applyToChildItems, applyToSalePrice, archivedChar,
			automaticallyAdded, combinableWithOtherOffers, offerDescription,
			offerDiscountType, endDate, marketingMessage, maxUsesPerCustomer,
			maxUses, maxUsesStrategy, minimumDaysPerUsage, offer.Name,
			offerItemQualifierRule, offerItemTargetRule, orderMinTotal,
			offer.OfferPriority, qualifyingItemMinTotal, requiresRelatedTarQual,
			offer.StartDate, targetMinTotal, targetSystem, totalitarianOffer,
			offer.OfferType, useListForDiscounts, offer.OfferValue, offer.CreatedAt, offer.UpdatedAt,
		).Scan(&offer.ID)
		if err != nil {
			return fmt.Errorf("failed to insert offer: %w", err)
		}
	} else {
		// Update existing offer
		query := `
			UPDATE blc_offer SET
				offer_adjustment_type = $1, apply_to_child_items = $2, apply_to_sale_price = $3, 
				archived = $4, automatically_added = $5, combinable_with_other_offers = $6, 
				offer_description = $7, offer_discount_type = $8, end_date = $9, 
				marketing_message = $10, max_uses_per_customer = $11, max_uses = $12, 
				max_uses_strategy = $13, minimum_days_per_usage = $14, offer_name = $15, 
				offer_item_qualifier_rule = $16, offer_item_target_rule = $17, 
				order_min_total = $18, offer_priority = $19, qualifying_item_min_total = $20, 
				requires_related_tar_qual = $21, start_date = $22, target_min_total = $23, 
				target_system = $24, totalitarian_offer = $25, offer_type = $26, 
				use_list_for_discounts = $27, offer_value = $28, updated_at = $29
			WHERE offer_id = $30`
		_, err = tx.ExecContext(ctx, query,
			offer.AdjustmentType, applyToChildItems, applyToSalePrice,
			archivedChar, automaticallyAdded, combinableWithOtherOffers,
			offerDescription, offerDiscountType, endDate,
			marketingMessage, maxUsesPerCustomer, maxUses,
			maxUsesStrategy, minimumDaysPerUsage, offer.Name,
			offerItemQualifierRule, offerItemTargetRule,
			orderMinTotal, offer.OfferPriority, qualifyingItemMinTotal,
			requiresRelatedTarQual, offer.StartDate, targetMinTotal,
			targetSystem, totalitarianOffer, offer.OfferType,
			useListForDiscounts, offer.OfferValue, offer.UpdatedAt, offer.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update offer: %w", err)
		}
	}

	return tx.Commit()
}

// FindByID retrieves an offer by its unique identifier.
func (r *OfferRepository) FindByID(ctx context.Context, id int64) (*domain.Offer, error) {
	query := `
		SELECT
			offer_id, offer_adjustment_type, apply_to_child_items, apply_to_sale_price, 
			archived, automatically_added, combinable_with_other_offers, offer_description, 
			offer_discount_type, end_date, marketing_message, max_uses_per_customer, 
			max_uses, max_uses_strategy, minimum_days_per_usage, offer_name, 
			offer_item_qualifier_rule, offer_item_target_rule, order_min_total, 
			offer_priority, qualifying_item_min_total, requires_related_tar_qual, 
			start_date, target_min_total, target_system, totalitarian_offer, 
			offer_type, use_list_for_discounts, offer_value, created_at, updated_at
		FROM blc_offer WHERE offer_id = $1`

	var offer domain.Offer
	var archivedChar string
	var applyToChildItems sql.NullBool
	var applyToSalePrice sql.NullBool
	var automaticallyAdded sql.NullBool
	var combinableWithOtherOffers sql.NullBool
	var offerDescription sql.NullString
	var offerDiscountType sql.NullString
	var endDate sql.NullTime
	var marketingMessage sql.NullString
	var maxUsesPerCustomer sql.NullInt64
	var maxUses sql.NullInt32
	var maxUsesStrategy sql.NullString
	var minimumDaysPerUsage sql.NullInt64
	var offerItemQualifierRule sql.NullString
	var offerItemTargetRule sql.NullString
	var orderMinTotal sql.NullFloat64
	var qualifyingItemMinTotal sql.NullFloat64
	var requiresRelatedTarQual sql.NullBool
	var targetMinTotal sql.NullFloat64
	var targetSystem sql.NullString
	var totalitarianOffer sql.NullBool
	var useListForDiscounts sql.NullBool

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&offer.ID, &offer.AdjustmentType, &applyToChildItems, &applyToSalePrice,
		&archivedChar, &automaticallyAdded, &combinableWithOtherOffers, &offerDescription,
		&offerDiscountType, &endDate, &marketingMessage, &maxUsesPerCustomer,
		&maxUses, &maxUsesStrategy, &minimumDaysPerUsage, &offer.Name,
		&offerItemQualifierRule, &offerItemTargetRule, &orderMinTotal,
		&offer.OfferPriority, &qualifyingItemMinTotal, &requiresRelatedTarQual,
		&offer.StartDate, &targetMinTotal, &targetSystem, &totalitarianOffer,
		&offer.OfferType, &useListForDiscounts, &offer.OfferValue, &offer.CreatedAt, &offer.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query offer by ID: %w", err)
	}

	offer.Archived = (archivedChar == "Y")
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
	if offerDescription.Valid {
		offer.OfferDescription = offerDescription.String
	}
	if offerDiscountType.Valid {
		offer.OfferDiscountType = domain.OfferDiscountType(offerDiscountType.String)
	}
	if endDate.Valid {
		offer.EndDate = &endDate.Time
	}
	if marketingMessage.Valid {
		offer.MarketingMessage = marketingMessage.String
	}
	if maxUsesPerCustomer.Valid {
		offer.MaxUsesPerCustomer = &maxUsesPerCustomer.Int64
	}
	if maxUses.Valid {
		intMaxUses := int(maxUses.Int32)
		offer.MaxUses = &intMaxUses
	}
	if maxUsesStrategy.Valid {
		offer.MaxUsesStrategy = maxUsesStrategy.String
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
	if orderMinTotal.Valid {
		offer.OrderMinTotal = orderMinTotal.Float64
	}
	if qualifyingItemMinTotal.Valid {
		offer.QualifyingItemMinTotal = qualifyingItemMinTotal.Float64
	}
	if requiresRelatedTarQual.Valid {
		offer.RequiresRelatedTarQual = requiresRelatedTarQual.Bool
	}
	if targetMinTotal.Valid {
		offer.TargetMinTotal = targetMinTotal.Float64
	}
	if targetSystem.Valid {
		offer.TargetSystem = targetSystem.String
	}
	if totalitarianOffer.Valid {
		offer.TotalitarianOffer = totalitarianOffer.Bool
	}
	if useListForDiscounts.Valid {
		offer.UseListForDiscounts = useListForDiscounts.Bool
	}

	return &offer, nil
}

// FindActiveOffers retrieves all currently active offers.
func (r *OfferRepository) FindActiveOffers(ctx context.Context) ([]*domain.Offer, error) {
	// Active implies archived = 'N' and current date is within start_date and end_date
	query := `
		SELECT
			offer_id, offer_adjustment_type, apply_to_child_items, apply_to_sale_price, 
			archived, automatically_added, combinable_with_other_offers, offer_description, 
			offer_discount_type, end_date, marketing_message, max_uses_per_customer, 
			max_uses, max_uses_strategy, minimum_days_per_usage, offer_name, 
			offer_item_qualifier_rule, offer_item_target_rule, order_min_total, 
			offer_priority, qualifying_item_min_total, requires_related_tar_qual, 
			start_date, target_min_total, target_system, totalitarian_offer, 
			offer_type, use_list_for_discounts, offer_value, created_at, updated_at
		FROM blc_offer 
		WHERE archived = 'N' AND start_date <= NOW() AND (end_date IS NULL OR end_date >= NOW())`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active offers: %w", err)
	}
	defer rows.Close()

	var offers []*domain.Offer
	for rows.Next() {
		var offer domain.Offer
		var archivedChar string
		var applyToChildItems sql.NullBool
		var applyToSalePrice sql.NullBool
		var automaticallyAdded sql.NullBool
		var combinableWithOtherOffers sql.NullBool
		var offerDescription sql.NullString
		var offerDiscountType sql.NullString
		var endDate sql.NullTime
		var marketingMessage sql.NullString
		var maxUsesPerCustomer sql.NullInt64
		var maxUses sql.NullInt32
		var maxUsesStrategy sql.NullString
		var minimumDaysPerUsage sql.NullInt64
		var offerItemQualifierRule sql.NullString
		var offerItemTargetRule sql.NullString
		var orderMinTotal sql.NullFloat64
		var qualifyingItemMinTotal sql.NullFloat64
		var requiresRelatedTarQual sql.NullBool
		var targetMinTotal sql.NullFloat64
		var targetSystem sql.NullString
		var totalitarianOffer sql.NullBool
		var useListForDiscounts sql.NullBool

		err := rows.Scan(
			&offer.ID, &offer.AdjustmentType, &applyToChildItems, &applyToSalePrice,
			&archivedChar, &automaticallyAdded, &combinableWithOtherOffers, &offerDescription,
			&offerDiscountType, &endDate, &marketingMessage, &maxUsesPerCustomer,
			&maxUses, &maxUsesStrategy, &minimumDaysPerUsage, &offer.Name,
			&offerItemQualifierRule, &offerItemTargetRule, &orderMinTotal,
			&offer.OfferPriority, &qualifyingItemMinTotal, &requiresRelatedTarQual,
			&offer.StartDate, &targetMinTotal, &targetSystem, &totalitarianOffer,
			&offer.OfferType, &useListForDiscounts, &offer.OfferValue, &offer.CreatedAt, &offer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan offer row: %w", err)
		}

		offer.Archived = (archivedChar == "Y")
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
		if offerDescription.Valid {
			offer.OfferDescription = offerDescription.String
		}
		if offerDiscountType.Valid {
			offer.OfferDiscountType = domain.OfferDiscountType(offerDiscountType.String)
		}
		if endDate.Valid {
			offer.EndDate = &endDate.Time
		}
		if marketingMessage.Valid {
			offer.MarketingMessage = marketingMessage.String
		}
		if maxUsesPerCustomer.Valid {
			offer.MaxUsesPerCustomer = &maxUsesPerCustomer.Int64
		}
		if maxUses.Valid {
			intMaxUses := int(maxUses.Int32)
			offer.MaxUses = &intMaxUses
		}
		if maxUsesStrategy.Valid {
			offer.MaxUsesStrategy = maxUsesStrategy.String
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
		if orderMinTotal.Valid {
			offer.OrderMinTotal = orderMinTotal.Float64
		}
		if qualifyingItemMinTotal.Valid {
			offer.QualifyingItemMinTotal = qualifyingItemMinTotal.Float64
		}
		if requiresRelatedTarQual.Valid {
			offer.RequiresRelatedTarQual = requiresRelatedTarQual.Bool
		}
		if targetMinTotal.Valid {
			offer.TargetMinTotal = targetMinTotal.Float64
		}
		if targetSystem.Valid {
			offer.TargetSystem = targetSystem.String
		}
		if totalitarianOffer.Valid {
			offer.TotalitarianOffer = totalitarianOffer.Bool
		}
		if useListForDiscounts.Valid {
			offer.UseListForDiscounts = useListForDiscounts.Bool
		}

		offers = append(offers, &offer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for active offers: %w", err)
	}

	return offers, nil
}

// FindByCode retrieves an offer by its promotional code.
// Note: This needs to join with blc_offer_code to find the offer.
func (r *OfferRepository) FindByCode(ctx context.Context, code string) (*domain.Offer, error) {
	query := `
		SELECT
			o.offer_id, o.offer_adjustment_type, o.apply_to_child_items, o.apply_to_sale_price, 
			o.archived, o.automatically_added, o.combinable_with_other_offers, o.offer_description, 
			o.offer_discount_type, o.end_date, o.marketing_message, o.max_uses_per_customer, 
			o.max_uses, o.max_uses_strategy, o.minimum_days_per_usage, o.offer_name, 
			o.offer_item_qualifier_rule, o.offer_item_target_rule, o.order_min_total, 
			o.offer_priority, o.qualifying_item_min_total, o.requires_related_tar_qual, 
			o.start_date, o.target_min_total, o.target_system, o.totalitarian_offer, 
			o.offer_type, o.use_list_for_discounts, o.offer_value, o.created_at, o.updated_at
		FROM blc_offer o
		JOIN blc_offer_code oc ON o.offer_id = oc.offer_id
		WHERE oc.offer_code = $1`

	var offer domain.Offer
	var archivedChar string
	var applyToChildItems sql.NullBool
	var applyToSalePrice sql.NullBool
	var automaticallyAdded sql.NullBool
	var combinableWithOtherOffers sql.NullBool
	var offerDescription sql.NullString
	var offerDiscountType sql.NullString
	var endDate sql.NullTime
	var marketingMessage sql.NullString
	var maxUsesPerCustomer sql.NullInt64
	var maxUses sql.NullInt32
	var maxUsesStrategy sql.NullString
	var minimumDaysPerUsage sql.NullInt64
	var offerItemQualifierRule sql.NullString
	var offerItemTargetRule sql.NullString
	var orderMinTotal sql.NullFloat64
	var qualifyingItemMinTotal sql.NullFloat64
	var requiresRelatedTarQual sql.NullBool
	var targetMinTotal sql.NullFloat64
	var targetSystem sql.NullString
	var totalitarianOffer sql.NullBool
	var useListForDiscounts sql.NullBool

	row := r.db.QueryRowContext(ctx, query, code)
	err := row.Scan(
		&offer.ID, &offer.AdjustmentType, &applyToChildItems, &applyToSalePrice,
		&archivedChar, &automaticallyAdded, &combinableWithOtherOffers, &offerDescription,
		&offerDiscountType, &endDate, &marketingMessage, &maxUsesPerCustomer,
		&maxUses, &maxUsesStrategy, &minimumDaysPerUsage, &offer.Name,
		&offerItemQualifierRule, &offerItemTargetRule, &orderMinTotal,
		&offer.OfferPriority, &qualifyingItemMinTotal, &requiresRelatedTarQual,
		&offer.StartDate, &targetMinTotal, &targetSystem, &totalitarianOffer,
		&offer.OfferType, &useListForDiscounts, &offer.OfferValue, &offer.CreatedAt, &offer.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query offer by code: %w", err)
	}

	offer.Archived = (archivedChar == "Y")
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
	if offerDescription.Valid {
		offer.OfferDescription = offerDescription.String
	}
	if offerDiscountType.Valid {
		offer.OfferDiscountType = domain.OfferDiscountType(offerDiscountType.String)
	}
	if endDate.Valid {
		offer.EndDate = &endDate.Time
	}
	if marketingMessage.Valid {
		offer.MarketingMessage = marketingMessage.String
	}
	if maxUsesPerCustomer.Valid {
		offer.MaxUsesPerCustomer = &maxUsesPerCustomer.Int64
	}
	if maxUses.Valid {
		intMaxUses := int(maxUses.Int32)
		offer.MaxUses = &intMaxUses
	}
	if maxUsesStrategy.Valid {
		offer.MaxUsesStrategy = maxUsesStrategy.String
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
	if orderMinTotal.Valid {
		offer.OrderMinTotal = orderMinTotal.Float64
	}
	if qualifyingItemMinTotal.Valid {
		offer.QualifyingItemMinTotal = qualifyingItemMinTotal.Float64
	}
	if requiresRelatedTarQual.Valid {
		offer.RequiresRelatedTarQual = requiresRelatedTarQual.Bool
	}
	if targetMinTotal.Valid {
		offer.TargetMinTotal = targetMinTotal.Float64
	}
	if targetSystem.Valid {
		offer.TargetSystem = targetSystem.String
	}
	if totalitarianOffer.Valid {
		offer.TotalitarianOffer = totalitarianOffer.Bool
	}
	if useListForDiscounts.Valid {
		offer.UseListForDiscounts = useListForDiscounts.Bool
	}

	return &offer, nil
}

// Delete removes an offer by its unique identifier.
func (r *OfferRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM blc_offer WHERE offer_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer: %w", err)
	}
	return nil
}
