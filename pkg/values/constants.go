package values

const (
	AuthorizationHeader                         = "Authorization"
	UserIdCtx                                   = "userId"
	RoleCtx                                     = "role"
	UserAccessTokenCtx                          = "accessToken"
	UserRefreshTokenCtx                         = "refreshToken"
)

const (
	UserRole = "user"
	AdminRole = "admin"
	OrganizerRole = "organizer"
)

const (
	EventStatusActive = "active"
	EventStatusCancelled = "cancelled"
	EventStatusFinished = "finished"
)