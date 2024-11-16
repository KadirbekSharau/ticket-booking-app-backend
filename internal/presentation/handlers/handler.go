package handlers

import (
	"github.com/gin-gonic/gin"
	"ticket-booking-app-backend/internal/application/service"
	"ticket-booking-app-backend/internal/presentation/middleware"
)

type Handler struct {
	services       *service.Services
	authMiddleware *middleware.AuthMiddleware
}

func NewHandler(services *service.Services, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		services:       services,
		authMiddleware: authMiddleware,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initEventsRoutes(v1)
	}
}