package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milly013/trello-project/back/task-service/model"
	"github.com/milly013/trello-project/back/task-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// Handler za kreiranje novog zadatka
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddTask(c, &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// Handler za dobijanje svih zadataka
func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.service.GetAllTasks(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

type AssignUserRequest struct {
	TaskID string `json:"taskId" binding:"required"`
	UserID string `json:"userId" binding:"required"`
}

// Handler za dodavanje korisnika na zadatak
func (h *TaskHandler) AssignMemberToTask(c *gin.Context) {
	var req AssignUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	taskID, err := primitive.ObjectIDFromHex(req.TaskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.service.AddMemberToTask(context.Background(), taskID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign user to task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to task successfully"})
}

type RemoveUserRequest struct {
	TaskID string `json:"taskId" binding:"required"`
	UserID string `json:"userId" binding:"required"`
}

// Handler za uklanjanje korisnika sa zadatka
func (h *TaskHandler) RemoveMemberFromTask(c *gin.Context) {
	var req RemoveUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	taskID, err := primitive.ObjectIDFromHex(req.TaskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.service.RemoveMemberFromTask(context.Background(), taskID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from task successfully"})
}

func (h *TaskHandler) GetTaskById(c *gin.Context) {
	taskID := c.Param("id")

	task, err := h.service.GetTaskById(context.Background(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	err := h.service.UpdateTask(context.Background(), &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}
