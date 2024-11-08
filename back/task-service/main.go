package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/milly013/trello-project/back/task-service/handler"
	"github.com/milly013/trello-project/back/task-service/repository"
	"github.com/milly013/trello-project/back/task-service/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var taskCollection *mongo.Collection

func main() {
	// Povezivanje na MongoDB
	client, err := connectToMongoDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Referenca na kolekciju
	taskCollection = client.Database("mydatabase").Collection("tasks")

	router := gin.Default()

	// CORS konfiguracija
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},            // Dozvoljava sve origene, možete ograničiti na specifične URL-ove
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},     // Dozvoljeni HTTP metodi
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"}, // Dozvoljeni HTTP headeri
		AllowCredentials: true,
	}))

	taskRepo := repository.NewTaskRepository(taskCollection)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	// API rute za zadatke
	router.POST("/tasks", taskHandler.CreateTask)
	router.GET("/tasks", taskHandler.GetTasks)
	// router.PUT("/tasks/:id", updateTask)
	// router.DELETE("/tasks/:id", deleteTask)

	router.Run(":8082")
}

func connectToMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Pokušajte da se povežete bez ponovnog povezivanja ako je već povezano
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Povezivanje sa MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Povežite se na MongoDB ako to nije već urađeno
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	// Pokušajte da pingujete bazu da biste proverili konekciju
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}
