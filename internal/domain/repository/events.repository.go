// domain/repository/events.repository.go
package repository

import (
    "context"
    "ticket-booking-app-backend/internal/domain/entities"
)

type EventsRepository interface {
    // Create operations
    CreateEvent(ctx context.Context, event *entities.Event, organizerID string) error
    
    // Read operations
    GetEventByID(ctx context.Context, eventID string) (*entities.Event, error)
    GetEventsByOrganizer(ctx context.Context, status string, organizerID string) ([]*entities.Event, error)
    GetEvents(ctx context.Context, status string) ([]*entities.Event, error)
    
    // Update operations
    UpdateEvent(ctx context.Context, organizerID string, event *entities.Event) error
    UpdateEventStatus(ctx context.Context, eventID, organizerID, status string) error
    UpdateEventCapacity(ctx context.Context, eventID string, capacity int) error
    IncrementTicketsSold(ctx context.Context, eventID string) error
    
    // Status management
    UpdateExpiredEvents(ctx context.Context) error
    
    // Capacity checks
    CheckEventCapacityIsFull(ctx context.Context, eventID string) (bool, error)
    
    // Delete operations (soft delete via status update)
    DeleteEvent(ctx context.Context, eventID, organizerID string) error
}