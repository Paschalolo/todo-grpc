package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/paschalolo/grpc/proto/todo/v2"
	"github.com/paschalolo/grpc/server/controller"
	grpcHandler "github.com/paschalolo/grpc/server/handler/grpc"
	"github.com/paschalolo/grpc/server/handler/interceptors"
	"github.com/paschalolo/grpc/server/repository/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
)

const (
	CA_CERT     = "../certs/ca_cert.pem"
	SERVER_CERT = "../certs/server_cert.pem"
	SERVER_KEY  = "../certs/server_key.pem"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("usage: server [IP_ADDR]")
	}
	addr := args[0]
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", addr))
	if err != nil {
		log.Fatalf("failed to listen : %v \n ", err)
	}
	defer func(lis net.Listener) {
		if err := lis.Close(); err != nil {
			log.Fatalf("unexpected error %v ", err)
		}

	}(lis)
	inMemory := memory.New()
	ctrl := controller.NewController(inMemory)
	h := grpcHandler.NewHandler(ctrl)
	creds, err := credentials.NewServerTLSFromFile(SERVER_CERT, SERVER_KEY)

	if err != nil {
		log.Fatalf("failed to create credentials :%v", err)
	}
	opt := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.UnaryInterceptor(interceptors.UnaryAuthinterceptors),
		grpc.StreamInterceptor(interceptors.StreamAuthInterceptor),
	}
	srv := grpc.NewServer(opt...)
	reflection.Register(srv)
	pb.RegisterTodoServiceServer(srv, h)
	log.Printf("listening at %s\n", addr)
	//registration of endpoints
	defer srv.Stop()

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}

}
