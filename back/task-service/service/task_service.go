package service

import (
	"context"
	
	"github.com/milly013/trello-project/back/task-service/model"
	"github.com/milly013/trello-project/back/task-service/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	taskRepo          *repository.TaskRepository
	projectServiceURL string
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (s *TaskService) AddTask(ctx context.Context, task *model.Task) error {
	return s.taskRepo.CreateTask(ctx, task)
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	return s.taskRepo.GetAllTasks(ctx)
}

func (s *TaskService) AddMemberToTask(ctx context.Context, taskID, userID primitive.ObjectID) error {
	return s.taskRepo.AddUserToTask(ctx, taskID, userID)
}

func (s *TaskService) RemoveMemberFromTask(ctx context.Context, taskID, userID primitive.ObjectID) error {
	return s.taskRepo.RemoveUserFromTask(ctx, taskID, userID)
}

func (s *TaskService) GetTaskById(ctx context.Context, taskId string) (*model.Task, error) {
	return s.taskRepo.GetTaskById(ctx, taskId)
}

func (s *TaskService) UpdateTask(ctx context.Context, task *model.Task) error {
	return s.taskRepo.UpdateTask(ctx, task)
}

// Provera da li je korisnik dodeljen zadatku
func (s *TaskService) IsUserAssignedToTask(ctx context.Context, userID string, taskID string) (bool, error) {
	task, err := s.GetTaskById(ctx, taskID)
	if err != nil {
		return false, err
	}
	if task == nil {
		return false, nil
	}

	// Pretpostavljamo da task ima listu korisnika u polju `AssignedTo`
	for _, assignedUserID := range task.AssignedTo {
		if assignedUserID.Hex() == userID {
			return true, nil
		}
	}
	return false, nil
}

// Provera da li postoje nezavršene zavisnosti za zadatak
func (s *TaskService) HasUnfinishedDependencies(ctx context.Context, taskID string) (bool, error) {
	dependencies, err := s.taskRepo.GetTaskDependencies(ctx, taskID)
	if err != nil {
		return false, err
	}

	for _, dependency := range dependencies {
		if dependency.Status != "completed" {
			return true, nil
		}
	}
	return false, nil
}