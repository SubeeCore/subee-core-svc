package service_v1

import (
	"context"
	"time"

	entities_subscriptions_v1 "github.com/subeecore/subee-core-svc/internal/entities/subscriptions/v1"
	entities_users_v1 "github.com/subeecore/subee-core-svc/internal/entities/users/v1"
)

type SubeeCoreService interface {
	//User
	CreateUser(ctx context.Context, req *entities_users_v1.CreateUserRequest) (*entities_users_v1.User, error)

	//Subscription
	CreateSubscription(ctx context.Context, req *entities_subscriptions_v1.CreateSubscriptionRequest) (*entities_subscriptions_v1.Subscription, error)
	GetSubscriptionByID(ctx context.Context, userID string, subscriptionID string) (*entities_subscriptions_v1.Subscription, error)
	FetchSubscriptions(ctx context.Context, userID string) ([]*entities_subscriptions_v1.Subscription, error)
	FinishSubscription(ctx context.Context, userID string, subscriptionID string, finishedAt time.Time) error
	DeleteSubscription(ctx context.Context, userID string, subscriptionID string) error
}
