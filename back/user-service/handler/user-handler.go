package handler

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/milly013/trello-project/back/user-service/model"
	"github.com/milly013/trello-project/back/user-service/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	repo *repository.UserRepository
}

// Kreiraj novi UserHandler
func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// Funkcija za generisanje slučajnog verifikacionog koda
func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000)) // Šestocifreni kod
}

// Funkcija za slanje verifikacionog emaila
func sendVerificationEmail(toEmail, verificationCode string) error {
	msg := []byte(fmt.Sprintf("Subject: Verifikacioni kod\n\nVas verifikacioni kod je: %s", verificationCode))
	auth := smtp.PlainAuth("", "teodosicmilos700@gmail.com", "azjjqovkylstwcjl", "smtp.gmail.com")
	err := smtp.SendMail(fmt.Sprintf("%s:%s", "smtp.gmail.com", "587"), auth, "teodosicmilos700@gmail.com", []string{toEmail}, msg)
	if err != nil {
		log.Println("jebiga")
	}

	return nil
}

// Handler za dodavanje novog korisnika uz proveru postojanja i slanje verifikacionog koda
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Provera da li korisnik već postoji
	exists, err := h.repo.CheckUserExists(c, user.Username, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User with given username or email already exists"})
		return
	}

	// Generisanje verifikacionog koda
	verificationCode := generateVerificationCode()

	// Slanje koda putem e-pošte
	err = sendVerificationEmail(user.Email, verificationCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
		return
	}

	// Čuvanje verifikacionog koda
	h.repo.SaveVerificationCode(c, user, verificationCode)

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
}

// Handler za verifikaciju koda
func (h *UserHandler) VerifyCode(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid, err := h.repo.VerifyCode(c, req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired code"})
		return
	}

	// Dobijanje korisničkih podataka na osnovu emaila
	var user model.User
	if err := h.repo.GetUserByEmail(c, req.Email, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	// Upisivanje korisnika u glavnu kolekciju nakon verifikacije
	if err := h.repo.CreateUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Handler za preuzimanje korisnika po ID-u
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	err := h.repo.GetUserByID(c, id, &user)
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
