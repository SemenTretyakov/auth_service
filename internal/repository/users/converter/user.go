package converter

import (
	"github.com/SemenTretyakov/auth_service/internal/model"
	modelRepo "github.com/SemenTretyakov/auth_service/internal/repository/users/model"
)

func RepoUserFieldsToDomain(userFields *modelRepo.UserFields) *model.UserFields {
	return &model.UserFields{
		Name:            userFields.Name,
		Email:           userFields.Email,
		Password:        userFields.Password,
		PasswordConfirm: userFields.PasswordConfirm,
		Role:            userFields.Role,
	}
}

// repo -> domain
func RepoUserToDomain(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
