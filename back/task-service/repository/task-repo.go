package repository

import (
	"context"

	"github.com/milly013/trello-project/back/task-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *TaskRepository) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task model.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, cursor.Err()
}

func (r *TaskRepository) AddUserToTask(ctx context.Context, taskID, userID primitive.ObjectID) error {
	filter := bson.M{"_id": taskID}
	update := bson.M{
		"$addToSet": bson.M{"assignedTo": userID},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
func (r *TaskRepository) RemoveUserFromTask(ctx context.Context, taskID, userID primitive.ObjectID) error {
	filter := bson.M{"_id": taskID}
	update := bson.M{
		"$pull": bson.M{"assignedTo": userID},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *TaskRepository) GetTaskById(ctx context.Context, taskId string) (*model.Task, error) {
	var task model.Task
	objID, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"id": objID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": task.ID},
		bson.M{"$set": task},
	)
	return err
}
