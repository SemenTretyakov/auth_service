package users

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/SemenTretyakov/auth_service/internal/model"
	"github.com/SemenTretyakov/auth_service/internal/repository"
	"github.com/SemenTretyakov/auth_service/internal/repository/users/converter"
	modelRepo "github.com/SemenTretyakov/auth_service/internal/repository/users/model"
)

const (
	tableName = "users"

	idColumn              = "id"
	fullnameColumn        = "fullname"
	emailColumn           = "email"
	passwordColumn        = "password"
	passwordConfirmColumn = "password_confirm"
	roleColumn            = "role"
	createdAtColumn       = "created_at"
	updatedAtColumn       = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) repository.UsersRepository {
	return &repo{db: pool}
}

func (r *repo) Create(ctx context.Context, info *model.UserFields) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(fullnameColumn, emailColumn, passwordColumn, passwordConfirmColumn, roleColumn).
		Values(
			info.Name,
			info.Email,
			info.Password,
			info.PasswordConfirm,
			info.Role,
		).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	log.Printf("Query: %s, args: %+v\n", query, args)

	var userID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Println("Error executing query:", err)
		return 0, err
	}

	return userID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	buildSelectOne := sq.
		Select(
			idColumn,
			fullnameColumn,
			emailColumn,
			roleColumn,
			createdAtColumn,
			updatedAtColumn,
		).
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := buildSelectOne.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	var user modelRepo.User

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Printf("failed to select users: %v", err)
	}

	return converter.RepoUserToDomain(&user), nil
}
