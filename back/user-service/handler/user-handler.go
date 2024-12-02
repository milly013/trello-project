package handler

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/milly013/trello-project/back/user-service/model"
	"github.com/milly013/trello-project/back/user-service/repository"
	"github.com/milly013/trello-project/back/user-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	repo       *repository.UserRepository
	jwtService *service.JWTService
	service    *service.UserService
}

// Kreiraj novi UserHandler
func NewUserHandler(repo *repository.UserRepository, jwtService *service.JWTService, service *service.UserService) *UserHandler {
	return &UserHandler{
		repo:       repo,
		service:    service,
		jwtService: jwtService}
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
		return err // Dodajte ovaj red
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

	// Provera da li lozinka pripada crnoj listi
	if service.IsPasswordBlacklisted(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lozinka koju ste uneli je previše slaba. Molimo vas koristite jaču lozinku."})
		return
	}

	// Provera da li korisnik već postoji
	exists, err := h.repo.CheckUserExists(c, user.Username, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Korisnik sa zadatim korisničkim imenom ili email-om već postoji"})
		return
	}

	// Heširanje lozinke
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom heširanja lozinke"})
		return
	}
	user.Password = string(hashedPassword)

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

	c.JSON(http.StatusCreated, gin.H{"message": "Korisnik uspešno kreiran"})
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
	_, err = h.repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	// Sada možeš nastaviti s procesom registracije
	c.JSON(http.StatusCreated, gin.H{"message": "User verified successfully"})

}

// Brisanje korisnika po ID-u
func (h *UserHandler) DeleteUserByID(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.DeleteUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Handler za preuzimanje korisnika po ID-u
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	// Pozovemo repo da dobijemo korisnika i eventualnu grešku
	user, err := h.repo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		// Ako je greška da nema dokumenata, vratiti odgovarajući status
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		// Ako je došlo do neke druge greške, vratiti internu grešku servera
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ako sve prođe kako treba, vratiti korisnika
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
func (h *UserHandler) GetUsersByIds(c *gin.Context) {
	var requestBody struct {
		UserIds []string `json:"userIds"`
	}

	// Parsiraj JSON telo zahteva
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Pretvaranje ID-ova iz string formata u ObjectID
	var userIDs []primitive.ObjectID
	for _, id := range requestBody.UserIds {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		userIDs = append(userIDs, objectID)
	}

	// Dobavljanje korisnika pomoću UserService-a
	users, err := h.service.GetUsersByIds(c.Request.Context(), userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Vraćanje liste korisnika
	c.JSON(http.StatusOK, users)
}

// Verifikacija korisnika
func (h *UserHandler) VerifyUser(c *gin.Context) {
	email := c.Param("email")
	code := c.Param("code")

	// Pozovi metodu iz repozitorijuma za verifikaciju korisnika
	success, err := h.repo.VerifyUserAndActivate(c, email, code)
	if err != nil {
		log.Printf("Error during verification: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Verification failed"})
		return
	}

	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully verified and activated"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode JSON iz zahteva
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Provera lozinke
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generiši JWT token
	token, err := h.jwtService.GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"userId":   user.ID.Hex(),
		"userRole": user.Role,
	})
}

// Handler za forgot password
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := h.repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Generisanje tokena za reset lozinke i njegovo vreme isteka
	resetToken := generateVerificationCode() // Možeš koristiti istu funkciju za generisanje tokena
	user.ResetToken = resetToken
	user.ResetTokenExpiresAt = time.Now().Add(15 * time.Minute)

	// Ažuriranje korisnika sa novim tokenom
	if err := h.repo.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user with reset token"})
		return
	}

	// Slanje email-a sa reset tokenom
	if err := sendVerificationEmail(user.Email, resetToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent successfully"})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := h.repo.GetUserByResetToken(c.Request.Context(), req.Token)
	if err != nil || user.ResetTokenExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Heširanje nove lozinke
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiresAt = time.Time{}

	if err := h.repo.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
	}
}

// Handler metoda za proveru da li je korisnik menadžer
func (h *UserHandler) CheckIfUserIsManager(c *gin.Context) {
	userID := c.Param("userId")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	isManager, err := h.service.IsUserManager(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isManager": isManager})
}

// Handler metoda za proveru da li je korisnik član
func (h *UserHandler) CheckIfUserIsMember(c *gin.Context) {
	userID := c.Param("userId")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	isMember, err := h.service.IsUserMember(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isMember": isMember})
}

func (h UserHandler) ChangePassword(c *gin.Context) {

	var req struct {
		UserID          string `json:"userId"`
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Proveri da li je userID validan ObjectID
	_, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = h.service.ChangePassword(c.Request.Context(), req.UserID, req.CurrentPassword, req.NewPassword)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// Secret key for JWT token (preporučujem da koristiš environment varijablu za ovo)
var jwtSecret = []byte("tvoja-tajna")

// Handler za slanje magic link-a korisniku
func (h *UserHandler) RequestMagicLinkHandler(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// Provera da li korisnik postoji
	user, err := h.repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding user"})
		return
	}

	// Generisanje JWT tokena sa istekon od 15 minuta
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := jwt.MapClaims{
		"email": user.Email,
		"exp":   expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Kreiraj magic link
	magicLink := fmt.Sprintf("http://localhost:4200/magic-login?token=%s", tokenString)

	// Pošalji magic link putem email-a
	err = sendVerificationEmail(user.Email, magicLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send magic link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Magic link sent successfully"})
}

// Handler za verifikaciju magic link tokena i prijavu korisnika
func (h *UserHandler) MagicLoginHandler(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	// Parsiraj i validiraj JWT token
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Proveri email iz tokena
	email := claims["email"].(string)
	user, err := h.repo.GetUserByEmail(context.Background(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	// Generiši JWT token za prijavu korisnika
	loginToken, err := h.jwtService.GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating login token"})
		return
	}

	// Vrati login token korisniku
	c.JSON(http.StatusOK, gin.H{
		"token":    loginToken,
		"userId":   user.ID.Hex(),
		"userRole": user.Role,
	})
}
