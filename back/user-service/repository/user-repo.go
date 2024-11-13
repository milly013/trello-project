package repository

import (
	"context"
	"time"

	"github.com/milly013/trello-project/back/user-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// Provera da li korisnik postoji prema korisničkom imenu ili emailu
func (r *UserRepository) CheckUserExists(ctx context.Context, username, email string) (bool, error) {
	filter := bson.M{"$or": []bson.M{{"username": username}, {"email": email}}}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Čuvanje verifikacionog koda za korisnika
func (r *UserRepository) SaveVerificationCode(ctx context.Context, user model.User, code string) error {
	verificationData := bson.M{
		"email":            user.Email,
		"verificationCode": code,
		"createdAt":        time.Now(),
	}
	_, err := r.collection.InsertOne(ctx, verificationData)
	return err
}

// Provera verifikacionog koda
func (r *UserRepository) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	filter := bson.M{"email": email, "verificationCode": code}
	var result bson.M
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return false, err
	}
	// Proveri da li je kod istekao (ako želite logiku isteka koda, dodajte proveru)
	return true, nil
}

// Preuzimanje korisnika na osnovu emaila
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string, user *model.User) error {
	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return err
	}
	return nil
}

// Kreiranje novog korisnika
func (r *UserRepository) CreateUser(ctx context.Context, user model.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Preuzimanje korisnika po ID-u
func (r *UserRepository) GetUserByID(ctx context.Context, id string, user *model.User) error {
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(user)
	return err
}

// Preuzimanje svih korisnika
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
