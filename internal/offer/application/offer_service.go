package application

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/offer/domain"
)

// OfferService defines the application service for offer-related operations.
type OfferService interface {
	// CreateOffer creates a new offer.
	CreateOffer(ctx context.Context, cmd *CreateOfferCommand) (*OfferDTO, error)

	// GetOfferByID retrieves an offer by its ID.
	GetOfferByID(ctx context.Context, id int64) (*OfferDTO, error)

	// UpdateOffer updates an existing offer.
	UpdateOffer(ctx context.Context, cmd *UpdateOfferCommand) (*OfferDTO, error)

	// DeleteOffer deletes an offer.
	DeleteOffer(ctx context.Context, id int64) error

	// CreateOfferCode creates a new offer code for an existing offer.
	CreateOfferCode(ctx context.Context, offerID int64, cmd *CreateOfferCodeCommand) (*OfferCodeDTO, error)

	// GetOfferCodeByID retrieves an offer code by its ID.
	GetOfferCodeByID(ctx context.Context, id int64) (*OfferCodeDTO, error)

	// UpdateOfferCode updates an existing offer code.
	UpdateOfferCode(ctx context.Context, id int64, cmd *UpdateOfferCodeCommand) (*OfferCodeDTO, error)

	// DeleteOfferCode deletes an offer code.
	DeleteOfferCode(ctx context.Context, id int64) error

	// CreateOfferItemCriteria creates new offer item criteria.
	CreateOfferItemCriteria(ctx context.Context, cmd *CreateOfferItemCriteriaCommand) (*OfferItemCriteriaDTO, error)

	// GetOfferItemCriteriaByID retrieves offer item criteria by ID.
	GetOfferItemCriteriaByID(ctx context.Context, id int64) (*OfferItemCriteriaDTO, error)

	// UpdateOfferItemCriteria updates existing offer item criteria.
	UpdateOfferItemCriteria(ctx context.Context, id int64, cmd *UpdateOfferItemCriteriaCommand) (*OfferItemCriteriaDTO, error)

	// DeleteOfferItemCriteria deletes offer item criteria.
	DeleteOfferItemCriteria(ctx context.Context, id int64) error

	// AddQualifyingItemCriteriaToOffer associates qualifying item criteria with an offer.
	AddQualifyingItemCriteriaToOffer(ctx context.Context, offerID, offerItemCriteriaID int64) (*QualCritOfferXrefDTO, error)

	// RemoveQualifyingItemCriteriaFromOffer removes qualifying item criteria association from an offer.
	RemoveQualifyingItemCriteriaFromOffer(ctx context.Context, offerID, offerItemCriteriaID int64) error

	// AddTargetItemCriteriaToOffer associates target item criteria with an offer.
	AddTargetItemCriteriaToOffer(ctx context.Context, offerID, offerItemCriteriaID int64) (*TarCritOfferXrefDTO, error)

	// RemoveTargetItemCriteriaFromOffer removes target item criteria association from an offer.
	RemoveTargetItemCriteriaFromOffer(ctx context.Context, offerID, offerItemCriteriaID int64) error

	// CreateOfferPriceData creates new offer price data for an existing offer.
	CreateOfferPriceData(ctx context.Context, offerID int64, cmd *CreateOfferPriceDataCommand) (*OfferPriceDataDTO, error)

	// GetOfferPriceDataByID retrieves offer price data by ID.
	GetOfferPriceDataByID(ctx context.Context, id int64) (*OfferPriceDataDTO, error)

	// UpdateOfferPriceData updates existing offer price data.
	UpdateOfferPriceData(ctx context.Context, id int64, cmd *UpdateOfferPriceDataCommand) (*OfferPriceDataDTO, error)

	// DeleteOfferPriceData deletes offer price data.
	DeleteOfferPriceData(ctx context.Context, id int64) error
}

