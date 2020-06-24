package main

import (
	gw "bangumi/pb"
	"context"
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	// command-line options:
	// gRPC server endpoint
	infoEP      = flag.String("info-endpoint", "localhost:50052", "gRPC server endpoint")
	recommendEP = flag.String("favorite-endpoint", "localhost:50053", "gRPC server endpoint")
	userEP      = flag.String("user-endpoint", "localhost:50054", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := gw.RegisterInfoServiceHandlerFromEndpoint(ctx, mux, *infoEP, opts); err != nil {
		return err
	}

	if err := gw.RegisterFavoriteServiceHandlerFromEndpoint(ctx, mux, *recommendEP, opts); err != nil {
		return err
	}

	if err := gw.RegisterUserServiceHandlerFromEndpoint(ctx, mux, *userEP, opts); err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
