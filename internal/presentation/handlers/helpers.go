package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ticket-booking-app-backend/internal/helpers"
	"ticket-booking-app-backend/internal/presentation/types"
)

func (h *Handler) validateQueryIDParam(c *gin.Context, key string) (string, error) {
	value := c.Query(key)
	if value == "" {
		helpers.NewErrorResponse(c, http.StatusBadRequest, key + " is required")
		return "", errors.New(key + " is required")
	}
	if err := h.validateUUIDParam(c, value); err != nil {
		return "", err
	}
	return value, nil
}

func (h *Handler) validateRequestParam(c *gin.Context, key string) (string, error) {
	value := c.Param(key)
	if value == "" {
		helpers.NewErrorResponse(c, http.StatusBadRequest, key + " is required")
		return "", errors.New(key + " is required")
	}
	return value, nil
}

func (h *Handler) validateRequestIDParam(c *gin.Context, key string) (string, error) {
	value := c.Param(key)
	if value == "" {
		helpers.NewErrorResponse(c, http.StatusBadRequest, key + " is required")
		return "", errors.New(key + " is required")
	}
	if err := h.validateUUIDParam(c, value); err != nil {
		return "", err
	}
	return value, nil
}

func (h *Handler) validateContextKey(c *gin.Context, key string) (string, error) {
	value, exists := c.Get(key)
	if !exists {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, key + " is required")
		return "", errors.New(key + " is required")
	}
	valueString := value.(string)

	if err := h.validateUUIDParam(c, valueString); err != nil {
		return "", err
	}
	return valueString, nil
}

func (h *Handler) validateUUIDParam(c *gin.Context, value string) error {
	if err := uuid.Validate(value); err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, types.ErrInvalidUUID.Error())
		return errors.New(types.ErrInvalidUUID.Error())
	}
	return nil
}
