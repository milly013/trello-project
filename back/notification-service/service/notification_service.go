package service

import (
	"github.com/milly013/trello-project/back/notification-service/model"
	"github.com/milly013/trello-project/back/notification-service/repository"
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
func (s *NotificationService) CreateNotification(notification *model.Notification) error {
	return s.repo.CreateNotification(notification)
}

// Dohvatanje obaveštenja po korisničkom ID-ju
func (s *NotificationService) GetNotificationsByUserID(userID string) ([]model.Notification, error) {
	return s.repo.GetNotificationsByUserID(userID)
}

// Označavanje obaveštenja kao pročitanog
func (s *NotificationService) MarkNotificationAsRead(notificationID string) error {
	return s.repo.MarkAsRead(notificationID)
}
func (s *NotificationService) GetAllNotifications() ([]model.Notification, error) {
	return s.repo.GetAllNotifications()
}
