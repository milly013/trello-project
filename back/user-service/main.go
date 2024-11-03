package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/milly013/trello-project/back/user-service/model"

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
	userRepo := model.NewUserRepository(client, "mydatabase")

	// Kreirajte UserHandler
	userHandler := model.NewUserHandler(userRepo)

	router := gin.Default()

	// API rute za korisnike
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users", userHandler.GetUsers)

	router.Run(":8081")
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
