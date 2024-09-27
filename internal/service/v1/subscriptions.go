package service_v1

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"

	entities_subscriptions_v1 "github.com/subeecore/subee-core-svc/internal/entities/subscriptions/v1"
)

func (s *Service) CreateSubscription(ctx context.Context, req *entities_subscriptions_v1.CreateSubscriptionRequest) (*entities_subscriptions_v1.Subscription, error) {
	subscription, err := s.store.CreateSubscription(ctx, req)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *Service) GetSubscriptionByID(ctx context.Context, userID string, subscriptionID string) (*entities_subscriptions_v1.Subscription, error) {
	key := generateSubscriptionCacheKeyByIDForUser(userID, subscriptionID)

	cachedSubscription, err := s.cache.Get(ctx, key)
	if err == nil {
		var subscription *entities_subscriptions_v1.Subscription
		err = json.Unmarshal([]byte(cachedSubscription), &subscription)
		if err != nil {
			log.Error().Err(err).
				Msg("service.v1.service.GetSubscriptionByID: unable to unmarshal subscription")
		} else {
			return subscription, nil
		}
	}

	subscription, err := s.store.GetSubscriptionByID(ctx, userID, subscriptionID)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(subscription)
	if err != nil {
		log.Error().Err(err).
			Msg("service.v1.service.GetSubscriptionByID: unable to marshal subscription")
	} else {
		s.cache.SetEx(ctx, key, bytes, subscriptionCacheDuration)
	}

	return subscription, nil
}

func (s *Service) FetchSubscriptions(ctx context.Context, userID string) ([]*entities_subscriptions_v1.Subscription, error) {
	key := generateSubscriptionsCacheKeyForUser(userID)

	cachedSubscriptions, err := s.cache.Get(ctx, key)
	if err == nil {
		var subscriptions []*entities_subscriptions_v1.Subscription
		err = json.Unmarshal([]byte(cachedSubscriptions), &subscriptions)
		if err != nil {
			log.Error().Err(err).
				Msg("service.v1.service.FetchSubscriptions: unable to unmarshal subscriptions")
		} else {
			return subscriptions, nil
		}
	}

	subscriptions, err := s.store.FetchSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(subscriptions)
	if err != nil {
		log.Error().Err(err).
			Msg("service.v1.service.FetchSubscriptions: unable to marshal subscriptions")
	} else {
		s.cache.SetEx(ctx, key, bytes, subscriptionCacheDuration)
	}

	return subscriptions, nil
}

func (s *Service) FinishSubscription(ctx context.Context, userID string, subscriptionID string, finishedAt time.Time) error {
	err := s.store.FinishSubscription(ctx, userID, subscriptionID, finishedAt)
	if err != nil {
		return err
	}

	s.cache.Del(ctx, generateSubscriptionsCacheKeyForUser(userID))
	s.cache.Del(ctx, generateSubscriptionCacheKeyByIDForUser(userID, subscriptionID))

	return nil
}

func (s *Service) DeleteSubscription(ctx context.Context, userID string, subscriptionID string) error {
	err := s.store.DeleteSubscription(ctx, userID, subscriptionID)
	if err != nil {
		return err
	}

	s.cache.Del(ctx, generateSubscriptionsCacheKeyForUser(userID))
	s.cache.Del(ctx, generateSubscriptionCacheKeyByIDForUser(userID, subscriptionID))

	return nil
}
