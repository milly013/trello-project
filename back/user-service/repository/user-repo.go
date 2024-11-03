package repository

import (
	"context"

	"github.com/milly013/trello-project/model" // Izmenite putanju prema vašem projektu
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName string) *UserRepository {
	collection := client.Database(dbName).Collection("users")
	return &UserRepository{collection: collection}
}

// Funkcija za dodavanje novog korisnika
func (r *UserRepository) CreateUser(ctx context.Context, user model.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Funkcija za dobijanje svih korisnika
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	cursor, err := r.collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
