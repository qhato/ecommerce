package domain

import "time"

// OfferType defines the type of offer (e.g., PERCENT_OFF, AMOUNT_OFF, BOGO)
type OfferType string

// OfferAdjustmentType defines how the adjustment is applied (e.g., ORDER_ITEM_OFFER, ORDER_OFFER)
type OfferAdjustmentType string

// OfferDiscountType defines how the discount is applied (e.g., FIX_PRICE, PERCENT_DISCOUNT)
type OfferDiscountType string

const (
	OfferTypePercentageOff OfferType = "PERCENTAGE_OFF"
	OfferTypeAmountOff     OfferType = "AMOUNT_OFF"
	OfferTypeBOGO          OfferType = "BOGO" // Buy One Get One

	OfferAdjustmentTypeOrderItem OfferAdjustmentType = "ORDER_ITEM_OFFER"
	OfferAdjustmentTypeOrder    OfferAdjustmentType = "ORDER_OFFER"

	OfferDiscountTypeFixPrice     OfferDiscountType = "FIX_PRICE"
	OfferDiscountTypePercentDiscount OfferDiscountType = "PERCENT_DISCOUNT"
	OfferDiscountTypeAmountOff    OfferDiscountType = "AMOUNT_OFF"
)

// Offer represents a promotional offer or discount
type Offer struct {
	ID                       int64
	Name                     string    // From blc_offer.offer_name
	OfferType                OfferType // From blc_offer.offer_type
	OfferValue               float64   // From blc_offer.offer_value (numeric(19,5))
	AdjustmentType           OfferAdjustmentType // From blc_offer.offer_adjustment_type
	ApplyToChildItems        bool      // From blc_offer.apply_to_child_items
	ApplyToSalePrice         bool      // From blc_offer.apply_to_sale_price
	Archived                 bool      // From blc_offer.archived (bpchar(1) 'Y'/'N')
	AutomaticallyAdded       bool      // From blc_offer.automatically_added
	CombinableWithOtherOffers bool      // From blc_offer.combinable_with_other_offers
	OfferDescription         string    // From blc_offer.offer_description
	OfferDiscountType        OfferDiscountType // From blc_offer.offer_discount_type
	EndDate                  *time.Time // From blc_offer.end_date
	MarketingMessage         string    // From blc_offer.marketing_message
	MaxUsesPerCustomer       *int64    // From blc_offer.max_uses_per_customer (int8)
	MaxUses                  *int      // From blc_offer.max_uses (int4)
	MaxUsesStrategy          string    // From blc_offer.max_uses_strategy
	MinimumDaysPerUsage      *int64    // From blc_offer.minimum_days_per_usage (int8)
	OfferItemQualifierRule   string    // From blc_offer.offer_item_qualifier_rule (text)
	OfferItemTargetRule      string    // From blc_offer.offer_item_target_rule (text)
	OrderMinTotal            float64   // From blc_offer.order_min_total (numeric(19,5))
	OfferPriority            int       // From blc_offer.offer_priority (int4)
	QualifyingItemMinTotal   float64   // From blc_offer.qualifying_item_min_total (numeric(19,5))
	RequiresRelatedTarQual   bool      // From blc_offer.requires_related_tar_qual
	StartDate                time.Time // From blc_offer.start_date
	TargetMinTotal           float64   // From blc_offer.target_min_total (numeric(19,5))
	TargetSystem             string    // From blc_offer.target_system
	TotalitarianOffer        bool      // From blc_offer.totalitarian_offer
	UseListForDiscounts      bool      // From blc_offer.use_list_for_discounts

	CreatedAt                time.Time
	UpdatedAt                time.Time
}

// NewOffer creates a new offer
func NewOffer(
	name string,
	offerType OfferType,
	offerValue float64,
	adjustmentType OfferAdjustmentType,
	startDate time.Time,
) (*Offer, error) {
	now := time.Now()
	if offerValue < 0 {
		return nil, NewDomainError("Offer value cannot be negative")
	}

	return &Offer{
		Name:                     name,
		OfferType:                offerType,
		OfferValue:               offerValue,
		AdjustmentType:           adjustmentType,
		StartDate:                startDate,
		Archived:                 false,
		ApplyToChildItems:        false,
		ApplyToSalePrice:         false,
		AutomaticallyAdded:       false,
		CombinableWithOtherOffers: true,
		OfferDescription:         "",
		OfferDiscountType:        "", // To be set explicitly if needed
		MarketingMessage:         "",
		MaxUsesStrategy:          "",
		OfferItemQualifierRule:   "",
		OfferItemTargetRule:      "",
		OrderMinTotal:            0.0,
		OfferPriority:            50, // Default priority
		QualifyingItemMinTotal:   0.0,
		RequiresRelatedTarQual:   false,
		TargetMinTotal:           0.0,
		TargetSystem:             "",
		TotalitarianOffer:        false,
		UseListForDiscounts:      false,
		CreatedAt:                now,
		UpdatedAt:                now,
	},
	nil
}

// Activate sets the offer to active (by unarchiving)
func (o *Offer) Activate() {
	o.Archived = false
	o.UpdatedAt = time.Now()
}

