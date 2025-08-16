package users

import (
	"context"

	"github.com/SemenTretyakov/auth_service/internal/model"
)

func (s *srv) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
