package application

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/qhato/ecommerce/internal/offer/domain"
	"github.com/shopspring/decimal"
)

// OfferProcessorService processes offers for orders
type OfferProcessorService interface {
	// ProcessOrderOffers processes all applicable offers for an order
	ProcessOrderOffers(ctx context.Context, request *ProcessOffersRequest) (*ProcessOffersResponse, error)

	// ApplyOfferCode applies a specific offer code to an order
	ApplyOfferCode(ctx context.Context, request *ApplyOfferCodeRequest) (*ApplyOfferCodeResponse, error)

	// RemoveOfferFromOrder removes an offer from an order
	RemoveOfferFromOrder(ctx context.Context, orderID, offerID int64) error
}

// ProcessOffersRequest contains the data needed to process offers for an order
type ProcessOffersRequest struct {
	OrderID        int64
	OrderSubtotal  decimal.Decimal
	OrderTotal     decimal.Decimal
	CustomerID     *string
	Items          []OrderItemData
	AppliedOfferIDs []int64 // Currently applied offer IDs (for re-calculation)
}

// OrderItemData represents an order item for offer processing
type OrderItemData struct {
	ItemID      string
	SKUID       string
	CategoryID  *string
	Price       decimal.Decimal
	SalePrice   *decimal.Decimal
	Quantity    int
	Subtotal    decimal.Decimal
	ProductID   *string
}

// ProcessOffersResponse contains the results of offer processing
type ProcessOffersResponse struct {
	OrderID              int64
	OriginalSubtotal     decimal.Decimal
	TotalDiscount        decimal.Decimal
	AdjustedSubtotal     decimal.Decimal
	AppliedOffers        []*AppliedOfferDTO
	OrderAdjustments     []*OrderAdjustmentDTO
	ItemAdjustments      []*OrderItemAdjustmentDTO
}

// ApplyOfferCodeRequest contains the data needed to apply an offer code
type ApplyOfferCodeRequest struct {
	OrderID       int64
	OfferCode     string
	OrderSubtotal decimal.Decimal
	OrderTotal    decimal.Decimal
	CustomerID    *string
	Items         []OrderItemData
}

// ApplyOfferCodeResponse contains the results of applying an offer code
type ApplyOfferCodeResponse struct {
	Success         bool
	Message         string
	Offer           *OfferDTO
	DiscountAmount  decimal.Decimal
}

// AppliedOfferDTO represents an offer that was successfully applied
type AppliedOfferDTO struct {
	OfferID        int64
	OfferName      string
	DiscountAmount decimal.Decimal
	Priority       int
}

// OrderAdjustmentDTO represents an order-level adjustment
type OrderAdjustmentDTO struct {
	OfferID          int64
	OfferName        string
	AdjustmentValue  decimal.Decimal
	AdjustmentReason string
}

// OrderItemAdjustmentDTO represents an item-level adjustment
type OrderItemAdjustmentDTO struct {
	ItemID          string
	OfferID         int64
	OfferName       string
	AdjustmentValue decimal.Decimal
	Quantity        int
}

type offerProcessorService struct {
	offerRepo            domain.OfferRepository
	offerCodeRepo        domain.OfferCodeRepository
	adjustmentRepo       domain.OrderAdjustmentRepository
	processor            *domain.OfferProcessor
}

// NewOfferProcessorService creates a new OfferProcessorService
func NewOfferProcessorService(
	offerRepo domain.OfferRepository,
	offerCodeRepo domain.OfferCodeRepository,
	adjustmentRepo domain.OrderAdjustmentRepository,
	processor *domain.OfferProcessor,
) OfferProcessorService {
	return &offerProcessorService{
		offerRepo:      offerRepo,
		offerCodeRepo:  offerCodeRepo,
		adjustmentRepo: adjustmentRepo,
		processor:      processor,
	}
}