// OfferDTO represents an offer data transfer object.
type OfferDTO struct {
	ID                       int64
	Name                     string
	OfferType                domain.OfferType
	OfferValue               float64
	AdjustmentType           domain.OfferAdjustmentType
	ApplyToChildItems        bool
	ApplyToSalePrice         bool
	Archived                 bool
	AutomaticallyAdded       bool
	CombinableWithOtherOffers bool
	OfferDescription         string
	OfferDiscountType        domain.OfferDiscountType
	EndDate                  *time.Time
	MarketingMessage         string
	MaxUsesPerCustomer       *int64
	MaxUses                  *int
	MaxUsesStrategy          string
	MinimumDaysPerUsage      *int64
	OfferItemQualifierRule   string
	OfferItemTargetRule      string
	OrderMinTotal            float64
	OfferPriority            int
	QualifyingItemMinTotal   float64
	RequiresRelatedTarQual   bool
	StartDate                time.Time
	TargetMinTotal           float64
	TargetSystem             string
	TotalitarianOffer        bool
	UseListForDiscounts      bool
	CreatedAt                time.Time
	UpdatedAt                time.Time
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

// CreateOfferCommand is a command to create a new offer.
type CreateOfferCommand struct {
	Name                     string
	OfferType                domain.OfferType
	OfferValue               float64
	AdjustmentType           domain.OfferAdjustmentType
	ApplyToChildItems        bool
	ApplyToSalePrice         bool
	AutomaticallyAdded       bool
	CombinableWithOtherOffers bool
	OfferDescription         string
	OfferDiscountType        domain.OfferDiscountType
	EndDate                  *time.Time
	MarketingMessage         string
	MaxUsesPerCustomer       *int64
	MaxUses                  *int
	MaxUsesStrategy          string
	MinimumDaysPerUsage      *int64
	OfferItemQualifierRule   string
	OfferItemTargetRule      string
	OrderMinTotal            float64
	OfferPriority            int
	QualifyingItemMinTotal   float64
	RequiresRelatedTarQual   bool
	StartDate                time.Time
	TargetMinTotal           float64
	TargetSystem             string
	TotalitarianOffer        bool
	UseListForDiscounts      bool
}

// UpdateOfferCommand is a command to update an existing offer.
type UpdateOfferCommand struct {
	ID                       int64
	Name                     *string
	OfferType                *domain.OfferType
	OfferValue               *float64
	AdjustmentType           *domain.OfferAdjustmentType
	ApplyToChildItems        *bool
	ApplyToSalePrice         *bool
	Archived                 *bool
	AutomaticallyAdded       *bool
	CombinableWithOtherOffers *bool
	OfferDescription         *string
	OfferDiscountType        *domain.OfferDiscountType
	EndDate                  *time.Time
	MarketingMessage         *string
	MaxUsesPerCustomer       *int64
	MaxUses                  *int
	MaxUsesStrategy          *string
	MinimumDaysPerUsage      *int64
	OfferItemQualifierRule   *string
	OfferItemTargetRule      *string
	OrderMinTotal            *float64
	OfferPriority            *int
	QualifyingItemMinTotal   *float64
	RequiresRelatedTarQual   *bool
	StartDate                *time.Time
	TargetMinTotal           *float64
	TargetSystem             *string
	TotalitarianOffer        *bool
	UseListForDiscounts      *bool
}

// CreateOfferCodeCommand is a command to create a new offer code.
type CreateOfferCodeCommand struct {
	Code         string
	MaxUses      *int
	EmailAddress *string
	StartDate    *time.Time
	EndDate      *time.Time
}

// UpdateOfferCodeCommand is a command to update an existing offer code.
type UpdateOfferCodeCommand struct {
	Code         *string
	MaxUses      *int
	Uses         *int
	EmailAddress *string
	StartDate    *time.Time
	EndDate      *time.Time
	Archived     *bool
}

// CreateOfferItemCriteriaCommand is a command to create new offer item criteria.
type CreateOfferItemCriteriaCommand struct {
	Quantity           int
	OrderItemMatchRule string
}

// UpdateOfferItemCriteriaCommand is a command to update existing offer item criteria.
type UpdateOfferItemCriteriaCommand struct {
	Quantity           *int
	OrderItemMatchRule *string
}

// CreateOfferPriceDataCommand is a command to create new offer price data.
type CreateOfferPriceDataCommand struct {
	Amount          float64
	DiscountType    string
	IdentifierType  string
	IdentifierValue string
	Quantity        int
	StartDate       *time.Time
	EndDate         *time.Time
}

// UpdateOfferPriceDataCommand is a command to update existing offer price data.
type UpdateOfferPriceDataCommand struct {
	Amount          *float64
	DiscountType    *string
	IdentifierType  *string
	IdentifierValue *string
	Quantity        *int
	StartDate       *time.Time
	EndDate         *time.Time
	Archived        *bool
}

type offerService struct {
	offerRepo               domain.OfferRepository
	offerCodeRepo           domain.OfferCodeRepository
	offerItemCriteriaRepo   domain.OfferItemCriteriaRepository
	offerRuleRepo           domain.OfferRuleRepository // Not used yet, but kept for future expansion
	offerPriceDataRepo      domain.OfferPriceDataRepository
	qualCritOfferXrefRepo   domain.QualCritOfferXrefRepository
	tarCritOfferXrefRepo    domain.TarCritOfferXrefRepository
}

// NewOfferService creates a new instance of OfferService.
func NewOfferService(
	offerRepo domain.OfferRepository,
	offerCodeRepo domain.OfferCodeRepository,
	offerItemCriteriaRepo domain.OfferItemCriteriaRepository,
	offerRuleRepo domain.OfferRuleRepository,
	offerPriceDataRepo domain.OfferPriceDataRepository,
	qualCritOfferXrefRepo domain.QualCritOfferXrefRepository,
	tarCritOfferXrefRepo domain.TarCritOfferXrefRepository,
) OfferService {
	return &offerService{
		offerRepo:               offerRepo,
		offerCodeRepo:           offerCodeRepo,
		offerItemCriteriaRepo:   offerItemCriteriaRepo,
		offerRuleRepo:           offerRuleRepo,
		offerPriceDataRepo:      offerPriceDataRepo,
		qualCritOfferXrefRepo:   qualCritOfferXrefRepo,
		tarCritOfferXrefRepo:    tarCritOfferXrefRepo,
	}
}

func (s *offerService) CreateOffer(ctx context.Context, cmd *CreateOfferCommand) (*OfferDTO, error) {
	offer, err := domain.NewOffer(cmd.Name, cmd.OfferType, cmd.OfferValue, cmd.AdjustmentType, cmd.StartDate)
	if err != nil {
		return nil, fmt.Errorf("failed to create offer domain entity: %w", err)
	}

	offer.ApplyToChildItems = cmd.ApplyToChildItems
	offer.ApplyToSalePrice = cmd.ApplyToSalePrice
	offer.AutomaticallyAdded = cmd.AutomaticallyAdded
	offer.CombinableWithOtherOffers = cmd.CombinableWithOtherOffers
	offer.OfferDescription = cmd.OfferDescription
	offer.OfferDiscountType = cmd.OfferDiscountType
	offer.EndDate = cmd.EndDate
	offer.MarketingMessage = cmd.MarketingMessage
	offer.MaxUsesPerCustomer = cmd.MaxUsesPerCustomer
	offer.MaxUses = cmd.MaxUses
	offer.MaxUsesStrategy = cmd.MaxUsesStrategy
	offer.MinimumDaysPerUsage = cmd.MinimumDaysPerUsage
	offer.OfferItemQualifierRule = cmd.OfferItemQualifierRule
	offer.OfferItemTargetRule = cmd.OfferItemTargetRule
	offer.OrderMinTotal = cmd.OrderMinTotal
	offer.OfferPriority = cmd.OfferPriority
	offer.QualifyingItemMinTotal = cmd.QualifyingItemMinTotal
	offer.RequiresRelatedTarQual = cmd.RequiresRelatedTarQual
	offer.TargetMinTotal = cmd.TargetMinTotal
	offer.TargetSystem = cmd.TargetSystem
	offer.TotalitarianOffer = cmd.TotalitarianOffer
	offer.UseListForDiscounts = cmd.UseListForDiscounts

	err = s.offerRepo.Save(ctx, offer)
	if err != nil {
		return nil, fmt.Errorf("failed to save offer: %w", err)
	}

	return toOfferDTO(offer), nil
}

func (s *offerService) GetOfferByID(ctx context.Context, id int64) (*OfferDTO, error) {
	offer, err := s.offerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer by ID: %w", err)
	}
	if offer == nil {
		return nil, fmt.Errorf("offer with ID %d not found", id)
	}
	return toOfferDTO(offer), nil
}

