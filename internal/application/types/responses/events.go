// internal/application/types/responses/events.go
package responses

type EventResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Event   interface{} `json:"event,omitempty"`
}