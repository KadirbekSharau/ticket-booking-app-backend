// internal/application/service/events.service.go
package service

import (
	"context"
	"time"

	types "ticket-booking-app-backend/internal/application/types/errors"
	"ticket-booking-app-backend/internal/application/types/requests"
	"ticket-booking-app-backend/internal/domain/entities"
	"ticket-booking-app-backend/internal/domain/repository"
	domainErrors "ticket-booking-app-backend/internal/domain/types"
	"ticket-booking-app-backend/pkg/values"
)

type Events interface {
	GetEvents(ctx context.Context, input *requests.GetEventsRequest) ([]*entities.Event, error)
	GetEventsByOrganizer(ctx context.Context, input *requests.GetEventsByOrganizerRequest) ([]*entities.Event, error)
	GetEventByID(ctx context.Context, input *requests.GetEventByIDRequest) (*entities.Event, error)
	CreateEvent(ctx context.Context, input *requests.CreateEventRequest) error
	UpdateEvent(ctx context.Context, input *requests.UpdateEventRequest) (*entities.Event, error)
	CancelEvent(ctx context.Context, input *requests.CancelEventRequest) error
	DeleteEvent(ctx context.Context, input *requests.DeleteEventRequest) error
}

type eventsService struct {
	repo       repository.EventsRepository
	commonRepo repository.CommonRepository
}

func NewEventsService(repo repository.EventsRepository, commonRepo repository.CommonRepository) *eventsService {
	return &eventsService{
		repo:       repo,
		commonRepo: commonRepo,
	}
}

func (s *eventsService) GetEvents(ctx context.Context, input *requests.GetEventsRequest) ([]*entities.Event, error) {
	// For regular users, only return active events
	if input.Role == values.UserRole {
		input.Status = values.EventStatusActive
	}

	return s.repo.GetEvents(ctx, input.Status)
}

func (s *eventsService) GetEventsByOrganizer(ctx context.Context, input *requests.GetEventsByOrganizerRequest) ([]*entities.Event, error) {
	// Verify permissions
	if input.Role != values.AdminRole && input.Role != values.OrganizerRole {
		return nil, types.ErrNotAuthorized
	}

	// For organizer, verify they exist
	if input.Role == values.OrganizerRole {
		if err := s.commonRepo.CheckIfUserExistsByIdAndRole(ctx, input.OrganizerID, values.OrganizerRole); err != nil {
			return nil, err
		}
	}

	return s.repo.GetEventsByOrganizer(ctx, input.Status, input.OrganizerID)
}

func (s *eventsService) GetEventByID(ctx context.Context, input *requests.GetEventByIDRequest) (*entities.Event, error) {
	event, err := s.repo.GetEventByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// Regular users can only view active events
	if input.Role == values.UserRole && event.Status != values.EventStatusActive {
		return nil, types.ErrNotAuthorized
	}

	// Organizers can only view their own events
	if input.Role == values.OrganizerRole {
		return nil, types.ErrNotAuthorized
	}

	return event, nil
}

func (s *eventsService) CreateEvent(ctx context.Context, input *requests.CreateEventRequest) error {
	// Verify permissions
	if input.Role != values.AdminRole && input.Role != values.OrganizerRole {
		return types.ErrNotAuthorized
	}

	// Verify organizer exists
	if input.Role == values.OrganizerRole {
		if err := s.commonRepo.CheckIfUserExistsByIdAndRole(ctx, input.OrganizerID, values.OrganizerRole); err != nil {
			return err
		}
	}

	// Validate event date
	if input.Body.Date.Before(time.Now()) {
		return domainErrors.ErrEventDateInvalid
	}

	event := &entities.Event{
		Title:       input.Body.Title,
		Description: input.Body.Description,
		Location:    input.Body.Location,
		Date:        input.Body.Date,
		Capacity:    input.Body.Capacity,
		Price:       input.Body.Price,
		Status:      values.EventStatusActive,
	}

	return s.repo.CreateEvent(ctx, event, input.OrganizerID)
}

func (s *eventsService) UpdateEvent(ctx context.Context, input *requests.UpdateEventRequest) (*entities.Event, error) {
	// Verify permissions
	if input.Role != values.AdminRole && input.Role != values.OrganizerRole {
		return nil, types.ErrNotAuthorized
	}

	// Get existing event
	existingEvent, err := s.repo.GetEventByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// Can't update finished or cancelled events
	if existingEvent.Status == values.EventStatusFinished || existingEvent.Status == values.EventStatusCancelled {
		return nil, domainErrors.ErrEventAlreadyFinished
	}

	// Validate event date
	if input.Body.Date.Before(time.Now()) {
		return nil, domainErrors.ErrEventDateInvalid
	}

	event := &entities.Event{
		ID:          input.ID,
		Title:       input.Body.Title,
		Description: input.Body.Description,
		Location:    input.Body.Location,
		Date:        input.Body.Date,
		Capacity:    input.Body.Capacity,
		Price:       input.Body.Price,
		Status:      existingEvent.Status,
	}

	err = s.repo.UpdateEvent(ctx, input.OrganizerID, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *eventsService) DeleteEvent(ctx context.Context, input *requests.DeleteEventRequest) error {
	// Verify permissions
	if input.Role != values.AdminRole && input.Role != values.OrganizerRole {
		return types.ErrNotAuthorized
	}

	// Get existing event
	existingEvent, err := s.repo.GetEventByID(ctx, input.ID)
	if err != nil {
		return err
	}

	// Can't delete finished events
	if existingEvent.Status == values.EventStatusFinished {
		return domainErrors.ErrEventAlreadyFinished
	}

	return s.repo.DeleteEvent(ctx, input.ID, input.OrganizerID)
}

func (s *eventsService) CancelEvent(ctx context.Context, input *requests.CancelEventRequest) error {
	// Verify permissions
	if input.Role != values.AdminRole && input.Role != values.OrganizerRole {
		return types.ErrNotAuthorized
	}

	// Get existing event
	existingEvent, err := s.repo.GetEventByID(ctx, input.ID)
	if err != nil {
		return err
	}

	// Can't cancel finished events
	if existingEvent.Status == values.EventStatusFinished {
		return domainErrors.ErrEventAlreadyFinished
	}

	return s.repo.UpdateEventStatus(ctx, input.ID, input.OrganizerID, values.EventStatusCancelled)
}
