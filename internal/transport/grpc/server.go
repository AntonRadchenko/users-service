package grpc

import (
	"net"

	"github.com/AntonRadchenko/users-service/internal/user"
	"google.golang.org/grpc"
	userpb "github.com/AntonRadchenko/project-protos/proto/user"
)


func RunGRPC(svc *user.UserService) error {
	// 1. net.Listen на ":50051"
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	// 2. grpc.NewServer()
	grpcServer := grpc.NewServer()

	// 3. userpb.RegisterUserServiceServer(grpcSrv, NewHandler(svc))
	handler := NewHandler(svc)
	userpb.RegisterUserServiceServer(grpcServer, handler)

	// 4. grpcSrv.Serve(listener) (блокируется)
	return grpcServer.Serve(listener)
}