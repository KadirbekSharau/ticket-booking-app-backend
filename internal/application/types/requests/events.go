// internal/application/types/requests/events.go
package requests

import (
    "time"
)

type CreateEventRequestBody struct {
    Title       string    `json:"title" binding:"required"`
    Description string    `json:"description" binding:"required"`
    Location    string    `json:"location" binding:"required"`
    Date        time.Time `json:"date" binding:"required,gtefield=time.Now"`
    Capacity    int       `json:"capacity" binding:"required,gt=0"`
    Price       float64   `json:"price" binding:"required,gte=0"`
}

type CreateEventRequest struct {
    Body        CreateEventRequestBody
    OrganizerID string
    Role        string
}

type UpdateEventRequestBody struct {
    Title       string    `json:"title" binding:"required"`
    Description string    `json:"description" binding:"required"`
    Location    string    `json:"location" binding:"required"`
    Date        time.Time `json:"date" binding:"required,gtefield=time.Now"`
    Capacity    int       `json:"capacity" binding:"required,gt=0"`
    Price       float64   `json:"price" binding:"required,gte=0"`
}

type UpdateEventRequest struct {
    ID          string
    OrganizerID string
    Role        string
    Body        UpdateEventRequestBody
}

type GetEventsByOrganizerRequest struct {
    OrganizerID string
    Role        string
    Status      string
}

type GetEventsRequest struct {
    Status string
    Role   string
}

type GetEventByIDRequest struct {
    ID          string
    OrganizerID string
    Role        string
}

type DeleteEventRequest struct {
    ID          string
    OrganizerID string
    Role        string
}