package repository

import (
	"context"
)

type CommonRepository interface {
	CheckIfUserExistsByEmail(ctx context.Context, email string) error
	CheckIfUserExistsByIdAndRole(ctx context.Context, userId, role string) error
	CheckIfEventIsActive(ctx context.Context, eventID string) error
	CheckIfEventExists(ctx context.Context, eventID string) error
	CheckIfEventBelongsToOrganizer(ctx context.Context, eventID, organizerID string) error
	CheckEventAvailableCapacity(ctx context.Context, eventID string) (int, error)
}
