package main

import (
	"log"

	"github.com/AntonRadchenko/users-service/internal/database"
	"github.com/AntonRadchenko/users-service/internal/transport/grpc"
	"github.com/AntonRadchenko/users-service/internal/user"
)

func main() {
	database.InitDB()

	repo := &user.UserRepo{}
	svc := user.NewUserService(repo)

	if err := grpc.RunGRPC(svc); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}