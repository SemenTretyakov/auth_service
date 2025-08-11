package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/SemenTretyakov/auth_service/internal/config"
	desc "github.com/SemenTretyakov/auth_service/pkg/user_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateReq) (*desc.CreateRes, error) {
	pass := gofakeit.Password(true, true, true, true, false, 5)
	buildInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("fullname", "email", "password", "password_confirm", "role").
		Values(gofakeit.Name(), gofakeit.Email(), pass, pass, 1).
		Suffix("RETURNING id")

	query, args, err := buildInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with ID: %d", userID)

	return &desc.CreateRes{
		Id: userID,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetReq) (*desc.GetRes, error) {
	buildSelectOne := sq.Select(
		"id",
		"fullname",
		"email",
		"role",
		"created_at",
		"updated_at",
	).
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := buildSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var (
		id        int64
		fullname  string
		email     string
		role      int32
		createdAt time.Time
		updatedAt sql.NullTime
	)

	err = s.pool.QueryRow(ctx, query, args...).Scan(
		&id,
		&fullname,
		&email,
		&role,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		log.Printf("failed to select users: %v", err)
	}

	var updatedAtProto *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtProto = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetRes{
		User: &desc.User{
			Id:        id,
			Name:      fullname,
			Email:     email,
			Role:      desc.Role(role),
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: updatedAtProto,
		},
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateReq) (*empty.Empty, error) {
	if req.GetId() == 0 {
		return &empty.Empty{}, status.Error(codes.InvalidArgument, "User ID is required")
	}

	log.Printf("Updating user with id: %d", req.GetId())
	return &empty.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteReq) (*empty.Empty, error) {
	if req.GetId() == 0 {
		return &empty.Empty{}, status.Error(codes.InvalidArgument, "User ID is required")
	}

	log.Printf("Deleting user with id: %d", req.GetId())
	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to load grpcConfig: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to load pgConfig: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}
}
