package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	userspb "github.com/go-crud-apis/services/grpc/gen/go/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func stratGRPCGatewayServer(ctx context.Context, proxyAdd, grpcGWAdd int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf(":%d", proxyAdd),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		return err
	}

	gwmux := runtime.NewServeMux()
	err = userspb.RegisterUserServiceHandler(ctx, gwmux, conn)
	if err != nil {
		return err
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", grpcGWAdd),
		Handler: gwmux,
	}

	log.Printf("Serving gRPC-Gateway on http://0.0.0.0:%d", grpcGWAdd)
	errch := make(chan error)
	select {
	case errch <- gwServer.ListenAndServe():
		return <-errch
	case <-ctx.Done():
		if err := gwServer.Close(); err != nil {
			log.Println(err)
		}
		return ctx.Err()
	default:
		return nil
	}
}