func (s *offerService) UpdateOffer(ctx context.Context, cmd *UpdateOfferCommand) (*OfferDTO, error) {
	offer, err := s.offerRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer by ID for update: %w", err)
	}
	if offer == nil {
		return nil, fmt.Errorf("offer with ID %d not found for update", cmd.ID)
	}

	if cmd.Name != nil {
		offer.Name = *cmd.Name
	}
	if cmd.OfferType != nil {
		offer.OfferType = *cmd.OfferType
	}
	if cmd.OfferValue != nil {
		offer.OfferValue = *cmd.OfferValue
	}
	if cmd.AdjustmentType != nil {
		offer.AdjustmentType = *cmd.AdjustmentType
	}
	if cmd.ApplyToChildItems != nil {
		offer.SetApplyToChildItems(*cmd.ApplyToChildItems)
	}
	if cmd.ApplyToSalePrice != nil {
		offer.SetApplyToSalePrice(*cmd.ApplyToSalePrice)
	}
	if cmd.Archived != nil {
		if *cmd.Archived {
			offer.Deactivate()
		} else {
			offer.Activate()
		}
	}
	if cmd.AutomaticallyAdded != nil {
		offer.SetAutomaticallyAdded(*cmd.AutomaticallyAdded)
	}
	if cmd.CombinableWithOtherOffers != nil {
		offer.SetCombinableWithOtherOffers(*cmd.CombinableWithOtherOffers)
	}
	if cmd.OfferDescription != nil {
		offer.SetOfferDescription(*cmd.OfferDescription)
	}
	if cmd.OfferDiscountType != nil {
		offer.OfferDiscountType = *cmd.OfferDiscountType
	}
	if cmd.EndDate != nil {
		offer.SetEndDate(*cmd.EndDate)
	}
	if cmd.MarketingMessage != nil {
		offer.SetMarketingMessage(*cmd.MarketingMessage)
	}
	if cmd.MaxUsesPerCustomer != nil {
		offer.SetMaxUsesPerCustomer(*cmd.MaxUsesPerCustomer)
	}
	if cmd.MaxUses != nil {
		offer.SetMaxUses(*cmd.MaxUses)
	}
	if cmd.MaxUsesStrategy != nil {
		offer.SetMaxUsesStrategy(*cmd.MaxUsesStrategy)
	}
	if cmd.MinimumDaysPerUsage != nil {
		offer.SetMinimumDaysPerUsage(*cmd.MinimumDaysPerUsage)
	}
	if cmd.OfferItemQualifierRule != nil {
		offer.SetOfferItemQualifierRule(*cmd.OfferItemQualifierRule)
	}
	if cmd.OfferItemTargetRule != nil {
		offer.SetOfferItemTargetRule(*cmd.OfferItemTargetRule)
	}
	if cmd.OrderMinTotal != nil {
		offer.SetOrderMinTotal(*cmd.OrderMinTotal)
	}
	if cmd.OfferPriority != nil {
		offer.SetOfferPriority(*cmd.OfferPriority)
	}
	if cmd.QualifyingItemMinTotal != nil {
		offer.SetQualifyingItemMinTotal(*cmd.QualifyingItemMinTotal)
	}
	if cmd.RequiresRelatedTarQual != nil {
		offer.SetRequiresRelatedTarQual(*cmd.RequiresRelatedTarQual)
	}
	if cmd.StartDate != nil {
		offer.StartDate = *cmd.StartDate
	}
	if cmd.TargetMinTotal != nil {
		offer.SetTargetMinTotal(*cmd.TargetMinTotal)
	}
	if cmd.TargetSystem != nil {
		offer.SetTargetSystem(*cmd.TargetSystem)
	}
	if cmd.TotalitarianOffer != nil {
		offer.SetTotalitarianOffer(*cmd.TotalitarianOffer)
	}
	if cmd.UseListForDiscounts != nil {
		offer.SetUseListForDiscounts(*cmd.UseListForDiscounts)
	}

	err = s.offerRepo.Save(ctx, offer)
	if err != nil {
		return nil, fmt.Errorf("failed to update offer: %w", err)
	}

	return toOfferDTO(offer), nil
}

