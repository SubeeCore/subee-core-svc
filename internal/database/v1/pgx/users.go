package database_v1_pgx

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/subeecore/pkg/constants"
	"github.com/subeecore/pkg/errors"

	entities_users_v1 "github.com/subeecore/subee-core-svc/internal/entities/users/v1"
)

func (d *dbClient) CreateUser(ctx context.Context, req *entities_users_v1.CreateUserRequest) (*entities_users_v1.User, error) {
	userID := constants.GenerateDataPrefixWithULID(constants.User)
	now := time.Now()

	_, err := d.connection.DB.ExecContext(ctx,
		`INSERT INTO 
			users (
				id,
				external_id,
				username,
				email, 
				created_at, 
				updated_at
			) 
			VALUES ($1, $2, $3, $4, $5, $6);
		`,
		userID, req.ExternalID, req.Username, req.Email, now, now)
	if err != nil {
		log.Error().Err(err).
			Msgf("failed to create user: %v", err.Error())
		return nil, errors.NewInternalServerError(fmt.Sprintf("failed to create user: %v", err.Error()))
	}

	return &entities_users_v1.User{
		ID:         userID,
		ExternalID: req.ExternalID,
		Username:   req.Username,
		Email:      req.Email,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
