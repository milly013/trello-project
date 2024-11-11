// repository/task_repository.go
package repository

import (
	"context"

	"github.com/milly013/trello-project/back/task-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskRepository - struktura za rad sa kolekcijom zadataka
type TaskRepository struct {
	collection *mongo.Collection
}

// NewTaskRepository - funkcija koja vraća novi TaskRepository
func NewTaskRepository(collection *mongo.Collection) *TaskRepository {
	return &TaskRepository{collection: collection}
}

// CreateTask - kreira novi zadatak u bazi
func (r *TaskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	_, err := r.collection.InsertOne(ctx, task)
	return err
}

// GetAllTasks - vraća sve zadatke
func (r *TaskRepository) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	for cursor.Next(ctx) {
		var task model.Task
		cursor.Decode(&task)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetTaskById - vraća zadatak prema ID-u
func (r *TaskRepository) GetTaskById(ctx context.Context, taskId string) (*model.Task, error) {
	var task model.Task
	objID, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Zadatak nije pronađen
		}
		return nil, err
	}

	return &task, nil
}

// UpdateTask - ažurira podatke o zadatku
func (r *TaskRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": task.ID},
		bson.M{"$set": task},
	)
	return err
}

// AddMemberToTask - dodaje člana na zadatak
func (r *TaskRepository) AddMemberToTask(ctx context.Context, taskID string, memberID primitive.ObjectID) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	// Dodavanje člana u zadatak
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$addToSet": bson.M{"assignedTo": memberID}}, // "assignedTo" je niz ID-eva članova
	)

	return err
}