func (s *offerService) DeleteOffer(ctx context.Context, id int64) error {
	// Delete associated offer codes
	err := s.offerCodeRepo.DeleteByOfferID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer codes for offer %d: %w", id, err)
	}
	// Delete associated price data
	err = s.offerPriceDataRepo.DeleteByOfferID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer price data for offer %d: %w", id, err)
	}
	// Delete associated qualifying criteria xrefs
	err = s.qualCritOfferXrefRepo.DeleteByOfferID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete qualifying criteria xrefs for offer %d: %w", id, err)
	}
	// Delete associated target criteria xrefs
	err = s.tarCritOfferXrefRepo.DeleteByOfferID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete target criteria xrefs for offer %d: %w", id, err)
	}

	err = s.offerRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer: %w", err)
	}
	return nil
}

func (s *offerService) CreateOfferCode(ctx context.Context, offerID int64, cmd *CreateOfferCodeCommand) (*OfferCodeDTO, error) {
	offerCode, err := domain.NewOfferCode(offerID, cmd.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to create offer code domain entity: %w", err)
	}
	if cmd.MaxUses != nil {
		offerCode.SetMaxUses(*cmd.MaxUses)
	}
	if cmd.EmailAddress != nil {
		offerCode.SetEmailAddress(*cmd.EmailAddress)
	}
	offerCode.SetValidityPeriod(cmd.StartDate, cmd.EndDate)

	err = s.offerCodeRepo.Save(ctx, offerCode)
	if err != nil {
		return nil, fmt.Errorf("failed to save offer code: %w", err)
	}
	return toOfferCodeDTO(offerCode), nil
}

