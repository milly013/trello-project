// handler/project_handler.go
package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/milly013/trello-project/back/project-service/model"
	"github.com/milly013/trello-project/back/project-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(service *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.CreateProject(context.Background(), &project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
	projects, err := h.service.GetProjects(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) AddMemberToProject(c *gin.Context) {
	projectId := c.Param("projectId")
	var request struct {
		MemberID primitive.ObjectID `json:"memberId"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddMemberToProject(context.Background(), projectId, request.MemberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *ProjectHandler) AddTaskToProject(c *gin.Context) {
	projectId := c.Param("projectId")
	var request struct {
		TaskID primitive.ObjectID `json:"taskId"`
	}

	// Proveravamo da li je JSON ispravno vezan
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddTaskToProject(context.Background(), projectId, request.TaskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Vraćamo status 201 i kreirani task
	c.JSON(http.StatusNoContent, nil)
}
