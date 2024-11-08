package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
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
	projectCollection = client.Database("mydatabase").Collection("projects")

	router := gin.Default()

	// API rute za projekte
	router.GET("/projects", getProjects)

	router.Run(":8080")
}

// Funkcija za povezivanje na MongoDB
func connectToMongoDB() (*mongo.Client, error) {
	// Opcije konekcije
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Povezujemo se na MongoDB server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Koristimo Connect umesto NewClient
	client, err := mongo.Connect(ctx, clientOptions)
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

// Handler za GET /projects
func getProjects(c *gin.Context) {
	// Dummy podaci za testiranje, kasnije ćemo ih zameniti podacima iz baze
	c.JSON(200, gin.H{"message": "Retrieve projects from MongoDB here"})
}
