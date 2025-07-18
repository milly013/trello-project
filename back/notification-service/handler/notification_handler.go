package handler

import (
	"net/http"

	"github.com/milly013/trello-project/back/notification-service/model"
	"github.com/milly013/trello-project/back/notification-service/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// Kreiranje novog obaveštenja
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var notification model.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.CreateNotification(&notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create notification"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Notification created successfully"})
}
func (h *NotificationHandler) GetAllNotifications(c *gin.Context) {

	notifications, err := h.service.GetAllNotifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve notifications"})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

// Dohvatanje obaveštenja po korisničkom ID-ju
func (h *NotificationHandler) GetNotificationsByUserID(c *gin.Context) {
	userID := c.Param("userID")

	notifications, err := h.service.GetNotificationsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve notifications"})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

// Označavanje obaveštenja kao pročitanog
func (h *NotificationHandler) MarkNotificationAsRead(c *gin.Context) {
	notificationID := c.Param("notificationID")

	err := h.service.MarkNotificationAsRead(notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to mark notification as read"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}