func (s *offerService) GetOfferCodeByID(ctx context.Context, id int64) (*OfferCodeDTO, error) {
	offerCode, err := s.offerCodeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer code by ID: %w", err)
	}
	if offerCode == nil {
		return nil, fmt.Errorf("offer code with ID %d not found", id)
	}
	return toOfferCodeDTO(offerCode), nil
}

func (s *offerService) UpdateOfferCode(ctx context.Context, id int64, cmd *UpdateOfferCodeCommand) (*OfferCodeDTO, error) {
	offerCode, err := s.offerCodeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer code by ID for update: %w", err)
	}
	if offerCode == nil {
		return nil, fmt.Errorf("offer code with ID %d not found for update", id)
	}

	if cmd.Code != nil {
		offerCode.Code = *cmd.Code
	}
	if cmd.MaxUses != nil {
		offerCode.SetMaxUses(*cmd.MaxUses)
	}
	if cmd.Uses != nil { // Directly updating uses is generally discouraged, but present in DTO
		offerCode.Uses = *cmd.Uses
	}
	if cmd.EmailAddress != nil {
		offerCode.SetEmailAddress(*cmd.EmailAddress)
	}
	if cmd.StartDate != nil || cmd.EndDate != nil {
		offerCode.SetValidityPeriod(cmd.StartDate, cmd.EndDate)
	}
	if cmd.Archived != nil {
		if *cmd.Archived {
			offerCode.Archived = true
		} else {
			offerCode.Archived = false
		}
	}

	err = s.offerCodeRepo.Save(ctx, offerCode)
	if err != nil {
		return nil, fmt.Errorf("failed to update offer code: %w", err)
	}
	return toOfferCodeDTO(offerCode), nil
}

