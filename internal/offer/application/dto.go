package application

import (
	"time"

	"github.com/qhato/ecommerce/internal/offer/domain"
)

// OfferDTO represents an offer data transfer object.
type OfferDTO struct {
	ID                        int64
	Name                      string
	OfferType                 domain.OfferType
	OfferValue                float64
	AdjustmentType            domain.OfferAdjustmentType
	ApplyToChildItems         bool
	ApplyToSalePrice          bool
	Archived                  bool
	AutomaticallyAdded        bool
	CombinableWithOtherOffers bool
	OfferDescription          string
	OfferDiscountType         domain.OfferDiscountType
	EndDate                   *time.Time
	MarketingMessage          string
	MaxUsesPerCustomer        *int64
	MaxUses                   *int
	MaxUsesStrategy           string
	MinimumDaysPerUsage       *int64
	OfferItemQualifierRule    string
	OfferItemTargetRule       string
	OrderMinTotal             float64
	OfferPriority             int
	QualifyingItemMinTotal    float64
	RequiresRelatedTarQual    bool
	StartDate                 time.Time
	TargetMinTotal            float64
	TargetSystem              string
	TotalitarianOffer         bool
	UseListForDiscounts       bool
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

// OfferCodeDTO represents an offer code data transfer object.
type OfferCodeDTO struct {
	ID           int64
	OfferID      int64
	Code         string
	MaxUses      *int
	Uses         int
	EmailAddress *string
	StartDate    *time.Time
	EndDate      *time.Time
	Archived     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// OfferItemCriteriaDTO represents offer item criteria data transfer object.
type OfferItemCriteriaDTO struct {
	ID                 int64
	Quantity           int
	OrderItemMatchRule string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// QualCritOfferXrefDTO represents a qualifying criteria xref data transfer object.
type QualCritOfferXrefDTO struct {
	ID                  int64
	OfferID             int64
	OfferItemCriteriaID int64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// TarCritOfferXrefDTO represents a target criteria xref data transfer object.
type TarCritOfferXrefDTO struct {
	ID                  int64
	OfferID             int64
	OfferItemCriteriaID int64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// OfferPriceDataDTO represents offer price data transfer object.
type OfferPriceDataDTO struct {
	ID              int64
	OfferID         int64
	Amount          float64
	DiscountType    string
	IdentifierType  string
	IdentifierValue string
	Quantity        int
	StartDate       *time.Time
	EndDate         *time.Time
	Archived        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ToOfferDTO converts a domain Offer to an OfferDTO.
func ToOfferDTO(offer *domain.Offer) *OfferDTO {
	return &OfferDTO{
		ID:                        offer.ID,
		Name:                      offer.Name,
		OfferType:                 offer.OfferType,
		OfferValue:                offer.OfferValue,
		AdjustmentType:            offer.AdjustmentType,
		ApplyToChildItems:         offer.ApplyToChildItems,
		ApplyToSalePrice:          offer.ApplyToSalePrice,
		Archived:                  offer.Archived,
		AutomaticallyAdded:        offer.AutomaticallyAdded,
		CombinableWithOtherOffers: offer.CombinableWithOtherOffers,
		OfferDescription:          offer.OfferDescription,
		OfferDiscountType:         offer.OfferDiscountType,
		EndDate:                   offer.EndDate,
		MarketingMessage:          offer.MarketingMessage,
		MaxUsesPerCustomer:        offer.MaxUsesPerCustomer,
		MaxUses:                   offer.MaxUses,
		MaxUsesStrategy:           offer.MaxUsesStrategy,
		MinimumDaysPerUsage:       offer.MinimumDaysPerUsage,
		OfferItemQualifierRule:    offer.OfferItemQualifierRule,
		OfferItemTargetRule:       offer.OfferItemTargetRule,
		OrderMinTotal:             offer.OrderMinTotal,
		OfferPriority:             offer.OfferPriority,
		QualifyingItemMinTotal:    offer.QualifyingItemMinTotal,
		RequiresRelatedTarQual:    offer.RequiresRelatedTarQual,
		StartDate:                 offer.StartDate,
		TargetMinTotal:            offer.TargetMinTotal,
		TargetSystem:              offer.TargetSystem,
		TotalitarianOffer:         offer.TotalitarianOffer,
		UseListForDiscounts:       offer.UseListForDiscounts,
		CreatedAt:                 offer.CreatedAt,
		UpdatedAt:                 offer.UpdatedAt,
	}
}

// ToOfferCodeDTO converts a domain OfferCode to an OfferCodeDTO.
func ToOfferCodeDTO(offerCode *domain.OfferCode) *OfferCodeDTO {
	return &OfferCodeDTO{
		ID:           offerCode.ID,
		OfferID:      offerCode.OfferID,
		Code:         offerCode.Code,
		MaxUses:      offerCode.MaxUses,
		Uses:         offerCode.Uses,
		EmailAddress: offerCode.EmailAddress,
		StartDate:    offerCode.StartDate,
		EndDate:      offerCode.EndDate,
		Archived:     offerCode.Archived,
		CreatedAt:    offerCode.CreatedAt,
		UpdatedAt:    offerCode.UpdatedAt,
	}
}

// ToOfferItemCriteriaDTO converts a domain OfferItemCriteria to an OfferItemCriteriaDTO.
func ToOfferItemCriteriaDTO(criteria *domain.OfferItemCriteria) *OfferItemCriteriaDTO {
	return &OfferItemCriteriaDTO{
		ID:                 criteria.ID,
		Quantity:           criteria.Quantity,
		OrderItemMatchRule: criteria.OrderItemMatchRule,
		CreatedAt:          criteria.CreatedAt,
		UpdatedAt:          criteria.UpdatedAt,
	}
}

// ToQualCritOfferXrefDTO converts a domain QualCritOfferXref to a QualCritOfferXrefDTO.
func ToQualCritOfferXrefDTO(xref *domain.QualCritOfferXref) *QualCritOfferXrefDTO {
	return &QualCritOfferXrefDTO{
		ID:                  xref.ID,
		OfferID:             xref.OfferID,
		OfferItemCriteriaID: xref.OfferItemCriteriaID,
		CreatedAt:           xref.CreatedAt,
		UpdatedAt:           xref.UpdatedAt,
	}
}

// ToTarCritOfferXrefDTO converts a domain TarCritOfferXref to a TarCritOfferXrefDTO.
func ToTarCritOfferXrefDTO(xref *domain.TarCritOfferXref) *TarCritOfferXrefDTO {
	return &TarCritOfferXrefDTO{
		ID:                  xref.ID,
		OfferID:             xref.OfferID,
		OfferItemCriteriaID: xref.OfferItemCriteriaID,
		CreatedAt:           xref.CreatedAt,
		UpdatedAt:           xref.UpdatedAt,
	}
}

// ToOfferPriceDataDTO converts a domain OfferPriceData to an OfferPriceDataDTO.
func ToOfferPriceDataDTO(priceData *domain.OfferPriceData) *OfferPriceDataDTO {
	return &OfferPriceDataDTO{
		ID:              priceData.ID,
		OfferID:         priceData.OfferID,
		Amount:          priceData.Amount,
		DiscountType:    priceData.DiscountType,
		IdentifierType:  priceData.IdentifierType,
		IdentifierValue: priceData.IdentifierValue,
		Quantity:        priceData.Quantity,
		StartDate:       priceData.StartDate,
		EndDate:         priceData.EndDate,
		Archived:        priceData.Archived,
		CreatedAt:       priceData.CreatedAt,
		UpdatedAt:       priceData.UpdatedAt,
	}
}

// ToOfferDomain converts an OfferDTO to a domain Offer.
func ToOfferDomain(offerDTO OfferDTO) *domain.Offer {
	return &domain.Offer{
		ID: offerDTO.ID,
		Name: offerDTO.Name,
		OfferType: offerDTO.OfferType,
		OfferValue: offerDTO.OfferValue,
		AdjustmentType: offerDTO.AdjustmentType,
		ApplyToChildItems: offerDTO.ApplyToChildItems,
		ApplyToSalePrice: offerDTO.ApplyToSalePrice,
		Archived: offerDTO.Archived,
		AutomaticallyAdded: offerDTO.AutomaticallyAdded,
		CombinableWithOtherOffers: offerDTO.CombinableWithOtherOffers,
		OfferDescription: offerDTO.OfferDescription,
		OfferDiscountType: offerDTO.OfferDiscountType,
		EndDate: offerDTO.EndDate,
		MarketingMessage: offerDTO.MarketingMessage,
		MaxUsesPerCustomer: offerDTO.MaxUsesPerCustomer,
		MaxUses: offerDTO.MaxUses,
		MaxUsesStrategy: offerDTO.MaxUsesStrategy,
		MinimumDaysPerUsage: offerDTO.MinimumDaysPerUsage,
		OfferItemQualifierRule: offerDTO.OfferItemQualifierRule,
		OfferItemTargetRule: offerDTO.OfferItemTargetRule,
		OrderMinTotal: offerDTO.OrderMinTotal,
		OfferPriority: offerDTO.OfferPriority,
		QualifyingItemMinTotal: offerDTO.QualifyingItemMinTotal,
		RequiresRelatedTarQual: offerDTO.RequiresRelatedTarQual,
		StartDate: offerDTO.StartDate,
		TargetMinTotal: offerDTO.TargetMinTotal,
		TargetSystem: offerDTO.TargetSystem,
		TotalitarianOffer: offerDTO.TotalitarianOffer,
		UseListForDiscounts: offerDTO.UseListForDiscounts,
		CreatedAt: offerDTO.CreatedAt,
		UpdatedAt: offerDTO.UpdatedAt,
	}
}
