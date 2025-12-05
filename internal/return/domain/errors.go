package domain

import "errors"

var (
	ErrReturnNotFound         = errors.New("return request not found")
	ErrReturnAlreadyProcessed = errors.New("return request already processed")
	ErrInvalidReturnStatus    = errors.New("invalid return status")
	ErrReturnNotApproved      = errors.New("return request not approved")
)