func (s *offerService) DeleteOfferCode(ctx context.Context, id int64) error {
	err := s.offerCodeRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer code: %w", err)
	}
	return nil
}

func (s *offerService) CreateOfferItemCriteria(ctx context.Context, cmd *CreateOfferItemCriteriaCommand) (*OfferItemCriteriaDTO, error) {
	criteria, err := domain.NewOfferItemCriteria(cmd.Quantity, cmd.OrderItemMatchRule)
	if err != nil {
		return nil, fmt.Errorf("failed to create offer item criteria domain entity: %w", err)
	}

	err = s.offerItemCriteriaRepo.Save(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to save offer item criteria: %w", err)
	}
	return toOfferItemCriteriaDTO(criteria), nil
}

func (s *offerService) GetOfferItemCriteriaByID(ctx context.Context, id int64) (*OfferItemCriteriaDTO, error) {
	criteria, err := s.offerItemCriteriaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer item criteria by ID: %w", err)
	}
	if criteria == nil {
		return nil, fmt.Errorf("offer item criteria with ID %d not found", id)
	}
	return toOfferItemCriteriaDTO(criteria), nil
}

func (s *offerService) UpdateOfferItemCriteria(ctx context.Context, id int64, cmd *UpdateOfferItemCriteriaCommand) (*OfferItemCriteriaDTO, error) {
	criteria, err := s.offerItemCriteriaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer item criteria by ID for update: %w", err)
	}
	if criteria == nil {
		return nil, fmt.Errorf("offer item criteria with ID %d not found for update", id)
	}

	quantity := criteria.Quantity
	if cmd.Quantity != nil {
		quantity = *cmd.Quantity
	}
	orderItemMatchRule := criteria.OrderItemMatchRule
	if cmd.OrderItemMatchRule != nil {
		orderItemMatchRule = *cmd.OrderItemMatchRule
	}
	criteria.UpdateCriteria(quantity, orderItemMatchRule)

	err = s.offerItemCriteriaRepo.Save(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to update offer item criteria: %w", err)
	}
	return toOfferItemCriteriaDTO(criteria), nil
}

func (s *offerService) DeleteOfferItemCriteria(ctx context.Context, id int64) error {
	// First, check if any xrefs still refer to this criteria
	qualXrefs, err := s.qualCritOfferXrefRepo.FindByOfferItemCriteriaID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check qualifying xrefs for criteria %d: %w", id, err)
	}
	if len(qualXrefs) > 0 {
		return fmt.Errorf("offer item criteria %d is still referenced by qualifying offer xrefs", id)
	}
	tarXrefs, err := s.tarCritOfferXrefRepo.FindByOfferItemCriteriaID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check target xrefs for criteria %d: %w", id, err)
	}
	if len(tarXrefs) > 0 {
		return fmt.Errorf("offer item criteria %d is still referenced by target offer xrefs", id)
	}

	err = s.offerItemCriteriaRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer item criteria: %w", err)
	}
	return nil
}

func (s *offerService) AddQualifyingItemCriteriaToOffer(ctx context.Context, offerID, offerItemCriteriaID int64) (*QualCritOfferXrefDTO, error) {
	xref, err := domain.NewQualCritOfferXref(offerID, offerItemCriteriaID)
	if err != nil {
		return nil, fmt.Errorf("failed to create qualifying criteria xref domain entity: %w", err)
	}
	err = s.qualCritOfferXrefRepo.Save(ctx, xref)
	if err != nil {
		return nil, fmt.Errorf("failed to save qualifying criteria xref: %w", err)
	}
	return toQualCritOfferXrefDTO(xref), nil
}

