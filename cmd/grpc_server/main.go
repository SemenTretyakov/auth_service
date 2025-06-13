package main

import (
	"context"
	"fmt"
	"log"
	"net"

	desc "github.com/SemenTretyakov/auth_service/pkg/user_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateReq) (*desc.CreateRes, error) {
	if req.GetInfo() == nil {
		log.Fatalf("Not provide data for create user")
	}

	info := req.GetInfo()

	if info.GetPassword() != info.GetPasswordConfirm() {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Password and confirmation do not match")
	}

	log.Printf("Creating user: %s, email: %s, role: %v",
		info.GetName(),
		info.GetEmail(),
		info.GetRole(),
	)

	return &desc.CreateRes{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetReq) (*desc.GetRes, error) {
	log.Printf("Getting user with id: %d", req.GetId())

	return &desc.GetRes{
		User: &desc.User{
			Id:        req.GetId(),
			Name:      gofakeit.Username(),
			Email:     gofakeit.Email(),
			Role:      desc.Role_USER,
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
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
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}
}
