package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

// Provera da li je korisnik član projekta
func (s *TaskService) IsUserPartOfProject(ctx context.Context, userID, projectID primitive.ObjectID) (bool, error) {

	url := fmt.Sprintf("%s/projects/%s/members/%s", s.projectServiceURL, projectID.Hex(), userID.Hex())

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var result struct {
		IsMember bool `json:"isMember"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false, err
	}

	return result.IsMember, nil
}
