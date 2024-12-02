// internal/application/service/tickets.service.go
package service

import (
	"context"
	"fmt"

	types "ticket-booking-app-backend/internal/application/types/errors"
	"ticket-booking-app-backend/internal/application/types/requests"
	"ticket-booking-app-backend/internal/domain/entities"
	"ticket-booking-app-backend/internal/domain/repository"
	domainErrors "ticket-booking-app-backend/internal/domain/types"
	"ticket-booking-app-backend/pkg/values"
)

type Tickets interface {
	ReserveTickets(ctx context.Context, input *requests.ReserveTicketsRequest) ([]*entities.Ticket, error)
	GetTicketByID(ctx context.Context, input *requests.GetTicketByIDRequest) (*entities.Ticket, error)
	GetUserTickets(ctx context.Context, input *requests.GetUserTicketsRequest) ([]*entities.Ticket, error)
	GetEventTickets(ctx context.Context, input *requests.GetEventTicketsRequest) ([]*entities.Ticket, error)
	CancelTicket(ctx context.Context, input *requests.CancelTicketRequest) error
}

type ticketsService struct {
	repo       repository.TicketsRepository
	commonRepo repository.CommonRepository
}

func NewTicketsService(repo repository.TicketsRepository, commonRepo repository.CommonRepository) *ticketsService {
	return &ticketsService{
		repo:       repo,
		commonRepo: commonRepo,
	}
}

func (s *ticketsService) ReserveTickets(ctx context.Context, input *requests.ReserveTicketsRequest) ([]*entities.Ticket, error) {
	// Validate ticket quantity
	if input.Body.Quantity > values.MaxTicketsPerPurchase {
		return nil, fmt.Errorf("cannot reserve more than %d tickets at once", values.MaxTicketsPerPurchase)
	}

	// Verify event is active
	if err := s.commonRepo.CheckIfEventIsActive(ctx, input.EventID); err != nil {
		return nil, err
	}

	// Check if there's enough capacity
	remainingCapacity, err := s.commonRepo.CheckEventAvailableCapacity(ctx, input.EventID)
	if err != nil {
		return nil, err
	}
	if remainingCapacity < input.Body.Quantity {
		return nil, domainErrors.ErrInsufficientTickets
	}

	// Create tickets now that we've verified everything
	tickets, err := s.repo.CreateTickets(ctx, input.EventID, input.UserID, input.Body.Quantity)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (s *ticketsService) GetTicketByID(ctx context.Context, input *requests.GetTicketByIDRequest) (*entities.Ticket, error) {
	ticket, err := s.repo.GetTicketWithEvent(ctx, input.TicketID)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if input.Role == values.UserRole {
		if err := s.repo.ValidateTicketOwnership(ctx, input.TicketID, input.UserID); err != nil {
			return nil, types.ErrNotAuthorized
		}
	}

	return ticket, nil
}

func (s *ticketsService) GetUserTickets(ctx context.Context, input *requests.GetUserTicketsRequest) ([]*entities.Ticket, error) {
	return s.repo.GetTicketsByUser(ctx, input.UserID)
}

func (s *ticketsService) GetEventTickets(ctx context.Context, input *requests.GetEventTicketsRequest) ([]*entities.Ticket, error) {
	// Verify permissions
	if input.Role != values.AdminRole && input.Role != values.OrganizerRole {
		return nil, types.ErrNotAuthorized
	}

	// For organizer, verify they own the event
	if input.Role == values.OrganizerRole {
		err := s.commonRepo.CheckIfEventBelongsToOrganizer(ctx, input.EventID, input.OrganizerID)
		if err != nil {
			return nil, err
		}
	}

	return s.repo.GetTicketsByEvent(ctx, input.EventID)
}

func (s *ticketsService) CancelTicket(ctx context.Context, input *requests.CancelTicketRequest) error {
	ticket, err := s.repo.GetTicketByID(ctx, input.TicketID)
	if err != nil {
		return err
	}

	// Verify permissions
	if input.Role == values.UserRole {
		if err := s.repo.ValidateTicketOwnership(ctx, input.TicketID, input.UserID); err != nil {
			return types.ErrNotAuthorized
		}
	}

	// Can only cancel reserved tickets
	if ticket.Status != values.TicketStatusReserved {
		return domainErrors.ErrInvalidTicketStatus
	}

	return s.repo.UpdateTicketStatus(ctx, input.TicketID, values.TicketStatusCancelled)
}
