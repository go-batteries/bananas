package middlewares

import (
	"context"
	"log"
	"runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCRecoverMiddleware is a gRPC middleware that recovers from panics.
func GRPCRecoverMiddleware() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic and the stack trace.
				log.Printf("panic recovered: %v", r)
				stackTrace := make([]byte, 4096)
				stackSize := runtime.Stack(stackTrace, true)
				log.Printf("stack trace:\n%s", stackTrace[:stackSize])

				// Return a gRPC error instead of crashing the server.
				err = status.Errorf(codes.Internal, "Internal server error")
			}
		}()

		// Call the handler
		resp, err = handler(ctx, req)
		return
	}
}
