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
var projectCollection *mongo.Collection

func main() {
	// Povezivanje na MongoDB
	client, err := connectToMongoDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Referenca na kolekcije
	taskCollection = client.Database("mydatabase").Collection("tasks")
	projectCollection = client.Database("mydatabase").Collection("projects")

	router := gin.Default()

	// API rute za zadatke
	router.POST("/tasks", createTask)
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:taskId", getTaskById)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)
	router.POST("/tasks/add-member", addMemberToTask)

	// CORS podešavanje
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}), // Set the correct origin
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Pokretanje servera sa CORS middleware-om
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
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title       string               `bson:"title" json:"title"`
	Description string               `bson:"description" json:"description"`
	Status      string               `bson:"status" json:"status"`
	CreatedAt   time.Time            `bson:"createdAt" json:"createdAt"`
	AssignedTo  []primitive.ObjectID `bson:"assignedTo" json:"assignedTo"` // Dodajemo AssignedTo za članove zadatka
}

// Struktura za zahtev za dodavanje člana na zadatak
type AddMemberRequest struct {
	TaskID   string             `json:"taskId"`
	MemberID primitive.ObjectID `json:"memberId"`
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

// Handler za dohvatanje zadatka po ID-u
func getTaskById(c *gin.Context) {
	taskId := c.Param("taskId")
	id, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid task ID"})
		return
	}

	var task Task
	err = taskCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&task)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, task)
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

// Handler za dodavanje člana na zadatak
func addMemberToTask(c *gin.Context) {
	var request AddMemberRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validacija da li zadatak postoji
	taskID, err := primitive.ObjectIDFromHex(request.TaskID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid task ID"})
		return
	}

	// Provera da li je član već u projektu
	var task Task
	err = taskCollection.FindOne(context.TODO(), bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		c.JSON(500, gin.H{"error": "Task not found"})
		return
	}

	// Dodavanje člana u zadatak (ako nije već prisutan)
	for _, memberID := range task.AssignedTo {
		if memberID == request.MemberID {
			c.JSON(400, gin.H{"error": "Member is already assigned to this task"})
			return
		}
	}

	// Dodavanje člana na zadatak
	_, err = taskCollection.UpdateByID(context.TODO(), taskID, bson.M{"$push": bson.M{"assignedTo": request.MemberID}})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Member added to task successfully"})
}
