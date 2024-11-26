package repository

import (
	"context"
	"fmt"
	"time"

	"regexp"

	"github.com/milly013/trello-project/back/user-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		"username":         user.Username,
		"email":            user.Email,
		"password":         user.Password,
		"verificationCode": code,
		"role":             user.Role,
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
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	filter := bson.M{"email": email}
	var user model.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Kreiranje novog korisnika
func (r *UserRepository) CreateUser(ctx context.Context, user model.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Brisanje korisnika po ID-u
func (r *UserRepository) DeleteUserByID(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectID}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}

// Preuzimanje korisnika po ID-u
func (r *UserRepository) GetUserByID(ctx context.Context, id string, user *model.User) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(user)
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

func (r *UserRepository) VerifyUserAndActivate(ctx context.Context, email, code string) (bool, error) {
	// Filtriraj korisnika prema email-u i verifikacionom kodu
	filter := bson.M{"email": email, "verificationCode": code}

	var user model.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		// Ako korisnik nije pronađen, vrati false
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err // Ako je došlo do druge greške
	}

	// Ako se korisnik pronađe, ažuriraj njegov status na aktiviran
	_, err = r.collection.UpdateOne(ctx, bson.M{"email": email}, bson.M{
		"$set": bson.M{
			"isActive": true,
		},
	})
	if err != nil {
		return false, err
	}

	// Opcionalno: Ukloni verifikacioni kod iz baze
	_, err = r.collection.UpdateOne(ctx, bson.M{"email": email}, bson.M{
		"$unset": bson.M{
			"verificationCode": "",
		},
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

// Ažuriranje lozinke korisnika
func (r *UserRepository) UpdatePassword(ctx context.Context, userID, newPassword string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format")
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"password": newPassword}}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Preuzimanje korisnika prema listi ID-eva
func (r *UserRepository) GetUsersByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.User, error) {
	// Filtriraj korisnike prema ID-evima koristeći $in operator
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []model.User
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
