package server

import (
	"context"
	"fmt"
	"net"

	userspb "github.com/go-crud-apis/services/grpc/gen/go/proto"
	users "github.com/go-crud-apis/services/grpc/users"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Test() {
	fmt.Println("test")
}

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
	g, ctx := errgroup.WithContext(context.Background())

	// Serve gRPC Server
	g.Go(func() error {
		userspb.RegisterUserServiceServer(grpcServer, service)
		grpcAddr := fmt.Sprintf(":%d", conf.GRPCPort)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			return err
		}
		if serverErr := grpcServer.Serve(lis); serverErr != nil {
			return serverErr
		}
		return nil
	})

	// Connect to the gRPC server and start gateway
	g.Go(func() error {
		return stratGRPCGatewayServer(ctx, conf.GRPCPort, conf.GRPCGatewayPort)
	})
	return g.Wait()
}
