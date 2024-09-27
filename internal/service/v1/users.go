package service_v1

import (
	"context"

	entities_users_v1 "github.com/subeecore/subee-core-svc/internal/entities/users/v1"
)

func (s *Service) CreateUser(ctx context.Context, req *entities_users_v1.CreateUserRequest) (*entities_users_v1.User, error) {
	user, err := s.store.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}
