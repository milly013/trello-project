package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
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
	sanitizedName := sanitizeProjectName(project.Name)
	if sanitizedName == "" {
		return fmt.Errorf("project name contains invalid characters or is empty")
	}
	project.Name = sanitizedName

	if !isValidEndDate(project.EndDate) {
		return fmt.Errorf("end date cannot be in the past")
	}

	if project.MaxMembers < project.MinMembers {
		return fmt.Errorf("maximum members cannot be less than minimum members")
	}

	project.CreatedAt = time.Now()
	project.IsActive = true

	_, err := s.repo.CreateProject(ctx, project)
	return err
}

func (s *ProjectService) GetProjects(ctx context.Context) ([]model.Project, error) {
	return s.repo.GetProjects(ctx)
}

func (s *ProjectService) GetProjectById(ctx context.Context, projectId string) (*model.Project, error) {
	return s.repo.GetProjectById(ctx, projectId)
}
func (s *ProjectService) GetProjectByManager(ctx context.Context, managerId string) ([]model.Project, error) {
	return s.repo.GetProjectsByManager(ctx, managerId)
}
func (s *ProjectService) GetProjectsByMember(ctx context.Context, memberId string) ([]model.Project, error) {
	return s.repo.GetProjectsByMember(ctx, memberId)
}
func (s *ProjectService) GetTaskIDsByProject(ctx context.Context, projectId string) ([]primitive.ObjectID, error) {
	return s.repo.GetTaskIDsByProject(ctx, projectId)
}

func (s *ProjectService) UserExists(ctx context.Context, memberId primitive.ObjectID) (bool, error) {
	url := fmt.Sprintf("http://api-gateway:8000/api/user/users/%s", memberId.Hex())

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
	log.Printf("Attempting to retrieve project with ID: %s", projectId)
	project, err := s.repo.GetProjectById(ctx, projectId)
	if err != nil {
		log.Printf("Error retrieving project with ID %s: %v", projectId, err)
		return err
	}

	log.Printf("Project retrieved successfully: %+v", project)

	// Proveri da li task već postoji u projektu
	for _, id := range project.TaskIDs {
		if id == taskId {
			log.Printf("Task with ID %s already exists in project %s", taskId.Hex(), projectId)
			return fmt.Errorf("task already exists in project")
		}
	}

	project.TaskIDs = append(project.TaskIDs, taskId)
	log.Printf("Adding task ID %s to project %s", taskId.Hex(), projectId)

	err = s.repo.UpdateProject(ctx, project)
	if err != nil {
		log.Printf("Error updating project with ID %s: %v", projectId, err)
		return err
	}

	log.Printf("Project updated successfully with new task ID %s", taskId.Hex())
	return nil
}

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

			log.Println("adsadsadsadsads")

			return s.repo.UpdateProject(ctx, project)
		}
	}

	return fmt.Errorf("member not found in project")
}

//===============Validacije=====================

func sanitizeProjectName(name string) string {
	// Dozvoljeni karakteri: slova, brojevi, razmaci, crtice, podvlake
	regex := regexp.MustCompile(`^[a-zA-Z0-9\s_-]+$`)
	if regex.MatchString(name) {
		return name
	}
	return ""

}
func isValidEndDate(endDate time.Time) bool {
	now := time.Now()
	return endDate.After(now)
}
