package grpcutil

import (
	"fmt"
	"log"

	"github.com/paschalolo/grpc/client/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

const (
	CA_CERT     = "../certs/ca_cert.pem"
	CLIENT_CERT = "../certs/server_cert.pem"
	CLIENT_KEY  = "../certs/server_key.pem"
)

func ServiceConnection(addr string) *grpc.ClientConn {
	creds, err := credentials.NewClientTLSFromFile(CA_CERT, "x.test.example.com")
	if err != nil {
		log.Fatalf("failed to load credentials : %v ", err)
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		grpc.WithUnaryInterceptor(interceptor.UnaryAuthInterceptor),
		grpc.WithStreamInterceptor(interceptor.StreamAuthInterceptor),
	}
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%v", addr), opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}
