package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/paschalolo/grpc/proto/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func addTask(c pb.TodoServiceClient, description string, dueDate time.Time) uint64 {
	ctx := context.Background()
	req := &pb.AddTaskRequest{
		Description: description,
		DueDate:     timestamppb.New(dueDate),
	}
	res, err := c.AddTask(ctx, req)
	if err != nil {
		panic(err)
	}
	return res.Id
}

func printTasks(c pb.TodoServiceClient) {
	req := &pb.ListTasksRequest{}
	stream, err := c.ListTasks(context.Background(), req)
	if err != nil {
		log.Fatalf("unexpected error %v", err)
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("unexpected error : %v", err)
		}
		fmt.Println(resp.Task.String(), "overdue:", resp.Overdue)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("usage client [IP_ADDR] ")
	}

	addr := args[0]

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%v", addr), opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewTodoServiceClient(conn)
	defer func(client *grpc.ClientConn) {
		if err := client.Close(); err != nil {
			log.Fatalf("unexpected error %v \n", err)
		}
	}(conn)
	fmt.Println("--------------ADD--------------")
	dueDate := time.Now().Add(5 * time.Second)
	id1 := addTask(c, "this is a task ", dueDate)

	fmt.Printf("added task %d", id1)
	fmt.Println("-------------------------------")
	fmt.Println("--------------LIST STREAM--------------")
	printTasks(c)
	fmt.Println("-------------------------------")

}
