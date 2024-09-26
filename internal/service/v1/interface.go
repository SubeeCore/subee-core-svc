package service_v1

import (
	"context"

	entities_user_v1 "github.com/subeecore/subee-core-svc/internal/entities/user/v1"
)

type AlyxCoreService interface {
	CreateUser(ctx context.Context, req *entities_user_v1.CreateUserRequest) (*entities_user_v1.User, error)
}
