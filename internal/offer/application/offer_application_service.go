package application

import (
	"context"
	"fmt"
	"sort"

	"github.com/qhato/ecommerce/internal/offer/domain"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/rules"
	"github.com/shopspring/decimal"
)

// OfferApplicationService handles real-time offer application to orders/carts
type OfferApplicationService struct {
	offerRepo  domain.OfferRepository
	processor  *domain.OfferProcessor
	ruleEngine *rules.RuleEngine
	log        logger.Logger
}

// NewOfferApplicationService creates a new OfferApplicationService
func NewOfferApplicationService(
	offerRepo domain.OfferRepository,
	ruleEngine *rules.RuleEngine,
	log logger.Logger,
) *OfferApplicationService {
	// Create rule evaluator adapter
	ruleEvaluator := &RuleEvaluatorAdapter{engine: ruleEngine}
	processor := domain.NewOfferProcessor(ruleEvaluator)

	return &OfferApplicationService{
		offerRepo:  offerRepo,
		processor:  processor,
		ruleEngine: ruleEngine,
		log:        log,
	}
}

// RuleEvaluatorAdapter adapts the rules.RuleEngine to domain.RuleEvaluator
type RuleEvaluatorAdapter struct {
	engine *rules.RuleEngine
}

// Evaluate evaluates a rule expression
func (r *RuleEvaluatorAdapter) Evaluate(ruleExpression string, context map[string]interface{}) (bool, error) {
	if ruleExpression == "" {
		return true, nil
	}

	// Create and compile rule on the fly
	rule, err := rules.NewRule("dynamic", ruleExpression, "Dynamic rule")
	if err != nil {
		return false, err
	}

	return rule.Evaluate(context)
}

// ApplyOffersToOrder applies available offers to an order
func (s *OfferApplicationService) ApplyOffersToOrder(ctx context.Context, orderCtx *domain.OfferContext) ([]*domain.OfferAdjustment, error) {
	s.log.Info(fmt.Sprintf("Applying offers to order with subtotal: %s", orderCtx.OrderSubtotal.String()))

	// Get all active offers
	offers, err := s.offerRepo.FindActiveOffers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active offers: %w", err)
	}

	s.log.Debug(fmt.Sprintf("Found %d active offers", len(offers)))

	// Qualify offers
	candidates := make([]*domain.CandidateOffer, 0)

	for _, offer := range offers {
		qualification, err := s.processor.QualifyOffer(offer, orderCtx)
		if err != nil {
			s.log.Error(fmt.Sprintf("Error qualifying offer %d: %v", offer.ID, err))
			continue
		}

		if !qualification.Qualifies {
			s.log.Debug(fmt.Sprintf("Offer %d (%s) did not qualify: %s", offer.ID, offer.Name, qualification.Reason))
			continue
		}

		// Calculate discount
		discountAmount, targetItems, err := s.processor.CalculateDiscount(offer, orderCtx)
		if err != nil {
			s.log.Error(fmt.Sprintf("Error calculating discount for offer %d: %v", offer.ID, err))
			continue
		}

		if discountAmount.GreaterThan(decimal.Zero) {
			candidate := &domain.CandidateOffer{
				Offer:          offer,
				Priority:       offer.OfferPriority,
				DiscountAmount: discountAmount,
				TargetItems:    targetItems,
			}
			candidates = append(candidates, candidate)

			s.log.Debug(fmt.Sprintf("Offer %d (%s) qualifies with discount: %s", offer.ID, offer.Name, discountAmount.String()))
		}
	}

	// Select best offers to apply
	selectedOffers := s.selectBestOffers(candidates)

	// Apply selected offers
	adjustments := make([]*domain.OfferAdjustment, 0, len(selectedOffers))
	for _, candidate := range selectedOffers {
		adjustment := s.processor.ApplyOffer(candidate.Offer, candidate.DiscountAmount)
		adjustments = append(adjustments, adjustment)

		s.log.Info(fmt.Sprintf("Applied offer %d (%s): %s discount", candidate.Offer.ID, candidate.Offer.Name, candidate.DiscountAmount.String()))
	}

	return adjustments, nil
}

// selectBestOffers selects the best combination of offers
func (s *OfferApplicationService) selectBestOffers(candidates []*domain.CandidateOffer) []*domain.CandidateOffer {
	if len(candidates) == 0 {
		return candidates
	}

	// Sort by priority (ascending) then by discount amount (descending)
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Priority != candidates[j].Priority {
			return candidates[i].Priority < candidates[j].Priority
		}
		return candidates[i].DiscountAmount.GreaterThan(candidates[j].DiscountAmount)
	})

	selected := make([]*domain.CandidateOffer, 0)
	totalitarianOfferApplied := false

	for _, candidate := range candidates {
		// If a totalitarian offer has been applied, no more offers can be added
		if totalitarianOfferApplied {
			break
		}

		// Check if this is a totalitarian offer
		if candidate.Offer.TotalitarianOffer {
			// Totalitarian offer takes precedence, clear previous selections
			selected = []*domain.CandidateOffer{candidate}
			totalitarianOfferApplied = true
			continue
		}

		// Check if offer can be combined with already selected offers
		canCombine := true
		if !candidate.Offer.CombinableWithOtherOffers && len(selected) > 0 {
			canCombine = false
		}

		for _, selectedOffer := range selected {
			if !selectedOffer.Offer.CombinableWithOtherOffers {
				canCombine = false
				break
			}
		}

		if canCombine {
			selected = append(selected, candidate)
		}
	}

	return selected
}

// ValidateOfferCode validates and returns an offer by code
func (s *OfferApplicationService) ValidateOfferCode(ctx context.Context, code string, orderCtx *domain.OfferContext) (*domain.Offer, error) {
	// This would fetch the offer by code from repository
	// For now, return error as it needs OfferCode repository method
	return nil, fmt.Errorf("offer code validation not yet implemented - needs FindByCode method in repository")
}
