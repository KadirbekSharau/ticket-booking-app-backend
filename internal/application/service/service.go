package service

import (
	"ticket-booking-app-backend/internal/domain/repository"
	"ticket-booking-app-backend/internal/helpers"
	"ticket-booking-app-backend/internal/infrastructure/jobs"
)

type Services struct {
	Users
	Events
	Tickets
	EventUpdater *jobs.EventStatusUpdater
}

func NewServices(repos *repository.Repository, jwt helpers.Jwt) *Services {
	return &Services{
		Users:        NewUsersService(repos.Users, repos.Common, jwt),
		Events:       NewEventsService(repos.Events, repos.Common),
		Tickets:      NewTicketsService(repos.Tickets, repos.Common),
		EventUpdater: jobs.NewEventStatusUpdater(repos.Events),
	}
}
