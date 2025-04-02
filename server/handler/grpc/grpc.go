package grpc

import (
	"context"
	"time"

	pb "github.com/paschalolo/grpc/proto/todo/v1"
	"github.com/paschalolo/grpc/server/controller"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	ctrl *controller.Controller
	pb.UnimplementedTodoServiceServer
}

func NewHandler(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) Addtask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
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
