// handler/project_handler.go
package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/milly013/trello-project/back/project-service/model"
	"github.com/milly013/trello-project/back/project-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	projectId := c.Param("id")

	// Proveri validnost ID-a
	if _, err := primitive.ObjectIDFromHex(projectId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// Dobavi projekat koristeći servisnu metodu
	project, err := h.service.GetProjectById(context.Background(), projectId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		}
		return
	}

	c.JSON(http.StatusOK, project)
}
func (h *ProjectHandler) GetProjectsByManager(c *gin.Context) {
	managerId := c.Param("managerId")

	if _, err := primitive.ObjectIDFromHex(managerId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manager ID format"})
		return
	}
	prjects, err := h.service.GetProjectByManager(context.Background(), managerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects"})
		return
	}
	c.JSON(http.StatusOK, prjects)
}
func (h *ProjectHandler) GetProjectsByMember(c *gin.Context) {
	memberId := c.Param("memberId")

	if _, err := primitive.ObjectIDFromHex(memberId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID format"})
		return
	}
	projects, err := h.service.GetProjectsByMember(context.Background(), memberId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects"})
		return
	}
	c.JSON(http.StatusOK, projects)
}
func (h *ProjectHandler) GetTaskIDsByProject(c *gin.Context) {
	projectId := c.Param("id")

	if _, err := primitive.ObjectIDFromHex(projectId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}
	TaskIDs, err := h.service.GetTaskIDsByProject(context.Background(), projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task IDs"})
		return
	}
	if TaskIDs == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No tasks found for this project"})
	}
	c.JSON(http.StatusOK, TaskIDs)
}

func (h *ProjectHandler) GetUsersByProjectId(c *gin.Context) {
	projectId := c.Param("projectId")

	// Proveri validnost ID-a projekta
	if _, err := primitive.ObjectIDFromHex(projectId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// Dobavi korisnike za projekat koristeći servisnu metodu
	users, err := h.service.GetUsersByProjectId(context.Background(), projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users for project"})
		return
	}

	c.JSON(http.StatusOK, users)
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
	// Logovanje podataka
	fmt.Printf("Adding member with ID %s to project with ID %s\n", request.MemberID, projectId)

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
func (h *ProjectHandler) RemoveMemberFromProject(c *gin.Context) {
	projectId := c.Param("projectId")
	var request struct {
		MemberID primitive.ObjectID `json:"memberId"`
	}

	// Verifikujemo da li je JSON ispravno vezan
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pozivamo servis za uklanjanje člana
	err := h.service.RemoveMemberFromProject(context.Background(), projectId, request.MemberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Vraćamo status 204 (No Content) kao potvrdu da je član uspešno uklonjen
	c.JSON(http.StatusNoContent, nil)
}
