package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/gorilla/handlers"
	"github.com/milly013/trello-project/back/notification-service/handler"
	"github.com/milly013/trello-project/back/notification-service/repository"
	"github.com/milly013/trello-project/back/notification-service/service"
)

var session *gocql.Session

func main() {

	var err error
	session, err = connectToCassandra()
	if err != nil {
		log.Fatal("Error connecting to Cassandra:", err)
	}
	defer session.Close()

	// Poziv funkcije za kreiranje tabele
	if err := createTableIfNotExists(session); err != nil {
		log.Fatal("Error creating Cassandra table:", err)
	}

	notificationRepo := repository.NewNotificationRepository(session)
	notificationService := service.NewNotificationService(notificationRepo)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	router := gin.Default()
	router.Use(CORSMiddleware())

	// Definisanje ruta
	router.POST("/notifications", notificationHandler.CreateNotification)
	router.GET("/notifications", notificationHandler.GetAllNotifications)
	router.GET("/notifications/:userID", notificationHandler.GetNotificationsByUserID)
	router.PUT("/notifications/:notificationID/read", notificationHandler.MarkNotificationAsRead)

	// Konfiguracija CORS-a
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{os.Getenv("CORS_ALLOWED_ORIGINS")}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" // Default port for notification-service
	}

	srv := &http.Server{
		Handler: corsHandler(router),
		Addr:    ":" + port,
	}

	log.Println("Notification service is running on port " + port)
	log.Fatal(srv.ListenAndServe())
}

// Funkcija za povezivanje na Cassandra bazu podataka
func connectToCassandra() (*gocql.Session, error) {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))
	cluster.Keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = 10 * time.Second

	var session *gocql.Session
	var err error
	for i := 0; i < 5; i++ { // Pokušaj 5 puta
		session, err = cluster.CreateSession()
		if err == nil {
			log.Printf("Connected to Cassandra keyspace %s!\n", cluster.Keyspace)
			return session, nil
		}
		log.Printf("Failed to connect to Cassandra, retrying in 5 seconds... (%d/5)\n", i+1)
		time.Sleep(5 * time.Second)
	}

	return nil, err
}

func createTableIfNotExists(session *gocql.Session) error {
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		id UUID,
		user_id TEXT,
		type TEXT,
		message TEXT,
		created_at TIMESTAMP,
		is_read BOOLEAN,
		PRIMARY KEY ((user_id), created_at)
	) WITH CLUSTERING ORDER BY (created_at DESC);
	`
	return session.Query(query).Exec()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://localhost:4200")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Ako je preflight (OPTIONS) zahtev
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
