package postgres

import (
	"ticket-booking-app-backend/internal/domain/entities"
	"ticket-booking-app-backend/internal/infrastructure/drivers/postgres/models"
)

// toDomainTicket maps the GORM Ticket model to the domain Ticket entity.
func toDomainTicket(ticketModel *models.Ticket) *entities.Ticket {
	return &entities.Ticket{
		ID:         ticketModel.ID.String(),
		Status:     ticketModel.Status,
		ReservedAt: ticketModel.ReservedAt,
		PaidAt:     ticketModel.PaidAt, // Converts *time.Time to time if not nil
		Price:      ticketModel.Price,
	}
}

// toGormTicket maps the domain Ticket entity to the GORM Ticket model.
func toGormTicket(ticket *entities.Ticket) *models.Ticket {
	ticketID, err := validateGormId(ticket.ID)
	if err != nil {
		return nil
	}

	gormTicket := &models.Ticket{
		ID:         ticketID,
		Status:     ticket.Status,
		ReservedAt: ticket.ReservedAt,
		Price:      ticket.Price,
	}

	// Set PaidAt if it exists
	if !ticket.PaidAt.IsZero() {
		gormTicket.PaidAt = ticket.PaidAt
	}

	return gormTicket
}
