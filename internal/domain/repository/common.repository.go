package repository

import (
	"context"
)

type CommonRepository interface {
	CheckIfUserExistsByEmail(ctx context.Context, email string) error
	CheckIfUserExistsByIdAndRole(ctx context.Context, userId, role string) error
}
