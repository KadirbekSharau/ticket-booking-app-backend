package repository

import (
	"ticket-booking-app-backend/internal/infrastructure/repositories/postgres"

	"gorm.io/gorm"
)

type Repository struct {
	Common  CommonRepository
	Users   UsersRepository
	Events  EventsRepository
	Tickets TicketsRepository
}

func NewRepositories(db *gorm.DB) *Repository {
	return &Repository{
		Common:  postgres.NewCommonRepository(db),
		Users:   postgres.NewUsersRepository(db),
		Events:  postgres.NewEventsRepository(db),
		Tickets: postgres.NewTicketsRepository(db),
	}
}