func (s *offerProcessorService) ProcessOrderOffers(ctx context.Context, request *ProcessOffersRequest) (*ProcessOffersResponse, error) {
	// Build offer context
	offerCtx := s.buildOfferContext(request)

	// Get all active offers
	activeOffers, err := s.offerRepo.FindActiveOffers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find active offers: %w", err)
	}

	offerCtx.AvailableOffers = activeOffers

	// Find qualifying offers and calculate discounts
	candidates := make([]*domain.CandidateOffer, 0)
	for _, offer := range activeOffers {
		// Check if offer should be automatically added
		if !offer.AutomaticallyAdded {
			continue
		}

		// Qualify the offer
		qualification, err := s.processor.QualifyOffer(offer, offerCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to qualify offer %d: %w", offer.ID, err)
		}

		if !qualification.Qualifies {
			continue
		}

		// Calculate discount
		discountAmount, targetItemIDs, err := s.processor.CalculateDiscount(offer, offerCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate discount for offer %d: %w", offer.ID, err)
		}

		if discountAmount.IsZero() {
			continue
		}

		candidates = append(candidates, &domain.CandidateOffer{
			Offer:           offer,
			Priority:        offer.OfferPriority,
			DiscountAmount:  discountAmount,
			TargetItems:     targetItemIDs,
		})
	}

	// Select best combination of offers
	selectedOffers := s.processor.SelectBestOffers(candidates)

	// Sort by priority
	sort.Slice(selectedOffers, func(i, j int) bool {
		return selectedOffers[i].Priority < selectedOffers[j].Priority
	})

	// Build response
	response := &ProcessOffersResponse{
		OrderID:          request.OrderID,
		OriginalSubtotal: request.OrderSubtotal,
		TotalDiscount:    decimal.Zero,
		AppliedOffers:    make([]*AppliedOfferDTO, 0),
		OrderAdjustments: make([]*OrderAdjustmentDTO, 0),
		ItemAdjustments:  make([]*OrderItemAdjustmentDTO, 0),
	}

	// Apply selected offers
	for _, candidate := range selectedOffers {
		offer := candidate.Offer
		discountAmount := candidate.DiscountAmount

		response.TotalDiscount = response.TotalDiscount.Add(discountAmount)

		response.AppliedOffers = append(response.AppliedOffers, &AppliedOfferDTO{
			OfferID:        offer.ID,
			OfferName:      offer.Name,
			DiscountAmount: discountAmount,
			Priority:       offer.OfferPriority,
		})

		// Create adjustments based on adjustment type
		switch offer.AdjustmentType {
		case domain.OfferAdjustmentTypeOrder:
			response.OrderAdjustments = append(response.OrderAdjustments, &OrderAdjustmentDTO{
				OfferID:          offer.ID,
				OfferName:        offer.Name,
				AdjustmentValue:  discountAmount,
				AdjustmentReason: "OFFER_DISCOUNT",
			})

		case domain.OfferAdjustmentTypeOrderItem:
			// Distribute discount across target items
			for _, itemID := range candidate.TargetItems {
				response.ItemAdjustments = append(response.ItemAdjustments, &OrderItemAdjustmentDTO{
					ItemID:          itemID,
					OfferID:         offer.ID,
					OfferName:       offer.Name,
					AdjustmentValue: discountAmount, // Simplified: full discount per item
					Quantity:        1,
				})
			}
		}
	}

	response.AdjustedSubtotal = response.OriginalSubtotal.Sub(response.TotalDiscount)

	return response, nil
}

func (s *offerProcessorService) ApplyOfferCode(ctx context.Context, request *ApplyOfferCodeRequest) (*ApplyOfferCodeResponse, error) {
	// Find the offer code
	offerCode, err := s.offerCodeRepo.FindByCode(ctx, request.OfferCode)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer code: %w", err)
	}

	if offerCode == nil {
		return &ApplyOfferCodeResponse{
			Success: false,
			Message: "Offer code not found",
		}, nil
	}

	// Check if offer code is active
	if !offerCode.IsActive() {
		return &ApplyOfferCodeResponse{
			Success: false,
			Message: "Offer code is not currently active",
		}, nil
	}

	// Get the associated offer
	offer, err := s.offerRepo.FindByID(ctx, offerCode.OfferID)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer: %w", err)
	}

	if offer == nil {
		return &ApplyOfferCodeResponse{
			Success: false,
			Message: "Associated offer not found",
		}, nil
	}

	// Build offer context
	offerCtx := &domain.OfferContext{
		OrderTotal:         request.OrderTotal,
		OrderSubtotal:      request.OrderSubtotal,
		CustomerID:         request.CustomerID,
		Items:              s.convertToOfferItems(request.Items),
		AppliedOffers:      make([]*domain.OfferAdjustment, 0),
		AvailableOffers:    []*domain.Offer{offer},
		CustomerUsageCount: make(map[int64]int),
	}

	// Qualify the offer
	qualification, err := s.processor.QualifyOffer(offer, offerCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to qualify offer: %w", err)
	}

	if !qualification.Qualifies {
		return &ApplyOfferCodeResponse{
			Success: false,
			Message: qualification.Reason,
		}, nil
	}

	// Calculate discount
	discountAmount, _, err := s.processor.CalculateDiscount(offer, offerCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate discount: %w", err)
	}

	if discountAmount.IsZero() {
		return &ApplyOfferCodeResponse{
			Success: false,
			Message: "Offer does not provide any discount for this order",
		}, nil
	}

	// Increment offer code usage
	offerCode.IncrementUses()
	err = s.offerCodeRepo.Save(ctx, offerCode)
	if err != nil {
		return nil, fmt.Errorf("failed to update offer code usage: %w", err)
	}

	return &ApplyOfferCodeResponse{
		Success:        true,
		Message:        "Offer code applied successfully",
		Offer:          ToOfferDTO(offer),
		DiscountAmount: discountAmount,
	}, nil
}

