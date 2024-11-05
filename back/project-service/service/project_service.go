package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/milly013/trello-project/back/project-service/model"
	"github.com/milly013/trello-project/back/project-service/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(ctx context.Context, project *model.Project) error {
	_, err := s.repo.CreateProject(ctx, project)
	return err
}

func (s *ProjectService) GetProjects(ctx context.Context) ([]model.Project, error) {
	return s.repo.GetProjects(ctx)
}

func (s *ProjectService) GetProjectById(ctx context.Context, projectId string) (*model.Project, error) {
	return s.repo.GetProjectById(ctx, projectId)
}

func (s *ProjectService) UserExists(ctx context.Context, memberId primitive.ObjectID) (bool, error) {
	url := fmt.Sprintf("http://localhost:8080/users/%s", memberId.Hex())

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("error checking user existence: %s", resp.Status)
}

func (s *ProjectService) AddMemberToProject(ctx context.Context, projectId string, memberId primitive.ObjectID) error {

	exists, err := s.UserExists(ctx, memberId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user does not exist")
	}

	project, err := s.repo.GetProjectById(ctx, projectId)
	if err != nil {
		return err
	}

	if len(project.MemberIDs) >= project.MaxMembers {
		return fmt.Errorf("maximum number of members reached")
	}

	for _, id := range project.MemberIDs {
		if id == memberId {
			return fmt.Errorf("member already exists in project")
		}
	}

	project.MemberIDs = append(project.MemberIDs, memberId)
	return s.repo.UpdateProject(ctx, project)
}

func (s *ProjectService) TaskExists(ctx context.Context, taskId primitive.ObjectID) (bool, error) {
	url := fmt.Sprintf("http://localhost:8080/tasks/%s", taskId.Hex())

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("error checking task existence: %s", resp.Status)
}

func (s *ProjectService) AddTaskToProject(ctx context.Context, projectId string, taskId primitive.ObjectID) error {

	exists, err := s.TaskExists(ctx, taskId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user does not exist")
	}

	project, err := s.repo.GetProjectById(ctx, projectId)
	if err != nil {
		return err
	}

	for _, id := range project.TaskIDs {
		if id == taskId {
			return fmt.Errorf("tasks already exists in project")
		}
	}

	project.TaskIDs = append(project.TaskIDs, taskId)
	return s.repo.UpdateProject(ctx, project)
}

// RemoveMemberFromProject uklanja člana iz projekta
// Prebaci ovu funkciju iznutra druge funkcije
func (s *ProjectService) RemoveMemberFromProject(ctx context.Context, projectId string, memberId primitive.ObjectID) error {
	project, err := s.repo.GetProjectById(ctx, projectId)
	if err != nil {
		return err
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}

	// Proveri da li član postoji u projektu
	for i, id := range project.MemberIDs {
		if id == memberId {
			// Ukloni člana
			project.MemberIDs = append(project.MemberIDs[:i], project.MemberIDs[i+1:]...)
			return s.repo.UpdateProject(ctx, project)
		}
	}

	return fmt.Errorf("member not found in project")
}
