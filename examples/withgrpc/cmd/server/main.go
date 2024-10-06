package main

import (
	"context"
	
	server "testgrpcsetup/cmd/server/grpcserve"
	
	"testgrpcsetup/pkg/config"
)

// Main entry point
func main() {
	cfg := config.Load("./config/")

	// Set up context for graceful shutdown
	appCtx := context.Background()
	ctx, cancel := context.WithCancel(appCtx)
	defer cancel()

	
	srv := &server.GrpcWithEchoServer{}
	// Start the gRPC server
	go srv.StartgRPCServer(ctx, &cfg)

	// Start the Echo server with gRPC Gateway
	srv.StartHTTPServer(ctx, &cfg)
	
}
