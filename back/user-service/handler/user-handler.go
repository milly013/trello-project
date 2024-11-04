// handler/user_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milly013/trello-project/back/user-service/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	repo *UserRepository
}

func NewUserHandler(repo *UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// Handler za dodavanje novog korisnika
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Handler za preuzimanje korisnika po ID-u
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	err := h.repo.GetUserByID(c, id, &user) // Moraš implementirati ovu funkciju u UserRepository
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Handler za dobijanje svih korisnika
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.repo.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
