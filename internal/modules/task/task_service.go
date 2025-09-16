package task

import (
	"context"

	"github.com/tamboto2000/otaqku-tasks/internal/dto"
)

type TaskService struct {
	taskRepo TaskRepository
}

func NewTaskService(taskRepo TaskRepository) TaskService {
	return TaskService{taskRepo: taskRepo}
}

func (svc TaskService) CreateTask(ctx context.Context, accId int, req dto.CreateTaskRequest) error {
	task, err := NewTask(accId, req)
	if err != nil {
		return err
	}

	err = svc.taskRepo.Save(ctx, task)

	return err
}
