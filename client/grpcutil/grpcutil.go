package grpcutil

import (
	"fmt"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/paschalolo/grpc/client/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

const (
	CA_CERT        = "../certs/ca_cert.pem"
	DOCKER_CA_CERT = "ca_cert.pem"
)

func ServiceConnection(addr string) *grpc.ClientConn {
	// run locally
	creds, err := credentials.NewClientTLSFromFile(CA_CERT, "x.test.example.com")
	// creds, err := credentials.NewClientTLSFromFile(DOCKER_CA_CERT, "x.test.example.com")

	if err != nil {
		log.Fatalf("failed to load credentials : %v ", err)
	}
	retryOpts := []retry.CallOption{
		retry.WithMax(3),
		retry.WithBackoff(retry.BackoffExponential(100 * time.Millisecond)),
		retry.WithCodes(codes.Unavailable),
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		grpc.WithChainUnaryInterceptor(
			retry.UnaryClientInterceptor(retryOpts...),
			interceptor.UnaryAuthInterceptor),
		grpc.WithStreamInterceptor(interceptor.StreamAuthInterceptor),
	}
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%v", addr), opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}