// Deactivate sets the offer to inactive (by archiving)
func (o *Offer) Deactivate() {
	o.Archived = true
	o.UpdatedAt = time.Now()
}

// SetEndDate sets the end date for the offer
func (o *Offer) SetEndDate(endDate time.Time) {
	o.EndDate = &endDate
	o.UpdatedAt = time.Now()
}

// SetMaxUses sets the maximum number of uses for the offer (total)
func (o *Offer) SetMaxUses(maxUses int) {
	maxUsesInt64 := int64(maxUses) // Convert to int64 for the nullable field
	o.MaxUses = &maxUses
	o.UpdatedAt = time.Now()
}

// SetMaxUsesPerCustomer sets the maximum number of uses per customer for the offer
func (o *Offer) SetMaxUsesPerCustomer(maxUses int64) {
	o.MaxUsesPerCustomer = &maxUses
	o.UpdatedAt = time.Now()
}

// SetMinimumDaysPerUsage sets the minimum days between usage for the offer
func (o *Offer) SetMinimumDaysPerUsage(days int64) {
	o.MinimumDaysPerUsage = &days
	o.UpdatedAt = time.Now()
}

// SetOrderMinTotal sets the minimum order subtotal required for the offer
func (o *Offer) SetOrderMinTotal(subtotal float64) {
	o.OrderMinTotal = subtotal
	o.UpdatedAt = time.Now()
}

// SetQualifyingItemMinTotal sets the minimum total for qualifying items required
func (o *Offer) SetQualifyingItemMinTotal(total float64) {
	o.QualifyingItemMinTotal = total
	o.UpdatedAt = time.Now()
}

// SetTargetMinTotal sets the minimum total for target items
func (o *Offer) SetTargetMinTotal(total float64) {
	o.TargetMinTotal = total
	o.UpdatedAt = time.Now()
}

// SetOfferPriority sets the priority of the offer (lower number = higher priority)
func (o *Offer) SetOfferPriority(priority int) {
	o.OfferPriority = priority
	o.UpdatedAt = time.Now()
}

// SetOfferDescription sets the offer description
func (o *Offer) SetOfferDescription(description string) {
	o.OfferDescription = description
	o.UpdatedAt = time.Now()
}

// SetMarketingMessage sets the marketing message
func (o *Offer) SetMarketingMessage(message string) {
	o.MarketingMessage = message
	o.UpdatedAt = time.Now()
}

// SetOfferItemQualifierRule sets the offer item qualifier rule
func (o *Offer) SetOfferItemQualifierRule(rule string) {
	o.OfferItemQualifierRule = rule
	o.UpdatedAt = time.Now()
}

// SetOfferItemTargetRule sets the offer item target rule
func (o *Offer) SetOfferItemTargetRule(rule string) {
	o.OfferItemTargetRule = rule
	o.UpdatedAt = time.Now()
}

// SetMaxUsesStrategy sets the maximum uses strategy
func (o *Offer) SetMaxUsesStrategy(strategy string) {
	o.MaxUsesStrategy = strategy
	o.UpdatedAt = time.Now()
}

// SetTargetSystem sets the target system
func (o *Offer) SetTargetSystem(system string) {
	o.TargetSystem = system
	o.UpdatedAt = time.Now()
}

// SetApplyToChildItems sets whether the offer applies to child items
func (o *Offer) SetApplyToChildItems(apply bool) {
	o.ApplyToChildItems = apply
	o.UpdatedAt = time.Now()
}

// SetApplyToSalePrice sets whether the offer applies to sale price
func (o *Offer) SetApplyToSalePrice(apply bool) {
	o.ApplyToSalePrice = apply
	o.UpdatedAt = time.Now()
}

// SetAutomaticallyAdded sets whether the offer is automatically added
func (o *Offer) SetAutomaticallyAdded(autoAdd bool) {
	o.AutomaticallyAdded = autoAdd
	o.UpdatedAt = time.Now()
}

// SetCombinableWithOtherOffers sets whether the offer is combinable with other offers
func (o *Offer) SetCombinableWithOtherOffers(combinable bool) {
	o.CombinableWithOtherOffers = combinable
	o.UpdatedAt = time.Now()
}

// SetRequiresRelatedTarQual sets whether the offer requires related target qualifier
func (o *Offer) SetRequiresRelatedTarQual(requires bool) {
	o.RequiresRelatedTarQual = requires
	o.UpdatedAt = time.Now()
}

// SetTotalitarianOffer sets whether the offer is a totalitarian offer
func (o *Offer) SetTotalitarianOffer(totalitarian bool) {
	o.TotalitarianOffer = totalitarian
	o.UpdatedAt = time.Now()
}

// SetUseListForDiscounts sets whether to use list for discounts
func (o *Offer) SetUseListForDiscounts(useList bool) {
	o.UseListForDiscounts = useList
	o.UpdatedAt = time.Now()
}

// DomainError represents a business rule validation error within the domain.
type DomainError struct {
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new DomainError.
func NewDomainError(message string) error {
	return &DomainError{Message: message}
}
