package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/SemenTretyakov/auth_service/internal/config"
	defConverter "github.com/SemenTretyakov/auth_service/internal/converter"
	usersRepo "github.com/SemenTretyakov/auth_service/internal/repository/users"
	"github.com/SemenTretyakov/auth_service/internal/service"
	userService "github.com/SemenTretyakov/auth_service/internal/service/users"
	desc "github.com/SemenTretyakov/auth_service/pkg/user_v1"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	usersService service.UsersService
	desc.UnimplementedUserV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateReq) (*desc.CreateRes, error) {
	userID, err := s.usersService.Create(ctx, defConverter.UserFieldsFromProto(req.GetInfo()))
	if err != nil {
		log.Printf("Error from repo.Create: %v\n", err)
		return nil, err
	}

	log.Printf("inserted user with ID: %d", userID)

	return &desc.CreateRes{
		Id: userID,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetReq) (*desc.GetRes, error) {
	user, err := s.usersService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetRes{
		User: defConverter.UserToProto(user),
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

	lis, err := net.Listen("tcp4", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	usersRepository := usersRepo.NewRepository(pool)
	userService := userService.NewService(usersRepository)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{usersService: userService})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}
}
