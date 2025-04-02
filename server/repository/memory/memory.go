package memory

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/paschalolo/grpc/proto/todo/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Repostiroy struct {
	sync.RWMutex
	Tasks []*pb.Task
}

func New() *Repostiroy {
	return &Repostiroy{}
}
func (r *Repostiroy) AddTask(_ context.Context, description string, dueDate time.Time) (uint64, error) {
	r.Lock()
	defer r.Unlock()
	log.Println("here")
	nextId := uint64(len(r.Tasks))
	task := &pb.Task{
		Id:          nextId,
		Description: description,
		DueDate:     timestamppb.New(dueDate),
	}
	r.Tasks = append(r.Tasks, task)
	return nextId, nil
}

func (r *Repostiroy) GetTasks(f func(any) error) error {
	for _, task := range r.Tasks {

		if err := f(task); err != nil {
			return err
		}
	}
	return nil
}
