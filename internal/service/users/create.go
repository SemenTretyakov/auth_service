package users

import (
	"context"

	"github.com/SemenTretyakov/auth_service/internal/model"
)

func (s *srv) Create(ctx context.Context, info *model.UserFields) (int64, error) {
	id, err := s.userRepo.Create(ctx, info)
	if err != nil {
		return 0, err
	}

	return id, nil
}
