package task

import (
	"context"
	"errors"

	"github.com/tamboto2000/otaqku-tasks/internal/common"
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

func (svc TaskService) GetTaskList(ctx context.Context, accId int, paginate dto.Pagination) (dto.TaskList, error) {
	taskList, err := svc.taskRepo.GetByAccountID(ctx, accId, paginate)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return dto.TaskList{}, nil
		}

		return dto.TaskList{}, err
	}

	var taskListDto dto.TaskList
	for _, task := range taskList.Tasks {
		taskListDto.Tasks = append(taskListDto.Tasks, taskToTaskDTO(task))
	}

	taskListDto.Pagination = taskList.Pagination

	return taskListDto, nil
}

func (svc TaskService) GetByID(ctx context.Context, accId int, id int) (dto.Task, error) {
	task, err := svc.taskRepo.GetByAccountIDAndID(ctx, accId, id)
	if err != nil {
		return dto.Task{}, err
	}

	return taskToTaskDTO(task), nil
}

func taskToTaskDTO(task Task) dto.Task {
	return dto.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func (svc TaskService) Update(ctx context.Context, accId int, req dto.Task) error {
	task, err := ValidateTaskForUpdate(accId, req)
	if err != nil {
		return err
	}

	return svc.taskRepo.UpdateByAccountIDAndID(ctx, task)
}

func (svc TaskService) Delete(ctx context.Context, accId, id int) error {
	return svc.taskRepo.DeleteByAccountIDAndID(ctx, accId, id)
}
