package repository

import (
	"context"

	desc "github.com/SemenTretyakov/auth_service/pkg/user_v1"
)

type UsersRepository interface {
	Create(ctx context.Context, info *desc.UserFields) (int64, error)
	Get(ctx context.Context, id int64) (*desc.User, error)
}