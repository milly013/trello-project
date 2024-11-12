package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
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

	taskRepo := repository.NewTaskRepository(taskCollection)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	// API rute za zadatke
	router.POST("/tasks", taskHandler.CreateTask)
	router.GET("/tasks", taskHandler.GetTasks)
	router.PUT("/tasks/:id", taskHandler.UpdateTask)
	router.GET("/tasks/:id", taskHandler.GetTaskById)
	//router.DELETE("/tasks/:id", taskHandler.DeleteTask)
	router.POST("/tasks/add-member", taskHandler.AssignMemberToTask)

	router.Run(":8082")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}), // Set the correct origin
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	srv := &http.Server{

		Handler: corsHandler(router),
		Addr:    ":8082",
	}

	log.Println("Server is running on port 8082")
	log.Fatal(srv.ListenAndServe())

}

// Funkcija za povezivanje na MongoDB
func connectToMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
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
