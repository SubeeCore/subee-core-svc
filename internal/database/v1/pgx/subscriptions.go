package database_v1_pgx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/subeecore/pkg/constants"
	"github.com/subeecore/pkg/errors"

	entities_categories_v1 "github.com/subeecore/subee-core-svc/internal/entities/categories/v1"
	entities_payments_v1 "github.com/subeecore/subee-core-svc/internal/entities/payments/v1"
	entities_recap_v1 "github.com/subeecore/subee-core-svc/internal/entities/recap/v1"
	entities_subscriptions_v1 "github.com/subeecore/subee-core-svc/internal/entities/subscriptions/v1"
)

func (d *dbClient) CreateSubscription(ctx context.Context, req *entities_subscriptions_v1.CreateSubscriptionRequest) (*entities_subscriptions_v1.Subscription, error) {
	subscriptionID := constants.GenerateDataPrefixWithULID(constants.Subscription)
	now := time.Now()

	date, err := time.Parse(time.DateOnly, req.StartedAt)
	if err != nil {
		log.Error().Err(err).
			Msgf("database.postgres.dbClient.CreateSubscription: failed to parse user subscription started_at: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.CreateSubscription: failed to parse user subscription started_at: %v", err.Error()))
	}

	_, err = d.connection.DB.ExecContext(ctx,
		`INSERT INTO 
			subscriptions (
				id,
				user_id,
				platform,
				category,
				reccurence,
				price,
				started_at,
				created_at,
				updated_at
			) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
		subscriptionID, req.UserID, req.Platform, req.Category, req.Reccurence, req.Price, date, now, now)
	if err != nil {
		log.Error().Err(err).
			Msgf("database.postgres.dbClient.CreateSubscription: failed to create user subscription: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.CreateSubscription: failed to create user subscription: %v", err.Error()))
	}

	return &entities_subscriptions_v1.Subscription{
		ID:         subscriptionID,
		UserID:     req.UserID,
		Platform:   req.Platform,
		Reccurence: req.Reccurence,
		Price:      req.Price,
		StartedAt:  date,
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
			category,
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
			&subscription.Category,
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
			category,
			reccurence,
			price,
			started_at,
			created_at,
			updated_at,
			finished_at
		FROM subscriptions
		WHERE 
			user_id = $1 AND 
			id = $2
	`, userID, subscriptionID).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.Platform,
		&subscription.Category,
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

func (d *dbClient) GetMonthlySubscriptionsRecap(ctx context.Context, userID string) (*entities_recap_v1.MonthlyRecap, error) {
	rows, err := d.connection.DB.QueryContext(ctx, `
		SELECT 
			id,
			platform,
			category,
			reccurence,
			price,
			started_at
		FROM 
			subscriptions
		WHERE
			(started_at <= CURRENT_DATE AND finished_at IS NULL) AND 
			user_id = $1
	`, userID)
	if err != nil {
		log.Error().Err(err).
			Str("user_id", userID).
			Msgf("database.postgres.dbClient.GetMonthlySubscriptionsPrice: failed to get subscriptions: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.GetMonthlySubscriptionsPrice: failed to get subscriptions: %v", err.Error()))
	}
	defer rows.Close()

	payments := make([]*entities_payments_v1.Payment, 0)
	for rows.Next() {
		payment := &entities_payments_v1.Payment{}

		err := rows.Scan(
			&payment.SubscriptionID,
			&payment.Platform,
			&payment.Category,
			&payment.Reccurence,
			&payment.Price,
			&payment.StartedAt,
		)
		if err != nil {
			log.Error().Err(err).
				Str("user_id", userID).
				Msgf("database.postgres.dbClient.GetMonthlySubscriptionsPrice: failed to scan payment: %v", err.Error())
			return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.GetMonthlySubscriptionsPrice: failed to scan payment: %v", err.Error()))
		}

		payments = append(payments, payment)
	}

	currentPayments, err := checkPayments(payments)
	if err != nil {
		log.Error().Err(err).
			Str("user_id", userID).
			Msgf("database.postgres.dbClient.GetMonthlySubscriptionsPrice: failed to check payments: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.GetMonthlySubscriptionsPrice: failed to check payments: %v", err.Error()))
	}

	price, err := getTotalPrice(currentPayments)
	if err != nil {
		log.Error().Err(err).
			Str("user_id", userID).
			Msgf("database.postgres.dbClient.GetMonthlySubscriptionsPrice: failed to get total price: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.GetMonthlySubscriptionsPrice: failed to get total price: %v", err.Error()))
	}

	categories, err := getCategoriesPercentage(payments)
	if err != nil {
		log.Error().Err(err).
			Str("user_id", userID).
			Msgf("database.postgres.dbClient.GetMonthlySubscriptionsPrice: failed to get categories: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.GetMonthlySubscriptionsPrice: failed to get categories: %v", err.Error()))
	}

	return &entities_recap_v1.MonthlyRecap{
		Price:      price,
		Payments:   payments,
		Categories: categories,
	}, nil
}

