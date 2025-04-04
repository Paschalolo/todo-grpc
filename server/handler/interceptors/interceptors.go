package interceptors

import (
	"context"

	"github.com/paschalolo/grpc/server/middleware"
	"google.golang.org/grpc"
)

func UnaryAuthinterceptors(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if err := middleware.ValidateAuthToken(ctx); err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func StreamAuthInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := middleware.ValidateAuthToken(ss.Context()); err != nil {
		return err
	}
	return handler(srv, ss)
}
