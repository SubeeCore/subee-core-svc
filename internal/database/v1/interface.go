package database_v1

import (
	"context"

	entities_recap_v1 "github.com/subeecore/subee-core-svc/internal/entities/recap/v1"
	entities_subscriptions_v1 "github.com/subeecore/subee-core-svc/internal/entities/subscriptions/v1"
	entities_users_v1 "github.com/subeecore/subee-core-svc/internal/entities/users/v1"
)

//go:generate mockgen -source interface.go -destination mocks/mock_database.go -package database_mocks
type Database interface {
	//User
	CreateUser(ctx context.Context, req *entities_users_v1.CreateUserRequest) (*entities_users_v1.User, error)

	//Subscription
	CreateSubscription(ctx context.Context, req *entities_subscriptions_v1.CreateSubscriptionRequest) (*entities_subscriptions_v1.Subscription, error)
	GetSubscriptionByID(ctx context.Context, userID string, subscriptionID string) (*entities_subscriptions_v1.Subscription, error)
	GetMonthlySubscriptionsRecap(ctx context.Context, userID string) (*entities_recap_v1.MonthlyRecap, error)
	GetGlobalSubscriptionsRecap(ctx context.Context, userID string) (*entities_recap_v1.GlobalRecap, error)
	FetchSubscriptions(ctx context.Context, userID string) ([]*entities_subscriptions_v1.Subscription, error)
	FinishSubscription(ctx context.Context, userID string, subscriptionID string, finishedAt string) error
	DeleteSubscription(ctx context.Context, userID string, subscriptionID string) error
}
