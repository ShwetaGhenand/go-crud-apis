package server

import (
	"fmt"
	"log"
	"net"

	userspb "github.com/go-crud-apis/services/grpc/gen/proto"
	users "github.com/go-crud-apis/services/grpc/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func run(conf *config) error {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	db, err := InitDB(conf.DBConfig)
	if err != nil {
		return err
	}

	var service *users.UserService
	service, err = users.NewUserService(db)
	if err != nil {
		return err
	}
	userspb.RegisterUserServiceServer(grpcServer, service)
	grpcAddr := fmt.Sprintf(":%d", 9000)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Println("Failed to listen on port ", err)
	}

	// Serve gRPC Server
	log.Println("Serving gRPC on ")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to connect to server %v", err)
	}

	return nil
}
