package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/paschalolo/grpc/proto/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func addTask(c pb.TodoServiceClient, description string, dueDate time.Time) (uint64, error) {
	ctx := context.Background()
	req := &pb.AddTaskRequest{
		Description: description,
		DueDate:     timestamppb.New(dueDate),
	}
	res, err := c.AddTask(ctx, req)
	if err != nil {
		return 0, err
	}
	return res.Id, nil
}

// func getTasks(c pb.TodoServiceClient) (*pb.ListTasksResponse, error) {
// 	ctx := context.Background()

// }

func main() {
	args := os.Args
	if len(args) == 0 {
		log.Fatalln("usage client [IP_ADDR] ")
	}
	addr := args[0]
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewTodoServiceClient(conn)
	fmt.Println("--------------ADD--------------")
	dueDate := time.Now().Add(5 * time.Second)
	id1, err := addTask(c, "this is a task ", dueDate)
	if err != nil {
		log.Fatalln("add task failed ")
	}
	fmt.Printf("added task %d", id1)
	fmt.Println("-------------------------------")

	defer func(client *grpc.ClientConn) {
		if err := client.Close(); err != nil {
			log.Fatalf("unexpected error %v \n", err)
		}
	}(conn)

}
