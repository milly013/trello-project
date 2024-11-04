package model

import (
	"context"
	"time"

	//provjeri za putanju

	"github.com/milly013/trello-project/back/task-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	collection *mongo.Collection
}

// NewTaskRepository kreira novi TaskRepository sa referencom na MongoDB kolekciju "tasks"
func NewTaskRepository(db *mongo.Database) *TaskRepository {
	return &TaskRepository{
		collection: db.Collection("tasks"),
	}
}

// CreateTask kreira novi task u kolekciji i postavlja CreatedAt na trenutni datum
func (repo *TaskRepository) CreateTask(ctx context.Context, task *model.Task) (*mongo.InsertOneResult, error) {
	task.CreatedAt = time.Now()
	return repo.collection.InsertOne(ctx, task)
}

// GetTasks vraća sve taskove iz kolekcije
func (repo *TaskRepository) GetTasks(ctx context.Context) ([]model.Task, error) {
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	for cursor.Next(ctx) {
		var task model.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
