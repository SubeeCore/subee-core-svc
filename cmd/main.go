package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/subeecore/pkg/cache/redis"
	pkg_config "github.com/subeecore/pkg/config"
	database_postgres "github.com/subeecore/pkg/database/postgres"

	"github.com/subeecore/subee-core-svc/internal/config"
	database_v1_pgx "github.com/subeecore/subee-core-svc/internal/database/v1/pgx"
	handlers_http "github.com/subeecore/subee-core-svc/internal/handlers/http"
	service "github.com/subeecore/subee-core-svc/internal/service/v1"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	cfg := &config.Config{}
	err := pkg_config.ParseConfig(cfg)
	if err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to parse config")
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cacheConnection := redis.GetConnection(ctx, &cfg.RedisConfig)
	cacheRedis := redis.NewRedisCache(ctx, cacheConnection)

	databaseConnection, err := database_postgres.NewDatabaseConnection(ctx, &cfg.PostgresConfig)
	if err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to create database connection")
	}
	databaseClient := database_v1_pgx.NewClient(ctx, databaseConnection)

	subeeCoreService, err := service.NewSubeeCoreService(ctx, databaseClient, cacheRedis)
	if err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to create subee core service")
	}

	httpServer, err := handlers_http.NewServer(ctx, cfg.HTTPServerConfig, subeeCoreService)
	if err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to create http server")
	}

	if err := httpServer.Setup(ctx); err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to setup http server")
	}

	if err := httpServer.Start(ctx); err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to start http server")
	}

	<-sigs
	cancel()

	if err := httpServer.Stop(ctx); err != nil {
		log.Fatal().Err(err).
			Msg("main: unable to stop http server")
	}

	os.Exit(0)
}
