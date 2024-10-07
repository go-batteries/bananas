package main

import (
	"context"
	
	server "openapivvv/cmd/server/httpserve"
	
	"openapivvv/pkg/config"
)

// Main entry point
func main() {
	cfg := config.Load("./config/")

	// Set up context for graceful shutdown
	appCtx := context.Background()
	ctx, cancel := context.WithCancel(appCtx)
	defer cancel()

	
	srv := &server.EchoServer{}
	srv.StartHTTPServer(ctx, &cfg)
	
}
