package domain

import "errors"

var (
	ErrMessageNotFound    = errors.New("promotion message not found")
	ErrInvalidMessageType = errors.New("invalid message type")
	ErrMessageExpired     = errors.New("message has expired")
	ErrMaxViewsReached    = errors.New("max views reached")
)
