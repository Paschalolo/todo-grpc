/*
The package works on the client handlers for grpc

The names of the handlers are:
  - AddTask
  - PrintTasks
  - UpdateTasks
  - DeleteTask

This is not rendered as code
*/
package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/paschalolo/grpc/proto/todo/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// this Client structure is for holding
// the gRPC TodoServiceClient
type Client struct {
	client pb.TodoServiceClient
}

/*
The NewClientCaller() makes a new client struct
it requires a grpc connection
returns a
*/
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

func (c *Client) PrintTasks(fm *fieldmaskpb.FieldMask) {
	req := &pb.ListTasksRequest{
		Mask: fm,
	}
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
			fmt.Printf("updated task with id: %d\n", req.Id)
		}
	}
	if _, err := stream.CloseAndRecv(); err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

}

func (c *Client) DeleteTask(reqs ...*pb.DeleteTasksRequest) {
	stream, err := c.client.DeleteTasks(context.Background())
	if err != nil {
		log.Fatalf("unexpected error : %v ", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while recieving : %v ", err)
			}
			log.Println("deleted tasks")

		}
	}()
	for _, req := range reqs {
		if err := stream.Send(req); err != nil {
			return
		}
	}
	if err := stream.CloseSend(); err != nil {
		return
	}
	wg.Wait()

}
