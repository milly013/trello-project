package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
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
func (r *UserRepository) CreateUser(ctx context.Context, user User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Funkcija za dobijanje svih korisnika
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]User, error) {
	var users []User
	filter := bson.D{} // Prazan filter za dobijanje svih korisnika

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Provera greške kursora nakon završetka iteracije
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
