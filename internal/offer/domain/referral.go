package domain

import (
	"errors"
	"time"
)

// ReferralProgram represents a referral program
type ReferralProgram struct {
	ID                      int64
	Name                    string
	Description             string
	StartDate               time.Time
	EndDate                 *time.Time
	IsActive                bool
	ReferrerRewardType      string // PERCENT, AMOUNT, POINTS
	ReferrerRewardValue     float64
	RefereeRewardType       string // PERCENT, AMOUNT, POINTS
	RefereeRewardValue      float64
	MinPurchaseAmount       float64
	MaxReferralsPerCustomer *int
	TotalReferrals          int64
	TotalConversions        int64
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// Referral represents a single referral
type Referral struct {
	ID                 int64
	ReferralProgramID  int64
	ReferrerCustomerID int64
	RefereeEmail       string
	RefereeCustomerID  *int64
	ReferralCode       string
	Status             string // SENT, REGISTERED, CONVERTED, EXPIRED
	ConvertedAt        *time.Time
	OrderID            *int64
	ReferrerRewarded   bool
	RefereeRewarded    bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// NewReferralProgram creates a new referral program
func NewReferralProgram(name, description string, startDate time.Time, referrerRewardType string, referrerRewardValue float64, refereeRewardType string, refereeRewardValue float64) (*ReferralProgram, error) {
	if referrerRewardValue <= 0 || refereeRewardValue <= 0 {
		return nil, errors.New("reward values must be greater than 0")
	}

	now := time.Now()
	return &ReferralProgram{
		Name:                name,
		Description:         description,
		StartDate:           startDate,
		IsActive:            true,
		ReferrerRewardType:  referrerRewardType,
		ReferrerRewardValue: referrerRewardValue,
		RefereeRewardType:   refereeRewardType,
		RefereeRewardValue:  refereeRewardValue,
		MinPurchaseAmount:   0,
		TotalReferrals:      0,
		TotalConversions:    0,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

// IsCurrentlyActive checks if the referral program is currently active
func (r *ReferralProgram) IsCurrentlyActive() bool {
	if !r.IsActive {
		return false
	}
	now := time.Now()
	if now.Before(r.StartDate) {
		return false
	}
	if r.EndDate != nil && now.After(*r.EndDate) {
		return false
	}
	return true
}

// CalculateReferrerReward calculates the referrer's reward
func (r *ReferralProgram) CalculateReferrerReward(orderTotal float64) float64 {
	if orderTotal < r.MinPurchaseAmount {
		return 0
	}
	if r.ReferrerRewardType == "PERCENT" {
		return orderTotal * (r.ReferrerRewardValue / 100)
	}
	return r.ReferrerRewardValue
}

// CalculateRefereeReward calculates the referee's reward
func (r *ReferralProgram) CalculateRefereeReward(orderTotal float64) float64 {
	if orderTotal < r.MinPurchaseAmount {
		return 0
	}
	if r.RefereeRewardType == "PERCENT" {
		return orderTotal * (r.RefereeRewardValue / 100)
	}
	return r.RefereeRewardValue
}

// NewReferral creates a new referral
func NewReferral(programID, referrerCustomerID int64, refereeEmail, referralCode string) *Referral {
	now := time.Now()
	return &Referral{
		ReferralProgramID:  programID,
		ReferrerCustomerID: referrerCustomerID,
		RefereeEmail:       refereeEmail,
		ReferralCode:       referralCode,
		Status:             "SENT",
		ReferrerRewarded:   false,
		RefereeRewarded:    false,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// MarkRegistered marks the referral as registered
func (r *Referral) MarkRegistered(refereeCustomerID int64) {
	r.RefereeCustomerID = &refereeCustomerID
	r.Status = "REGISTERED"
	r.UpdatedAt = time.Now()
}

// MarkConverted marks the referral as converted
func (r *Referral) MarkConverted(orderID int64) {
	now := time.Now()
	r.OrderID = &orderID
	r.Status = "CONVERTED"
	r.ConvertedAt = &now
	r.UpdatedAt = now
}

// MarkRewarded marks rewards as given
func (r *Referral) MarkRewarded(referrerRewarded, refereeRewarded bool) {
	r.ReferrerRewarded = referrerRewarded
	r.RefereeRewarded = refereeRewarded
	r.UpdatedAt = time.Now()
}