func (s *offerService) RemoveQualifyingItemCriteriaFromOffer(ctx context.Context, offerID, offerItemCriteriaID int64) error {
	err := s.qualCritOfferXrefRepo.RemoveQualCritOfferXref(ctx, offerID, offerItemCriteriaID)
	if err != nil {
		return fmt.Errorf("failed to remove qualifying criteria xref: %w", err)
	}
	return nil
}

func (s *offerService) AddTargetItemCriteriaToOffer(ctx context.Context, offerID, offerItemCriteriaID int64) (*TarCritOfferXrefDTO, error) {
	xref, err := domain.NewTarCritOfferXref(offerID, offerItemCriteriaID)
	if err != nil {
		return nil, fmt.Errorf("failed to create target criteria xref domain entity: %w", err)
	}
	err = s.tarCritOfferXrefRepo.Save(ctx, xref)
	if err != nil {
		return nil, fmt.Errorf("failed to save target criteria xref: %w", err)
	}
	return toTarCritOfferXrefDTO(xref), nil
}

func (s *offerService) RemoveTargetItemCriteriaFromOffer(ctx context.Context, offerID, offerItemCriteriaID int64) error {
	err := s.tarCritOfferXrefRepo.RemoveTarCritOfferXref(ctx, offerID, offerItemCriteriaID)
	if err != nil {
		return fmt.Errorf("failed to remove target criteria xref: %w", err)
	}
	return nil
}

func (s *offerService) CreateOfferPriceData(ctx context.Context, offerID int64, cmd *CreateOfferPriceDataCommand) (*OfferPriceDataDTO, error) {
	priceData, err := domain.NewOfferPriceData(offerID, cmd.Amount, cmd.DiscountType, cmd.IdentifierType, cmd.IdentifierValue, cmd.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to create offer price data domain entity: %w", err)
	}
	priceData.SetValidityPeriod(cmd.StartDate, cmd.EndDate)

	err = s.offerPriceDataRepo.Save(ctx, priceData)
	if err != nil {
		return nil, fmt.Errorf("failed to save offer price data: %w", err)
	}
	return toOfferPriceDataDTO(priceData), nil
}

func (s *offerService) GetOfferPriceDataByID(ctx context.Context, id int64) (*OfferPriceDataDTO, error) {
	priceData, err := s.offerPriceDataRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer price data by ID: %w", err)
	}
	if priceData == nil {
		return nil, fmt.Errorf("offer price data with ID %d not found", id)
	}
	return toOfferPriceDataDTO(priceData), nil
}

func (s *offerService) UpdateOfferPriceData(ctx context.Context, id int64, cmd *UpdateOfferPriceDataCommand) (*OfferPriceDataDTO, error) {
	priceData, err := s.offerPriceDataRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer price data by ID for update: %w", err)
	}
	if priceData == nil {
		return nil, fmt.Errorf("offer price data with ID %d not found for update", id)
	}

	amount := priceData.Amount
	if cmd.Amount != nil {
		amount = *cmd.Amount
	}
	discountType := priceData.DiscountType
	if cmd.DiscountType != nil {
		discountType = *cmd.DiscountType
	}
	identifierType := priceData.IdentifierType
	if cmd.IdentifierType != nil {
		identifierType = *cmd.IdentifierType
	}
	identifierValue := priceData.IdentifierValue
	if cmd.IdentifierValue != nil {
		identifierValue = *cmd.IdentifierValue
	}
	quantity := priceData.Quantity
	if cmd.Quantity != nil {
		quantity = *cmd.Quantity
	}
	priceData.UpdateData(amount, discountType, identifierType, identifierValue, quantity)

	if cmd.StartDate != nil || cmd.EndDate != nil {
		priceData.SetValidityPeriod(cmd.StartDate, cmd.EndDate)
	}
	if cmd.Archived != nil {
		if *cmd.Archived {
			priceData.Archive()
		} else {
			priceData.Unarchive()
		}
	}

	err = s.offerPriceDataRepo.Save(ctx, priceData)
	if err != nil {
		return nil, fmt.Errorf("failed to update offer price data: %w", err)
	}
	return toOfferPriceDataDTO(priceData), nil
}

