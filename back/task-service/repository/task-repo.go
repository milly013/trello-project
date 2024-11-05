package repository

import (
	"context"

	"github.com/milly013/trello-project/back/task-service/model"
	"go.mongodb.org/mongo-driver/bson"
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
	cursor, err := r.collection.Find(ctx, bson.M{}) // Pronađite sve zadatke
	if err != nil {
		return nil, err // Vraćamo grešku ako dođe do problema sa upitom
	}
	defer cursor.Close(ctx) // Zatvorite kursor kada završite

	var tasks []model.Task
	for cursor.Next(ctx) {
		var task model.Task
		if err := cursor.Decode(&task); err != nil { // Dekodirajte zadatak
			return nil, err // Vraćamo grešku ako dekodiranje ne uspe
		}
		tasks = append(tasks, task) // Dodajte zadatak u slice
	}

	// Proverite da li je bilo grešaka tokom iteracije
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil // Vraćamo slice sa svim zadacima
}
