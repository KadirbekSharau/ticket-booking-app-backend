package repository

import (
	"context"

	"ticket-booking-app-backend/internal/domain/entities"
)

type UsersRepository interface {
	Create(ctx context.Context, organizationId string, user *entities.User) error
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	// Update(ctx context.Context, user domain.User) error
}
