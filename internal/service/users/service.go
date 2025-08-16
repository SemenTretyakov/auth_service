package users

import (
	"github.com/SemenTretyakov/auth_service/internal/repository"
	"github.com/SemenTretyakov/auth_service/internal/service"
)

type srv struct {
	userRepo repository.UsersRepository
}

func NewService(userRepo repository.UsersRepository) service.UsersService {
	return &srv{userRepo: userRepo}
}
