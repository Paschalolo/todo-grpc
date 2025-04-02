package controller

import (
	"context"
	"time"
)

type Repository interface {
	AddTask(ctx context.Context, description string, dueDate time.Time) (uint64, error)
	GetTasks(f func(any) error) error
}
type Controller struct {
	Repo Repository
}

func NewController(Repo Repository) *Controller {
	return &Controller{Repo: Repo}
}
func (c *Controller) AddTask(ctx context.Context, description string, duedate time.Time) (uint64, error) {
	id, err := c.Repo.AddTask(ctx, description, duedate)
	if err != nil {
		return 0, err
	}
	return id, nil

}

func (c *Controller) GetTasks(f func(any) error) error {
	if err := c.Repo.GetTasks(f); err != nil {
		return err
	}
	return nil
}
