package service

import (
	"context"

	"github.com/milly013/trello-project/back/notification-service/model"
	"github.com/milly013/trello-project/back/notification-service/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

// Kreiranje obaveštenja
func (s *NotificationService) CreateNotification(ctx context.Context, notification *model.Notification) error {
	return s.repo.CreateNotification(ctx, notification)
}

// Dohvatanje obaveštenja po korisničkom ID-ju
func (s *NotificationService) GetNotificationsByUserID(ctx context.Context, userID primitive.ObjectID) ([]model.Notification, error) {
	return s.repo.GetNotificationsByUserID(ctx, userID)
}

// Označavanje obaveštenja kao pročitanog
func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, notificationID primitive.ObjectID) error {
	return s.repo.MarkAsRead(ctx, notificationID)
}
