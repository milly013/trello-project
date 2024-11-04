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

// GetProjectById vraća projekat na osnovu ID-a
func (s *ProjectService) GetProjectById(ctx context.Context, projectId string) (*model.Project, error) {
	return s.repo.GetProjectById(ctx, projectId)
}

// UserExists proverava da li korisnik postoji u user-service
func (s *ProjectService) UserExists(ctx context.Context, memberId primitive.ObjectID) (bool, error) {
	url := fmt.Sprintf("http://localhost:8080/users/%s", memberId.Hex()) // Zameni sa stvarnim URL-om user-service

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil // Korisnik postoji
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil // Korisnik ne postoji
	}

	return false, fmt.Errorf("error checking user existence: %s", resp.Status)
}

// AddMemberToProject dodaje člana u projekat
func (s *ProjectService) AddMemberToProject(ctx context.Context, projectId string, memberId primitive.ObjectID) error {
	// Proveri da li korisnik postoji
	exists, err := s.UserExists(ctx, memberId)
	if err != nil {
		return err // Vraćamo grešku ako se desila
	}
	if !exists {
		return fmt.Errorf("user does not exist") // Korisnik ne postoji
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

// Možete dodati i druge funkcije kao što su DeleteProject, UpdateProject, FindById...
