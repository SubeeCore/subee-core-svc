package handlers_http

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	pkg_http "github.com/subeecore/pkg/http"

	"github.com/subeecore/subee-core-svc/internal/handlers"
	handlers_http_internal_subscriptions_v1 "github.com/subeecore/subee-core-svc/internal/handlers/http/internal/subscriptions/v1"
	handlers_http_internal_users_v1 "github.com/subeecore/subee-core-svc/internal/handlers/http/internal/users/v1"
	service_v1 "github.com/subeecore/subee-core-svc/internal/service/v1"
)

type httpServer struct {
	router  *echo.Echo
	config  pkg_http.HTTPServerConfig
	service *service_v1.Service
}

func NewServer(ctx context.Context, cfg pkg_http.HTTPServerConfig, service *service_v1.Service) (handlers.Server, error) {
	return &httpServer{
		router:  echo.New(),
		config:  cfg,
		service: service,
	}, nil
}

func (s *httpServer) Setup(ctx context.Context) error {
	log.Info().
		Msg("handlers.http.httpServer.Setup: Setting up HTTP server...")

	// setup handlers
	internalUsersV1Handlers := handlers_http_internal_users_v1.NewHandler(ctx, s.service)
	internalSubscriptionsV1Handlers := handlers_http_internal_subscriptions_v1.NewHandler(ctx, s.service)

	// setup middlewares
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recover())
	s.router.Use(middleware.CORS())

	// setup endpoints

	// internal endpoints
	privateV1 := s.router.Group("/private/v1")

	// users related endpoints
	usersV1 := privateV1.Group("/users")
	usersV1.POST("/", internalUsersV1Handlers.Create)

	// subscriptions related endpoints
	subscriptionsV1 := privateV1.Group("/subscriptions")
	subscriptionsV1.POST("/", internalSubscriptionsV1Handlers.CreateSubscription)
	subscriptionsV1.GET("/:user_id", internalSubscriptionsV1Handlers.FetchSubscriptions)
	subscriptionsV1.GET("/:user_id/recap", internalSubscriptionsV1Handlers.GetMonthlySubscriptionsRecap)
	subscriptionsV1.GET("/:user_id/global_recap", internalSubscriptionsV1Handlers.GetGlobalSubscriptionsRecap)
	subscriptionsV1.GET("/:user_id/:subscription_id", internalSubscriptionsV1Handlers.GetSubscriptionByIDForUser)
	subscriptionsV1.PATCH("/:user_id/:subscription_id/finish", internalSubscriptionsV1Handlers.FinishSubscription)
	subscriptionsV1.DELETE("/:user_id/:subscription_id/delete", internalSubscriptionsV1Handlers.DeleteSubscription)

	return nil
}

func (s *httpServer) Start(ctx context.Context) error {
	log.Info().
		Uint16("port", s.config.Port).
		Msg("handlers.http.httpServer.Start: Starting HTTP server...")

	return s.router.Start(fmt.Sprintf(":%d", s.config.Port))
}

func (s *httpServer) Stop(ctx context.Context) error {
	log.Info().
		Msg("handlers.http.httpServer.Stop: Stopping HTTP server...")

	return s.router.Shutdown(ctx)
}
