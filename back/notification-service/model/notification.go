package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // Jedinstveni identifikator za svako obaveštenje
	UserID    primitive.ObjectID `bson:"user_id"`       // ID korisnika kome je namenjeno obaveštenje
	Type      string             `bson:"type"`          // Tip obaveštenja (npr. "added_to_project", "task_status_changed")
	Message   string             `bson:"message"`       // Poruka koja se šalje korisniku
	CreatedAt time.Time          `bson:"created_at"`    // Datum i vreme kreiranja obaveštenja
	IsRead    bool               `bson:"is_read"`       // Da li je obaveštenje pročitano
}
