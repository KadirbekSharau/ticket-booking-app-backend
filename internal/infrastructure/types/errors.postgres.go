package types

import "errors"


var (
	ErrUserNotFound = errors.New("user doesn't exists")
	ErrEventNotFound = errors.New("event doesn't exists")
	ErrInvalidUUID = errors.New("invalid UUID")
)
