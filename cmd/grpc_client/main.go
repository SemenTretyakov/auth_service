package main

import (
	"context"
	"log"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	desc "github.com/SemenTretyakov/auth_service/pkg/user_v1"
)

const (
	address = "localhost:50051"
	userID  = 32
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetReq{Id: userID})
	if err != nil {
		log.Fatalf("failed get user by id: %v", err)
	}

	log.Printf(color.RedString("User info:\n"), color.GreenString("%+v", r.GetUser()))
}
