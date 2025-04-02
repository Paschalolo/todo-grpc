package memory

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	pb "github.com/paschalolo/grpc/proto/todo/v2"

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
	r.RLock()
	defer r.RUnlock()
	for _, task := range r.Tasks {

		if err := f(task); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repostiroy) UpdateTasks(id uint64, description string, dueDate time.Time, done bool) error {
	r.Lock()
	defer r.Unlock()
	for i, task := range r.Tasks {
		if task.Id == id {
			t := r.Tasks[i]
			t.Description = description
			t.Done = done
			t.DueDate = timestamppb.New(dueDate)
			return nil
		}
	}
	return fmt.Errorf("task with id %d not found ", id)

}

func (r *Repostiroy) DeleteTask(id uint64) error {
	r.Lock()
	defer r.Unlock()
	for i, task := range r.Tasks {
		if task.Id == id {
			r.Tasks = slices.Delete(r.Tasks, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("task with id not found %v ", id)
}
