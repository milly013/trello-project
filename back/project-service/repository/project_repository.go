// repository/project_repository.go
package repository

import (
	"context"
	"time"

	"github.com/milly013/trello-project/back/project-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectRepository struct {
	collection *mongo.Collection
}

func NewProjectRepository(db *mongo.Database) *ProjectRepository {
	return &ProjectRepository{
		collection: db.Collection("projects"),
	}
}

func (repo *ProjectRepository) CreateProject(ctx context.Context, project *model.Project) (*mongo.InsertOneResult, error) {
	project.CreatedAt = time.Now()
	project.IsActive = true
	return repo.collection.InsertOne(ctx, project)
}

func (repo *ProjectRepository) GetProjects(ctx context.Context) ([]model.Project, error) {
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []model.Project
	for cursor.Next(ctx) {
		var project model.Project
		cursor.Decode(&project)
		projects = append(projects, project)
	}
	return projects, nil
}

// Dodati ostale funkcije kao što su UpdateProject, DeleteProject, FindById...
