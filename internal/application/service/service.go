package service

import (
	"ticket-booking-app-backend/internal/domain/repository"
	"ticket-booking-app-backend/internal/helpers"
)

type Services struct {
	Users
	Events
}

func NewServices(repos *repository.Repository, jwt helpers.Jwt) *Services {
	return &Services{
		Users:  NewUsersService(repos.Users, repos.Common, jwt),
		Events: NewEventsService(repos.Events, repos.Common),
	}
}