// service/project_service.go
package service

import (
	"context"

	"github.com/milly013/trello-project/back/project-service/model"
	"github.com/milly013/trello-project/back/project-service/repository"
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

// Možete dodati i druge funkcije kao što su DeleteProject, UpdateProject, FindById...
