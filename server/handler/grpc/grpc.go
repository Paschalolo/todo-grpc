package grpc

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/paschalolo/grpc/proto/todo/v2"
	"github.com/paschalolo/grpc/server/controller"
	mask "github.com/paschalolo/grpc/utils/masks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	ctrl *controller.Controller
	pb.UnimplementedTodoServiceServer
}

func NewHandler(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			log.Printf("request deadline exceeded : %s ", ctx.Err())
		default:
		}
		return nil, ctx.Err()
		// case <-time.After(1 * time.Millisecond):
	default:
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// if len(req.Description) == 0 {
	// 	return nil, status.Error(codes.InvalidArgument, "expected a task description got an empty string")
	// }
	// if req.DueDate.AsTime().Before(time.Now().UTC()) {
	// 	return nil, status.Error(codes.InvalidArgument, "expected a task due dates that is in the future ")
	// }
	id, err := h.ctrl.Repo.AddTask(ctx, req.Description, req.DueDate.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AddTaskResponse{
		Id: id,
	}, nil
}

func (h *Handler) ListTasks(req *pb.ListTasksRequest, stream pb.TodoService_ListTasksServer) error {
	ctx := stream.Context()
	return h.ctrl.Repo.GetTasks(func(a any) error {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				log.Printf("request deadline exceeded : %s ", ctx.Err())
			default:
			}
			return ctx.Err()
			// case <-time.After(1 * time.Millisecond):
		default:
		}
		task, ok := a.(*pb.Task)
		mask.Filter(task, req.Mask)
		if !ok {
			return status.Error(codes.InvalidArgument, "could not convert ")
		}
		overdue := task.DueDate != nil && !task.Done && task.DueDate.AsTime().Before(time.Now().UTC())
		err := stream.Send(&pb.ListTasksResponse{
			Task:    task,
			Overdue: overdue,
		})
		return err
	})

}

func (h *Handler) UpdateTasks(stream pb.TodoService_UpdateTasksServer) error {
	totalLength := 0
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				log.Printf("request deadline exceeded : %s ", ctx.Err())
			default:
			}
			return ctx.Err()
			// case <-time.After(1 * time.Millisecond):
		default:
		}
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("total length:", totalLength)
			return stream.SendAndClose(&pb.UpdateTasksResponse{})
		}
		if err != nil {
			return err
		}

		out, _ := proto.Marshal(req)
		totalLength += len(out)

		err = h.ctrl.Repo.UpdateTasks(req.Id, req.Description, req.DueDate.AsTime(), req.Done)
		if err != nil {
			return err
		}
	}
}
func (h *Handler) DeleteTasks(stream pb.TodoService_DeleteTasksServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				log.Printf("request deadline exceeded : %s ", ctx.Err())
			default:
			}
			return ctx.Err()
		default:
		}
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := h.ctrl.Repo.DeleteTask(req.Id); err != nil {
			return err
		}
		if err := stream.Send(&pb.UpdateTasksResponse{}); err != nil {
			return err
		}
	}
}
