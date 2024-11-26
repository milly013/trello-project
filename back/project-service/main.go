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

	"github.com/milly013/trello-project/back/project-service/handler"
	"github.com/milly013/trello-project/back/project-service/repository"
	"github.com/milly013/trello-project/back/project-service/service"
)

var projectCollection *mongo.Collection

func main() {
	// Učitajte .env fajl
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Povezivanje na MongoDB
	client, err := connectToMongoDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Kreiramo instancu baze podataka
	db := client.Database(os.Getenv("MONGODB_DATABASE"))
	projectCollection = db.Collection("projects")

	projectRepo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	router := gin.Default()

	// Definisanje ruta
	router.POST("/projects", projectHandler.CreateProject)
	router.GET("/projects", projectHandler.GetProjects)
	router.GET("/projects/:id", projectHandler.GetProjectByID)
	router.POST("/projects/:projectId/members", projectHandler.AddMemberToProject)
	router.DELETE("/projects/:projectId/members", projectHandler.RemoveMemberFromProject)
	router.POST("/projects/:projectId/tasks", projectHandler.AddTaskToProject)
	router.GET("/projects/:id/tasks", projectHandler.GetTaskIDsByProject)
	router.GET("/projects/manager/:managerId", projectHandler.GetProjectsByManager)
	router.GET("/projects/member/:memberId", projectHandler.GetProjectsByMember)

	// Konfiguracija CORS-a
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{os.Getenv("CORS_ALLOWED_ORIGINS")}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Default port
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
	dbName := os.Getenv("MONGODB_DATABASE")
	if mongoURI == "" || dbName == "" {
		log.Fatal("MongoDB URI or Database environment variable not set")
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

	fmt.Printf("Connected to MongoDB database %s!\n", dbName)
	return client, nil
}
