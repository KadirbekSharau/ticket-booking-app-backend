// infrastructure/repositories/postgres/events.postgres.go
package postgres

import (
    "context"
    "errors"
    "time"

    "ticket-booking-app-backend/internal/domain/entities"
    "ticket-booking-app-backend/internal/infrastructure/drivers/postgres/models"
    "ticket-booking-app-backend/internal/infrastructure/types"
    "ticket-booking-app-backend/pkg/values"

    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"
)

type eventsRepository struct {
    db *gorm.DB
}

func NewEventsRepository(db *gorm.DB) *eventsRepository {
    return &eventsRepository{db: db}
}

func (r *eventsRepository) CreateEvent(ctx context.Context, event *entities.Event, organizerID string) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        gormEvent := toGormEvent(event)
        orgID, err := validateGormId(organizerID)
        if err != nil {
            return err
        }
        
        gormEvent.OrganizerID = orgID
        gormEvent.Status = values.EventStatusActive
        
        if err := tx.Create(gormEvent).Error; err != nil {
            return err
        }
        
        return nil
    })
}

func (r *eventsRepository) GetEventByID(ctx context.Context, eventID string) (*entities.Event, error) {
    // First update any expired events
    if err := r.UpdateExpiredEvents(ctx); err != nil {
        logrus.Errorf("Failed to update expired events: %v", err)
    }

    var event models.Event
    err := r.db.WithContext(ctx).
        Where("id = ?", eventID).
        First(&event).Error
        
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, types.ErrEventNotFound
    }
    if err != nil {
        return nil, err
    }
    
    return toDomainEvent(&event), nil
}

func (r *eventsRepository) GetEventsByOrganizer(ctx context.Context, status string, organizerID string) ([]*entities.Event, error) {
    // First update any expired events
    if err := r.UpdateExpiredEvents(ctx); err != nil {
        logrus.Errorf("Failed to update expired events: %v", err)
    }

    var events []models.Event
    query := r.db.WithContext(ctx).Where("organizer_id = ?", organizerID)
    
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    if err := query.Find(&events).Error; err != nil {
        return nil, err
    }

    return toDomainEvents(events), nil
}

func (r *eventsRepository) GetEvents(ctx context.Context, status string) ([]*entities.Event, error) {
    // First update any expired events
    if err := r.UpdateExpiredEvents(ctx); err != nil {
        logrus.Errorf("Failed to update expired events: %v", err)
    }

    var events []models.Event
    query := r.db.WithContext(ctx)
    
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    if err := query.Find(&events).Error; err != nil {
        return nil, err
    }

    return toDomainEvents(events), nil
}

func (r *eventsRepository) UpdateEvent(ctx context.Context, organizerID string, event *entities.Event) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var existingEvent models.Event
        err := tx.Where("id = ? AND organizer_id = ?", event.ID, organizerID).
            First(&existingEvent).Error
            
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return types.ErrEventNotFound
        }
        if err != nil {
            return err
        }

        // Don't allow updating if event is finished or cancelled
        if existingEvent.Status == values.EventStatusFinished || 
           existingEvent.Status == values.EventStatusCancelled {
            return errors.New("cannot update finished or cancelled event")
        }
        
        return tx.Model(&existingEvent).Updates(toGormEvent(event)).Error
    })
}

func (r *eventsRepository) UpdateEventStatus(ctx context.Context, eventID, organizerID, status string) error {
    result := r.db.WithContext(ctx).
        Model(&models.Event{}).
        Where("id = ? AND organizer_id = ?", eventID, organizerID).
        Update("status", status)
        
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return types.ErrEventNotFound
    }
    return nil
}

func (r *eventsRepository) UpdateExpiredEvents(ctx context.Context) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        result := tx.Model(&models.Event{}).
            Where("date < ? AND status = ?", time.Now(), values.EventStatusActive).
            Update("status", values.EventStatusFinished)
        
        if result.Error != nil {
            return result.Error
        }
        
        return nil
    })
}

func (r *eventsRepository) CheckEventCapacityIsFull(ctx context.Context, eventID string) (bool, error) {
    var event models.Event
    err := r.db.WithContext(ctx).
        Where("id = ?", eventID).
        First(&event).Error
        
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return false, types.ErrEventNotFound
    }
    if err != nil {
        return false, err
    }

    return event.TicketsSold >= event.Capacity, nil
}

func (r *eventsRepository) UpdateEventCapacity(ctx context.Context, eventID string, capacity int) error {
    result := r.db.WithContext(ctx).
        Model(&models.Event{}).
        Where("id = ?", eventID).
        Update("capacity", capacity)
        
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return types.ErrEventNotFound
    }
    return nil
}

func (r *eventsRepository) IncrementTicketsSold(ctx context.Context, eventID string) error {
    result := r.db.WithContext(ctx).
        Model(&models.Event{}).
        Where("id = ?", eventID).
        Update("tickets_sold", gorm.Expr("tickets_sold + ?", 1))
        
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return types.ErrEventNotFound
    }
    return nil
}

func (r *eventsRepository) DeleteEvent(ctx context.Context, eventID, organizerID string) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var event models.Event
        err := tx.Where("id = ? AND organizer_id = ?", eventID, organizerID).
            First(&event).Error
            
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return types.ErrEventNotFound
        }
        if err != nil {
            return err
        }

        // Don't allow deleting if event is already finished
        if event.Status == values.EventStatusFinished {
            return errors.New("cannot delete finished event")
        }

        // Update status to cancelled instead of deleting
        return tx.Model(&event).Update("status", values.EventStatusCancelled).Error
    })
}

// Helper functions
func toDomainEvents(events []models.Event) []*entities.Event {
    result := make([]*entities.Event, len(events))
    for i, event := range events {
        result[i] = toDomainEvent(&event)
    }
    return result
}

func toDomainEvent(eventModel *models.Event) *entities.Event {
    return &entities.Event{
        ID:          eventModel.ID.String(),
        Title:       eventModel.Title,
        Description: eventModel.Description,
        Location:    eventModel.Location,
        Date:        eventModel.Date,
        Capacity:    eventModel.Capacity,
        TicketsSold: eventModel.TicketsSold,
        Price:       eventModel.Price,
        Status:      eventModel.Status,
        CreatedAt:   eventModel.CreatedAt,
    }
}

func toGormEvent(event *entities.Event) *models.Event {
    var eventID uuid.UUID
    var err error
    if event.ID != "" {
        eventID, err = validateGormId(event.ID)
        if err != nil {
            return nil
        }
    }

    return &models.Event{
        ID:          eventID,
        Title:       event.Title,
        Description: event.Description,
        Location:    event.Location,
        Date:        event.Date,
        Capacity:    event.Capacity,
        TicketsSold: event.TicketsSold,
        Price:       event.Price,
        Status:      event.Status,
    }
}