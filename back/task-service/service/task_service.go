package service

import (
	"context"

	"github.com/milly013/trello-project/back/task-service/model"
	"github.com/milly013/trello-project/back/task-service/repository"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (s *TaskService) AddTask(ctx context.Context, task *model.Task) error {
	return s.taskRepo.CreateTask(ctx, task)
}

func (s *TaskService) GetTasks(ctx context.Context) ([]model.Task, error) {
	return s.taskRepo.GetAllTasks(ctx)
}
