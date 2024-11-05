package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// API rute za zadatke
	router.POST("/tasks", createTask)
	router.GET("/tasks", getTasks)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)

	//router.Run(":8082")

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

// Struktura za zadatak
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      string             `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
}

// Handler za kreiranje novog zadatka
func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	task.CreatedAt = time.Now()

	result, err := taskCollection.InsertOne(context.TODO(), task)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"taskId": result.InsertedID})
}

// Handler za dohvatanje svih zadataka
func getTasks(c *gin.Context) {
	cursor, err := taskCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var tasks []Task
	if err = cursor.All(context.TODO(), &tasks); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tasks)
}

// Handler za ažuriranje zadatka
func updateTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid task ID"})
		return
	}

	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{"$set": task}
	_, err = taskCollection.UpdateByID(context.TODO(), id, update)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Task updated"})
}

// Handler za brisanje zadatka
func deleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid task ID"})
		return
	}

	_, err = taskCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Task deleted"})
}
