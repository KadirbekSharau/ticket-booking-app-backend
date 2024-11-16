package postgres

import (
	"context"
	"fmt"

	"ticket-booking-app-backend/internal/infrastructure/drivers/postgres/models"
	"ticket-booking-app-backend/internal/infrastructure/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type commonRepository struct {
	db *gorm.DB
}

func NewCommonRepository(db *gorm.DB) *commonRepository {
	return &commonRepository{db: db}
}

func (r *commonRepository) CheckIfUserExistsByEmail(ctx context.Context, email string) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("user with email %s does not exist", email)
	}
	return nil
}

func (r *commonRepository) CheckIfUserExistsByIdAndRole(ctx context.Context, userId, role string) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ? AND role = ?", userId, role).Count(&count).Error; err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("user with id %s does not exist", userId)
	}
	return nil
}

// validateGormId validates the GORM ID.
func validateGormId(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, nil
	}
	newId, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, types.ErrInvalidUUID
	}
	return newId, nil
}