package entities

import (
	"time"
)

// Ticket represents a ticket for an event.
type Ticket struct {
	ID         string     `json:"id"`
	Status     string     `json:"status"` // Status: 'reserved', 'paid', 'cancelled', 'expired'
	ReservedAt time.Time  `json:"reserved_at"`
	PaidAt     time.Time `json:"paid_at"`
	Price      float64    `json:"price"`
	CreatedAt  time.Time  `json:"created_at"`
}
