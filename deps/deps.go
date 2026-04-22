// deps/deps.go
package deps

import (
	"stationhub-api/config"
	"stationhub-api/handler"
	"stationhub-api/middleware"
	"stationhub-api/repository"
	"stationhub-api/service"

	"gorm.io/gorm"
)

type Dependencies struct {
	UserRepo *repository.UserRepository

	AuthenticateService *service.AuthenticateService
	WebhookClerkService *service.WebhookClerkService
	ClerkService        *service.ClerkService
	UserService         *service.UserService

	WebhookClerkHandler *handler.WebhookClerkHandler
	UserHandler         *handler.UserHandler
	StationHandler      *handler.StationHandler

	AuthenticateMiddleware *middleware.AuthenticateMiddleware
	ClerkWebhookMiddleware *middleware.ClerkWebhookMiddleware
}

func New(db *gorm.DB, cfg *config.Config) *Dependencies {
	userRepo := repository.NewUserRepository(db)
	stationRepo := repository.NewStationRepository(db)

	authenticateService := service.NewAuthenticateService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	webhookClerkService := service.NewWebhookClerkService(userRepo, userService)
	clerkService := service.NewClerkService(cfg)
	stationService := service.NewStationService(stationRepo)

	webhookClerkHandler := handler.NewWebhookClerkHandler(webhookClerkService)
	userHandler := handler.NewUserHandler(userService)
	stationHandler := handler.NewStationHandler(stationService)

	clerkWebhookMiddleware := middleware.NewClerkWebhookMiddleware(cfg.ClerkWebhookSecret)
	authenticateMiddleware := middleware.NewAuthenticateMiddleware(authenticateService, clerkService, userService, userRepo)

	return &Dependencies{
		UserRepo:               userRepo,
		AuthenticateService:    authenticateService,
		WebhookClerkService:    webhookClerkService,
		UserService:            userService,
		WebhookClerkHandler:    webhookClerkHandler,
		UserHandler:            userHandler,
		AuthenticateMiddleware: authenticateMiddleware,
		ClerkWebhookMiddleware: clerkWebhookMiddleware,
		ClerkService:           clerkService,
		StationHandler:         stationHandler,
	}
}
