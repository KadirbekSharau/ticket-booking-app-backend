package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"ticket-booking-app-backend/internal/helpers"
	"ticket-booking-app-backend/pkg/values"
)

type AuthMiddleware struct {
	Jwt               helpers.Jwt
}

func NewAuthMiddleware(jwt helpers.Jwt) *AuthMiddleware {
	return &AuthMiddleware{
		Jwt:               jwt,
	}
}

func (m *AuthMiddleware) RefreshTokenIdentity(c *gin.Context) {
	header := c.GetHeader(values.AuthorizationHeader)
	if header == "" {
		helpers.NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		helpers.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userClaims, err := m.Jwt.VerifyRefreshToken(headerParts[1])
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(values.UserIdCtx, userClaims.UserId)
	c.Set(values.UserRefreshTokenCtx, header)
}

func (m *AuthMiddleware) UserIdentity(c *gin.Context) {
	header := c.GetHeader(values.AuthorizationHeader)
	if header == "" {
		helpers.NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		helpers.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userClaims, err := m.Jwt.Verify(headerParts[1])
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(values.UserIdCtx, userClaims.UserId)
	c.Set(values.RoleCtx, userClaims.Role)
	c.Set(values.UserRefreshTokenCtx, header)
}

// RoleMiddleware checks if a user has the required role to access a route
func (m *AuthMiddleware) RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user's role from the context (set by UserIdentity middleware)
		role, exists := c.Get(values.RoleCtx)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No role found in token"})
			c.Abort()
			return
		}

		// Check if the role is in the list of allowed roles
		userRole := role.(string)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next() // Role is allowed, continue
				return
			}
		}

		// If the role is not allowed, return Forbidden
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
		c.Abort()
	}
}
