package config

import (
	"github.com/subeecore/pkg/cache/redis"
	"github.com/subeecore/pkg/config"
	database_postgres "github.com/subeecore/pkg/database/postgres"
	pkg_http "github.com/subeecore/pkg/http"
)

type Config struct {
	ServiceConfig config.Config

	HTTPServerConfig pkg_http.HTTPServerConfig
	PostgresConfig   database_postgres.Config
	RedisConfig      redis.Config
}
