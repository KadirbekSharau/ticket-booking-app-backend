package types

import "errors"

var (
	ErrInvalidUUID             = errors.New("invalid UUID")
	ErrInvalidIDQueryParameter = errors.New("invalid ID query parameter")
	ErrInvalidInputBody        = errors.New("invalid input body")
)