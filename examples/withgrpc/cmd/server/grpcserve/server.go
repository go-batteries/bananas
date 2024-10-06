package grpcserve

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"testgrpcsetup/app/core/hellow/protos"
	"testgrpcsetup/app/web/middlewares"
	"testgrpcsetup/pkg/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcWithEchoServer struct{}

func (srv *GrpcWithEchoServer) StartHTTPServer(ctx context.Context, cfg *config.AppConfig) {
	port, grpcPort := cfg.AppPort, cfg.GrpcPort

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{fmt.Sprintf("http://localhost:%s", port)},
		AllowCredentials: true,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderContentDisposition,
			echo.HeaderConnection,
			echo.HeaderCacheControl,
			// Access Token Headers,
		},
		ExposeHeaders: []string{
			echo.HeaderContentLength,
			echo.HeaderContentDisposition,
			echo.HeaderContentEncoding,
			echo.HeaderContentType,
			echo.HeaderCacheControl,
			echo.HeaderConnection,
			// Access Token Headers,
		},
	}))

	// Start gRPC Gateway and Echo HTTP server concurrently
	go func() {
		// Dial the gRPC server
		conn, err := grpc.NewClient(
			fmt.Sprintf("localhost:%s", grpcPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to gRPC server")
		}
		defer conn.Close()

		// Create the HTTP gRPC Gateway mux
		gwMux := runtime.NewServeMux()

		// TODO: register the client in the app/web/router.go
		// This is an example
		client := protos.NewHellowServiceClient(conn)
		if err := protos.RegisterHellowServiceHandlerClient(ctx, gwMux, client); err != nil {
			log.Fatal().Err(err).Msg("Failed to register gRPC Gateway")
		}

		// Combine Echo with the gRPC Gateway
		e.Any("/*", echo.WrapHandler(gwMux))

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: e,
		}

		// Start the Echo HTTP server
		log.Info().Str("port", port).Msg("HTTP server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to serve HTTP")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Info().Msg("Shutting down HTTP server...")

	if err := e.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to gracefully shutdown HTTP server")
	}
}

// Start the gRPC server
func (srv *GrpcWithEchoServer) StartgRPCServer(ctx context.Context, cfg *config.AppConfig) {

	grpcPort := cfg.GrpcPort
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen on gRPC port")
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middlewares.GRPCRecoverMiddleware()))

	// Register all services in app/web/routers.go
	// This is an example
	protos.RegisterHellowServiceServer(grpcServer, &HellowServiceImpl{})

	go func() {
		log.Info().Str("port", grpcPort).Msg("gRPC server started")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("Failed to serve gRPC")
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown of gRPC server
	grpcServer.GracefulStop()
	log.Info().Msg("gRPC server stopped")
}

// Implement your gRPC HellowServiceImpl
type HellowServiceImpl struct {
	protos.UnimplementedHellowServiceServer
}

// TODO: these go inside app/core/hellow/hellow.controllers.go
// Implement the gRPC HellowService methods...
func (s *HellowServiceImpl) GetHellow(ctx context.Context, req *protos.HellowGetRequest) (*protos.HellowResponse, error) {
	name := req.GetName()
	if name == "" {
		return &protos.HellowResponse{
			Success: false,
			Data:    "",
			Errors:  []string{"Name is empty"},
		}, nil
	}
	return &protos.HellowResponse{
		Success: true,
		Data:    fmt.Sprintf("Hi, %s", name),
		Errors:  nil,
	}, nil
}

func (s *HellowServiceImpl) DeleteHellow(ctx context.Context, req *protos.HellowGetRequest) (*protos.HellowResponse, error) {
	name := req.GetName()
	if name == "" {
		return &protos.HellowResponse{
			Success: false,
			Data:    "",
			Errors:  []string{"Name is empty"},
		}, nil
	}
	return &protos.HellowResponse{
		Success: true,
		Data:    fmt.Sprintf("Bye, %s", name),
		Errors:  nil,
	}, nil
}

func (s *HellowServiceImpl) PostHellow(ctx context.Context, req *protos.HellowRequest) (*protos.HellowResponse, error) {
	name := req.Name
	if name == "" {
		return &protos.HellowResponse{
			Success: false,
			Data:    "",
			Errors:  []string{"Name is empty"},
		}, nil
	}
	return &protos.HellowResponse{
		Success: true,
		Data:    fmt.Sprintf("Hello, %s", name),
		Errors:  nil,
	}, nil
}
