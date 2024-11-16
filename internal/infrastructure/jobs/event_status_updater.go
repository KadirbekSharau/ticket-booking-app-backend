// internal/infrastructure/jobs/event_status_updater.go
package jobs

import (
	"context"
	"time"

	"ticket-booking-app-backend/internal/domain/repository"

	"github.com/sirupsen/logrus"
)

type EventStatusUpdater struct {
	repo repository.EventsRepository
}

func NewEventStatusUpdater(repo repository.EventsRepository) *EventStatusUpdater {
	return &EventStatusUpdater{
		repo: repo,
	}
}

func (u *EventStatusUpdater) Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		// Run once at startup
		logrus.Warn("Running initial expired events update")
		if err := u.repo.UpdateExpiredEvents(ctx); err != nil {
			logrus.Errorf("Initial expired events update failed: %v", err)
		}

		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if err := u.repo.UpdateExpiredEvents(ctx); err != nil {
					logrus.Errorf("Error updating expired events: %v", err)
				} else {
					logrus.Debug("Successfully updated expired events")
				}
			}
		}
	}()
}
