package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "ticket-booking-app-backend/docs"
	"ticket-booking-app-backend/internal/infrastructure/configs"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"ticket-booking-app-backend/internal/application/service"
	"ticket-booking-app-backend/internal/presentation/handlers"
	"ticket-booking-app-backend/internal/presentation/middleware"
)

type Router struct {
	services       *service.Services
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(services *service.Services, authMiddleware *middleware.AuthMiddleware) *Router {
	return &Router{
		services:       services,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) Init(cfg *configs.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		corsMiddleware,
	)

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.initAPI(router)

	return router
}

func (r *Router) initAPI(router *gin.Engine) {
	handlerV1 := handlers.NewHandler(r.services, r.authMiddleware)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
