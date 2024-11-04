// repository/user_repository.go
package repository

import (
	"context"

	"github.com/milly013/trello-project/back/user-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Funkcija za dobijanje korisnika po ID-u
func (r *UserRepository) GetUserByID(ctx context.Context, id string, user *model.User) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(user)
}

// Funkcija za dobijanje svih korisnika
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	cursor, err := r.collection.Find(ctx, bson.D{})
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
