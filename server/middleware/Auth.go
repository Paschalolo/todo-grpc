package middleware

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const authTokenKey string = "auth_token"
const authTokenValue string = "authd"

func ValidateAuthToken(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	if t, ok := md["auth_token"]; ok {
		switch {
		case len(t) != 1:
			return status.Errorf(codes.InvalidArgument, "%s should contain only 1 value", authTokenKey)
		case t[0] != "authd":
			return status.Errorf(codes.Unauthenticated, "inncorect %s", authTokenValue)
		default:
			log.Println("authenticated")
		}
	} else {
		return status.Errorf(codes.Unauthenticated, "failed to get auth token ")
	}
	return nil
}
