package types

import "errors"


var (
	ErrUserNotFound = errors.New("user doesn't exists")
	ErrEventNotFound = errors.New("event doesn't exists")
	ErrInvalidUUID = errors.New("invalid UUID")
)

var (
	ErrTicketNotFound = errors.New("ticket not found")
	ErrTicketLimitExceeded = errors.New("ticket limit exceeded")
)