// infrastructure/repositories/postgres/tickets.postgres.go
package postgres

import (
	"context"
	"errors"
	"time"

	"ticket-booking-app-backend/internal/domain/entities"
	"ticket-booking-app-backend/internal/infrastructure/drivers/postgres/models"
	"ticket-booking-app-backend/internal/infrastructure/types"
	"ticket-booking-app-backend/pkg/values"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ticketsRepository struct {
	db *gorm.DB
}

func NewTicketsRepository(db *gorm.DB) *ticketsRepository {
	return &ticketsRepository{db: db}
}

func (r *ticketsRepository) CreateTicket(ctx context.Context, eventID, userID string, ticket *entities.Ticket) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First get event to get the price
		var event models.Event
		if err := tx.Where("id = ?", eventID).First(&event).Error; err != nil {
			return err
		}

		gormTicket := toGormTicket(ticket)

		eventUUID, err := validateGormId(eventID)
		if err != nil {
			return err
		}
		userUUID, err := validateGormId(userID)
		if err != nil {
			return err
		}

		gormTicket.EventID = eventUUID
		gormTicket.UserID = userUUID
		gormTicket.Status = values.TicketStatusReserved
		gormTicket.ReservedAt = time.Now()
		gormTicket.Price = event.Price

		if err := tx.Create(gormTicket).Error; err != nil {
			return err
		}

		// Update event's tickets_sold count
		if err := tx.Model(&event).
			Update("tickets_sold", gorm.Expr("tickets_sold + ?", 1)).
			Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *ticketsRepository) CreateTickets(ctx context.Context, eventID, userID string, count int) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First get event to get the price
		var event models.Event
		if err := tx.Where("id = ?", eventID).First(&event).Error; err != nil {
			return err
		}

		eventUUID, err := validateGormId(eventID)
		if err != nil {
			return err
		}
		userUUID, err := validateGormId(userID)
		if err != nil {
			return err
		}

		// Create tickets
		for i := 0; i < count; i++ {
			gormTicket := &models.Ticket{
				EventID:    eventUUID,
				UserID:     userUUID,
				Status:     values.TicketStatusReserved,
				ReservedAt: time.Now(),
				Price:      event.Price,
			}

			if err := tx.Create(gormTicket).Error; err != nil {
				return err
			}

			tickets = append(tickets, toDomainTicket(gormTicket))
		}

		// Update event's tickets_sold count
		if err := tx.Model(&event).
			Update("tickets_sold", gorm.Expr("tickets_sold + ?", count)).
			Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *ticketsRepository) GetTicketByID(ctx context.Context, ticketID string) (*entities.Ticket, error) {
	var ticket models.Ticket
	err := r.db.WithContext(ctx).
		Where("id = ?", ticketID).
		First(&ticket).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, types.ErrTicketNotFound
	}
	if err != nil {
		return nil, err
	}

	return toDomainTicket(&ticket), nil
}

func (r *ticketsRepository) GetTicketsByEvent(ctx context.Context, eventID, status string) ([]*entities.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND status = ?", eventID, status).
		Find(&tickets).Error

	if err != nil {
		return nil, err
	}

	return toDomainTickets(tickets), nil
}

func (r *ticketsRepository) GetTicketsByUser(ctx context.Context, userID, status string) ([]*entities.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, status).
		Find(&tickets).Error

	if err != nil {
		return nil, err
	}

	return toDomainTickets(tickets), nil
}

func (r *ticketsRepository) GetTicketWithEvent(ctx context.Context, ticketID string) (*entities.Ticket, error) {
	var ticket models.Ticket
	err := r.db.WithContext(ctx).
		Preload("Event").
		Where("id = ?", ticketID).
		First(&ticket).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, types.ErrTicketNotFound
	}
	if err != nil {
		return nil, err
	}

	return toDomainTicket(&ticket), nil
}

func (r *ticketsRepository) UpdateTicketStatus(ctx context.Context, ticketID string, status string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Where("id = ?", ticketID).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return types.ErrTicketNotFound
	}
	return nil
}

func (r *ticketsRepository) UpdateTicketPayment(ctx context.Context, ticketID string, paidAt time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Ticket{}).
			Where("id = ?", ticketID).
			Updates(map[string]interface{}{
				"status":  values.TicketStatusPaid,
				"paid_at": paidAt,
			})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return types.ErrTicketNotFound
		}
		return nil
	})
}

func (r *ticketsRepository) UpdateExpiredTickets(ctx context.Context) error {
	// Get reservation expiration time (e.g., 15 minutes)
	expirationTime := time.Now().Add(-15 * time.Minute)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update expired reserved tickets
		result := tx.Model(&models.Ticket{}).
			Where("status = ? AND reserved_at <= ?", values.TicketStatusReserved, expirationTime).
			Update("status", values.TicketStatusExpired)

		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

func (r *ticketsRepository) CancelEventTickets(ctx context.Context, eventID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Ticket{}).
			Where("event_id = ? AND status = ?", eventID, values.TicketStatusReserved).
			Update("status", values.TicketStatusCancelled)

		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

func (r *ticketsRepository) CheckTicketAvailability(ctx context.Context, eventID string, count int) (bool, error) {
	var event models.Event
	err := r.db.WithContext(ctx).
		Select("capacity, tickets_sold").
		Where("id = ?", eventID).
		First(&event).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, types.ErrEventNotFound
	}
	if err != nil {
		return false, err
	}

	return (event.Capacity - event.TicketsSold) >= count, nil
}

func (r *ticketsRepository) ValidateTicketOwnership(ctx context.Context, ticketID, userID string) error {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Ticket{}).
		Where("id = ? AND user_id = ?", ticketID, userID).
		Count(&count).Error

	if err != nil {
		return err
	}
	if count == 0 {
		return types.ErrTicketNotFound
	}
	return nil
}

// Helper functions for mapping between domain and GORM models
func toDomainTickets(tickets []models.Ticket) []*entities.Ticket {
	result := make([]*entities.Ticket, len(tickets))
	for i, ticket := range tickets {
		result[i] = toDomainTicket(&ticket)
	}
	return result
}

func toDomainTicket(ticketModel *models.Ticket) *entities.Ticket {
	return &entities.Ticket{
		ID:         ticketModel.ID.String(),
		Status:     ticketModel.Status,
		ReservedAt: ticketModel.ReservedAt,
		PaidAt:     ticketModel.PaidAt,
		Price:      ticketModel.Price,
		CreatedAt:  ticketModel.CreatedAt,
	}
}

func toGormTicket(ticket *entities.Ticket) *models.Ticket {
	var ticketID uuid.UUID
	var err error
	if ticket.ID != "" {
		ticketID, err = validateGormId(ticket.ID)
		if err != nil {
			return nil
		}
	}

	return &models.Ticket{
		ID:         ticketID,
		Status:     ticket.Status,
		ReservedAt: ticket.ReservedAt,
		PaidAt:     ticket.PaidAt,
		Price:      ticket.Price,
	}
}
