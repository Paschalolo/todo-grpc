package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/paschalolo/grpc/proto/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	client pb.TodoServiceClient
}

func NewClientCaller(conn *grpc.ClientConn) *Client {
	return &Client{
		client: pb.NewTodoServiceClient(conn),
	}
}
func (c *Client) AddTask(description string, dueDate time.Time) uint64 {
	ctx := context.Background()
	req := &pb.AddTaskRequest{
		Description: description,
		DueDate:     timestamppb.New(dueDate),
	}
	res, err := c.client.AddTask(ctx, req)
	if err != nil {
		panic(err)
	}
	return res.Id
}

func (c *Client) PrintTasks() {
	req := &pb.ListTasksRequest{}
	stream, err := c.client.ListTasks(context.Background(), req)
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

func (c *Client) UpdateTasks(reqs ...*pb.UpdateTasksRequest) {
	stream, err := c.client.UpdateTasks(context.Background())
	if err != nil {
		log.Fatalf("unexpected error : %v", err)
	}
	for _, req := range reqs {
		if err := stream.Send(req); err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
		if req != nil {
			fmt.Printf("updated task with id: %d\n", req.Task.Id)
		}
	}
	if _, err := stream.CloseAndRecv(); err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

}
