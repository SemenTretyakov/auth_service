package service

import (
	"context"

	"github.com/SemenTretyakov/auth_service/internal/model"
)

type UsersService interface {
	Create(ctx context.Context, info *model.UserFields) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
}
