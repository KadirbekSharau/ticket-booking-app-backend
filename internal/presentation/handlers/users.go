package handlers

import (
	"errors"
	"net/http"

	"ticket-booking-app-backend/internal/application/types/requests"
	domainErrors "ticket-booking-app-backend/internal/domain/types"
	"ticket-booking-app-backend/internal/helpers"

	"github.com/gin-gonic/gin"
)

// initUsersRoutes initializes the user routes.
func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-in", h.userSignIn)
		users.POST("/sign-up", h.userSignUp)
	}
	admin := api.Group("/admin")
	{
		admin.POST("/sign-in", h.adminSignIn)
	}
	organizer := api.Group("/organizer")
	{
		organizer.POST("/sign-in", h.organizerSignIn)
		organizer.POST("/sign-up", h.organizerSignUp)
	}
}

// userSignUp handles the user sign up request.
// @Summary User SignUp
// @Tags users-auth
// @Description Register a new user
// @Accept json
// @Produce json
// @Param input body requests.UserSignUpRequest true "User sign up info"
// @Success 201 {object} helpers.Response
// @Failure 400 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/users/sign-up [post]
func (h *Handler) userSignUp(c *gin.Context) {
	var inp requests.UserSignUpRequest
	if err := c.BindJSON(&inp); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.services.Users.UserSignUp(c.Request.Context(), &inp)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, helpers.NewResponse("user created"))
}

// userSignIn handles the user sign in request.
// @Summary User SignIn
// @Tags users-auth
// @Description Authenticate an existing user
// @Accept json
// @Produce json
// @Param input body requests.UserSignInRequest true "User sign in info"
// @Success 200 {object} responses.TokenResponse
// @Failure 400 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/users/sign-in [post]
func (h *Handler) userSignIn(c *gin.Context) {
	var inp requests.UserSignInRequest
	if err := c.BindJSON(&inp); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	res, err := h.services.Users.UserSignIn(c.Request.Context(), &inp)
	if err != nil {
		if errors.Is(err, domainErrors.ErrUserNotFound) {
			helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, &res)
}

// adminSignIn handles the admin sign in request.
// @Summary Admin SignIn
// @Tags admin-auth
// @Description Authenticate an admin user
// @Accept json
// @Produce json
// @Param input body requests.AdminSignInRequest true "Admin sign in info"
// @Success 200 {object} responses.TokenResponse
// @Failure 400 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/admin/sign-in [post]
func (h *Handler) adminSignIn(c *gin.Context) {
	var inp requests.AdminSignInRequest
	if err := c.BindJSON(&inp); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	res, err := h.services.Users.AdminSignIn(c.Request.Context(), &inp)
	if err != nil {
		if errors.Is(err, domainErrors.ErrAdminNotFound) {
			helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, &res)
}

// organizerSignUp handles the organizer sign up request.
// @Summary Organizer SignUp
// @Tags organizer-auth
// @Description Register a new organizer
// @Accept json
// @Produce json
// @Param input body requests.OrganizerSignUpRequest true "Organizer sign up info"
// @Success 201 {object} helpers.Response
// @Failure 400 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/organizer/sign-up [post]
func (h *Handler) organizerSignUp(c *gin.Context) {
	var inp requests.OrganizerSignUpRequest
	if err := c.BindJSON(&inp); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.services.Users.OrganizerSignUp(c.Request.Context(), &inp)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, helpers.NewResponse("organizer created"))
}

// organizerSignIn handles the organizer sign in request.
// @Summary Organizer SignIn
// @Tags organizer-auth
// @Description Authenticate an organizer user
// @Accept json
// @Produce json
// @Param input body requests.OrganizerSignInRequest true "Organizer sign in info"
// @Success 200 {object} responses.TokenResponse
// @Failure 400 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /api/v1/organizer/sign-in [post]
func (h *Handler) organizerSignIn(c *gin.Context) {
	var inp requests.OrganizerSignInRequest
	if err := c.BindJSON(&inp); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	res, err := h.services.Users.OrganizerSignIn(c.Request.Context(), &inp)
	if err != nil {
		if errors.Is(err, domainErrors.ErrOrganizerNotFound) {
			helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, &res)
}
