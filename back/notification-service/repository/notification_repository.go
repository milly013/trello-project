package repository

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/milly013/trello-project/back/notification-service/model"
)

type NotificationRepository struct {
	session *gocql.Session
}

func NewNotificationRepository(session *gocql.Session) *NotificationRepository {
	return &NotificationRepository{
		session: session,
	}
}

// Metoda za kreiranje obaveštenja
func (r *NotificationRepository) CreateNotification(notification *model.Notification) error {
	// Proveri i generiši novi UUID ako nije postavljen
	if notification.ID == (gocql.UUID{}) {
		newID, err := gocql.RandomUUID()
		if err != nil {
			return fmt.Errorf("failed to generate UUID: %w", err)
		}
		notification.ID = newID
	}

	// Proveri i postavi tip notifikacije ako nije dodeljen
	if notification.Type == "" {
		notification.Type = "added_to_project" // ili neki podrazumevani tip
	}

	notification.CreatedAt = time.Now()

	query := `INSERT INTO notifications (id, user_id, type, message, created_at, is_read) VALUES (?, ?, ?, ?, ?, ?)`
	err := r.session.Query(query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.Message,
		notification.CreatedAt,
		notification.IsRead).Exec()

	if err != nil {
		return fmt.Errorf("failed to insert notification: %w", err)
	}

	return nil
}

func (r *NotificationRepository) GetNotificationsByUserID(userID string) ([]model.Notification, error) {
	query := `SELECT id, user_id, type, message, created_at, is_read FROM notifications WHERE user_id = ?`
	iter := r.session.Query(query, userID).Iter()

	var notifications []model.Notification

	var (
		id        gocql.UUID
		uID       string
		ntype     string
		message   string
		createdAt time.Time
		isRead    bool
	)

	for iter.Scan(&id, &uID, &ntype, &message, &createdAt, &isRead) {
		notifications = append(notifications, model.Notification{
			ID:        id,
			UserID:    uID,
			Type:      ntype,
			Message:   message,
			CreatedAt: createdAt,
			IsRead:    isRead,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to retrieve notifications: %w", err)
	}

	return notifications, nil
}

// Označavanje obaveštenja kao pročitanog
func (r *NotificationRepository) MarkAsRead(notificationID string) error {
	query := `UPDATE notifications SET is_read = true WHERE id = ?`
	return r.session.Query(query, notificationID).Exec()
}

// Metoda za dohvatanje svih obaveštenja
func (r *NotificationRepository) GetAllNotifications() ([]model.Notification, error) {
	query := `SELECT id, user_id, message, created_at, is_read FROM notifications`
	iter := r.session.Query(query).Iter()

	var notifications []model.Notification
	var notification model.Notification
	for iter.Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.CreatedAt, &notification.IsRead) {
		notifications = append(notifications, notification)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return notifications, nil
}
