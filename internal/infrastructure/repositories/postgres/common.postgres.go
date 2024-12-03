package postgres

import (
	"context"
	"fmt"

	"ticket-booking-app-backend/internal/infrastructure/drivers/postgres/models"
	"ticket-booking-app-backend/internal/infrastructure/types"
	"ticket-booking-app-backend/pkg/values"

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

func (r *commonRepository) CheckIfEventExists(ctx context.Context, eventID string) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Event{}).Where("id = ?", eventID).Count(&count).Error; err != nil {
		return fmt.Errorf("error checking event existence: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("event with id %s does not exist", eventID)
	}
	return nil
}

func (r *commonRepository) CheckIfEventIsActive(ctx context.Context, eventID string) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Event{}).Where("id = ? AND status = ?", eventID, values.EventStatusActive).Count(&count).Error; err != nil {
		return fmt.Errorf("error checking event status: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("event with id %s is not active", eventID)
	}
	return nil
}

func (r *commonRepository) CheckIfEventBelongsToOrganizer(ctx context.Context, eventID, organizerID string) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Event{}).Where("id = ? AND organizer_id = ?", eventID, organizerID).Count(&count).Error; err != nil {
		return fmt.Errorf("error checking event ownership: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("event with id %s does not belong to organizer with id %s", eventID, organizerID)
	}
	return nil
}

func (r *commonRepository) CheckEventAvailableCapacity(ctx context.Context, eventID string) (int, error) {
	var event models.Event
	if err := r.db.WithContext(ctx).Where("id = ?", eventID).First(&event).Error; err != nil {
		return 0, fmt.Errorf("error getting event: %w", err)
	}
	return event.Capacity - event.TicketsSold, nil
}

func (r *commonRepository) CheckIfUserExceededCapacityForEvent(ctx context.Context, eventID, userID string, ticketCount int) error {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Where(
			"event_id = ? AND user_id = ? AND status IN (?)",
			eventID,
			userID,
			[]string{values.TicketStatusReserved, values.TicketStatusPaid},
		).
		Count(&count).Error

	if err != nil {
		return fmt.Errorf("error checking user capacity: %w", err)
	}

	if ticketCount >= values.MaxTicketsPerPurchase-int(count) {
		return types.ErrTicketLimitExceeded
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
