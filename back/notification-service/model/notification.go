package model

import (
	"time"

	"github.com/gocql/gocql"
)

type Notification struct {
	ID        gocql.UUID `json:"id"`         // Jedinstveni identifikator za svako obaveštenje
	UserID    string     `json:"user_id"`    // ID korisnika kome je namenjeno obaveštenje
	Type      string     `json:"type"`       // Tip obaveštenja (npr. "added_to_project", "task_status_changed")
	Message   string     `json:"message"`    // Poruka koja se šalje korisniku
	CreatedAt time.Time  `json:"created_at"` // Datum i vreme kreiranja obaveštenja
	IsRead    bool       `json:"is_read"`    // Da li je obaveštenje pročitano
}
