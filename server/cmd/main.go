package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	pb "github.com/paschalolo/grpc/proto/todo/v2"
	"github.com/paschalolo/grpc/server/controller"
	grpcHandler "github.com/paschalolo/grpc/server/handler/grpc"
	"github.com/paschalolo/grpc/server/middleware"
	"github.com/paschalolo/grpc/server/repository/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
)

const (
	SERVER_CERT        = "../certs/server_cert.pem"
	SERVER_KEY         = "../certs/server_key.pem"
	DOCKER_SERVER_CERT = "server_cert.pem"
	DOCKER_SERVER_KEY  = "server_key.pem"
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
	//run this line for docekr build
	// creds, err := credentials.NewServerTLSFromFile(DOCKER_SERVER_CERT, DOCKER_SERVER_KEY)
	creds, err := credentials.NewServerTLSFromFile(SERVER_CERT, SERVER_KEY)

	if err != nil {
		log.Fatalf("failed to create credentials :%v", err)
	}
	opt := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor(middleware.ValidateAuthToken)),
		grpc.StreamInterceptor(auth.StreamServerInterceptor(middleware.ValidateAuthToken)),
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
