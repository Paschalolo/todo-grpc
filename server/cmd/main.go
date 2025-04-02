package main

import (
	"log"
	"net"
	"os"

	pb "github.com/paschalolo/grpc/proto/todo/v1"
	"github.com/paschalolo/grpc/server/controller"
	grpcHandler "github.com/paschalolo/grpc/server/handler/grpc"
	"github.com/paschalolo/grpc/server/repository/memory"
	"google.golang.org/grpc"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("usage: server [IP_ADDR]")
	}
	addr := args[0]
	lis, err := net.Listen("tcp", addr)
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
	todoHandler := grpcHandler.NewHandler(ctrl)
	log.Printf("listening at %s\n", addr)
	opt := []grpc.ServerOption{}
	s := grpc.NewServer(opt...)
	pb.RegisterTodoServiceServer(s, todoHandler)
	//registration of endpoints
	defer s.Stop()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}

}
