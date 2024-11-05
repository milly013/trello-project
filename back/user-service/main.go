package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"

	"github.com/milly013/trello-project/back/user-service/handler"
	"github.com/milly013/trello-project/back/user-service/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var projectCollection *mongo.Collection

func main() {
	// Povezivanje na MongoDB
	client, err := connectToMongoDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Referenca na kolekciju
	projectCollection = client.Database("mydatabase").Collection("users")

	// Kreirajte repozitorijum
	userRepo := repository.NewUserRepository(client, "mydatabase")

	// Kreirajte UserHandler
	userHandler := handler.NewUserHandler(userRepo)

	router := gin.Default()

	// API rute za korisnike
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users", userHandler.GetUsers)
	router.GET("/users/:id", userHandler.GetUserByID)
	router.DELETE("/users/:id",userHandler) // Dodana ruta za preuzimanje korisnika po ID-u

	//router.Run(":8080")


	corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"http://localhost:4200"}), // Set the correct origin
        handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS","DELETE"}),
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
    )

    srv := &http.Server{

        Handler: corsHandler(router),
        Addr:    ":8080",
    }

	log.Println("Server is running on port 8080")
    log.Fatal(srv.ListenAndServe())
}

// Funkcija za povezivanje na MongoDB
func connectToMongoDB() (*mongo.Client, error) {
	// Opcije konekcije
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Kreiramo novi MongoDB klijent
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	// Povezujemo se na MongoDB server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Testiramo konekciju
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}
