package service_v1

import (
	"context"
	"fmt"
	"time"

	"github.com/subeecore/pkg/cache"
	database_v1 "github.com/subeecore/subee-core-svc/internal/database/v1"
)

const (
	userCacheDuration         = time.Hour * 24
	subscriptionCacheDuration = time.Hour * 24
)

func generateSubscriptionCacheKeyByIDForUser(userID string, subscriptionID string) string {
	return fmt.Sprintf("subee-core-svc:subscription:user_id:%v:subscription_id:%v", userID, subscriptionID)
}

func generateMonthlyRecapSubscriptionCacheKeyByID(userID string) string {
	return fmt.Sprintf("subee-core-svc:recap:user_id:%v", userID)
}

func generateSubscriptionsCacheKeyForUser(userID string) string {
	return fmt.Sprintf("subee-core-svc:subscription:user_id:%v", userID)
}

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
