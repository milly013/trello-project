package repository

import (
	"context"

	"github.com/milly013/trello-project/back/task-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) *TaskRepository {
	return &TaskRepository{collection: collection}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	_, err := r.collection.InsertOne(ctx, task)
	return err
}

func (r *TaskRepository) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
