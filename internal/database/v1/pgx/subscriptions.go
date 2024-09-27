package database_v1_pgx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/subeecore/pkg/constants"
	"github.com/subeecore/pkg/errors"

	entities_subscriptions_v1 "github.com/subeecore/subee-core-svc/internal/entities/subscriptions/v1"
)

func (d *dbClient) CreateSubscription(ctx context.Context, req *entities_subscriptions_v1.CreateSubscriptionRequest) (*entities_subscriptions_v1.Subscription, error) {
	subscriptionID := constants.GenerateDataPrefixWithULID(constants.Subscription)
	now := time.Now()

	_, err := d.connection.DB.ExecContext(ctx,
		`INSERT INTO 
			subscriptions (
				id,
				user_id,
				platform,
				reccurence,
				price,
				started_at,
				created_at,
			) 
			VALUES ($1, $2, $3, $4, $5, $6, $7);
		`,
		subscriptionID, req.UserID, req.Platform, req.Reccurence, req.Price, req.StartedAt, now)
	if err != nil {
		log.Error().Err(err).
			Msgf("failed to create user subscription: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("failed to create user subscription: %v", err.Error()))
	}

	return &entities_subscriptions_v1.Subscription{
		ID:         subscriptionID,
		UserID:     req.UserID,
		Platform:   req.Platform,
		Reccurence: req.Reccurence,
		Price:      req.Price,
		StartedAt:  req.StartedAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (d *dbClient) FetchSubscriptions(ctx context.Context, userID string) ([]*entities_subscriptions_v1.Subscription, error) {
	rows, err := d.connection.DB.QueryContext(ctx, `
		SELECT
			id,
			user_id,
			platform,
			reccurence,
			price,
			started_at,
			created_at,
			updated_at,
			finished_at
		FROM
			subscriptions
		WHERE
			user_id = $1
	`, userID)
	if err != nil {
		log.Error().Err(err).
			Str("user_id", userID).
			Msgf("database.postgres.dbClient.FetchSubscriptions: failed to get subscriptions: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.FetchSubscriptions: failed to get subscriptions: %v", err.Error()))
	}
	defer rows.Close()

	subscriptions := make([]*entities_subscriptions_v1.Subscription, 0)
	for rows.Next() {
		subscription := &entities_subscriptions_v1.Subscription{}

		err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.Platform,
			&subscription.Reccurence,
			&subscription.Price,
			&subscription.StartedAt,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
			&subscription.FinishedAt,
		)
		if err != nil {
			log.Error().Err(err).
				Str("user_id", userID).
				Msgf("database.postgres.dbClient.FetchSubscriptions: failed to scan subscription: %v", err.Error())
			return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.FetchSubscriptions: failed to scan subscription: %v", err.Error()))
		}

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

func (d *dbClient) GetSubscriptionByID(ctx context.Context, userID string, subscriptionID string) (*entities_subscriptions_v1.Subscription, error) {
	subscription := &entities_subscriptions_v1.Subscription{}

	err := d.connection.DB.QueryRowContext(ctx, `
		SELECT 
			id,
			user_id,
			platform,
			reccurence,
			price,
			started_at,
			created_at,
			updated_at,
			finished_at
		FROM subscriptions
		WHERE 
			user_id = $1 AND 
			game_id = $2
	`, userID, subscriptionID).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.Platform,
		&subscription.Reccurence,
		&subscription.Price,
		&subscription.StartedAt,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&subscription.FinishedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).
				Str("user_id", userID).
				Msgf("database.postgres.dbClient.GetSubscriptionByID: subscription not found")
			return nil, errors.NewNotFoundError("database.postgres.dbClient.GetSubscriptionByID: subscription not found")
		}

		log.Error().Err(err).
			Str("user_id", userID).
			Str("subscription_id", subscriptionID).
			Msgf("database.postgres.dbClient.GetSubscriptionByID: failed to get subscription: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.GetSubscriptionByID: failed to get subscription: %v", err.Error()))
	}

	return subscription, nil
}

func (d *dbClient) FinishSubscription(ctx context.Context, userID string, subscriptionID string, finishedAt time.Time) error {
	_, err := d.connection.DB.ExecContext(ctx, `
		UPDATE 
			subscriptions
		SET 
			finished_at = $3
		WHERE 
			id = $1 AND 
			user_id = $2
	`, subscriptionID, userID, finishedAt)
	if err != nil {
		log.Error().Err(err).
			Str("subscription_id", subscriptionID).
			Str("user_id", userID).
			Msgf("database.postgres.dbClient.FinishSubscription: failed to finish subscription: %v", err.Error())
		return errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.FinishSubscription: failed to finish subscription: %v", err.Error()))
	}

	return nil
}

func (d *dbClient) DeleteSubscription(ctx context.Context, userID string, subscriptionID string) error {
	_, err := d.connection.DB.ExecContext(ctx, `
		DELETE FROM 
			subscriptions
		WHERE 
			id = $1 AND 
			user_id = $2
	`, subscriptionID, userID)
	if err != nil {
		log.Error().Err(err).
			Str("subscription_id", subscriptionID).
			Str("user_id", userID).
			Msgf("database.postgres.dbClient.DeleteSubscription: failed to delete subscription: %v", err.Error())
		return errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.DeleteSubscription: failed to delete subscription: %v", err.Error()))
	}

	return nil
}
