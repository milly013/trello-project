package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/milly013/trello-project/back/user-service/handler"
	"github.com/milly013/trello-project/back/user-service/repository"
	"github.com/milly013/trello-project/back/user-service/service"
)

var userCollection *mongo.Collection

func main() {
	// Učitavanje .env fajla
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Učitavanje crne liste lozinki
	service.LoadBlacklistedPasswords("blacklist_passwords.txt") // Dodaj ovo da učitaš crnu listu lozinki

	// Povezivanje na MongoDB
	client, err := connectToMongoDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Kreiramo instancu baze podataka koristeći MONGODB_DATABASE
	db := client.Database(os.Getenv("MONGODB_DATABASE"))
	userCollection = db.Collection("users")

	userRepo := repository.NewUserRepository(db)
	jwtService := service.NewJWTService()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userRepo, jwtService, userService)

	router := gin.Default()
	router.Use(CORSMiddleware())

	// API rute za korisnike bez autentifikacije (npr. registracija i verifikacija)
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users", userHandler.GetUsers)
	router.POST("/verify/:email/:code", userHandler.VerifyUser)
	router.POST("/users/login", userHandler.Login)
	router.DELETE("/users/:id", userHandler.DeleteUserByID)
	router.POST("/users/forgot-password", userHandler.ForgotPassword)
	router.POST("/users/reset-password", userHandler.ResetPassword)
	router.GET("/users/:id", userHandler.GetUserByID)
	router.POST("/users/getByIds", userHandler.GetUsersByIds)
	router.GET("/users/isManager/:userId", userHandler.CheckIfUserIsManager)
	router.GET("/users/isMember/:userId", userHandler.CheckIfUserIsMember)
	router.POST("/users/change-password", userHandler.ChangePassword)

	router.POST("/users/request-magic-link", userHandler.RequestMagicLinkHandler)
	router.POST("/users/magic-login", userHandler.MagicLoginHandler)

	// Middleware za zaštitu ruta
	// authMiddleware := middleware.JWTAuth(jwtService)

	// Zaštićene rute
	// authRoutes := router.Group("/")
	// authRoutes.Use(authMiddleware)
	// {

	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Println("Server is running on port " + port)
	log.Fatal(srv.ListenAndServe())
}

// Funkcija za povezivanje na MongoDB
func connectToMongoDB() (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DATABASE")
	if mongoURI == "" || dbName == "" {
		log.Fatal("MONGODB_URI or MONGODB_DATABASE environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Povezivanje na MongoDB
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Connected to MongoDB database %s!\n", dbName)
	return client, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://localhost:4200")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Ako je preflight (OPTIONS) zahtev
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
