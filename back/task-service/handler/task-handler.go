package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milly013/trello-project/back/task-service/model"
	"github.com/milly013/trello-project/back/task-service/service"
)

type TaskHandler struct {
	service *service.TaskService // Menjamo repo u service
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service} // Dodeljujemo servis
}

// Handler za dodavanje novog zadatka
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddTask(c, &task); err != nil { // Koristimo servis
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// Handler za dobijanje svih zadataka
func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.service.GetAllTasks(c) // Koristimo servis
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
