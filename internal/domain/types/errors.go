package errors

import "errors"

var (
	ErrUserNotFound          = errors.New("user doesn't exists")
	ErrUserAlreadyExists     = errors.New("user with such email already exists")
	ErrUserPasswordIncorrect = errors.New("password incorrect")
)

var (
	ErrAdminNotFound           = errors.New("admin doesn't exists")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrOrganizerNotFound       = errors.New("organizer doesn't exists")
)

var (
	ErrEventNotFound           = errors.New("event not found")
	ErrEventAlreadyFinished    = errors.New("event already finished")
	ErrEventAlreadyCancelled   = errors.New("event already cancelled")
	ErrEventCapacityExceeded   = errors.New("event capacity exceeded")
	ErrEventDateInvalid        = errors.New("event date must be in the future")
	ErrUnauthorizedEventAccess = errors.New("unauthorized access to event")
	ErrInvalidEventStatus      = errors.New("invalid event status")
)

var (
	ErrEventNotActive      = errors.New("event is not active")
	ErrInsufficientTickets = errors.New("insufficient tickets")
	ErrInvalidTicketStatus = errors.New("invalid ticket status")
	ErrTicketNotFound      = errors.New("ticket not found")
)
