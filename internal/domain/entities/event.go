package entities

import (
	"time"
)

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Date        time.Time `json:"date"`
	Capacity    int       `json:"capacity"`
	TicketsSold int       `json:"tickets_sold"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	Tickets     []*Ticket `json:"tickets"`
	CreatedAt   time.Time `json:"created_at"`
}