func (s *offerProcessorService) RemoveOfferFromOrder(ctx context.Context, orderID, offerID int64) error {
	// Find all adjustments for this order and offer
	orderAdjustments, err := s.adjustmentRepo.FindByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("failed to find order adjustments: %w", err)
	}

	// Filter and delete adjustments for this specific offer
	for _, adj := range orderAdjustments {
		if adj.OfferID == offerID {
			// In a real implementation, you'd have a Delete method
			// For now, this is a placeholder
		}
	}

	return nil
}

func (s *offerProcessorService) buildOfferContext(request *ProcessOffersRequest) *domain.OfferContext {
	return &domain.OfferContext{
		OrderTotal:         request.OrderTotal,
		OrderSubtotal:      request.OrderSubtotal,
		CustomerID:         request.CustomerID,
		Items:              s.convertToOfferItems(request.Items),
		AppliedOffers:      make([]*domain.OfferAdjustment, 0),
		AvailableOffers:    make([]*domain.Offer, 0),
		CustomerUsageCount: make(map[int64]int),
	}
}

func (s *offerProcessorService) convertToOfferItems(items []OrderItemData) []domain.OfferItem {
	offerItems := make([]domain.OfferItem, len(items))
	for i, item := range items {
		offerItems[i] = domain.OfferItem{
			ItemID:      item.ItemID,
			SKUID:       item.SKUID,
			CategoryID:  item.CategoryID,
			Price:       item.Price,
			SalePrice:   item.SalePrice,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
			ProductID:   item.ProductID,
			Adjustments: make([]domain.OfferAdjustment, 0),
		}
	}
	return offerItems
}

// PersistAdjustments persists the calculated adjustments to the database
func (s *offerProcessorService) PersistAdjustments(ctx context.Context, orderID int64, response *ProcessOffersResponse) error {
	now := time.Now()

	// Delete existing adjustments for this order
	err := s.adjustmentRepo.DeleteByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("failed to delete existing adjustments: %w", err)
	}

	// Create order-level adjustments
	for _, adjDTO := range response.OrderAdjustments {
		adj := &domain.OrderAdjustment{
			OrderID:          orderID,
			OfferID:          adjDTO.OfferID,
			OfferName:        adjDTO.OfferName,
			AdjustmentValue:  adjDTO.AdjustmentValue,
			AdjustmentReason: adjDTO.AdjustmentReason,
			AppliedDate:      now,
			CreatedAt:        now,
		}

		err := s.adjustmentRepo.CreateOrderAdjustment(adj)
		if err != nil {
			return fmt.Errorf("failed to create order adjustment: %w", err)
		}
	}

	// Create item-level adjustments
	for _, adjDTO := range response.ItemAdjustments {
		// In a real system, you'd need to map itemID to order_item_id
		// For now, this is a simplified version
		adj := &domain.OrderItemAdjustment{
			OrderItemID:     0, // TODO: Map itemID to order_item_id
			OfferID:         adjDTO.OfferID,
			OfferName:       adjDTO.OfferName,
			AdjustmentValue: adjDTO.AdjustmentValue,
			Quantity:        adjDTO.Quantity,
			AppliedDate:     now,
			CreatedAt:       now,
		}

		err := s.adjustmentRepo.CreateOrderItemAdjustment(adj)
		if err != nil {
			return fmt.Errorf("failed to create order item adjustment: %w", err)
		}
	}

	return nil
}
