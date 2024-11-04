package repository

import (
    "context"
    "github.com/milly013/trello-project/back/task-service/model"
    "go.mongodb.org/mongo-driver/mongo"
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