func (d *dbClient) GetGlobalSubscriptionsRecap(ctx context.Context, userID string) (*entities_recap_v1.GlobalRecap, error) {
	rows, err := d.connection.DB.QueryContext(ctx, `
		SELECT 
			price,
			started_at
		FROM 
			subscriptions
		WHERE
			user_id = $1
	`, userID)
	if err != nil {
		log.Error().Err(err).
			Str("user_id", userID).
			Msgf("database.postgres.dbClient.GetGlobalSubscriptionsRecap: failed to get subscriptions: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.GetGlobalSubscriptionsRecap: failed to get subscriptions: %v", err.Error()))
	}
	defer rows.Close()

	payments := make([]*entities_payments_v1.Payment_Light, 0)
	for rows.Next() {
		payment := &entities_payments_v1.Payment_Light{}

		err := rows.Scan(
			&payment.Price,
			&payment.StartedAt,
		)
		if err != nil {
			log.Error().Err(err).
				Str("user_id", userID).
				Msgf("database.postgres.dbClient.GetGlobalSubscriptionsRecap: failed to scan payment: %v", err.Error())
			return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.GetGlobalSubscriptionsRecap: failed to scan payment: %v", err.Error()))
		}

		payments = append(payments, payment)
	}

	_, err = groupByYear(payments)
	if err != nil {
		log.Error().Err(err).
			Str("user_id", userID).
			Msgf("database.postgres.dbClient.GetGlobalSubscriptionsRecap: failed to group by year: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.GetGlobalSubscriptionsRecap: failed to group by year: %v", err.Error()))
	}

	return nil, nil
}

func (d *dbClient) FinishSubscription(ctx context.Context, userID string, subscriptionID string, finishedAt string) error {
	date, err := time.Parse(time.DateOnly, finishedAt)
	if err != nil {
		log.Error().Err(err).
			Msgf("database.postgres.dbClient.FinishSubscription: failed to parse finished_at: %v", err.Error())
		return errors.NewInternalServerError(fmt.Sprintf("database.postgres.dbClient.FinishSubscription: failed to parse finished_at: %v", err.Error()))
	}

	_, err = d.connection.DB.ExecContext(ctx, `
		UPDATE 
			subscriptions
		SET 
			finished_at = $3
		WHERE 
			id = $1 AND 
			user_id = $2
	`, subscriptionID, userID, date)
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

func checkPayments(allPayments []*entities_payments_v1.Payment) ([]*entities_payments_v1.Payment, error) {
	payments := make([]*entities_payments_v1.Payment, 0)

	for _, payment := range allPayments {
		date := payment.StartedAt

		for date.Month() <= time.Now().Month() && date.Year() <= time.Now().Year() {
			if date.Month() == time.Now().Month() && date.Year() == time.Now().Year() {
				payments = append(payments, payment)
			}
			date = date.Add(time.Hour * time.Duration(24*payment.Reccurence))
		}
	}

	return payments, nil
}

func getTotalPrice(currentPayments []*entities_payments_v1.Payment) (float64, error) {
	price := 0.0

	for _, payment := range currentPayments {
		price += payment.Price
	}

	return price, nil
}

func groupByYear(allPayments []*entities_payments_v1.Payment_Light) (map[int][]*entities_recap_v1.MonthlyRecap_Light, error) {
	groupedByYear := make(map[int][]*entities_recap_v1.MonthlyRecap_Light)

	for _, payment := range allPayments {
		year := payment.StartedAt.Year()
		groupedByYear[year] = append(groupedByYear[year], &entities_recap_v1.MonthlyRecap_Light{
			Month: payment.StartedAt.Month().String(),
		})
	}

	log.Info().Msgf("groupedByYear: %v\n", groupedByYear)
	fmt.Printf("groupedByYear: %v\n", groupedByYear)

	return groupedByYear, nil
}

func getCategoriesPercentage(allPayments []*entities_payments_v1.Payment) (*entities_categories_v1.CategoriesRecap, error) {
	billsCategoriesLength := 0
	entertainmentCategoriesLength := 0

	for _, payment := range allPayments {
		if payment.Category == "bills" {
			billsCategoriesLength++
		} else if payment.Category == "entertainment" {
			entertainmentCategoriesLength++
		}
	}

	return &entities_categories_v1.CategoriesRecap{
		Categories: []*entities_categories_v1.CategoryRecap{
			{
				Name:       "bills",
				Percentage: fmt.Sprintf("%.2f", float64(billsCategoriesLength)/float64(len(allPayments))*100),
			},
			{
				Name:       "entertainment",
				Percentage: fmt.Sprintf("%.2f", float64(entertainmentCategoriesLength)/float64(len(allPayments))*100),
			},
		},
	}, nil
}
