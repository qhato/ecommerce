package domain

import (
	"time"
)

// OfferRule represents a rule that can be applied to an offer (e.g., conditions for eligibility)
type OfferRule struct {
	ID        int64
	MatchRule string // From blc_offer_rule.match_rule (text)
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewOfferRule creates a new OfferRule
func NewOfferRule(matchRule string) (*OfferRule, error) {
	if matchRule == "" {
		return nil, NewDomainError("MatchRule cannot be empty for OfferRule")
	}

	now := time.Now()
	return &OfferRule{
		MatchRule: matchRule,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateMatchRule updates the match rule string
func (or *OfferRule) UpdateMatchRule(matchRule string) {
	or.MatchRule = matchRule
	or.UpdatedAt = time.Now()
}
