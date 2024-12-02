// internal/application/handlers/events.go
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"ticket-booking-app-backend/internal/application/types/requests"
	domainErrors "ticket-booking-app-backend/internal/domain/types"
	"ticket-booking-app-backend/internal/helpers"
	"ticket-booking-app-backend/pkg/values"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// initEventsRoutes initializes the event routes
func (h *Handler) initEventsRoutes(api *gin.RouterGroup) {
	events := api.Group("/events", h.authMiddleware.UserIdentity)
	{
		// Public routes
		events.GET("/", h.getActiveEvents) // For all users to see active events

		// Protected routes
		// Organizer routes
		organizer := events.Group("/organizer", h.authMiddleware.RoleMiddleware(values.OrganizerRole, values.AdminRole))
		{
			organizer.GET("/", h.getOrganizerEvents)    // Get organizer's own events
			organizer.POST("/", h.createEvent)          // Create new event
			organizer.PUT("/:id", h.updateEvent)        // Update own event
			organizer.DELETE("/:id", h.deleteEvent)     // Delete own event
			organizer.PUT("/cancel/:id", h.cancelEvent) // Cancel own event
		}

		// Admin routes
		admin := events.Group("/admin", h.authMiddleware.RoleMiddleware(values.AdminRole))
		{
			admin.GET("/", h.getAllEvents) // Get all events
		}
	}
}

