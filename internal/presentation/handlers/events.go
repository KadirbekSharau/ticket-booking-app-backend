// internal/application/handlers/events.go
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

// initEventsRoutes initializes the event routes
func (h *Handler) initEventsRoutes(api *gin.RouterGroup) {
	events := api.Group("/events")
	{
		// Public routes
		events.GET("/", h.getActiveEvents) // For all users to see active events

		// Protected routes
		authorized := events.Group("/", h.authMiddleware.UserIdentity)
		{
			// Organizer routes
			organizer := authorized.Group("/organizer", h.authMiddleware.RoleMiddleware(values.OrganizerRole, values.AdminRole))
			{
				organizer.GET("/events", h.getOrganizerEvents) // Get organizer's own events
				organizer.POST("/events", h.createEvent)       // Create new event
				organizer.PUT("/events/:id", h.updateEvent)    // Update own event
				organizer.DELETE("/events/:id", h.deleteEvent) // Delete own event
			}

			// Admin routes
			admin := authorized.Group("/admin", h.authMiddleware.RoleMiddleware(values.AdminRole))
			{
				admin.GET("/events", h.getAllEvents) // Get all events
			}
		}
	}
}

// @Summary Get Active Events
// @Tags events
// @Description Fetch all active events (public access)
// @Accept json
// @Produce json
// @Success 200 {array} entities.Event "List of active events"
// @Failure 500 {object} helpers.Response "Internal server error"
// @Router /api/v1/events [get]
func (h *Handler) getActiveEvents(c *gin.Context) {
	inp := requests.GetEventsRequest{
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

// @Summary Get Organizer Events
// @Tags events
// @Description Fetch all events for the authenticated organizer
// @Accept json
// @Produce json
// @Param status query string false "Event status" Enums(active, cancelled, finished)
// @Security ApiKeyAuth
// @Success 200 {array} entities.Event "List of organizer's events"
// @Failure 401 {object} helpers.Response "Unauthorized"
// @Failure 403 {object} helpers.Response "Forbidden"
// @Failure 500 {object} helpers.Response "Internal server error"
// @Router /api/v1/events/organizer/events [get]
func (h *Handler) getOrganizerEvents(c *gin.Context) {
	organizerID := c.GetString(values.UserIdCtx)
	status := c.DefaultQuery("status", values.EventStatusActive)

	inp := requests.GetEventsByOrganizerRequest{
		OrganizerID: organizerID,
		Role:        values.OrganizerRole,
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
// @Description Create a new event (Organizer only)
// @Accept json
// @Produce json
// @Param input body requests.CreateEventRequestBody true "Event details"
// @Security ApiKeyAuth
// @Success 201 {object} helpers.Response "Event created successfully"
// @Failure 400 {object} helpers.Response "Bad request"
// @Failure 401 {object} helpers.Response "Unauthorized"
// @Failure 403 {object} helpers.Response "Forbidden"
// @Failure 500 {object} helpers.Response "Internal server error"
// @Router /api/v1/events/organizer/events [post]
func (h *Handler) createEvent(c *gin.Context) {
	var inp requests.CreateEventRequest
	if err := c.BindJSON(&inp.Body); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	inp.OrganizerID = c.GetString(values.UserIdCtx)
	inp.Role = values.OrganizerRole

	err := h.services.Events.CreateEvent(c.Request.Context(), &inp)
	if err != nil {
		logrus.Errorf("Error creating event: %s", err)
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, helpers.NewResponse("event created successfully"))
}

// @Summary Update Event
// @Tags events
// @Description Update organizer's own event
// @Accept json
// @Produce json
// @Param id path string true "Event ID" format(uuid)
// @Param input body requests.UpdateEventRequestBody true "Event details"
// @Security ApiKeyAuth
// @Success 200 {object} entities.Event "Updated event"
// @Failure 400 {object} helpers.Response "Bad request"
// @Failure 401 {object} helpers.Response "Unauthorized"
// @Failure 403 {object} helpers.Response "Forbidden"
// @Failure 404 {object} helpers.Response "Event not found"
// @Failure 500 {object} helpers.Response "Internal server error"
// @Router /api/v1/events/organizer/events/{id} [put]
func (h *Handler) updateEvent(c *gin.Context) {
	eventID, err := h.validateRequestIDParam(c, values.IdQueryParam)
	if err != nil {
		return
	}

	var inp requests.UpdateEventRequest
	if err := c.BindJSON(&inp.Body); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	inp.ID = eventID
	inp.OrganizerID = c.GetString(values.UserIdCtx)
	inp.Role = values.OrganizerRole

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
// @Description Delete organizer's own event
// @Accept json
// @Produce json
// @Param id path string true "Event ID" format(uuid)
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Failure 400 {object} helpers.Response "Bad request"
// @Failure 401 {object} helpers.Response "Unauthorized"
// @Failure 403 {object} helpers.Response "Forbidden"
// @Failure 404 {object} helpers.Response "Event not found"
// @Failure 500 {object} helpers.Response "Internal server error"
// @Router /api/v1/events/organizer/events/{id} [delete]
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

	c.Status(http.StatusNoContent)
}

// @Summary Get All Events (Admin)
// @Tags events
// @Description Fetch all events (Admin only)
// @Accept json
// @Produce json
// @Param status query string false "Event status" Enums(active, cancelled, finished)
// @Security ApiKeyAuth
// @Success 200 {array} entities.Event "List of all events"
// @Failure 401 {object} helpers.Response "Unauthorized"
// @Failure 403 {object} helpers.Response "Forbidden"
// @Failure 500 {object} helpers.Response "Internal server error"
// @Router /api/v1/events/admin/events [get]
func (h *Handler) getAllEvents(c *gin.Context) {
	status := c.DefaultQuery("status", "")

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

// Admin handlers for updating and deleting any event follow the same pattern
// Would you like me to include those as well?
