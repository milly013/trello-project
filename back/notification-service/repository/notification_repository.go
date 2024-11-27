package repository

import (
	"context"
	"time"

	"github.com/milly013/trello-project/back/notification-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	return &NotificationRepository{
		collection: db.Collection("notifications"),
	}
}

// Metoda za kreiranje obaveštenja
func (r *NotificationRepository) CreateNotification(ctx context.Context, notification *model.Notification) error {
	notification.ID = primitive.NewObjectID()
	notification.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, notification)
	return err
}

// Metoda za dohvatanje obaveštenja po korisničkom ID-ju
func (r *NotificationRepository) GetNotificationsByUserID(ctx context.Context, userID primitive.ObjectID) ([]model.Notification, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var notifications []model.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

// Označavanje obaveštenja kao pročitanog
func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID primitive.ObjectID) error {
	filter := bson.M{"_id": notificationID}
	update := bson.M{"$set": bson.M{"is_read": true}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
