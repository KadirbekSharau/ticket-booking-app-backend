package internal

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ticket-booking-app-backend/internal/application/service"
	"ticket-booking-app-backend/internal/application/types/requests"
	"ticket-booking-app-backend/internal/domain/repository"
	"ticket-booking-app-backend/internal/helpers"
	"ticket-booking-app-backend/internal/infrastructure/configs"
	postgres "ticket-booking-app-backend/internal/infrastructure/drivers/postgres/connection"
	infrastructure "ticket-booking-app-backend/internal/infrastructure/http"
	"ticket-booking-app-backend/internal/presentation/middleware"

	"github.com/sirupsen/logrus"
	//"github.com/xuri/excelize/v2"
)

// @title TicketBooking API
// @version 1.0
// @description REST API endpoints for Munai Plan App

// @host localhost:8000
// @BasePath /api/v1/

func Run(configPath string) {
	cfg, err := configs.Init(configPath)
	if err != nil {
		logrus.Error(err)
		return
	}

	// Dependencies
	db := postgres.NewDatabase()
	if db == nil {
		logrus.Error("failed to initialize database connection")
		return
	}

	jwt, err := helpers.NewJwt()
	if err != nil {
		logrus.Error(err)
		return
	}

	// Initializing repositories
	repos := repository.NewRepositories(db.Conn)

	// Initializing services
	services := service.NewServices(repos, jwt)

	adminEmail, err := helpers.GetEnv("ADMIN_EMAIL")
	if err != nil {
		logrus.Error(err)
		return
	}

	adminPassword, err := helpers.GetEnv("ADMIN_PASSWORD")
	if err != nil {
		logrus.Error(err)
		return
	}

	err = services.AdminSignUp(context.Background(), &requests.AdminSignUpRequest{
		Email:    adminEmail,
		Password: adminPassword,
	})

	if err != nil {
		return
	}

	// Initializing middleware
	authMiddleware := middleware.NewAuthMiddleware(jwt)

	// Initializing router and handlers
	router := infrastructure.NewRouter(services, authMiddleware)

	// HTTP Server
	srv := infrastructure.NewServer(cfg, router.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logrus.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logrus.Errorf("failed to stop server: %v", err)
	}

	sqlDB, err := db.Conn.DB()
	if err != nil {
		logrus.Error(err.Error())
	}

	if err := sqlDB.Close(); err != nil {
		logrus.Error(err.Error())
	}
}
