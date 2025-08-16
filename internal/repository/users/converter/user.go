package converter

import (
	"github.com/SemenTretyakov/auth_service/internal/repository/users/model"
	desc "github.com/SemenTretyakov/auth_service/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// repo -> proto
func RepoUserFieldsToDesc(userFields *model.UserFields) *desc.UserFields {
	role := desc.Role_USER
	switch userFields.Role {
	case 0:
		role = desc.Role_USER
	case 1:
		role = desc.Role_ADMIN
	}

	return &desc.UserFields{
		Name:            userFields.Name,
		Email:           userFields.Email,
		Password:        userFields.Password,
		PasswordConfirm: userFields.PasswordConfirm,
		Role:            role,
	}
}

// repo -> proto
func RepoUserToDesc(user *model.User) *desc.User {
	role := desc.Role_USER
	switch user.Role {
	case 0:
		role = desc.Role_USER
	case 1:
		role = desc.Role_ADMIN
	}

	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      role,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