// @Summary List Active Events
// @Tags events
// @Description Get a list of all active events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} entities.Event
// @Failure 401 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/events [get]
func (h *Handler) getActiveEvents(c *gin.Context) {
	role, err := h.validateContextKey(c, values.RoleCtx)
	if err != nil {
		logrus.Warn("Error getting role from context")
		return
	}

	inp := requests.GetEventsRequest{
		Role:   role,
		Status: values.EventStatusActive,
	}

	events, err := h.services.Events.GetEvents(c.Request.Context(), &inp)
	if err != nil {
		logrus.Errorf("Error getting active events: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, events)
}

// @Summary List Organizer Events
// @Tags events
// @Description Get list of events for authenticated organizer
// @Accept json
// @Produce json
// @Param status query string false "Event status filter (active/cancelled/finished)"
// @Security ApiKeyAuth
// @Success 200 {array} entities.Event
// @Failure 400 {object} helpers.Response
// @Failure 401 {object} helpers.Response
// @Failure 403 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/events/organizer [get]
func (h *Handler) getOrganizerEvents(c *gin.Context) {
	organizerID, err := h.validateContextIDKey(c, values.UserIdCtx)
	if err != nil {
		return
	}
	role, err := h.validateContextKey(c, values.RoleCtx)
	if err != nil {
		logrus.Warn("Error getting role from context")
		return
	}
	status, err := h.validateQueryParam(c, values.StatusQueryParam)
	if err != nil {
		logrus.Warn("Error getting role from context")
		return
	}

	inp := requests.GetEventsByOrganizerRequest{
		OrganizerID: organizerID,
		Role:        role,
		Status:      status,
	}

	events, err := h.services.Events.GetEventsByOrganizer(c.Request.Context(), &inp)
	if err != nil {
		logrus.Errorf("Error getting organizer events: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Create Event
// @Tags events
// @Description Create a new event
// @Accept json
// @Produce json
// @Param input body requests.CreateEventRequestBody true "Event data"
// @Security ApiKeyAuth
// @Success 201 {object} helpers.Response
// @Failure 400 {object} helpers.Response
// @Failure 401 {object} helpers.Response
// @Failure 403 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/events/organizer [post]
func (h *Handler) createEvent(c *gin.Context) {
	var inp requests.CreateEventRequest
	if err := c.BindJSON(&inp.Body); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid input body: "+err.Error())
		return
	}

	organizerID, err := h.validateContextIDKey(c, values.UserIdCtx)
	if err != nil {
		return
	}
	role, err := h.validateContextKey(c, values.RoleCtx)
	fmt.Println(role)
	if err != nil {
		logrus.Warn("Error getting role from context")
		return
	}

	inp.OrganizerID = organizerID
	inp.Role = role

	err = h.services.Events.CreateEvent(c.Request.Context(), &inp)
	if err != nil {
		logrus.Errorf("Error creating event: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, helpers.NewResponse("event created successfully"))
}

// @Summary Update Event
// @Tags events
// @Description Update an existing event
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param input body requests.UpdateEventRequestBody true "Event data"
// @Security ApiKeyAuth
// @Success 200 {object} entities.Event
// @Failure 400 {object} helpers.Response
// @Failure 401 {object} helpers.Response
// @Failure 403 {object} helpers.Response
// @Failure 404 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/events/organizer/{id} [put]
func (h *Handler) updateEvent(c *gin.Context) {
	eventID, err := h.validateRequestIDParam(c, values.IdQueryParam)
	if err != nil {
		return
	}
	organizerID, err := h.validateContextIDKey(c, values.UserIdCtx)
	if err != nil {
		return
	}
	role, err := h.validateContextKey(c, values.RoleCtx)
	if err != nil {
		logrus.Warn("Error getting role from context")
		return
	}

	var inp requests.UpdateEventRequest
	if err := c.BindJSON(&inp.Body); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	inp.ID = eventID
	inp.OrganizerID = organizerID
	inp.Role = role

	event, err := h.services.Events.UpdateEvent(c.Request.Context(), &inp)
	if err != nil {
		if errors.Is(err, domainErrors.ErrEventNotFound) {
			helpers.NewErrorResponse(c, http.StatusNotFound, "event not found")
			return
		}
		logrus.Errorf("Error updating event: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, event)
}

// @Summary Delete Event
// @Tags events
// @Description Delete an event
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Security ApiKeyAuth
// @Success 200 {object} helpers.Response
// @Failure 400 {object} helpers.Response
// @Failure 401 {object} helpers.Response
// @Failure 403 {object} helpers.Response
// @Failure 404 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/events/organizer/{id} [delete]
func (h *Handler) deleteEvent(c *gin.Context) {
	eventID, err := h.validateRequestIDParam(c, values.IdQueryParam)
	if err != nil {
		return
	}

	inp := requests.DeleteEventRequest{
		ID:          eventID,
		OrganizerID: c.GetString(values.UserIdCtx),
		Role:        values.OrganizerRole,
	}

	if err := h.services.Events.DeleteEvent(c.Request.Context(), &inp); err != nil {
		if errors.Is(err, domainErrors.ErrEventNotFound) {
			helpers.NewErrorResponse(c, http.StatusNotFound, "event not found")
			return
		}
		logrus.Errorf("Error deleting event: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, helpers.NewResponse("event deleted successfully"))
}

// @Summary List All Events
// @Tags events
// @Description Get all events (admin only)
// @Accept json
// @Produce json
// @Param status query string false "Event status filter (active/cancelled/finished)"
// @Security ApiKeyAuth
// @Success 200 {array} entities.Event
// @Failure 401 {object} helpers.Response
// @Failure 403 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/events/admin [get]
func (h *Handler) getAllEvents(c *gin.Context) {
	status, err := h.validateQueryParam(c, values.StatusQueryParam)
	if err != nil {
		logrus.Warn("Error getting status from query")
		return
	}

	inp := requests.GetEventsRequest{
		Role:   values.AdminRole,
		Status: status,
	}

	events, err := h.services.Events.GetEvents(c.Request.Context(), &inp)
	if err != nil {
		logrus.Errorf("Error getting all events: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Cancel Event
// @Tags events
// @Description Cancel an event
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Security ApiKeyAuth
// @Success 200 {object} helpers.Response
// @Failure 400 {object} helpers.Response
// @Failure 401 {object} helpers.Response
// @Failure 403 {object} helpers.Response
// @Failure 404 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/events/organizer/cancel/{id} [put]
func (h *Handler) cancelEvent(c *gin.Context) {
	eventID, err := h.validateRequestIDParam(c, values.IdQueryParam)
	if err != nil {
		return
	}
	role, err := h.validateContextKey(c, values.RoleCtx)
	if err != nil {
		logrus.Warn("Error getting role from context")
		return
	}
	organizerID, err := h.validateContextIDKey(c, values.UserIdCtx)
	if err != nil {
		return
	}

	inp := requests.CancelEventRequest{
		ID:          eventID,
		Role:        role,
		OrganizerID: organizerID,
	}

	if err := h.services.Events.CancelEvent(c.Request.Context(), &inp); err != nil {
		if errors.Is(err, domainErrors.ErrEventNotFound) {
			helpers.NewErrorResponse(c, http.StatusNotFound, "event not found")
			return
		}
		logrus.Errorf("Error cancelling event: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, helpers.NewResponse("event cancelled successfully"))
}
