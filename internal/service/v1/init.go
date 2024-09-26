package service_v1

import (
	"context"
	"time"

	"github.com/subeecore/pkg/cache"
	database_v1 "github.com/subeecore/subee-core-svc/internal/database/v1"
)

const (
	userCacheDuration = time.Hour * 24
)

type Service struct {
	store database_v1.Database
	cache cache.Cache
}

func NewSubeeCoreService(ctx context.Context, store database_v1.Database, cache cache.Cache) (*Service, error) {
	return &Service{
		store: store,
		cache: cache,
	}, nil
}
