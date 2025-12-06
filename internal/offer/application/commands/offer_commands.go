package commands

import "time"

// Offer Management Commands
type CreateOfferCommand struct {
	Name                      string
	OfferType                 string
	OfferValue                float64
	AdjustmentType            string
	ApplyToChildItems         bool
	ApplyToSalePrice          bool
	AutomaticallyAdded        bool
	CombinableWithOtherOffers bool
	OfferDescription          string
	OfferDiscountType         string
	EndDate                   *time.Time
	MarketingMessage          string
	MaxUsesPerCustomer        *int64
	MaxUses                   *int
	OrderMinTotal             float64
	OfferPriority             int
	StartDate                 time.Time
	
	// Advanced features
	CustomerSegmentIDs        []int64
	RequiredProductIDs        []int64
	ExcludedProductIDs        []int64
	RequiredCategoryIDs       []int64
	ExcludedCategoryIDs       []int64
}

type UpdateOfferCommand struct {
	ID                        int64
	Name                      string
	OfferValue                float64
	OfferDescription          string
	MarketingMessage          string
	EndDate                   *time.Time
	MaxUsesPerCustomer        *int64
	MaxUses                   *int
	OrderMinTotal             float64
	OfferPriority             int
	CombinableWithOtherOffers bool
	CustomerSegmentIDs        []int64
}

type ArchiveOfferCommand struct {
	ID int64
}

type DeleteOfferCommand struct {
	ID int64
}

// Offer Code Commands
type CreateOfferCodeCommand struct {
	OfferID           int64
	Code              string
	MaxUses           *int
	MaxUsesPerCustomer *int64
	StartDate         time.Time
	EndDate           *time.Time
}

type UpdateOfferCodeCommand struct {
	ID                int64
	MaxUses           *int
	MaxUsesPerCustomer *int64
	EndDate           *time.Time
}

type DeactivateOfferCodeCommand struct {
	ID int64
}

// Advanced Promotion Commands

// Tiered Offer Command
type CreateTieredOfferCommand struct {
	Name             string
	Description      string
	StartDate        time.Time
	EndDate          *time.Time
	Tiers            []OfferTier
	CustomerSegmentIDs []int64
}

type OfferTier struct {
	MinSpend     float64
	DiscountType string // PERCENT, AMOUNT
	DiscountValue float64
}

// Bundle Offer Command
type CreateBundleOfferCommand struct {
	Name               string
	Description        string
	StartDate          time.Time
	EndDate            *time.Time
	BundleProducts     []BundleProduct
	DiscountType       string
	DiscountValue      float64
	CustomerSegmentIDs []int64
}

type BundleProduct struct {
	ProductID int64
	Quantity  int
}

// Referral Program Command
type CreateReferralProgramCommand struct {
	Name                  string
	Description           string
	StartDate             time.Time
	EndDate               *time.Time
	ReferrerRewardType    string // PERCENT, AMOUNT, POINTS
	ReferrerRewardValue   float64
	RefereeRewardType     string
	RefereeRewardValue    float64
	MinPurchaseAmount     float64
	MaxReferralsPerCustomer *int
}

// Gift Card Commands
type CreateGiftCardCommand struct {
	Code          string
	Amount        float64
	PurchasedBy   *int64
	RecipientEmail string
	Message       string
	ExpiresAt     *time.Time
}

type RedeemGiftCardCommand struct {
	Code    string
	Amount  float64
	OrderID int64
}

// Loyalty Points Commands
type CreateLoyaltyRuleCommand struct {
	Name              string
	Description       string
	PointsPerDollar   float64
	MinPurchaseAmount float64
	ProductIDs        []int64
	CategoryIDs       []int64
	StartDate         time.Time
	EndDate           *time.Time
}

type AwardLoyaltyPointsCommand struct {
	CustomerID int64
	OrderID    int64
	Points     int64
	Reason     string
}

type RedeemLoyaltyPointsCommand struct {
	CustomerID      int64
	Points          int64
	OrderID         int64
	DiscountAmount  float64
}

// Customer Segmentation Commands
type CreateCustomerSegmentCommand struct {
	Name        string
	Description string
	Rules       []SegmentRule
}

type SegmentRule struct {
	Field    string // total_spent, order_count, last_order_date, etc.
	Operator string // >, <, >=, <=, =, IN
	Value    interface{}
}

type UpdateCustomerSegmentCommand struct {
	ID          int64
	Name        string
	Description string
	Rules       []SegmentRule
}