func (s *offerService) DeleteOfferPriceData(ctx context.Context, id int64) error {
	err := s.offerPriceDataRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer price data: %w", err)
	}
	return nil
}

func toOfferDTO(offer *domain.Offer) *OfferDTO {
	return &OfferDTO{
		ID:                       offer.ID,
		Name:                     offer.Name,
		OfferType:                offer.OfferType,
		OfferValue:               offer.OfferValue,
		AdjustmentType:           offer.AdjustmentType,
		ApplyToChildItems:        offer.ApplyToChildItems,
		ApplyToSalePrice:         offer.ApplyToSalePrice,
		Archived:                 offer.Archived,
		AutomaticallyAdded:       offer.AutomaticallyAdded,
		CombinableWithOtherOffers: offer.CombinableWithOtherOffers,
		OfferDescription:         offer.OfferDescription,
		OfferDiscountType:        offer.OfferDiscountType,
		EndDate:                  offer.EndDate,
		MarketingMessage:         offer.MarketingMessage,
		MaxUsesPerCustomer:       offer.MaxUsesPerCustomer,
		MaxUses:                  offer.MaxUses,
		MaxUsesStrategy:          offer.MaxUsesStrategy,
		MinimumDaysPerUsage:      offer.MinimumDaysPerUsage,
		OfferItemQualifierRule:   offer.OfferItemQualifierRule,
		OfferItemTargetRule:      offer.OfferItemTargetRule,
		OrderMinTotal:            offer.OrderMinTotal,
		OfferPriority:            offer.OfferPriority,
		QualifyingItemMinTotal:   offer.QualifyingItemMinTotal,
		RequiresRelatedTarQual:   offer.RequiresRelatedTarQual,
		StartDate:                offer.StartDate,
		TargetMinTotal:           offer.TargetMinTotal,
		TargetSystem:             offer.TargetSystem,
		TotalitarianOffer:        offer.TotalitarianOffer,
		UseListForDiscounts:      offer.UseListForDiscounts,
		CreatedAt:                offer.CreatedAt,
		UpdatedAt:                offer.UpdatedAt,
	}
}

func toOfferCodeDTO(offerCode *domain.OfferCode) *OfferCodeDTO {
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

func toOfferItemCriteriaDTO(criteria *domain.OfferItemCriteria) *OfferItemCriteriaDTO {
	return &OfferItemCriteriaDTO{
		ID:                 criteria.ID,
		Quantity:           criteria.Quantity,
		OrderItemMatchRule: criteria.OrderItemMatchRule,
		CreatedAt:          criteria.CreatedAt,
		UpdatedAt:          criteria.UpdatedAt,
	}
}

func toQualCritOfferXrefDTO(xref *domain.QualCritOfferXref) *QualCritOfferXrefDTO {
	return &QualCritOfferXrefDTO{
		ID:                  xref.ID,
		OfferID:             xref.OfferID,
		OfferItemCriteriaID: xref.OfferItemCriteriaID,
		CreatedAt:           xref.CreatedAt,
		UpdatedAt:           xref.UpdatedAt,
	}
}

func toTarCritOfferXrefDTO(xref *domain.TarCritOfferXref) *TarCritOfferXrefDTO {
	return &TarCritOfferXrefDTO{
		ID:                  xref.ID,
		OfferID:             xref.OfferID,
		OfferItemCriteriaID: xref.OfferItemCriteriaID,
		CreatedAt:           xref.CreatedAt,
		UpdatedAt:           xref.UpdatedAt,
	}
}

func toOfferPriceDataDTO(priceData *domain.OfferPriceData) *OfferPriceDataDTO {
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