package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/milly013/trello-project/back/task-service/model"
	"github.com/milly013/trello-project/back/task-service/repository"
	userModel "github.com/milly013/trello-project/back/user-service/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

	if task.ID.IsZero() {
		task.ID = primitive.NewObjectID()
	}
	sanitizedTitle := sanitizeInput(task.Title)
	if sanitizedTitle == "" {
		return fmt.Errorf("task title is invalid or empty")
	}
	task.Title = sanitizedTitle

	sanitizedDescription := sanitizeInput(task.Description)
	task.Description = sanitizedDescription

	if !isValidDateRange(task.StartDate, task.EndDate) {
		return fmt.Errorf("end date cannot be before start date or start date cannot be in the past")
	}

	// Create the task in the task repository
	err := s.taskRepo.CreateTask(ctx, task)
	if err != nil {
		return err
	}

	// Make a request to the project service to add the task ID to the project
	projectID := task.ProjectID
	if err := s.addTaskToProject(ctx, projectID, task.ID); err != nil {
		return fmt.Errorf("failed to add task to project: %w", err)
	}

	return nil
}
func (s *TaskService) addTaskToProject(ctx context.Context, projectID, taskID primitive.ObjectID) error {
	url := fmt.Sprintf("http://project-service:8081/projects/%s/tasks", projectID.Hex())

	requestBody, err := json.Marshal(map[string]string{
		"taskId": taskID.Hex(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to add task to project, received status code: %d", resp.StatusCode)
	}

	return nil
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
func (s *TaskService) GetTaskIDsByProject(ctx context.Context, projectId string) ([]primitive.ObjectID, error) {
	url := fmt.Sprintf("http://project-service:8081/projects/%s/tasks", projectId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Proverite statusni kod
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get task IDs by project, received status code: %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Čitanje tela odgovora
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("Response body: %s\n", string(bodyBytes)) // Dodajte logovanje odgovora

	// Dekodiranje odgovora
	var taskIDStrings []string
	if err := json.Unmarshal(bodyBytes, &taskIDStrings); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("Decoded task IDs:", taskIDStrings)

	// Konvertovanje u `ObjectID`
	var taskIDs []primitive.ObjectID
	for _, taskIDStr := range taskIDStrings {
		taskID, err := primitive.ObjectIDFromHex(taskIDStr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert task ID from string: %w", err)
		}
		taskIDs = append(taskIDs, taskID)
	}
	return taskIDs, nil
}

func (s *TaskService) GetTasksByProject(ctx context.Context, projectId string) ([]model.Task, error) {
	taskIDs, err := s.GetTaskIDsByProject(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to get task IDs from project-service: %w", err)
	}

	if projectId == "" {
		return nil, fmt.Errorf("projectId cannot be empty")
	}

	fmt.Println("Retrieved Task IDs:", taskIDs)

	tasks := []model.Task{}

	for _, taskID := range taskIDs {
		task, err := s.GetTaskById(ctx, taskID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to get task by id %s: %w", taskID.Hex(), err)
		}
		if task != nil {
			tasks = append(tasks, *task)
		}
	}
	return tasks, nil
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

// Funkcija koja dohvata sve korisnike dodeljene zadatku
func (s *TaskService) GetUsersByTaskId(ctx context.Context, taskId string) ([]*userModel.User, error) {
	// Prvo dobavljamo zadatak po ID-u
	task, err := s.GetTaskById(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Ako nema dodeljenih korisnika, vratimo prazan niz
	if len(task.AssignedTo) == 0 {
		return []*userModel.User{}, nil
	}

	// Dobavljamo informacije o korisnicima iz `user-service`
	users, err := s.getUsersByIDs(ctx, task.AssignedTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

// // Funkcija za dobijanje korisnika prema listi ID-ova iz `user-service`
func (s *TaskService) getUsersByIDs(ctx context.Context, userIDs []primitive.ObjectID) ([]*userModel.User, error) {
	// Pretvaranje ID-ova u niz stringova
	userIDStrings := make([]string, len(userIDs))
	for i, id := range userIDs {
		userIDStrings[i] = id.Hex()
	}

	// Kreiramo URL za pozivanje `user-service`
	url := fmt.Sprintf("http://user-service:8080/users/getByIds")
	requestBody, err := json.Marshal(map[string][]string{
		"userIds": userIDStrings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get users, status code: %d", resp.StatusCode)
	}

	var users []*userModel.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return users, nil
}

func (s *TaskService) GetTaskStatus(ctx context.Context, taskID string) (string, error) {
	task, err := s.GetTaskById(ctx, taskID)
	if err != nil {
		return "", err
	}
	if task == nil {
		return "", nil
	}

	return task.Status, nil
}

// Method to check if there are incomplete tasks for a given project
func (s *TaskService) HasIncompleteTasksByProject(ctx context.Context, projectID string) (bool, error) {
	tasks, err := s.GetTasksByProject(ctx, projectID)
	if err != nil {
		return false, fmt.Errorf("failed to get tasks for project: %w", err)
	}

	for _, task := range tasks {
		if task.Status != "Completed" {
			return true, nil
		}
	}

	return false, nil
}

//==========Validacije=========================

func sanitizeInput(input string) string {
	// Dozvoljeni karakteri: slova, brojevi, razmaci, osnovni interpunkcijski znakovi
	regex := regexp.MustCompile(`[^a-zA-Z0-9\s.,!?_-]`)
	return regex.ReplaceAllString(input, "")
}

func isValidDateRange(startDate, endDate time.Time) bool {
	now := time.Now()
	if startDate.Before(now) {
		return false
	}
	if endDate.Before(startDate) {
		return false
	}
	return true
}
