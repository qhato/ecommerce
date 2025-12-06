package queries

import "time"

// Offer Queries
type GetOfferQuery struct {
	ID int64
}

type GetActiveOffersQuery struct {
	CustomerID      *int64
	ProductIDs      []int64
	CategoryIDs     []int64
	OrderTotal      float64
	IncludeArchived bool
	Page            int
	PageSize        int
}

type SearchOffersQuery struct {
	SearchTerm      string
	OfferType       *string
	ActiveOnly      bool
	CustomerSegmentID *int64
	StartDate       *time.Time
	EndDate         *time.Time
	Page            int
	PageSize        int
	SortBy          string // name, priority, start_date, created_at
	SortOrder       string // asc, desc
}

// Offer Code Queries
type GetOfferCodeQuery struct {
	ID int64
}

type GetOfferByCodeQuery struct {
	Code string
}

type ValidateOfferCodeQuery struct {
	Code       string
	CustomerID *int64
	OrderTotal float64
}

// Analytics Queries
type GetOfferPerformanceQuery struct {
	OfferID   int64
	StartDate time.Time
	EndDate   time.Time
}

type GetOfferRevenueImpactQuery struct {
	OfferID   int64
	StartDate time.Time
	EndDate   time.Time
}

type GetTopPerformingOffersQuery struct {
	StartDate time.Time
	EndDate   time.Time
	Limit     int
	SortBy    string // usage_count, revenue_impact, avg_discount
}

// Customer Segment Queries
type GetCustomerSegmentQuery struct {
	ID int64
}

type GetAllCustomerSegmentsQuery struct {
	Page     int
	PageSize int
}

type GetCustomersInSegmentQuery struct {
	SegmentID int64
	Page      int
	PageSize  int
}

// Gift Card Queries
type GetGiftCardQuery struct {
	Code string
}

type GetGiftCardBalanceQuery struct {
	Code string
}

type GetCustomerGiftCardsQuery struct {
	CustomerID int64
}

// Loyalty Points Queries
type GetCustomerLoyaltyPointsQuery struct {
	CustomerID int64
}

type GetLoyaltyPointsHistoryQuery struct {
	CustomerID int64
	StartDate  *time.Time
	EndDate    *time.Time
	Page       int
	PageSize   int
}

type GetLoyaltyRulesQuery struct {
	ActiveOnly bool
}

// Referral Program Queries
type GetReferralProgramQuery struct {
	ID int64
}

type GetCustomerReferralsQuery struct {
	CustomerID int64
}

type GetReferralStatsQuery struct {
	CustomerID int64
}

// Tiered Offer Queries
type GetTieredOfferQuery struct {
	ID int64
}

type CalculateTieredDiscountQuery struct {
	OfferID    int64
	OrderTotal float64
}

// Bundle Offer Queries
type GetBundleOfferQuery struct {
	ID int64
}

type ValidateBundleQuery struct {
	OfferID    int64
	ProductIDs []int64
}
