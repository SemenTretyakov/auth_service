package users

import (
	"database/sql"
	"time"

	"github.com/SemenTretyakov/auth_service/internal/model"
	desc "github.com/SemenTretyakov/auth_service/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// proto -> domain
func UserFieldsFromProto(in *desc.UserFields) *model.UserFields {
	return &model.UserFields{
		Name:            in.Name,
		Email:           in.Email,
		Role:            int32(in.Role),
		Password:        in.Password,
		PasswordConfirm: in.PasswordConfirm,
	}
}

func UserFromProto(in *desc.User) *model.User {
	var updatedAt sql.NullTime
	var createdAt time.Time

	if in.UpdatedAt != nil {
		updatedAt = sql.NullTime{
			Time:  in.UpdatedAt.AsTime(),
			Valid: true,
		}
	}

	if in.CreatedAt.IsValid() {
		createdAt = in.CreatedAt.AsTime()
	}

	return &model.User{
		ID:        in.Id,
		Name:      in.Name,
		Email:     in.Email,
		Role:      int8(in.Role),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// domain -> proto
func UserFieldsToProto(in *model.UserFields) *desc.UserFields {
	return &desc.UserFields{
		Name:            in.Name,
		Email:           in.Email,
		Role:            desc.Role(in.Role),
		Password:        in.Password,
		PasswordConfirm: in.PasswordConfirm,
	}
}

func UserToProto(in *model.User) *desc.User {
	var createdAt *timestamppb.Timestamp
	var updatedAt *timestamppb.Timestamp

	if !in.CreatedAt.IsZero() {
		createdAt = timestamppb.New(in.CreatedAt)
	}

	if in.UpdatedAt.Valid {
		updatedAt = timestamppb.New(in.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        in.ID,
		Name:      in.Name,
		Email:     in.Email,
		Role:      desc.Role(in.Role),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
