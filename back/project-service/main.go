// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/milly013/trello-project/back/project-service/handler"
	"github.com/milly013/trello-project/back/project-service/repository"
	"github.com/milly013/trello-project/back/project-service/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Povezivanje na MongoDB
	client, err := connectToMongoDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())
	db := client.Database("mydatabase")

	// Inicijalizacija repozitorijuma, servisa i handlera
	projectRepo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	// Kreiranje Gin routera
	router := gin.Default()

	// CORS konfiguracija
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Definisanje ruta
	router.POST("/projects", projectHandler.CreateProject)
	router.GET("/projects", projectHandler.GetProjects)

	// Pokretanje servera
	router.Run(":8081")
}

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
