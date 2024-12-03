// domain/repository/tickets.repository.go
package repository

import (
	"context"
	"ticket-booking-app-backend/internal/domain/entities"
	"time"
)

type TicketsRepository interface {
	// Create operations
	CreateTicket(ctx context.Context, eventID, userID string, ticket *entities.Ticket) error
	CreateTickets(ctx context.Context, eventID, userID string, count int) ([]*entities.Ticket, error)

	// Read operations
	GetTicketByID(ctx context.Context, ticketID string) (*entities.Ticket, error)
	GetTicketsByEvent(ctx context.Context, eventID, status string) ([]*entities.Ticket, error)
	GetTicketsByUser(ctx context.Context, userID, status string) ([]*entities.Ticket, error)
	GetTicketWithEvent(ctx context.Context, ticketID string) (*entities.Ticket, error)

	// Update operations
	UpdateTicketStatus(ctx context.Context, ticketID string, status string) error
	UpdateTicketPayment(ctx context.Context, ticketID string, paidAt time.Time) error

	// Batch operations
	UpdateExpiredTickets(ctx context.Context) error
	CancelEventTickets(ctx context.Context, eventID string) error

	// Validation operations
	CheckTicketAvailability(ctx context.Context, eventID string, count int) (bool, error)
	ValidateTicketOwnership(ctx context.Context, ticketID, userID string) error
}
