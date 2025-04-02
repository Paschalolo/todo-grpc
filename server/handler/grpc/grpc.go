package grpc

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/paschalolo/grpc/proto/todo/v1"
	"github.com/paschalolo/grpc/server/controller"
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
	id, err := h.ctrl.Repo.AddTask(ctx, req.Description, req.DueDate.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AddTaskResponse{
		Id: id,
	}, nil
}

func (h *Handler) ListTasks(req *pb.ListTasksRequest, stream pb.TodoService_ListTasksServer) error {

	return h.ctrl.Repo.GetTasks(func(a any) error {
		task, ok := a.(*pb.Task)
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
	for {
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

		err = h.ctrl.Repo.UpdateTasks(req.Task.Id, req.Task.Description, req.Task.DueDate.AsTime(), req.Task.Done)
		if err != nil {
			return err
		}
	}
}
func (h *Handler) DeleteTasks(stream pb.TodoService_DeleteTasksServer) error {
	for {
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
