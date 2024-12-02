package values

const (
	AuthorizationHeader = "Authorization"
	UserIdCtx           = "userId"
	RoleCtx             = "role"
	UserAccessTokenCtx  = "accessToken"
	UserRefreshTokenCtx = "refreshToken"
)

const (
	UserRole      = "user"
	AdminRole     = "admin"
	OrganizerRole = "organizer"
)

const (
	EventStatusActive    = "active"
	EventStatusCancelled = "cancelled"
	EventStatusFinished  = "finished"
)

const (
	StatusQueryParam = "status"
)

const (
	TicketStatusReserved  = "reserved"
	TicketStatusPaid      = "paid"
	TicketStatusCancelled = "cancelled"
	TicketStatusExpired   = "expired"
)

const (
	PaymentStatusPending   = "pending"
	PaymentStatusCompleted = "completed"
	PaymentStatusFailed    = "failed"
	PaymentStatusRefunded  = "refunded"
)

// Ticket limits and timeouts
const (
	MaxTicketsPerPurchase    = 5
	TicketReservationMinutes = 15
)
