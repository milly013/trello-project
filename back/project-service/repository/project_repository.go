// repository/project_repository.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/milly013/trello-project/back/project-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func (repo *ProjectRepository) GetProjectsByManager(ctx context.Context, managerId string) ([]model.Project, error) {
	objID, err := primitive.ObjectIDFromHex(managerId)
	if err != nil {
		return nil, err
	}
	cursor, err := repo.collection.Find(ctx, bson.M{"managerId": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []model.Project
	for cursor.Next(ctx) {
		var project model.Project
		if err := cursor.Decode(&project); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return projects, nil

}
func (repo *ProjectRepository) GetProjectsByMember(ctx context.Context, memberId string) ([]model.Project, error) {
	objID, err := primitive.ObjectIDFromHex(memberId)
	if err != nil {
		return nil, err
	}
	cursor, err := repo.collection.Find(ctx, bson.M{"memberIds": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []model.Project
	for cursor.Next(ctx) {
		var project model.Project
		if err := cursor.Decode(&project); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if cursor.Err(); err != nil {
		return nil, err
	}

	return projects, nil

}

func (repo *ProjectRepository) GetProjectById(ctx context.Context, projectId string) (*model.Project, error) {
	var project model.Project
	objID, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return nil, err
	}

	err = repo.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Projekat nije pronađen
		}
		return nil, err
	}

	return &project, nil
}

func (repo *ProjectRepository) UpdateProject(ctx context.Context, project *model.Project) error {
	_, err := repo.collection.UpdateOne(
		ctx,
		bson.M{"_id": project.ID},
		bson.M{"$set": project},
	)
	return err
}

func (repo *ProjectRepository) GetTaskIDsByProject(ctx context.Context, projectId string) ([]primitive.ObjectID, error) {
	var project model.Project
	objID, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return nil, err
	}
	err = repo.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Projekat nije pronađen
		}
		return nil, err
	}
	return project.TaskIDs, nil
}

func (repo *ProjectRepository) GetUserIDsByProject(ctx context.Context, projectId string) ([]primitive.ObjectID, error) {
	var project model.Project

	objID, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return nil, err
	}

	err = repo.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return project.MemberIDs, nil
}
func (repo *ProjectRepository) DeleteProject(ctx context.Context, projectId string) error {
	objID, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return err
	}

	result, err := repo.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("project with ID %s not found", projectId)
	}

	return nil
}
