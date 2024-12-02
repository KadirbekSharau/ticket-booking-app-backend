package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User model with UUID primary key.
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Email     string         `gorm:"type:varchar(255);not null;unique" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"password"`
	Name      string         `gorm:"type:varchar(100)" json:"name"`
	Address   string         `gorm:"type:varchar(255)" json:"address"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Role      string         `gorm:"type:varchar(50);not null;default:'user'" json:"role"` // Roles: 'user', 'organizer'
	Events    []Event        `gorm:"foreignKey:OrganizerID" json:"events"`
	Tickets   []Ticket       `gorm:"constraint:OnDelete:SET NULL;" json:"tickets"`
	Payments  []Payment      `gorm:"constraint:OnDelete:CASCADE;" json:"payments"`
}

// Event model with UUID primary key.
type Event struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	OrganizerID uuid.UUID      `gorm:"type:uuid;not null" json:"organizer_id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Location    string         `gorm:"type:varchar(255)" json:"location"`
	Date        time.Time      `gorm:"type:timestamptz;not null" json:"date"`
	Capacity    int            `gorm:"not null" json:"capacity"`
	TicketsSold int            `gorm:"not null;default:0" json:"tickets_sold"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Status      string         `gorm:"type:varchar(50);not null;default:'upcoming'" json:"status"` // Status: 'upcoming', 'ongoing', 'completed', 'cancelled'
	Tickets     []Ticket       `gorm:"constraint:OnDelete:CASCADE;" json:"tickets"`
}

// Ticket model with UUID primary key.
type Ticket struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	EventID    uuid.UUID      `gorm:"type:uuid;not null" json:"event_id"`
	UserID     uuid.UUID      `gorm:"type:uuid" json:"user_id"`                                   // Nullable for unclaimed tickets
	Status     string         `gorm:"type:varchar(50);not null;default:'reserved'" json:"status"` // Status: 'reserved', 'paid', 'cancelled', 'expired'
	ReservedAt time.Time      `gorm:"autoCreateTime" json:"reserved_at"`
	PaidAt     time.Time      `json:"paid_at"`
	Price      float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Event      Event          `gorm:"foreignKey:EventID" json:"event"`
}

// Payment model with UUID primary key.
type Payment struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	TicketID        uuid.UUID      `gorm:"type:uuid;not null" json:"ticket_id"`
	StripePaymentID string         `gorm:"type:varchar(255);unique" json:"stripe_payment_id"`
	Amount          float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status          string         `gorm:"type:varchar(50);not null;default:'pending'" json:"status"` // Status: 'pending', 'completed', 'failed', 'refunded'
}
