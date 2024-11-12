package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`        // ID zadatka
	ProjectID   primitive.ObjectID   `bson:"projectId" json:"projectId"`     // ID projekta kojem pripada
	Title       string               `bson:"title" json:"title"`             // Naziv zadatka
	Description string               `bson:"description" json:"description"` // Opis zadatka
	StartDate   time.Time            `bson:"startDate" json:"startDate"`     // Datum početka zadatka
	EndDate     time.Time            `bson:"endDate" json:"endDate"`         // Datum završetka zadatka
	AssignedTo  []primitive.ObjectID `bson:"assignedTo" json:"assignedTo"`   // Lista ID-ova korisnika kojima je zadatak dodeljen
	Status      string               `bson:"status" json:"status"`           // Status zadatka (npr. "Pending", "InProgress", "Completed")
	CreatedAt   time.Time            `bson:"createdAt" json:"createdAt"`     // Datum kreiranja zadatka
}
