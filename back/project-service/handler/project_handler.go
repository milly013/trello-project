// handler/project_handler.go
package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milly013/trello-project/back/project-service/model"
	"github.com/milly013/trello-project/back/project-service/service"
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
