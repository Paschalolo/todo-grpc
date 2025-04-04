package grpcutil

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

func ServiceConnection(addr string) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%v", addr), opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}
