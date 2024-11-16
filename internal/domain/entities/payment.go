package entities

import (
	"time"
)

// Payment represents a payment made for a ticket.
type Payment struct {
	ID              string    `json:"id"`
	StripePaymentID string    `json:"stripe_payment_id"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`           // Status: 'pending', 'completed', 'failed', 'refunded'
	CreatedAt       time.Time `json:"created_at"`
}
