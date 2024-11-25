package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/milly013/trello-project/back/user-service/handler"
	"github.com/milly013/trello-project/back/user-service/middleware"
	"github.com/milly013/trello-project/back/user-service/repository"
	"github.com/milly013/trello-project/back/user-service/service"
)

var userCollection *mongo.Collection

func main() {
	// Ucitavanje .env fajla
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Povezivanje na MongoDB
	client, err := connectToMongoDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("mydatabase")
	userCollection = db.Collection("users")

	userRepo := repository.NewUserRepository(db)
	jwtService := service.NewJWTService()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userRepo, jwtService, userService)

	router := gin.Default()

	// API rute za korisnike bez autentifikacije (npr. registracija i verifikacija)
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users", userHandler.GetUsers)
	router.POST("/verify/:email/:code", userHandler.VerifyUser)
	router.POST("/login", userHandler.Login)
	router.DELETE("/users/:id", userHandler.DeleteUserByID)
	router.POST("/users/change-password", userHandler.ChangePassword)

	// Middleware za zaštitu ruta
	authMiddleware := middleware.JWTAuth(jwtService)

	// Zaštićene rute
	authRoutes := router.Group("/")
	authRoutes.Use(authMiddleware)
	{

		authRoutes.GET("/users/:id", userHandler.GetUserByID)

	}
	// Konfiguracija CORS-a
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{os.Getenv("CORS_ALLOWED_ORIGINS")}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	srv := &http.Server{
		Handler: corsHandler(router),
		Addr:    ":" + port,
	}

	log.Println("Server is running on port " + port)
	log.Fatal(srv.ListenAndServe())
}

// Funkcija za povezivanje na MongoDB
func connectToMongoDB() (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}
