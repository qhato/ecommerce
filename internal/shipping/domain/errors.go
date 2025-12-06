package domain

import "errors"

var (
	ErrCarrierNotFound      = errors.New("carrier not found")
	ErrMethodNotFound       = errors.New("shipping method not found")
	ErrRuleNotFound         = errors.New("shipping rule not found")
	ErrInvalidBandRange     = errors.New("invalid band range")
	ErrNoShippingAvailable  = errors.New("no shipping methods available")
)
