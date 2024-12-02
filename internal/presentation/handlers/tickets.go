// internal/application/handlers/tickets.go
package handlers

import (
	"errors"
	"net/http"

	"ticket-booking-app-backend/internal/application/types/requests"
	domainErrors "ticket-booking-app-backend/internal/domain/types"
	"ticket-booking-app-backend/internal/helpers"
	"ticket-booking-app-backend/pkg/values"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// initTicketsRoutes initializes the ticket routes
func (h *Handler) initTicketsRoutes(api *gin.RouterGroup) {
	tickets := api.Group("/tickets", h.authMiddleware.UserIdentity)
	{
		// User routes
		tickets.POST("/reserve", h.reserveTickets)
		tickets.GET("/my", h.getUserTickets)
		tickets.GET("/my/:id", h.getTicketByID)
		tickets.PUT("/my/:id/cancel", h.cancelTicket)

		// Organizer routes
		organizer := tickets.Group("/organizer", h.authMiddleware.RoleMiddleware(values.OrganizerRole))
		{
			organizer.GET("/events/:id", h.getEventTickets)
		}

		// Admin routes
		admin := tickets.Group("/admin", h.authMiddleware.RoleMiddleware(values.AdminRole))
		{
			admin.GET("/events/:id", h.getEventTickets)
		}
	}
}

// @Summary Reserve Tickets
// @Tags tickets
// @Description Reserve tickets for an event
// @Accept json
// @Produce json
// @Param input body requests.ReserveTicketsRequest true "Reservation details"
// @Security ApiKeyAuth
// @Success 201 {array} entities.Ticket
// @Failure 400 {object} helpers.Response
// @Failure 401 {object} helpers.Response
// @Failure 403 {object} helpers.Response
// @Failure 404 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/tickets/reserve [post]
func (h *Handler) reserveTickets(c *gin.Context) {
	var inp requests.ReserveTicketsRequest
	if err := c.BindJSON(&inp.Body); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid input body: "+err.Error())
		return
	}

	userID, err := h.validateContextIDKey(c, values.UserIdCtx)
	if err != nil {
		return
	}
	eventID, err := h.validateQueryIDParam(c, values.EventIdQueryParam)
	if err != nil {
		return
	}
	role, err := h.validateContextKey(c, values.RoleCtx)
	if err != nil {
		logrus.Warn("Error getting role from context")
		return
	}

	inp.UserID = userID
	inp.EventID = eventID
	inp.Role = role

	tickets, err := h.services.Tickets.ReserveTickets(c.Request.Context(), &inp)
	if err != nil {
		if errors.Is(err, domainErrors.ErrInsufficientTickets) {
			helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		logrus.Errorf("Error reserving tickets: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, tickets)
}

// @Summary Get User Tickets
// @Tags tickets
// @Description Get all tickets for the authenticated user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} entities.Ticket
// @Failure 401 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/tickets/my [get]
func (h *Handler) getUserTickets(c *gin.Context) {
	userID, err := h.validateContextIDKey(c, values.UserIdCtx)
	if err != nil {
		return
	}
	role, err := h.validateContextKey(c, values.RoleCtx)
	if err != nil {
		return
	}
	status, err := h.validateQueryParam(c, values.StatusQueryParam)
	if err != nil {
		logrus.Warn("Error getting status from query")
		return
	}

	inp := requests.GetUserTicketsRequest{
		UserID: userID,
		Role:   role,
		Status: status,
	}

	tickets, err := h.services.Tickets.GetUserTickets(c.Request.Context(), &inp)
	if err != nil {
		logrus.Errorf("Error getting user tickets: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// @Summary Get Ticket Details
// @Tags tickets
// @Description Get details of a specific ticket
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Security ApiKeyAuth
// @Success 200 {object} entities.Ticket
// @Failure 401 {object} helpers.Response
// @Failure 403 {object} helpers.Response
// @Failure 404 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/tickets/my/{id} [get]
func (h *Handler) getTicketByID(c *gin.Context) {
	ticketID, err := h.validateRequestIDParam(c, values.IdQueryParam)
	if err != nil {
		return
	}
	userID, err := h.validateContextIDKey(c, values.UserIdCtx)
	if err != nil {
		return
	}
	role, err := h.validateContextKey(c, values.RoleCtx)
	if err != nil {
		return
	}

	inp := requests.GetTicketByIDRequest{
		TicketID: ticketID,
		UserID:   userID,
		Role:     role,
	}

	ticket, err := h.services.Tickets.GetTicketByID(c.Request.Context(), &inp)
	if err != nil {
		if errors.Is(err, domainErrors.ErrTicketNotFound) {
			helpers.NewErrorResponse(c, http.StatusNotFound, "ticket not found")
			return
		}
		logrus.Errorf("Error getting ticket: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// @Summary Get Event Tickets
// @Tags tickets
// @Description Get all tickets for a specific event (Organizer/Admin only)
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Security ApiKeyAuth
// @Success 200 {array} entities.Ticket
// @Failure 401 {object} helpers.Response
// @Failure 403 {object} helpers.Response
// @Failure 404 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/tickets/organizer/events/{id} [get]
func (h *Handler) getEventTickets(c *gin.Context) {
	eventID, err := h.validateQueryIDParam(c, values.EventIdQueryParam)
	if err != nil {
		return
	}
	userID, err := h.validateContextIDKey(c, values.UserIdCtx)
	if err != nil {
		return
	}
	role, err := h.validateContextKey(c, values.RoleCtx)
	if err != nil {
		return
	}
	status, err := h.validateQueryParam(c, values.StatusQueryParam)
	if err != nil {
		logrus.Warn("Error getting status from query")
		return
	}

	inp := requests.GetEventTicketsRequest{
		EventID:     eventID,
		OrganizerID: userID,
		Role:        role,
		Status:      status,
	}

	tickets, err := h.services.Tickets.GetEventTickets(c.Request.Context(), &inp)
	if err != nil {
		logrus.Errorf("Error getting event tickets: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// @Summary Cancel Ticket
// @Tags tickets
// @Description Cancel a reserved ticket
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Security ApiKeyAuth
// @Success 200 {object} helpers.Response
// @Failure 400 {object} helpers.Response
// @Failure 401 {object} helpers.Response
// @Failure 403 {object} helpers.Response
// @Failure 404 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/tickets/my/{id}/cancel [put]
func (h *Handler) cancelTicket(c *gin.Context) {
	ticketID, err := h.validateRequestIDParam(c, values.IdQueryParam)
	if err != nil {
		return
	}
	userID, err := h.validateContextIDKey(c, values.UserIdCtx)
	if err != nil {
		return
	}
	role, err := h.validateContextKey(c, values.RoleCtx)
	if err != nil {
		return
	}

	inp := requests.CancelTicketRequest{
		TicketID: ticketID,
		UserID:   userID,
		Role:     role,
	}

	if err := h.services.Tickets.CancelTicket(c.Request.Context(), &inp); err != nil {
		if errors.Is(err, domainErrors.ErrInvalidTicketStatus) {
			helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, domainErrors.ErrTicketNotFound) {
			helpers.NewErrorResponse(c, http.StatusNotFound, "ticket not found")
			return
		}
		logrus.Errorf("Error cancelling ticket: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, helpers.NewResponse("ticket cancelled successfully"))
}
