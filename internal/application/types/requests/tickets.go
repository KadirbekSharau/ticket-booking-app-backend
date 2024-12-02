// internal/application/types/requests/tickets.go
package requests

type ReserveTicketsRequest struct {
	Body    ReserveTicketsRequestBody
	EventID  string `json:"event_id" binding:"required"`
	UserID   string
	Role     string
}

type ReserveTicketsRequestBody struct {
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}


type GetUserTicketsRequest struct {
	UserID string
	Role   string
	Status string
}

type GetEventTicketsRequest struct {
	EventID     string
	OrganizerID string
	Role        string
	Status      string
}

type CancelTicketRequest struct {
	TicketID string
	UserID   string
	Role     string
}

type GetTicketByIDRequest struct {
	TicketID string
	UserID   string
	Role     string
}
